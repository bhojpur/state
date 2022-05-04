package remote

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
	"errors"
	"fmt"

	v1 "github.com/bhojpur/state/pkg/api/v1/database"
	db "github.com/bhojpur/state/pkg/database"
	"github.com/bhojpur/state/pkg/database/remote/rpc"
)

type RemoteDB struct {
	ctx context.Context
	dc  v1.DBClient
}

func NewRemoteDB(serverAddr string, serverKey string) (*RemoteDB, error) {
	return newRemoteDB(rpc.NewClient(serverAddr, serverKey))
}

func newRemoteDB(gdc v1.DBClient, err error) (*RemoteDB, error) {
	if err != nil {
		return nil, err
	}
	return &RemoteDB{dc: gdc, ctx: context.Background()}, nil
}

type Init struct {
	Dir  string
	Name string
	Type string
}

func (rd *RemoteDB) InitRemote(in *Init) error {
	_, err := rd.dc.Init(rd.ctx, &v1.Init{Dir: in.Dir, Type: in.Type, Name: in.Name})
	return err
}

var _ db.DB = (*RemoteDB)(nil)

// Close is a noop currently
func (rd *RemoteDB) Close() error {
	return nil
}

func (rd *RemoteDB) Delete(key []byte) error {
	if _, err := rd.dc.Delete(rd.ctx, &v1.Entity{Key: key}); err != nil {
		return fmt.Errorf("remoteDB.Delete: %w", err)
	}
	return nil
}

func (rd *RemoteDB) DeleteSync(key []byte) error {
	if _, err := rd.dc.DeleteSync(rd.ctx, &v1.Entity{Key: key}); err != nil {
		return fmt.Errorf("remoteDB.DeleteSync: %w", err)
	}
	return nil
}

func (rd *RemoteDB) Set(key, value []byte) error {
	if _, err := rd.dc.Set(rd.ctx, &v1.Entity{Key: key, Value: value}); err != nil {
		return fmt.Errorf("remoteDB.Set: %w", err)
	}
	return nil
}

func (rd *RemoteDB) SetSync(key, value []byte) error {
	if _, err := rd.dc.SetSync(rd.ctx, &v1.Entity{Key: key, Value: value}); err != nil {
		return fmt.Errorf("remoteDB.SetSync: %w", err)
	}
	return nil
}

func (rd *RemoteDB) Get(key []byte) ([]byte, error) {
	res, err := rd.dc.Get(rd.ctx, &v1.Entity{Key: key})
	if err != nil {
		return nil, fmt.Errorf("remoteDB.Get error: %w", err)
	}
	return res.Value, nil
}

func (rd *RemoteDB) Has(key []byte) (bool, error) {
	res, err := rd.dc.Has(rd.ctx, &v1.Entity{Key: key})
	if err != nil {
		return false, err
	}
	return res.Exists, nil
}

func (rd *RemoteDB) ReverseIterator(start, end []byte) (db.Iterator, error) {
	dic, err := rd.dc.ReverseIterator(rd.ctx, &v1.Entity{Start: start, End: end})
	if err != nil {
		return nil, fmt.Errorf("RemoteDB.Iterator error: %w", err)
	}
	return makeReverseIterator(dic), nil
}

func (rd *RemoteDB) NewBatch() db.Batch {
	return newBatch(rd)
}

// TODO: Implement Print when db.DB implements a method
// to print to a string and not db.Print to stdout.
func (rd *RemoteDB) Print() error {
	return errors.New("remoteDB.Print: unimplemented")
}

func (rd *RemoteDB) Stats() map[string]string {
	stats, err := rd.dc.Stats(rd.ctx, &v1.Nothing{})
	if err != nil || stats == nil {
		return nil
	}
	return stats.Data
}

func (rd *RemoteDB) Iterator(start, end []byte) (db.Iterator, error) {
	dic, err := rd.dc.Iterator(rd.ctx, &v1.Entity{Start: start, End: end})
	if err != nil {
		return nil, fmt.Errorf("RemoteDB.Iterator error: %w", err)
	}
	return makeIterator(dic), nil
}
