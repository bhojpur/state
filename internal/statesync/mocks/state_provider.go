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
	state "github.com/bhojpur/state/internal/state"

	testing "testing"

	types "github.com/bhojpur/state/pkg/types"
)

// StateProvider is an autogenerated mock type for the StateProvider type
type StateProvider struct {
	mock.Mock
}

// AppHash provides a mock function with given fields: ctx, height
func (_m *StateProvider) AppHash(ctx context.Context, height uint64) ([]byte, error) {
	ret := _m.Called(ctx, height)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, uint64) []byte); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Commit provides a mock function with given fields: ctx, height
func (_m *StateProvider) Commit(ctx context.Context, height uint64) (*types.Commit, error) {
	ret := _m.Called(ctx, height)

	var r0 *types.Commit
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *types.Commit); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Commit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// State provides a mock function with given fields: ctx, height
func (_m *StateProvider) State(ctx context.Context, height uint64) (state.State, error) {
	ret := _m.Called(ctx, height)

	var r0 state.State
	if rf, ok := ret.Get(0).(func(context.Context, uint64) state.State); ok {
		r0 = rf(ctx, height)
	} else {
		r0 = ret.Get(0).(state.State)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStateProvider creates a new instance of StateProvider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewStateProvider(t testing.TB) *StateProvider {
	mock := &StateProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}