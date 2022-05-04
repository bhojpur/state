package p2p

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
	"testing"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/bhojpur/state/pkg/libs/log"
)

type testMessage = gogotypes.StringValue

func TestCloseWhileDequeueFull(t *testing.T) {
	enqueueLength := 5
	chDescs := []*ChannelDescriptor{
		{ID: 0x01, Priority: 1},
	}
	pqueue := newPQScheduler(log.NewNopLogger(), NopMetrics(), chDescs, uint(enqueueLength), 1, 120)

	for i := 0; i < enqueueLength; i++ {
		pqueue.enqueue() <- Envelope{
			ChannelID: 0x01,
			Message:   &testMessage{Value: "foo"}, // 5 bytes
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pqueue.process(ctx)

	// sleep to allow context switch for process() to run
	time.Sleep(10 * time.Millisecond)
	doneCh := make(chan struct{})
	go func() {
		pqueue.close()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("pqueue failed to close")
	}
}
