package service

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
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/libs/log"
)

type testService struct {
	started      bool
	stopped      bool
	multiStopped bool
	mu           sync.Mutex
	BaseService
}

func (t *testService) OnStop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped == true {
		t.multiStopped = true
	}
	t.stopped = true
}
func (t *testService) OnStart(context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.started = true
	return nil
}

func (t *testService) isStarted() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.started
}

func (t *testService) isStopped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.stopped
}

func (t *testService) isMultiStopped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.multiStopped
}

func TestBaseService(t *testing.T) {
	t.Cleanup(leaktest.Check(t))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewNopLogger()

	t.Run("Wait", func(t *testing.T) {
		wctx, wcancel := context.WithCancel(ctx)
		defer wcancel()
		ts := &testService{}
		ts.BaseService = *NewBaseService(logger, t.Name(), ts)
		err := ts.Start(wctx)
		require.NoError(t, err)
		require.True(t, ts.isStarted())

		waitFinished := make(chan struct{})
		wcancel()
		go func() {
			ts.Wait()
			close(waitFinished)
		}()

		select {
		case <-waitFinished:
			assert.True(t, ts.isStopped(), "failed to stop")
			assert.False(t, ts.IsRunning(), "is not running")

		case <-time.After(100 * time.Millisecond):
			t.Fatal("expected Wait() to finish within 100 ms.")
		}
	})
	t.Run("ManualStop", func(t *testing.T) {
		ts := &testService{}
		ts.BaseService = *NewBaseService(logger, t.Name(), ts)
		require.False(t, ts.IsRunning())
		require.False(t, ts.isStarted())
		require.NoError(t, ts.Start(ctx))

		require.True(t, ts.isStarted())

		ts.Stop()
		require.True(t, ts.isStopped())
		require.False(t, ts.IsRunning())
	})
	t.Run("MultiStop", func(t *testing.T) {
		t.Run("SingleThreaded", func(t *testing.T) {
			ts := &testService{}
			ts.BaseService = *NewBaseService(logger, t.Name(), ts)

			require.NoError(t, ts.Start(ctx))
			require.True(t, ts.isStarted())
			ts.Stop()
			require.True(t, ts.isStopped())
			require.False(t, ts.isMultiStopped())
			ts.Stop()
			require.False(t, ts.isMultiStopped())
		})
		t.Run("MultiThreaded", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ts := &testService{}
			ts.BaseService = *NewBaseService(logger, t.Name(), ts)

			require.NoError(t, ts.Start(ctx))
			require.True(t, ts.isStarted())

			go ts.Stop()
			go cancel()

			ts.Wait()

			require.True(t, ts.isStopped())
			require.False(t, ts.isMultiStopped())
		})

	})

}
