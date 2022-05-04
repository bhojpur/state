package types

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
	"encoding/json"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"

	typespb "github.com/bhojpur/state/pkg/api/v1/types"
)

func TestMarshalJSON(t *testing.T) {
	b, err := json.Marshal(&ExecTxResult{Code: 1})
	assert.NoError(t, err)
	// include empty fields.
	assert.True(t, strings.Contains(string(b), "code"))
	r1 := ResponseCheckTx{
		Code:      1,
		Data:      []byte("hello"),
		GasWanted: 43,
		Events: []Event{
			{
				Type: "testEvent",
				Attributes: []EventAttribute{
					{Key: "pho", Value: "bo"},
				},
			},
		},
	}
	b, err = json.Marshal(&r1)
	assert.NoError(t, err)

	var r2 ResponseCheckTx
	err = json.Unmarshal(b, &r2)
	assert.NoError(t, err)
	assert.Equal(t, r1, r2)
}

func TestWriteReadMessageSimple(t *testing.T) {
	cases := []proto.Message{
		&RequestEcho{
			Message: "Hello",
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		err := WriteMessage(c, buf)
		assert.NoError(t, err)

		msg := new(RequestEcho)
		err = ReadMessage(buf, msg)
		assert.NoError(t, err)

		assert.True(t, proto.Equal(c, msg))
	}
}

func TestWriteReadMessage(t *testing.T) {
	cases := []proto.Message{
		&typespb.Header{
			Height:  4,
			ChainID: "test",
		},
		// TODO: add the rest
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		err := WriteMessage(c, buf)
		assert.NoError(t, err)

		msg := new(typespb.Header)
		err = ReadMessage(buf, msg)
		assert.NoError(t, err)

		assert.True(t, proto.Equal(c, msg))
	}
}

func TestWriteReadMessage2(t *testing.T) {
	phrase := "hello-world"
	cases := []proto.Message{
		&ResponseCheckTx{
			Data:      []byte(phrase),
			Log:       phrase,
			GasWanted: 10,
			Events: []Event{
				{
					Type: "testEvent",
					Attributes: []EventAttribute{
						{Key: "abc", Value: "def"},
					},
				},
			},
		},
		// TODO: add the rest
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		err := WriteMessage(c, buf)
		assert.NoError(t, err)

		msg := new(ResponseCheckTx)
		err = ReadMessage(buf, msg)
		assert.NoError(t, err)

		assert.True(t, proto.Equal(c, msg))
	}
}
