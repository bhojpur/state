//go:build gofuzz || go1.18

package tests

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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhojpur/state/pkg/libs/log"
	rpcserver "github.com/bhojpur/state/pkg/rpc/jsonrpc/server"
	"github.com/bhojpur/state/pkg/rpc/jsonrpc/types"
)

func FuzzRPCJSONRPCServer(f *testing.F) {
	type args struct {
		S string `json:"s"`
		I int    `json:"i"`
	}
	var rpcFuncMap = map[string]*rpcserver.RPCFunc{
		"c": rpcserver.NewRPCFunc(func(context.Context, *args) (string, error) {
			return "foo", nil
		}),
	}

	mux := http.NewServeMux()
	rpcserver.RegisterRPCFuncs(mux, rpcFuncMap, log.NewNopLogger())
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}

		req, err := http.NewRequest("POST", "http://localhost/", bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		res := rec.Result()
		blob, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
		if len(blob) == 0 {
			return
		}

		if outputJSONIsSlice(blob) {
			var recv []types.RPCResponse
			if err := json.Unmarshal(blob, &recv); err != nil {
				panic(err)
			}
			return
		}
		var recv types.RPCResponse
		if err := json.Unmarshal(blob, &recv); err != nil {
			panic(err)
		}
	})
}

func outputJSONIsSlice(input []byte) bool {
	var slice []json.RawMessage
	return json.Unmarshal(input, &slice) == nil
}
