package rpc

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	v1 "github.com/bhojpur/state/pkg/api/v1/database"
	db "github.com/bhojpur/state/pkg/database"
)

// ListenAndServe is a blocking function that sets up a gRPC based
// server at the address supplied, with the gRPC options passed in.
// Normally in usage, invoke it in a goroutine like you would for http.ListenAndServe.
func ListenAndServe(addr, cert, key string, opts ...grpc.ServerOption) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv, err := NewServer(cert, key, opts...)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func NewServer(cert, key string, opts ...grpc.ServerOption) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		return nil, err
	}
	opts = append(opts, grpc.Creds(creds))
	srv := grpc.NewServer(opts...)
	v1.RegisterDBServer(srv, new(server))
	return srv, nil
}

type server struct {
	v1.UnimplementedDBServer
	mu sync.Mutex
	db db.DB
}

var _ v1.DBServer = (*server)(nil)

// Init initializes the server's database. Only one type of database
// can be initialized per server.
//
// Dir is the directory on the file system in which the DB will be stored(if backed by disk) (TODO: remove)
//
// Name is representative filesystem entry's basepath
//
// Type can be either one of:
//  * cleveldb (if built with gcc enabled)
//  * fsdb
//  * memdB
//  * goleveldb
func (s *server) Init(ctx context.Context, in *v1.Init) (*v1.Entity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	s.db, err = db.NewDB(in.Name, db.BackendType(in.Type), in.Dir)
	if err != nil {
		return nil, err
	}
	return &v1.Entity{CreatedAt: time.Now().Unix()}, nil
}

func (s *server) Delete(ctx context.Context, in *v1.Entity) (*v1.Nothing, error) {
	err := s.db.Delete(in.Key)
	if err != nil {
		return nil, err
	}
	return nothing, nil
}

var nothing = new(v1.Nothing)

func (s *server) DeleteSync(ctx context.Context, in *v1.Entity) (*v1.Nothing, error) {
	err := s.db.DeleteSync(in.Key)
	if err != nil {
		return nil, err
	}
	return nothing, nil
}

func (s *server) Get(ctx context.Context, in *v1.Entity) (*v1.Entity, error) {
	value, err := s.db.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &v1.Entity{Value: value}, nil
}

func (s *server) GetStream(ds v1.DB_GetStreamServer) error {
	// Receive routine
	responsesChan := make(chan *v1.Entity)
	go func() {
		defer close(responsesChan)
		ctx := context.Background()
		for {
			in, err := ds.Recv()
			if err != nil {
				responsesChan <- &v1.Entity{Err: err.Error()}
				return
			}
			out, err := s.Get(ctx, in)
			if err != nil {
				if out == nil {
					out = new(v1.Entity)
					out.Key = in.Key
				}
				out.Err = err.Error()
				responsesChan <- out
				return
			}

			// Otherwise continue on
			responsesChan <- out
		}
	}()

	// Send routine, block until we return
	for out := range responsesChan {
		if err := ds.Send(out); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Has(ctx context.Context, in *v1.Entity) (*v1.Entity, error) {
	exists, err := s.db.Has(in.Key)
	if err != nil {
		return nil, err
	}
	return &v1.Entity{Exists: exists}, nil
}

func (s *server) Set(ctx context.Context, in *v1.Entity) (*v1.Nothing, error) {
	err := s.db.Set(in.Key, in.Value)
	if err != nil {
		return nil, err
	}
	return nothing, nil
}

func (s *server) SetSync(ctx context.Context, in *v1.Entity) (*v1.Nothing, error) {
	err := s.db.SetSync(in.Key, in.Value)
	if err != nil {
		return nil, err
	}
	return nothing, nil
}

func (s *server) Iterator(query *v1.Entity, dis v1.DB_IteratorServer) error {
	it, err := s.db.Iterator(query.Start, query.End)
	if err != nil {
		return err
	}
	defer it.Close()
	return s.handleIterator(it, dis.Send)
}

func (s *server) handleIterator(it db.Iterator, sendFunc func(*v1.Iterator) error) error {
	for it.Valid() {
		start, end := it.Domain()
		key := it.Key()
		value := it.Value()

		out := &v1.Iterator{
			Domain: &v1.Domain{Start: start, End: end},
			Valid:  it.Valid(),
			Key:    key,
			Value:  value,
		}
		if err := sendFunc(out); err != nil {
			return err
		}

		// Finally move the iterator forward,
		it.Next()

	}
	return nil
}

func (s *server) ReverseIterator(query *v1.Entity, dis v1.DB_ReverseIteratorServer) error {
	it, err := s.db.ReverseIterator(query.Start, query.End)
	if err != nil {
		return err
	}
	defer it.Close()
	return s.handleIterator(it, dis.Send)
}

func (s *server) Stats(context.Context, *v1.Nothing) (*v1.Stats, error) {
	stats := s.db.Stats()
	return &v1.Stats{Data: stats, TimeAt: time.Now().Unix()}, nil
}

func (s *server) BatchWrite(c context.Context, b *v1.Batch) (*v1.Nothing, error) {
	return s.batchWrite(c, b, false)
}

func (s *server) BatchWriteSync(c context.Context, b *v1.Batch) (*v1.Nothing, error) {
	return s.batchWrite(c, b, true)
}

func (s *server) batchWrite(c context.Context, b *v1.Batch, sync bool) (*v1.Nothing, error) {
	bat := s.db.NewBatch()
	defer bat.Close()
	for _, op := range b.Ops {
		switch op.Type {
		case v1.Operation_SET:
			err := bat.Set(op.Entity.Key, op.Entity.Value)
			if err != nil {
				return nil, err
			}
		case v1.Operation_DELETE:
			err := bat.Delete(op.Entity.Key)
			if err != nil {
				return nil, err
			}
		}
	}
	if sync {
		err := bat.WriteSync()
		if err != nil {
			return nil, err
		}
	} else {
		err := bat.Write()
		if err != nil {
			return nil, err
		}
	}
	return nothing, nil
}
