package clist

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

import "testing"

func BenchmarkDetaching(b *testing.B) {
	lst := New()
	for i := 0; i < b.N+1; i++ {
		lst.PushBack(i)
	}
	start := lst.Front()
	nxt := start.Next()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start.removed = true
		start.detachNext()
		start.DetachPrev()
		tmp := nxt
		nxt = nxt.Next()
		start = tmp
	}
}

// This is used to benchmark the time of RMutex.
func BenchmarkRemoved(b *testing.B) {
	lst := New()
	for i := 0; i < b.N+1; i++ {
		lst.PushBack(i)
	}
	start := lst.Front()
	nxt := start.Next()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start.Removed()
		tmp := nxt
		nxt = nxt.Next()
		start = tmp
	}
}

func BenchmarkPushBack(b *testing.B) {
	lst := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lst.PushBack(i)
	}
}
