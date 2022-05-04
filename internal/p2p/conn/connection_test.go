package conn

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
	"encoding/hex"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/libs/protoio"
	v1 "github.com/bhojpur/state/pkg/api/v1/p2p"
	"github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/service"
)

const maxPingPongPacketSize = 1024 // bytes

func createTestMConnection(logger log.Logger, conn net.Conn) *MConnection {
	return createMConnectionWithCallbacks(logger, conn,
		// onRecieve
		func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		},
		// onError
		func(ctx context.Context, r interface{}) {
		})
}

func createMConnectionWithCallbacks(
	logger log.Logger,
	conn net.Conn,
	onReceive func(ctx context.Context, chID ChannelID, msgBytes []byte),
	onError func(ctx context.Context, r interface{}),
) *MConnection {
	cfg := DefaultMConnConfig()
	cfg.PingInterval = 250 * time.Millisecond
	cfg.PongTimeout = 500 * time.Millisecond
	chDescs := []*ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
	c := NewMConnection(logger, conn, chDescs, onReceive, onError, cfg)
	return c
}

func TestMConnectionSendFlushStop(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientConn := createTestMConnection(log.NewNopLogger(), client)
	err := clientConn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(clientConn))

	msg := []byte("abc")
	assert.True(t, clientConn.Send(0x01, msg))

	msgLength := 14

	// start the reader in a new routine, so we can flush
	errCh := make(chan error)
	go func() {
		msgB := make([]byte, msgLength)
		_, err := server.Read(msgB)
		if err != nil {
			t.Error(err)
			return
		}
		errCh <- err
	}()

	timer := time.NewTimer(3 * time.Second)
	select {
	case <-errCh:
	case <-timer.C:
		t.Error("timed out waiting for msgs to be read")
	}
}

func TestMConnectionSend(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createTestMConnection(log.NewNopLogger(), client)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	msg := []byte("Ant-Man")
	assert.True(t, mconn.Send(0x01, msg))
	// Note: subsequent Send/TrySend calls could pass because we are reading from
	// the send queue in a separate goroutine.
	_, err = server.Read(make([]byte, len(msg)))
	if err != nil {
		t.Error(err)
	}

	msg = []byte("Spider-Man")
	assert.True(t, mconn.Send(0x01, msg))
	_, err = server.Read(make([]byte, len(msg)))
	if err != nil {
		t.Error(err)
	}

	assert.False(t, mconn.Send(0x05, []byte("Absorbing Man")), "Send should return false because channel is unknown")
}

func TestMConnectionReceive(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case receivedCh <- msgBytes:
		case <-ctx.Done():
		}
	}
	onError := func(ctx context.Context, r interface{}) {
		select {
		case errorsCh <- r:
		case <-ctx.Done():
		}
	}
	logger := log.NewNopLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn1 := createMConnectionWithCallbacks(logger, client, onReceive, onError)
	err := mconn1.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn1))

	mconn2 := createTestMConnection(logger, server)
	err = mconn2.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn2))

	msg := []byte("Cyclops")
	assert.True(t, mconn2.Send(0x01, msg))

	select {
	case receivedBytes := <-receivedCh:
		assert.Equal(t, msg, receivedBytes)
	case err := <-errorsCh:
		t.Fatalf("Expected %s, got %+v", msg, err)
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("Did not receive %s message in 500ms", msg)
	}
}

func TestMConnectionWillEventuallyTimeout(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createMConnectionWithCallbacks(log.NewNopLogger(), client, nil, nil)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))
	require.True(t, mconn.IsRunning())

	go func() {
		// read the send buffer so that the send receive
		// doesn't get blocked.
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, _ = io.ReadAll(server)
			case <-ctx.Done():
				return
			}
		}
	}()

	// wait for the send routine to die because it doesn't
	select {
	case <-mconn.doneSendRoutine:
		require.True(t, time.Since(mconn.getLastMessageAt()) > mconn.config.PongTimeout,
			"the connection state reflects that we've passed the pong timeout")
		// since we hit the timeout, things should be shutdown
		require.False(t, mconn.IsRunning())
	case <-time.After(2 * mconn.config.PongTimeout):
		t.Fatal("connection did not hit timeout", mconn.config.PongTimeout)
	}
}

func TestMConnectionMultiplePongsInTheBeginning(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case receivedCh <- msgBytes:
		case <-ctx.Done():
		}
	}
	onError := func(ctx context.Context, r interface{}) {
		select {
		case errorsCh <- r:
		case <-ctx.Done():
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createMConnectionWithCallbacks(log.NewNopLogger(), client, onReceive, onError)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	// sending 3 pongs in a row (abuse)
	protoWriter := protoio.NewDelimitedWriter(server)

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	// read ping (one byte)
	var packet v1.Packet
	_, err = protoio.NewDelimitedReader(server, maxPingPongPacketSize).ReadMsg(&packet)
	require.NoError(t, err)

	// respond with pong
	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	pongTimerExpired := mconn.config.PongTimeout + 20*time.Millisecond
	select {
	case msgBytes := <-receivedCh:
		t.Fatalf("Expected no data, but got %v", msgBytes)
	case err := <-errorsCh:
		t.Fatalf("Expected no error, but got %v", err)
	case <-time.After(pongTimerExpired):
		assert.True(t, mconn.IsRunning())
	}
}

func TestMConnectionMultiplePings(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case receivedCh <- msgBytes:
		case <-ctx.Done():
		}
	}
	onError := func(ctx context.Context, r interface{}) {
		select {
		case errorsCh <- r:
		case <-ctx.Done():
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createMConnectionWithCallbacks(log.NewNopLogger(), client, onReceive, onError)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	// sending three pings in a row (abuse)
	protoReader := protoio.NewDelimitedReader(server, maxPingPongPacketSize)
	protoWriter := protoio.NewDelimitedWriter(server)
	var pkt v1.Packet

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPing{}))
	require.NoError(t, err)

	_, err = protoReader.ReadMsg(&pkt)
	require.NoError(t, err)

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPing{}))
	require.NoError(t, err)

	_, err = protoReader.ReadMsg(&pkt)
	require.NoError(t, err)

	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPing{}))
	require.NoError(t, err)

	_, err = protoReader.ReadMsg(&pkt)
	require.NoError(t, err)

	assert.True(t, mconn.IsRunning())
}

func TestMConnectionPingPongs(t *testing.T) {
	// check that we are not leaking any go-routines
	t.Cleanup(leaktest.CheckTimeout(t, 10*time.Second))

	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case receivedCh <- msgBytes:
		case <-ctx.Done():
		}
	}
	onError := func(ctx context.Context, r interface{}) {
		select {
		case errorsCh <- r:
		case <-ctx.Done():
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createMConnectionWithCallbacks(log.NewNopLogger(), client, onReceive, onError)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	protoReader := protoio.NewDelimitedReader(server, maxPingPongPacketSize)
	protoWriter := protoio.NewDelimitedWriter(server)
	var pkt v1.PacketPing

	// read ping
	_, err = protoReader.ReadMsg(&pkt)
	require.NoError(t, err)

	// respond with pong
	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	time.Sleep(mconn.config.PingInterval)

	// read ping
	_, err = protoReader.ReadMsg(&pkt)
	require.NoError(t, err)

	// respond with pong
	_, err = protoWriter.WriteMsg(mustWrapPacket(&v1.PacketPong{}))
	require.NoError(t, err)

	pongTimerExpired := (mconn.config.PongTimeout + 20*time.Millisecond) * 4
	select {
	case msgBytes := <-receivedCh:
		t.Fatalf("Expected no data, but got %v", msgBytes)
	case err := <-errorsCh:
		t.Fatalf("Expected no error, but got %v", err)
	case <-time.After(2 * pongTimerExpired):
		assert.True(t, mconn.IsRunning())
	}
}

func TestMConnectionStopsAndReturnsError(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))

	receivedCh := make(chan []byte)
	errorsCh := make(chan interface{})
	onReceive := func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case receivedCh <- msgBytes:
		case <-ctx.Done():
		}
	}
	onError := func(ctx context.Context, r interface{}) {
		select {
		case errorsCh <- r:
		case <-ctx.Done():
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createMConnectionWithCallbacks(log.NewNopLogger(), client, onReceive, onError)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	if err := client.Close(); err != nil {
		t.Error(err)
	}

	select {
	case receivedBytes := <-receivedCh:
		t.Fatalf("Expected error, got %v", receivedBytes)
	case err := <-errorsCh:
		assert.NotNil(t, err)
		assert.False(t, mconn.IsRunning())
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Did not receive error in 500ms")
	}
}

func newClientAndServerConnsForReadErrors(
	ctx context.Context,
	t *testing.T,
	chOnErr chan struct{},
) (*MConnection, *MConnection) {
	server, client := net.Pipe()

	onReceive := func(context.Context, ChannelID, []byte) {}
	onError := func(context.Context, interface{}) {}

	// create client conn with two channels
	chDescs := []*ChannelDescriptor{
		{ID: 0x01, Priority: 1, SendQueueCapacity: 1},
		{ID: 0x02, Priority: 1, SendQueueCapacity: 1},
	}
	logger := log.NewNopLogger()

	mconnClient := NewMConnection(logger.With("module", "client"), client, chDescs, onReceive, onError, DefaultMConnConfig())
	err := mconnClient.Start(ctx)
	require.NoError(t, err)

	// create server conn with 1 channel
	// it fires on chOnErr when there's an error
	serverLogger := logger.With("module", "server")
	onError = func(ctx context.Context, r interface{}) {
		select {
		case <-ctx.Done():
		case chOnErr <- struct{}{}:
		}
	}

	mconnServer := createMConnectionWithCallbacks(serverLogger, server, onReceive, onError)
	err = mconnServer.Start(ctx)
	require.NoError(t, err)
	return mconnClient, mconnServer
}

func expectSend(ch chan struct{}) bool {
	after := time.After(time.Second * 5)
	select {
	case <-ch:
		return true
	case <-after:
		return false
	}
}

func TestMConnectionReadErrorBadEncoding(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chOnErr := make(chan struct{})
	mconnClient, mconnServer := newClientAndServerConnsForReadErrors(ctx, t, chOnErr)

	client := mconnClient.conn

	// Write it.
	_, err := client.Write([]byte{1, 2, 3, 4, 5})
	require.NoError(t, err)
	assert.True(t, expectSend(chOnErr), "badly encoded msgPacket")
	t.Cleanup(waitAll(mconnClient, mconnServer))
}

func TestMConnectionReadErrorUnknownChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chOnErr := make(chan struct{})
	mconnClient, mconnServer := newClientAndServerConnsForReadErrors(ctx, t, chOnErr)

	msg := []byte("Ant-Man")

	// fail to send msg on channel unknown by client
	assert.False(t, mconnClient.Send(0x03, msg))

	// send msg on channel unknown by the server.
	// should cause an error
	assert.True(t, mconnClient.Send(0x02, msg))
	assert.True(t, expectSend(chOnErr), "unknown channel")
	t.Cleanup(waitAll(mconnClient, mconnServer))
}

func TestMConnectionReadErrorLongMessage(t *testing.T) {
	chOnErr := make(chan struct{})
	chOnRcv := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconnClient, mconnServer := newClientAndServerConnsForReadErrors(ctx, t, chOnErr)
	t.Cleanup(waitAll(mconnClient, mconnServer))

	mconnServer.onReceive = func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case <-ctx.Done():
		case chOnRcv <- struct{}{}:
		}
	}

	client := mconnClient.conn
	protoWriter := protoio.NewDelimitedWriter(client)

	// send msg thats just right
	var packet = v1.PacketMsg{
		ChannelID: 0x01,
		EOF:       true,
		Data:      make([]byte, mconnClient.config.MaxPacketMsgPayloadSize),
	}

	_, err := protoWriter.WriteMsg(mustWrapPacket(&packet))
	require.NoError(t, err)
	assert.True(t, expectSend(chOnRcv), "msg just right")

	// send msg thats too long
	packet = v1.PacketMsg{
		ChannelID: 0x01,
		EOF:       true,
		Data:      make([]byte, mconnClient.config.MaxPacketMsgPayloadSize+100),
	}

	_, err = protoWriter.WriteMsg(mustWrapPacket(&packet))
	require.Error(t, err)
	assert.True(t, expectSend(chOnErr), "msg too long")
}

func TestMConnectionReadErrorUnknownMsgType(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chOnErr := make(chan struct{})
	mconnClient, mconnServer := newClientAndServerConnsForReadErrors(ctx, t, chOnErr)
	t.Cleanup(waitAll(mconnClient, mconnServer))

	// send msg with unknown msg type
	_, err := protoio.NewDelimitedWriter(mconnClient.conn).WriteMsg(&types.Header{ChainID: "x"})
	require.NoError(t, err)
	assert.True(t, expectSend(chOnErr), "unknown msg type")
}

func TestMConnectionTrySend(t *testing.T) {
	server, client := net.Pipe()
	t.Cleanup(closeAll(t, client, server))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconn := createTestMConnection(log.NewNopLogger(), client)
	err := mconn.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(waitAll(mconn))

	msg := []byte("Semicolon-Woman")
	resultCh := make(chan string, 2)
	assert.True(t, mconn.Send(0x01, msg))
	_, err = server.Read(make([]byte, len(msg)))
	require.NoError(t, err)
	assert.True(t, mconn.Send(0x01, msg))
	go func() {
		mconn.Send(0x01, msg)
		resultCh <- "TrySend"
	}()
	assert.False(t, mconn.Send(0x01, msg))
	assert.Equal(t, "TrySend", <-resultCh)
}

func TestConnVectors(t *testing.T) {

	testCases := []struct {
		testName string
		msg      proto.Message
		expBytes string
	}{
		{"PacketPing", &v1.PacketPing{}, "0a00"},
		{"PacketPong", &v1.PacketPong{}, "1200"},
		{"PacketMsg", &v1.PacketMsg{ChannelID: 1, EOF: false, Data: []byte("data transmitted over the wire")}, "1a2208011a1e64617461207472616e736d6974746564206f766572207468652077697265"},
	}

	for _, tc := range testCases {
		tc := tc

		pm := mustWrapPacket(tc.msg)
		bz, err := pm.Marshal()
		require.NoError(t, err, tc.testName)

		require.Equal(t, tc.expBytes, hex.EncodeToString(bz), tc.testName)
	}
}

func TestMConnectionChannelOverflow(t *testing.T) {
	chOnErr := make(chan struct{})
	chOnRcv := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mconnClient, mconnServer := newClientAndServerConnsForReadErrors(ctx, t, chOnErr)
	t.Cleanup(waitAll(mconnClient, mconnServer))

	mconnServer.onReceive = func(ctx context.Context, chID ChannelID, msgBytes []byte) {
		select {
		case <-ctx.Done():
		case chOnRcv <- struct{}{}:
		}
	}

	client := mconnClient.conn
	protoWriter := protoio.NewDelimitedWriter(client)

	var packet = v1.PacketMsg{
		ChannelID: 0x01,
		EOF:       true,
		Data:      []byte(`42`),
	}
	_, err := protoWriter.WriteMsg(mustWrapPacket(&packet))
	require.NoError(t, err)
	assert.True(t, expectSend(chOnRcv))

	packet.ChannelID = int32(1025)
	_, err = protoWriter.WriteMsg(mustWrapPacket(&packet))
	require.NoError(t, err)
	assert.False(t, expectSend(chOnRcv))

}

func waitAll(waiters ...service.Service) func() {
	return func() {
		switch len(waiters) {
		case 0:
			return
		case 1:
			waiters[0].Wait()
			return
		default:
			wg := &sync.WaitGroup{}

			for _, w := range waiters {
				wg.Add(1)
				go func(s service.Service) {
					defer wg.Done()
					s.Wait()
				}(w)
			}

			wg.Wait()
		}
	}
}

type closer interface {
	Close() error
}

func closeAll(t *testing.T, closers ...closer) func() {
	return func() {
		for _, s := range closers {
			if err := s.Close(); err != nil {
				t.Log(err)
			}
		}
	}
}
