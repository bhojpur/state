package server

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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bhojpur/state/pkg/libs/log"
	rpctypes "github.com/bhojpur/state/pkg/rpc/jsonrpc/types"
)

// uriReqID is a placeholder ID used for GET requests, which do not receive a
// JSON-RPC request ID from the caller.
const uriReqID = -1

// convert from a function name to the http handler
func makeHTTPHandler(rpcFunc *RPCFunc, logger log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := rpctypes.WithCallInfo(req.Context(), &rpctypes.CallInfo{
			HTTPRequest: req,
		})
		args, err := parseURLParams(rpcFunc.args, req)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err.Error())
			return
		}
		jreq := rpctypes.NewRequest(uriReqID)
		result, err := rpcFunc.Call(ctx, args)
		if err == nil {
			writeHTTPResponse(w, logger, jreq.MakeResponse(result))
		} else {
			writeHTTPResponse(w, logger, jreq.MakeError(err))
		}
	}
}

func parseURLParams(args []argInfo, req *http.Request) ([]byte, error) {
	if err := req.ParseForm(); err != nil {
		return nil, fmt.Errorf("invalid HTTP request: %w", err)
	}
	getArg := func(name string) (string, bool) {
		if req.Form.Has(name) {
			return req.Form.Get(name), true
		}
		return "", false
	}

	params := make(map[string]interface{})
	for _, arg := range args {
		v, ok := getArg(arg.name)
		if !ok {
			continue
		}
		if z, err := decodeInteger(v); err == nil {
			params[arg.name] = z
		} else if b, err := strconv.ParseBool(v); err == nil {
			params[arg.name] = b
		} else if lc := strings.ToLower(v); strings.HasPrefix(lc, "0x") {
			dec, err := hex.DecodeString(lc[2:])
			if err != nil {
				return nil, fmt.Errorf("invalid hex string: %w", err)
			} else if len(dec) == 0 {
				return nil, errors.New("invalid empty hex string")
			}
			if arg.isBinary {
				params[arg.name] = dec
			} else {
				params[arg.name] = string(dec)
			}
		} else if isQuotedString(v) {
			var dec string
			if err := json.Unmarshal([]byte(v), &dec); err != nil {
				return nil, fmt.Errorf("invalid quoted string: %w", err)
			}
			if arg.isBinary {
				params[arg.name] = []byte(dec)
			} else {
				params[arg.name] = dec
			}
		} else {
			params[arg.name] = v
		}
	}
	return json.Marshal(params)
}

// isQuotedString reports whether s is enclosed in double quotes.
func isQuotedString(s string) bool {
	return len(s) >= 2 && strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)
}

// decodeInteger decodes s into an int64. If s is "double quoted" the quotes
// are removed; otherwise s must be a base-10 digit string.
func decodeInteger(s string) (int64, error) {
	if isQuotedString(s) {
		s = s[1 : len(s)-1]
	}
	return strconv.ParseInt(s, 10, 64)
}
