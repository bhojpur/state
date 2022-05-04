package database

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

type prefixDBBatch struct {
	prefix []byte
	source Batch
}

var _ Batch = (*prefixDBBatch)(nil)

func newPrefixBatch(prefix []byte, source Batch) prefixDBBatch {
	return prefixDBBatch{
		prefix: prefix,
		source: source,
	}
}

// Set implements Batch.
func (pb prefixDBBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	if value == nil {
		return errValueNil
	}
	pkey := append(cp(pb.prefix), key...)
	return pb.source.Set(pkey, value)
}

// Delete implements Batch.
func (pb prefixDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	pkey := append(cp(pb.prefix), key...)
	return pb.source.Delete(pkey)
}

// Write implements Batch.
func (pb prefixDBBatch) Write() error {
	return pb.source.Write()
}

// WriteSync implements Batch.
func (pb prefixDBBatch) WriteSync() error {
	return pb.source.WriteSync()
}

// Close implements Batch.
func (pb prefixDBBatch) Close() error {
	return pb.source.Close()
}
