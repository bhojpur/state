package math

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

	"github.com/stretchr/testify/assert"
)

func TestParseFraction(t *testing.T) {

	testCases := []struct {
		f   string
		exp Fraction
		err bool
	}{
		{
			f:   "2/3",
			exp: Fraction{2, 3},
			err: false,
		},
		{
			f:   "15/5",
			exp: Fraction{15, 5},
			err: false,
		},
		// test divide by zero error
		{
			f:   "2/0",
			exp: Fraction{},
			err: true,
		},
		// test negative
		{
			f:   "-1/2",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "1/-2",
			exp: Fraction{},
			err: true,
		},
		// test overflow
		{
			f:   "9223372036854775808/2",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "2/9223372036854775808",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "2/3/4",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "123",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "1a2/4",
			exp: Fraction{},
			err: true,
		},
		{
			f:   "1/3bc4",
			exp: Fraction{},
			err: true,
		},
	}

	for idx, tc := range testCases {
		output, err := ParseFraction(tc.f)
		if tc.err {
			assert.Error(t, err, idx)
		} else {
			assert.NoError(t, err, idx)
		}
		assert.Equal(t, tc.exp, output, idx)
	}

}
