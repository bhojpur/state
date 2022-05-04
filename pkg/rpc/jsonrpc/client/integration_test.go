//go:build release
// +build release

package client

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
	"bytes"
	"context"
	"errors"
	"net"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWSClientReconnectWithJitter(t *testing.T) {
	const numClients = 8
	const maxReconnectAttempts = 3
	const maxSleepTime = time.Duration(((1<<maxReconnectAttempts)-1)+maxReconnectAttempts) * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	failDialer := func(net, addr string) (net.Conn, error) {
		return nil, errors.New("not connected")
	}

	clientMap := make(map[int]*WSClient)
	buf := new(bytes.Buffer)
	for i := 0; i < numClients; i++ {
		c, err := NewWS("tcp://foo", "/websocket")
		require.NoError(t, err)
		c.Dialer = failDialer
		c.maxReconnectAttempts = maxReconnectAttempts
		c.Start(ctx)

		// Not invoking defer c.Stop() because
		// after all the reconnect attempts have been
		// exhausted, c.Stop is implicitly invoked.
		clientMap[i] = c
		// Trigger the reconnect routine that performs exponential backoff.
		go c.reconnect(ctx)
	}

	// Next we have to examine the logs to ensure that no single time was repeated
	backoffDurRegexp := regexp.MustCompile(`backoff_duration=(.+)`)
	matches := backoffDurRegexp.FindAll(buf.Bytes(), -1)
	seenMap := make(map[string]int)
	for i, match := range matches {
		if origIndex, seen := seenMap[string(match)]; seen {
			t.Errorf("match #%d (%q) was seen originally at log entry #%d", i, match, origIndex)
		} else {
			seenMap[string(match)] = i
		}
	}
}
