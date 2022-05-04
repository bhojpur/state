package mock

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

	"github.com/bhojpur/state/pkg/rpc/client"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
)

// StatusMock returns the result specified by the Call
type StatusMock struct {
	Call
}

var (
	_ client.StatusClient = (*StatusMock)(nil)
	_ client.StatusClient = (*StatusRecorder)(nil)
)

func (m *StatusMock) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	res, err := m.GetResponse(nil)
	if err != nil {
		return nil, err
	}
	return res.(*coretypes.ResultStatus), nil
}

// StatusRecorder can wrap another type (StatusMock, full client)
// and record the status calls
type StatusRecorder struct {
	Client client.StatusClient
	Calls  []Call
}

func NewStatusRecorder(client client.StatusClient) *StatusRecorder {
	return &StatusRecorder{
		Client: client,
		Calls:  []Call{},
	}
}

func (r *StatusRecorder) addCall(call Call) {
	r.Calls = append(r.Calls, call)
}

func (r *StatusRecorder) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	res, err := r.Client.Status(ctx)
	r.addCall(Call{
		Name:     "status",
		Response: res,
		Error:    err,
	})
	return res, err
}
