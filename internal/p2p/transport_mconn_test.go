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
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/p2p"
	"github.com/bhojpur/state/internal/p2p/conn"
	"github.com/bhojpur/state/pkg/libs/log"
)

// Transports are mainly tested by common tests in transport_test.go, we
// register a transport factory here to get included in those tests.
func init() {
	testTransports["mconn"] = func(t *testing.T) p2p.Transport {
		transport := p2p.NewMConnTransport(
			log.NewNopLogger(),
			conn.DefaultMConnConfig(),
			[]*p2p.ChannelDescriptor{{ID: chID, Priority: 1}},
			p2p.MConnTransportOptions{},
		)
		err := transport.Listen(&p2p.Endpoint{
			Protocol: p2p.MConnProtocol,
			IP:       net.IPv4(127, 0, 0, 1),
			Port:     0, // assign a random port
		})
		require.NoError(t, err)

		t.Cleanup(func() { _ = transport.Close() })

		return transport
	}
}

func TestMConnTransport_AcceptBeforeListen(t *testing.T) {
	transport := p2p.NewMConnTransport(
		log.NewNopLogger(),
		conn.DefaultMConnConfig(),
		[]*p2p.ChannelDescriptor{{ID: chID, Priority: 1}},
		p2p.MConnTransportOptions{
			MaxAcceptedConnections: 2,
		},
	)
	t.Cleanup(func() {
		_ = transport.Close()
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := transport.Accept(ctx)
	require.Error(t, err)
	require.NotEqual(t, io.EOF, err) // io.EOF should be returned after Close()
}

func TestMConnTransport_AcceptMaxAcceptedConnections(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transport := p2p.NewMConnTransport(
		log.NewNopLogger(),
		conn.DefaultMConnConfig(),
		[]*p2p.ChannelDescriptor{{ID: chID, Priority: 1}},
		p2p.MConnTransportOptions{
			MaxAcceptedConnections: 2,
		},
	)
	t.Cleanup(func() {
		_ = transport.Close()
	})
	err := transport.Listen(&p2p.Endpoint{
		Protocol: p2p.MConnProtocol,
		IP:       net.IPv4(127, 0, 0, 1),
	})
	require.NoError(t, err)
	endpoint, err := transport.Endpoint()
	require.NoError(t, err)
	require.NotNil(t, endpoint)

	// Start a goroutine to just accept any connections.
	acceptCh := make(chan p2p.Connection, 10)
	go func() {
		for {
			conn, err := transport.Accept(ctx)
			if err != nil {
				return
			}
			acceptCh <- conn
		}
	}()

	// The first two connections should be accepted just fine.
	dial1, err := transport.Dial(ctx, endpoint)
	require.NoError(t, err)
	defer dial1.Close()
	accept1 := <-acceptCh
	defer accept1.Close()
	require.Equal(t, dial1.LocalEndpoint(), accept1.RemoteEndpoint())

	dial2, err := transport.Dial(ctx, endpoint)
	require.NoError(t, err)
	defer dial2.Close()
	accept2 := <-acceptCh
	defer accept2.Close()
	require.Equal(t, dial2.LocalEndpoint(), accept2.RemoteEndpoint())

	// The third connection will be dialed successfully, but the accept should
	// not go through.
	dial3, err := transport.Dial(ctx, endpoint)
	require.NoError(t, err)
	defer dial3.Close()
	select {
	case <-acceptCh:
		require.Fail(t, "unexpected accept")
	case <-time.After(time.Second):
	}

	// However, once either of the other connections are closed, the accept
	// goes through.
	require.NoError(t, accept1.Close())
	accept3 := <-acceptCh
	defer accept3.Close()
	require.Equal(t, dial3.LocalEndpoint(), accept3.RemoteEndpoint())
}

func TestMConnTransport_Listen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testcases := []struct {
		endpoint *p2p.Endpoint
		ok       bool
	}{
		// Valid v4 and v6 addresses, with mconn and tcp protocols.
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, IP: net.IPv4zero}, true},
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, IP: net.IPv4(127, 0, 0, 1)}, true},
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, IP: net.IPv6zero}, true},
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, IP: net.IPv6loopback}, true},
		{&p2p.Endpoint{Protocol: p2p.TCPProtocol, IP: net.IPv4zero}, true},

		// Invalid endpoints.
		{&p2p.Endpoint{}, false},
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, Path: "foo"}, false},
		{&p2p.Endpoint{Protocol: p2p.MConnProtocol, IP: net.IPv4zero, Path: "foo"}, false},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.endpoint.String(), func(t *testing.T) {
			t.Cleanup(leaktest.Check(t))

			transport := p2p.NewMConnTransport(
				log.NewNopLogger(),
				conn.DefaultMConnConfig(),
				[]*p2p.ChannelDescriptor{{ID: chID, Priority: 1}},
				p2p.MConnTransportOptions{},
			)

			// Transport should not listen on any endpoints yet.
			endpoint, err := transport.Endpoint()
			require.Error(t, err)
			require.Nil(t, endpoint)

			// Start listening, and check any expected errors.
			err = transport.Listen(tc.endpoint)
			if !tc.ok {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check the endpoint.
			endpoint, err = transport.Endpoint()
			require.NoError(t, err)
			require.NotNil(t, endpoint)

			require.Equal(t, p2p.MConnProtocol, endpoint.Protocol)
			if tc.endpoint.IP.IsUnspecified() {
				require.True(t, endpoint.IP.IsUnspecified(),
					"expected unspecified IP, got %v", endpoint.IP)
			} else {
				require.True(t, tc.endpoint.IP.Equal(endpoint.IP),
					"expected %v, got %v", tc.endpoint.IP, endpoint.IP)
			}
			require.NotZero(t, endpoint.Port)
			require.Empty(t, endpoint.Path)

			dialedChan := make(chan struct{})

			var peerConn p2p.Connection
			go func() {
				// Dialing the endpoint should work.
				var err error
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				peerConn, err = transport.Dial(ctx, endpoint)
				require.NoError(t, err)
				close(dialedChan)
			}()

			conn, err := transport.Accept(ctx)
			require.NoError(t, err)
			_ = conn.Close()
			<-dialedChan

			// closing the connection should not error
			require.NoError(t, peerConn.Close())

			// try to read from the connection should error
			_, _, err = peerConn.ReceiveMessage(ctx)
			require.Error(t, err)

			// Trying to listen again should error.
			err = transport.Listen(tc.endpoint)
			require.Error(t, err)

			// close the transport
			_ = transport.Close()

			// Dialing the closed endpoint should error
			_, err = transport.Dial(ctx, endpoint)
			require.Error(t, err)
		})
	}
}
