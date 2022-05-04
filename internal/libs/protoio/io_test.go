package protoio_test

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
	"math/rand"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/test"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/libs/protoio"
)

func iotest(t *testing.T, writer protoio.WriteCloser, reader protoio.ReadCloser) error {
	t.Helper()
	varint := make([]byte, binary.MaxVarintLen64)
	size := 1000
	msgs := make([]*test.NinOptNative, size)
	lens := make([]int, size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range msgs {
		msgs[i] = test.NewPopulatedNinOptNative(r, true)
		// issue 31
		if i == 5 {
			msgs[i] = &test.NinOptNative{}
		}
		// issue 31
		if i == 999 {
			msgs[i] = &test.NinOptNative{}
		}
		// FIXME Check size
		bz, err := proto.Marshal(msgs[i])
		if err != nil {
			return err
		}
		visize := binary.PutUvarint(varint, uint64(len(bz)))
		n, err := writer.WriteMsg(msgs[i])
		if err != nil {
			return err
		}
		if n != len(bz)+visize {
			return fmt.Errorf("WriteMsg() wrote %v bytes, expected %v", n, len(bz)+visize)
		}
		lens[i] = n
	}
	if err := writer.Close(); err != nil {
		return err
	}
	i := 0
	for {
		msg := &test.NinOptNative{}
		if n, err := reader.ReadMsg(msg); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else if n != lens[i] {
			return fmt.Errorf("read %v bytes, expected %v", n, lens[i])
		}
		if err := msg.VerboseEqual(msgs[i]); err != nil {
			return err
		}
		i++
	}
	require.Equal(t, size, i, "messages read ≠ messages written")
	if err := reader.Close(); err != nil {
		return err
	}
	return nil
}

type buffer struct {
	*bytes.Buffer
	closed bool
}

func (b *buffer) Close() error {
	b.closed = true
	return nil
}

func newBuffer() *buffer {
	return &buffer{bytes.NewBuffer(nil), false}
}

func TestVarintNormal(t *testing.T) {
	buf := newBuffer()
	writer := protoio.NewDelimitedWriter(buf)
	reader := protoio.NewDelimitedReader(buf, 1024*1024)
	err := iotest(t, writer, reader)
	require.NoError(t, err)
	require.True(t, buf.closed, "did not close buffer")
}

func TestVarintNoClose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := protoio.NewDelimitedWriter(buf)
	reader := protoio.NewDelimitedReader(buf, 1024*1024)
	err := iotest(t, writer, reader)
	require.NoError(t, err)
}

// issue 32
func TestVarintMaxSize(t *testing.T) {
	buf := newBuffer()
	writer := protoio.NewDelimitedWriter(buf)
	reader := protoio.NewDelimitedReader(buf, 20)
	err := iotest(t, writer, reader)
	require.Error(t, err)
}

func TestVarintError(t *testing.T) {
	buf := newBuffer()
	buf.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	reader := protoio.NewDelimitedReader(buf, 1024*1024)
	msg := &test.NinOptNative{}
	n, err := reader.ReadMsg(msg)
	require.Error(t, err)
	require.Equal(t, 10, n)
}

func TestVarintTruncated(t *testing.T) {
	buf := newBuffer()
	buf.Write([]byte{0xff, 0xff})
	reader := protoio.NewDelimitedReader(buf, 1024*1024)
	msg := &test.NinOptNative{}
	n, err := reader.ReadMsg(msg)
	require.Error(t, err)
	require.Equal(t, 2, n)
}

func TestShort(t *testing.T) {
	buf := newBuffer()

	varintBuf := make([]byte, binary.MaxVarintLen64)
	varintLen := binary.PutUvarint(varintBuf, 100)
	_, err := buf.Write(varintBuf[:varintLen])
	require.NoError(t, err)

	bz, err := proto.Marshal(&test.NinOptNative{Field15: []byte{0x01, 0x02, 0x03}})
	require.NoError(t, err)
	buf.Write(bz)

	reader := protoio.NewDelimitedReader(buf, 1024*1024)
	require.NoError(t, err)
	msg := &test.NinOptNative{}
	n, err := reader.ReadMsg(msg)
	require.Error(t, err)
	require.Equal(t, varintLen+len(bz), n)
}
