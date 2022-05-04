package privval

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
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/crypto/ed25519"
)

// helper funcs

func newPrivKey() ed25519.PrivKey {
	return ed25519.GenPrivKey()
}

// tests

type listenerTestCase struct {
	description string // For test reporting purposes.
	listener    net.Listener
	dialer      SocketDialer
}

// testUnixAddr will attempt to obtain a platform-independent temporary file
// name for a Unix socket
func testUnixAddr(t *testing.T) (string, error) {
	// N.B. We can't use t.TempDir here because socket filenames have a
	// restrictive length limit (~100 bytes) for silly historical reasons.
	f, err := os.CreateTemp("", "bhojpur-privval-test-*.sock")
	if err != nil {
		return "", err
	}
	addr := f.Name()
	f.Close()
	os.Remove(addr)                       // remove so the test can bind it
	t.Cleanup(func() { os.Remove(addr) }) // clean up after the test
	return addr, nil
}

func tcpListenerTestCase(t *testing.T, timeoutAccept, timeoutReadWrite time.Duration) listenerTestCase {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	tcpLn := NewTCPListener(ln, newPrivKey())
	TCPListenerTimeoutAccept(timeoutAccept)(tcpLn)
	TCPListenerTimeoutReadWrite(timeoutReadWrite)(tcpLn)
	return listenerTestCase{
		description: "TCP",
		listener:    tcpLn,
		dialer:      DialTCPFn(ln.Addr().String(), testTimeoutReadWrite, newPrivKey()),
	}
}

func unixListenerTestCase(t *testing.T, timeoutAccept, timeoutReadWrite time.Duration) listenerTestCase {
	addr, err := testUnixAddr(t)
	if err != nil {
		t.Fatal(err)
	}
	ln, err := net.Listen("unix", addr)
	if err != nil {
		t.Fatal(err)
	}

	unixLn := NewUnixListener(ln)
	UnixListenerTimeoutAccept(timeoutAccept)(unixLn)
	UnixListenerTimeoutReadWrite(timeoutReadWrite)(unixLn)
	return listenerTestCase{
		description: "Unix",
		listener:    unixLn,
		dialer:      DialUnixFn(addr),
	}
}

func listenerTestCases(t *testing.T, timeoutAccept, timeoutReadWrite time.Duration) []listenerTestCase {
	return []listenerTestCase{
		tcpListenerTestCase(t, timeoutAccept, timeoutReadWrite),
		unixListenerTestCase(t, timeoutAccept, timeoutReadWrite),
	}
}

func TestListenerTimeoutAccept(t *testing.T) {
	for _, tc := range listenerTestCases(t, time.Millisecond, time.Second) {
		_, err := tc.listener.Accept()
		opErr, ok := err.(*net.OpError)
		if !ok {
			t.Fatalf("for %s listener, have %v, want *net.OpError", tc.description, err)
		}

		if have, want := opErr.Op, "accept"; have != want {
			t.Errorf("for %s listener,  have %v, want %v", tc.description, have, want)
		}
	}
}

func TestListenerTimeoutReadWrite(t *testing.T) {
	const (
		// This needs to be long enough s.t. the Accept will definitely succeed:
		timeoutAccept = time.Second
		// This can be really short but in the TCP case, the accept can
		// also trigger a timeoutReadWrite. Hence, we need to give it some time.
		// Note: this controls how long this test actually runs.
		timeoutReadWrite = 10 * time.Millisecond
	)
	for _, tc := range listenerTestCases(t, timeoutAccept, timeoutReadWrite) {
		go func(dialer SocketDialer) {
			_, err := dialer()
			require.NoError(t, err)
		}(tc.dialer)

		c, err := tc.listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		// this will timeout because we don't write anything:
		msg := make([]byte, 200)
		_, err = c.Read(msg)
		opErr, ok := err.(*net.OpError)
		if !ok {
			t.Fatalf("for %s listener, have %v, want *net.OpError", tc.description, err)
		}

		if have, want := opErr.Op, "read"; have != want {
			t.Errorf("for %s listener, have %v, want %v", tc.description, have, want)
		}

		if !opErr.Timeout() {
			t.Errorf("for %s listener, got unexpected error: have %v, want Timeout error", tc.description, opErr)
		}
	}
}
