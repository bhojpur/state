package sync

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

// Waker is used to wake up a sleeper when some event occurs. It debounces
// multiple wakeup calls occurring between each sleep, and wakeups are
// non-blocking to avoid having to coordinate goroutines.
type Waker struct {
	wakeCh chan struct{}
}

// NewWaker creates a new Waker.
func NewWaker() *Waker {
	return &Waker{
		wakeCh: make(chan struct{}, 1), // buffer used for debouncing
	}
}

// Sleep returns a channel that blocks until Wake() is called.
func (w *Waker) Sleep() <-chan struct{} {
	return w.wakeCh
}

// Wake wakes up the sleeper.
func (w *Waker) Wake() {
	// A non-blocking send with a size 1 buffer ensures that we never block, and
	// that we queue up at most a single wakeup call between each Sleep().
	select {
	case w.wakeCh <- struct{}{}:
	default:
	}
}
