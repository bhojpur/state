package p2p_test

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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/p2p"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"
)

// Transports are mainly tested by common tests in transport_test.go, we
// register a transport factory here to get included in those tests.
func init() {
	var network *p2p.MemoryNetwork // shared by transports in the same test

	testTransports["memory"] = func(t *testing.T) p2p.Transport {
		if network == nil {
			network = p2p.NewMemoryNetwork(log.NewNopLogger(), 1)
		}
		i := byte(network.Size())
		nodeID, err := types.NewNodeID(hex.EncodeToString(bytes.Repeat([]byte{i<<4 + i}, 20)))
		require.NoError(t, err)
		transport := network.CreateTransport(nodeID)

		t.Cleanup(func() {
			require.NoError(t, transport.Close())
			network = nil // set up a new memory network for the next test
		})

		return transport
	}
}
