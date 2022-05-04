package proxy

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
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gotest.tools/assert"

	"github.com/bhojpur/state/example/kvstore"
	abciclient "github.com/bhojpur/state/pkg/abci/client"
	abcimocks "github.com/bhojpur/state/pkg/abci/client/mocks"
	"github.com/bhojpur/state/pkg/abci/server"
	"github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/libs/log"
	librand "github.com/bhojpur/state/pkg/libs/rand"
)

type appConnTestI interface {
	Echo(context.Context, string) (*types.ResponseEcho, error)
	Flush(context.Context) error
	Info(context.Context, *types.RequestInfo) (*types.ResponseInfo, error)
}

type appConnTest struct {
	appConn abciclient.Client
}

func newAppConnTest(appConn abciclient.Client) appConnTestI {
	return &appConnTest{appConn}
}

func (app *appConnTest) Echo(ctx context.Context, msg string) (*types.ResponseEcho, error) {
	return app.appConn.Echo(ctx, msg)
}

func (app *appConnTest) Flush(ctx context.Context) error {
	return app.appConn.Flush(ctx)
}

func (app *appConnTest) Info(ctx context.Context, req *types.RequestInfo) (*types.ResponseInfo, error) {
	return app.appConn.Info(ctx, req)
}

var SOCKET = "socket"

func TestEcho(t *testing.T) {
	sockPath := fmt.Sprintf("unix:///tmp/echo_%v.sock", librand.Str(6))
	logger := log.NewNopLogger()
	client, err := abciclient.NewClient(logger, sockPath, SOCKET, true)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	s := server.NewSocketServer(logger.With("module", "abci-server"), sockPath, kvstore.NewApplication())
	require.NoError(t, s.Start(ctx), "error starting socket server")
	t.Cleanup(func() { cancel(); s.Wait() })

	// Start client
	require.NoError(t, client.Start(ctx), "Error starting ABCI client")

	proxy := newAppConnTest(client)
	t.Log("Connected")

	for i := 0; i < 1000; i++ {
		_, err = proxy.Echo(ctx, fmt.Sprintf("echo-%v", i))
		if err != nil {
			t.Error(err)
		}
		// flush sometimes
		if i%128 == 0 {
			if err := proxy.Flush(ctx); err != nil {
				t.Error(err)
			}
		}
	}
	if err := proxy.Flush(ctx); err != nil {
		t.Error(err)
	}
}

func BenchmarkEcho(b *testing.B) {
	b.StopTimer() // Initialize
	sockPath := fmt.Sprintf("unix:///tmp/echo_%v.sock", librand.Str(6))
	logger := log.NewNopLogger()
	client, err := abciclient.NewClient(logger, sockPath, SOCKET, true)
	if err != nil {
		b.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	s := server.NewSocketServer(logger.With("module", "abci-server"), sockPath, kvstore.NewApplication())
	require.NoError(b, s.Start(ctx), "Error starting socket server")
	b.Cleanup(func() { cancel(); s.Wait() })

	// Start client
	require.NoError(b, client.Start(ctx), "Error starting ABCI client")

	proxy := newAppConnTest(client)
	b.Log("Connected")
	echoString := strings.Repeat(" ", 200)
	b.StartTimer() // Start benchmarking tests

	for i := 0; i < b.N; i++ {
		_, err = proxy.Echo(ctx, echoString)
		if err != nil {
			b.Error(err)
		}
		// flush sometimes
		if i%128 == 0 {
			if err := proxy.Flush(ctx); err != nil {
				b.Error(err)
			}
		}
	}
	if err := proxy.Flush(ctx); err != nil {
		b.Error(err)
	}

	b.StopTimer()
	// info := proxy.Info(types.RequestInfo{""})
	// b.Log("N: ", b.N, info)
}

func TestInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sockPath := fmt.Sprintf("unix:///tmp/echo_%v.sock", librand.Str(6))
	logger := log.NewNopLogger()
	client, err := abciclient.NewClient(logger, sockPath, SOCKET, true)
	if err != nil {
		t.Fatal(err)
	}

	// Start server
	s := server.NewSocketServer(logger.With("module", "abci-server"), sockPath, kvstore.NewApplication())
	require.NoError(t, s.Start(ctx), "Error starting socket server")
	t.Cleanup(func() { cancel(); s.Wait() })

	// Start client
	require.NoError(t, client.Start(ctx), "Error starting ABCI client")

	proxy := newAppConnTest(client)
	t.Log("Connected")

	resInfo, err := proxy.Info(ctx, &RequestInfo)
	require.NoError(t, err)

	if resInfo.Data != "{\"size\":0}" {
		t.Error("Expected ResponseInfo with one element '{\"size\":0}' but got something else")
	}
}

type noopStoppableClientImpl struct {
	abciclient.Client
	count int
}

func (c *noopStoppableClientImpl) Stop() { c.count++ }

func TestAppConns_Start_Stop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientMock := &abcimocks.Client{}
	clientMock.On("Start", mock.Anything).Return(nil)
	clientMock.On("Error").Return(nil)
	clientMock.On("IsRunning").Return(true)
	clientMock.On("Wait").Return(nil).Times(1)
	cl := &noopStoppableClientImpl{Client: clientMock}

	appConns := New(cl, log.NewNopLogger(), NopMetrics())

	err := appConns.Start(ctx)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	cancel()
	appConns.Wait()

	clientMock.AssertExpectations(t)
	assert.Equal(t, 1, cl.count)
}

// Upon failure, we call bos.Kill
func TestAppConns_Failure(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGABRT)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientMock := &abcimocks.Client{}
	clientMock.On("SetLogger", mock.Anything).Return()
	clientMock.On("Start", mock.Anything).Return(nil)
	clientMock.On("IsRunning").Return(true)
	clientMock.On("Wait").Return(nil)
	clientMock.On("Error").Return(errors.New("EOF"))
	cl := &noopStoppableClientImpl{Client: clientMock}

	appConns := New(cl, log.NewNopLogger(), NopMetrics())

	err := appConns.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { cancel(); appConns.Wait() })

	select {
	case sig := <-c:
		t.Logf("signal %q successfully received", sig)
	case <-ctx.Done():
		t.Fatal("expected process to receive SIGTERM signal")
	}
}
