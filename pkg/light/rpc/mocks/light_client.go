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

	testing "testing"

	time "time"

	types "github.com/bhojpur/state/pkg/types"
)

// LightClient is an autogenerated mock type for the LightClient type
type LightClient struct {
	mock.Mock
}

// ChainID provides a mock function with given fields:
func (_m *LightClient) ChainID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Status provides a mock function with given fields: ctx
func (_m *LightClient) Status(ctx context.Context) *types.LightClientInfo {
	ret := _m.Called(ctx)

	var r0 *types.LightClientInfo
	if rf, ok := ret.Get(0).(func(context.Context) *types.LightClientInfo); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LightClientInfo)
		}
	}

	return r0
}

// TrustedLightBlock provides a mock function with given fields: height
func (_m *LightClient) TrustedLightBlock(height int64) (*types.LightBlock, error) {
	ret := _m.Called(height)

	var r0 *types.LightBlock
	if rf, ok := ret.Get(0).(func(int64) *types.LightBlock); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LightBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, now
func (_m *LightClient) Update(ctx context.Context, now time.Time) (*types.LightBlock, error) {
	ret := _m.Called(ctx, now)

	var r0 *types.LightBlock
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) *types.LightBlock); ok {
		r0 = rf(ctx, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LightBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time) error); ok {
		r1 = rf(ctx, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyLightBlockAtHeight provides a mock function with given fields: ctx, height, now
func (_m *LightClient) VerifyLightBlockAtHeight(ctx context.Context, height int64, now time.Time) (*types.LightBlock, error) {
	ret := _m.Called(ctx, height, now)

	var r0 *types.LightBlock
	if rf, ok := ret.Get(0).(func(context.Context, int64, time.Time) *types.LightBlock); ok {
		r0 = rf(ctx, height, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LightBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, time.Time) error); ok {
		r1 = rf(ctx, height, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLightClient creates a new instance of LightClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLightClient(t testing.TB) *LightClient {
	mock := &LightClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}