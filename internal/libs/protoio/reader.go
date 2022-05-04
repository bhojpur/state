package protoio

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
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
)

// NewDelimitedReader reads varint-delimited Protobuf messages from a reader.
// Unlike the gogoproto NewDelimitedReader, this does not buffer the reader,
// which may cause poor performance but is necessary when only reading single
// messages (e.g. in the p2p package). It also returns the number of bytes
// read, which is necessary for the p2p package.
func NewDelimitedReader(r io.Reader, maxSize int) ReadCloser {
	var closer io.Closer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}
	return &varintReader{r, nil, maxSize, closer}
}

type varintReader struct {
	r       io.Reader
	buf     []byte
	maxSize int
	closer  io.Closer
}

func (r *varintReader) ReadMsg(msg proto.Message) (int, error) {
	// ReadUvarint needs an io.ByteReader, and we also need to keep track of the
	// number of bytes read, so we use our own byteReader. This can't be
	// buffered, so the caller should pass a buffered io.Reader to avoid poor
	// performance.
	byteReader := newByteReader(r.r)
	l, err := binary.ReadUvarint(byteReader)
	n := byteReader.bytesRead
	if err != nil {
		return n, err
	}

	// Make sure length doesn't overflow the native int size (e.g. 32-bit),
	// and that the returned sum of n+length doesn't overflow either.
	length := int(l)
	if l >= uint64(^uint(0)>>1) || length < 0 || n+length < 0 {
		return n, fmt.Errorf("invalid out-of-range message length %v", l)
	}
	if length > r.maxSize {
		return n, fmt.Errorf("message exceeds max size (%v > %v)", length, r.maxSize)
	}

	if len(r.buf) < length {
		r.buf = make([]byte, length)
	}
	buf := r.buf[:length]
	nr, err := io.ReadFull(r.r, buf)
	n += nr
	if err != nil {
		return n, err
	}
	return n, proto.Unmarshal(buf, msg)
}

func (r *varintReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}

func UnmarshalDelimited(data []byte, msg proto.Message) error {
	_, err := NewDelimitedReader(bytes.NewReader(data), len(data)).ReadMsg(msg)
	return err
}
