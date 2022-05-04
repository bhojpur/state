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
	"io"

	"github.com/gogo/protobuf/proto"
)

type Writer interface {
	WriteMsg(proto.Message) (int, error)
}

type WriteCloser interface {
	Writer
	io.Closer
}

type Reader interface {
	ReadMsg(msg proto.Message) (int, error)
}

type ReadCloser interface {
	Reader
	io.Closer
}

type marshaler interface {
	MarshalTo(data []byte) (n int, err error)
}

func getSize(v interface{}) (int, bool) {
	if sz, ok := v.(interface {
		Size() (n int)
	}); ok {
		return sz.Size(), true
	} else if sz, ok := v.(interface {
		ProtoSize() (n int)
	}); ok {
		return sz.ProtoSize(), true
	} else {
		return 0, false
	}
}

// byteReader wraps an io.Reader and implements io.ByteReader, required by
// binary.ReadUvarint(). Reading one byte at a time is extremely slow, but this
// is what Amino did previously anyway, and the caller can wrap the underlying
// reader in a bufio.Reader if appropriate.
type byteReader struct {
	reader    io.Reader
	buf       []byte
	bytesRead int // keeps track of bytes read via ReadByte()
}

func newByteReader(r io.Reader) *byteReader {
	return &byteReader{
		reader: r,
		buf:    make([]byte, 1),
	}
}

func (r *byteReader) ReadByte() (byte, error) {
	n, err := r.reader.Read(r.buf)
	r.bytesRead += n
	if err != nil {
		return 0x00, err
	}
	return r.buf[0], nil
}
