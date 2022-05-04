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
	"errors"
	"fmt"

	v1 "github.com/bhojpur/state/pkg/api/v1/database"
	db "github.com/bhojpur/state/pkg/database"
)

var errBatchClosed = errors.New("batch has been written or closed")

type batch struct {
	db  *RemoteDB
	ops []*v1.Operation
}

var _ db.Batch = (*batch)(nil)

func newBatch(rdb *RemoteDB) *batch {
	return &batch{
		db:  rdb,
		ops: []*v1.Operation{},
	}
}

// Set implements Batch.
func (b *batch) Set(key, value []byte) error {
	if b.ops == nil {
		return errBatchClosed
	}
	op := &v1.Operation{
		Entity: &v1.Entity{Key: key, Value: value},
		Type:   v1.Operation_SET,
	}
	b.ops = append(b.ops, op)
	return nil
}

// Delete implements Batch.
func (b *batch) Delete(key []byte) error {
	if b.ops == nil {
		return errBatchClosed
	}
	op := &v1.Operation{
		Entity: &v1.Entity{Key: key},
		Type:   v1.Operation_DELETE,
	}
	b.ops = append(b.ops, op)
	return nil
}

// Write implements Batch.
func (b *batch) Write() error {
	if b.ops == nil {
		return errBatchClosed
	}
	_, err := b.db.dc.BatchWrite(b.db.ctx, &v1.Batch{Ops: b.ops})
	if err != nil {
		return fmt.Errorf("remoteDB.BatchWrite: %w", err)
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// WriteSync implements Batch.
func (b *batch) WriteSync() error {
	if b.ops == nil {
		return errBatchClosed
	}
	_, err := b.db.dc.BatchWriteSync(b.db.ctx, &v1.Batch{Ops: b.ops})
	if err != nil {
		return fmt.Errorf("RemoteDB.BatchWriteSync: %w", err)
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	return b.Close()
}

// Close implements Batch.
func (b *batch) Close() error {
	b.ops = nil
	return nil
}
