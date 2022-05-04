package strings

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
)

func TestSplitAndTrimEmpty(t *testing.T) {
	testCases := []struct {
		s        string
		sep      string
		cutset   string
		expected []string
	}{
		{"a,b,c", ",", " ", []string{"a", "b", "c"}},
		{" a , b , c ", ",", " ", []string{"a", "b", "c"}},
		{" a, b, c ", ",", " ", []string{"a", "b", "c"}},
		{" a, ", ",", " ", []string{"a"}},
		{"   ", ",", " ", []string{}},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, SplitAndTrimEmpty(tc.s, tc.sep, tc.cutset), "%s", tc.s)
	}
}

func assertCorrectTrim(t *testing.T, input, expected string) {
	t.Helper()
	output, err := ASCIITrim(input)
	require.NoError(t, err)
	require.Equal(t, expected, output)
}

func TestASCIITrim(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		t.Run("NonASCII", func(t *testing.T) {
			notASCIIText := []string{
				"\xC2", "\xC2\xA2", "\xFF", "\x80", "\xF0", "\n", "\t",
			}
			for _, v := range notASCIIText {
				_, err := ASCIITrim(v)
				require.Error(t, err, "%q is not ascii-text", v)
			}
		})
		t.Run("EmptyString", func(t *testing.T) {
			out, err := ASCIITrim("")
			require.NoError(t, err)
			require.Zero(t, out)
		})
		t.Run("ASCIIText", func(t *testing.T) {
			asciiText := []string{
				" ", ".", "x", "$", "_", "abcdefg;", "-", "0x00", "0", "123",
			}
			for _, v := range asciiText {
				_, err := ASCIITrim(v)
				require.NoError(t, err, "%q is  ascii-text", v)
			}
		})
		_, err := ASCIITrim("\xC2\xA2")
		require.Error(t, err)
	})
	t.Run("Trimming", func(t *testing.T) {
		assertCorrectTrim(t, " ", "")
		assertCorrectTrim(t, " a", "a")
		assertCorrectTrim(t, "a ", "a")
		assertCorrectTrim(t, " a ", "a")
	})

}

func TestStringSliceEqual(t *testing.T) {
	tests := []struct {
		a    []string
		b    []string
		want bool
	}{
		{[]string{"hello", "world"}, []string{"hello", "world"}, true},
		{[]string{"test"}, []string{"test"}, true},
		{[]string{"test1"}, []string{"test2"}, false},
		{[]string{"hello", "world."}, []string{"hello", "world!"}, false},
		{[]string{"only 1 word"}, []string{"two", "words!"}, false},
		{[]string{"two", "words!"}, []string{"only 1 word"}, false},
	}
	for i, tt := range tests {
		require.Equal(t, tt.want, StringSliceEqual(tt.a, tt.b),
			"StringSliceEqual failed on test %d", i)
	}
}
