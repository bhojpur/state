// Code generated by mockery. DO NOT EDIT.

package mocks

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
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/bhojpur/state/pkg/abci/types"
)

// AppConnQuery is an autogenerated mock type for the AppConnQuery type
type AppConnQuery struct {
	mock.Mock
}

// Echo provides a mock function with given fields: _a0, _a1
func (_m *AppConnQuery) Echo(_a0 context.Context, _a1 string) (*types.ResponseEcho, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.ResponseEcho
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.ResponseEcho); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseEcho)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Error provides a mock function with given fields:
func (_m *AppConnQuery) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Info provides a mock function with given fields: _a0, _a1
func (_m *AppConnQuery) Info(_a0 context.Context, _a1 types.RequestInfo) (*types.ResponseInfo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.ResponseInfo
	if rf, ok := ret.Get(0).(func(context.Context, types.RequestInfo) *types.ResponseInfo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.RequestInfo) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: _a0, _a1
func (_m *AppConnQuery) Query(_a0 context.Context, _a1 types.RequestQuery) (*types.ResponseQuery, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.ResponseQuery
	if rf, ok := ret.Get(0).(func(context.Context, types.RequestQuery) *types.ResponseQuery); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseQuery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.RequestQuery) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}