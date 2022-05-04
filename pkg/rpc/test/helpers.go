package rpctest

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
	"fmt"
	"os"
	"testing"
	"time"

	abciclient "github.com/bhojpur/state/pkg/abci/client"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	libnet "github.com/bhojpur/state/pkg/libs/net"
	"github.com/bhojpur/state/pkg/libs/service"
	"github.com/bhojpur/state/pkg/node"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	rpcclient "github.com/bhojpur/state/pkg/rpc/jsonrpc/client"
)

// Options helps with specifying some parameters for our RPC testing for greater
// control.
type Options struct {
	suppressStdout bool
}

// waitForRPC connects to the RPC service and blocks until a /status call succeeds.
func waitForRPC(ctx context.Context, conf *config.Config) {
	laddr := conf.RPC.ListenAddress
	client, err := rpcclient.New(laddr)
	if err != nil {
		panic(err)
	}
	result := new(coretypes.ResultStatus)
	for {
		err := client.Call(ctx, "status", map[string]interface{}{}, result)
		if err == nil {
			return
		}

		fmt.Println("error", err)
		time.Sleep(time.Millisecond)
	}
}

func randPort() int {
	port, err := libnet.GetFreePort()
	if err != nil {
		panic(err)
	}
	return port
}

// makeAddrs constructs local listener addresses for node services.  This
// implementation uses random ports so test instances can run concurrently.
func makeAddrs() (p2pAddr, rpcAddr string) {
	const addrTemplate = "tcp://127.0.0.1:%d"
	return fmt.Sprintf(addrTemplate, randPort()), fmt.Sprintf(addrTemplate, randPort())
}

func CreateConfig(t *testing.T, testName string) (*config.Config, error) {
	c, err := config.ResetTestRoot(t.TempDir(), testName)
	if err != nil {
		return nil, err
	}

	p2pAddr, rpcAddr := makeAddrs()
	c.P2P.ListenAddress = p2pAddr
	c.RPC.ListenAddress = rpcAddr
	c.RPC.EventLogWindowSize = 5 * time.Minute
	c.Consensus.WalPath = "rpc-test"
	c.RPC.CORSAllowedOrigins = []string{"https://bhojpur.net/"}
	return c, nil
}

type ServiceCloser func(context.Context) error

func StartBhojpurState(
	ctx context.Context,
	conf *config.Config,
	app abci.Application,
	opts ...func(*Options),
) (service.Service, ServiceCloser, error) {
	ctx, cancel := context.WithCancel(ctx)

	nodeOpts := &Options{}
	for _, opt := range opts {
		opt(nodeOpts)
	}
	var logger log.Logger
	if nodeOpts.suppressStdout {
		logger = log.NewNopLogger()
	} else {
		var err error
		logger, err = log.NewDefaultLogger(log.LogFormatPlain, log.LogLevelInfo)
		if err != nil {
			return nil, func(_ context.Context) error { cancel(); return nil }, err
		}

	}
	papp := abciclient.NewLocalClient(logger, app)
	tmNode, err := node.New(ctx, conf, logger, papp, nil)
	if err != nil {
		return nil, func(_ context.Context) error { cancel(); return nil }, err
	}

	err = tmNode.Start(ctx)
	if err != nil {
		return nil, func(_ context.Context) error { cancel(); return nil }, err
	}

	waitForRPC(ctx, conf)

	if !nodeOpts.suppressStdout {
		fmt.Println("Bhojpur State machine running!")
	}

	return tmNode, func(ctx context.Context) error {
		cancel()
		tmNode.Wait()
		os.RemoveAll(conf.RootDir)
		return nil
	}, nil
}

// SuppressStdout is an option that tries to make sure the RPC test Bhojpur State
// node doesn't log anything to stdout.
func SuppressStdout(o *Options) {
	o.suppressStdout = true
}
