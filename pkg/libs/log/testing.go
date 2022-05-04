package log

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

	"github.com/rs/zerolog"
)

// NewTestingLogger converts a testing.T into a logging interface to
// make test failures and verbose provide better feedback associated
// with test failures. This logging instance is safe for use from
// multiple threads, but in general you should create one of these
// loggers ONCE for each *testing.T instance that you interact with.
//
// By default it collects only ERROR messages, or DEBUG messages in
// verbose mode, and relies on the underlying behavior of
// testing.T.Log()
//
// Users should be careful to ensure that no calls to this logger are
// made in goroutines that are running after (which, by the rules of
// testing.TB will panic.)
func NewTestingLogger(t testing.TB) Logger {
	level := LogLevelError
	if testing.Verbose() {
		level = LogLevelDebug
	}

	return NewTestingLoggerWithLevel(t, level)
}

// NewTestingLoggerWithLevel creates a testing logger instance at a
// specific level that wraps the behavior of testing.T.Log().
func NewTestingLoggerWithLevel(t testing.TB, level string) Logger {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		t.Fatalf("failed to parse log level (%s): %v", level, err)
	}

	return defaultLogger{
		Logger: zerolog.New(newSyncWriter(testingWriter{t})).Level(logLevel),
	}
}

type testingWriter struct {
	t testing.TB
}

func (tw testingWriter) Write(in []byte) (int, error) {
	tw.t.Log(string(in))
	return len(in), nil
}
