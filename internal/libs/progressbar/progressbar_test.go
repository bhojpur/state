package progressbar

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
	"time"

	"github.com/stretchr/testify/require"
)

func TestProgressBar(t *testing.T) {
	zero := int64(0)
	hundred := int64(100)

	var bar Bar
	bar.NewOption(zero, hundred)

	require.Equal(t, zero, bar.start)
	require.Equal(t, zero, bar.cur)
	require.Equal(t, hundred, bar.total)
	require.Equal(t, zero, bar.percent)
	require.Equal(t, "█", bar.graph)
	require.Equal(t, "", bar.rate)

	defer bar.Finish()
	for i := zero; i <= hundred; i++ {
		time.Sleep(1 * time.Millisecond)
		bar.Play(i)
	}

	require.Equal(t, zero, bar.start)
	require.Equal(t, hundred, bar.cur)
	require.Equal(t, hundred, bar.total)
	require.Equal(t, hundred, bar.percent)

	var rate string
	for i := zero; i < hundred/2; i++ {
		rate += "█"
	}

	require.Equal(t, rate, bar.rate)
}
