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
	"io"
	"sync"

	"github.com/gogo/protobuf/proto"
)

// NewDelimitedWriter writes a varint-delimited Protobuf message to a writer. It is
// equivalent to the gogoproto NewDelimitedWriter, except WriteMsg() also returns the
// number of bytes written, which is necessary in the p2p package.
func NewDelimitedWriter(w io.Writer) WriteCloser {
	return &varintWriter{w, make([]byte, binary.MaxVarintLen64), nil}
}

type varintWriter struct {
	w      io.Writer
	lenBuf []byte
	buffer []byte
}

func (w *varintWriter) WriteMsg(msg proto.Message) (int, error) {
	if m, ok := msg.(marshaler); ok {
		n, ok := getSize(m)
		if ok {
			if n+binary.MaxVarintLen64 >= len(w.buffer) {
				w.buffer = make([]byte, n+binary.MaxVarintLen64)
			}
			lenOff := binary.PutUvarint(w.buffer, uint64(n))
			_, err := m.MarshalTo(w.buffer[lenOff:])
			if err != nil {
				return 0, err
			}
			_, err = w.w.Write(w.buffer[:lenOff+n])
			return lenOff + n, err
		}
	}

	// fallback
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, err
	}
	length := uint64(len(data))
	n := binary.PutUvarint(w.lenBuf, length)
	_, err = w.w.Write(w.lenBuf[:n])
	if err != nil {
		return 0, err
	}
	_, err = w.w.Write(data)
	return len(data) + n, err
}

func (w *varintWriter) Close() error {
	if closer, ok := w.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func varintWrittenBytes(m marshaler, size int) ([]byte, error) {
	buf := make([]byte, size+binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(size))
	nw, err := m.MarshalTo(buf[n:])
	if err != nil {
		return nil, err
	}
	return buf[:n+nw], nil
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func MarshalDelimited(msg proto.Message) ([]byte, error) {
	// The goal here is to write proto message as is knowning already if
	// the exact size can be retrieved and if so just use that.
	if m, ok := msg.(marshaler); ok {
		size, ok := getSize(msg)
		if ok {
			return varintWrittenBytes(m, size)
		}
	}

	// Otherwise, go down the route of using proto.Marshal,
	// and use the buffer pool to retrieve a writer.
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()
	_, err := NewDelimitedWriter(buf).WriteMsg(msg)
	if err != nil {
		return nil, err
	}
	// Given that we are reusing buffers, we should
	// make a copy of the returned bytes.
	bytesCopy := make([]byte, buf.Len())
	copy(bytesCopy, buf.Bytes())
	return bytesCopy, nil
}
