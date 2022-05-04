package timer

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
	"testing"
	"time"

	// make govet noshadow happy...

	asrt "github.com/stretchr/testify/assert"
)

type thCounter struct {
	input chan struct{}
	mtx   sync.Mutex
	count int
}

func (c *thCounter) Increment() {
	c.mtx.Lock()
	c.count++
	c.mtx.Unlock()
}

func (c *thCounter) Count() int {
	c.mtx.Lock()
	val := c.count
	c.mtx.Unlock()
	return val
}

// Read should run in a go-routine and
// updates count by one every time a packet comes in
func (c *thCounter) Read() {
	for range c.input {
		c.Increment()
	}
}

func TestThrottle(test *testing.T) {
	assert := asrt.New(test)

	ms := 50
	delay := time.Duration(ms) * time.Millisecond
	longwait := time.Duration(2) * delay
	t := NewThrottleTimer("foo", delay)

	// start at 0
	c := &thCounter{input: t.Ch}
	assert.Equal(0, c.Count())
	go c.Read()

	// waiting does nothing
	time.Sleep(longwait)
	assert.Equal(0, c.Count())

	// send one event adds one
	t.Set()
	time.Sleep(longwait)
	assert.Equal(1, c.Count())

	// send a burst adds one
	for i := 0; i < 5; i++ {
		t.Set()
	}
	time.Sleep(longwait)
	assert.Equal(2, c.Count())

	// send 12, over 2 delay sections, adds 3 or more. It
	// is possible for more to be added if the overhead
	// in executing the loop is large
	short := time.Duration(ms/5) * time.Millisecond
	for i := 0; i < 13; i++ {
		t.Set()
		time.Sleep(short)
	}
	time.Sleep(longwait)
	assert.LessOrEqual(5, c.Count())

	close(t.Ch)
}
