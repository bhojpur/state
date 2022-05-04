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
	"sync"
)

// default capacity for the size of a queue
const defaultCapacity uint = 16e6 // ~16MB

// queue does QoS scheduling for Envelopes, enqueueing and dequeueing according
// to some policy. Queues are used at contention points, i.e.:
//
// - Receiving inbound messages to a single channel from all peers.
// - Sending outbound messages to a single peer from all channels.
type queue interface {
	// enqueue returns a channel for submitting envelopes.
	enqueue() chan<- Envelope

	// dequeue returns a channel ordered according to some queueing policy.
	dequeue() <-chan Envelope

	// close closes the queue. After this call enqueue() will block, so the
	// caller must select on closed() as well to avoid blocking forever. The
	// enqueue() and dequeue() channels will not be closed.
	close()

	// closed returns a channel that's closed when the scheduler is closed.
	closed() <-chan struct{}
}

// fifoQueue is a simple unbuffered lossless queue that passes messages through
// in the order they were received, and blocks until message is received.
type fifoQueue struct {
	queueCh chan Envelope
	closeFn func()
	closeCh <-chan struct{}
}

func newFIFOQueue(size int) queue {
	closeCh := make(chan struct{})
	once := &sync.Once{}

	return &fifoQueue{
		queueCh: make(chan Envelope, size),
		closeFn: func() { once.Do(func() { close(closeCh) }) },
		closeCh: closeCh,
	}
}

func (q *fifoQueue) enqueue() chan<- Envelope { return q.queueCh }
func (q *fifoQueue) dequeue() <-chan Envelope { return q.queueCh }
func (q *fifoQueue) close()                   { q.closeFn() }
func (q *fifoQueue) closed() <-chan struct{}  { return q.closeCh }
