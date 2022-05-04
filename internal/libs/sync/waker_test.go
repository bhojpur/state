package sync_test

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
	"testing"

	"github.com/stretchr/testify/require"

	tmsync "github.com/bhojpur/state/internal/libs/sync"
)

func TestWaker(t *testing.T) {

	// A new waker should block when sleeping.
	waker := tmsync.NewWaker()

	select {
	case <-waker.Sleep():
		require.Fail(t, "unexpected wakeup")
	default:
	}

	// Wakeups should not block, and should cause the next sleeper to awaken.
	waker.Wake()

	select {
	case <-waker.Sleep():
	default:
		require.Fail(t, "expected wakeup, but sleeping instead")
	}

	// Multiple wakeups should only wake a single sleeper.
	waker.Wake()
	waker.Wake()
	waker.Wake()

	select {
	case <-waker.Sleep():
	default:
		require.Fail(t, "expected wakeup, but sleeping instead")
	}

	select {
	case <-waker.Sleep():
		require.Fail(t, "unexpected wakeup")
	default:
	}
}
