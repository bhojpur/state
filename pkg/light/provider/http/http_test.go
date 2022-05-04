package http_test

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/example/kvstore"
	"github.com/bhojpur/state/pkg/light/provider"
	lighthttp "github.com/bhojpur/state/pkg/light/provider/http"
	rpcclient "github.com/bhojpur/state/pkg/rpc/client"
	rpchttp "github.com/bhojpur/state/pkg/rpc/client/http"
	rpctest "github.com/bhojpur/state/pkg/rpc/test"
	"github.com/bhojpur/state/pkg/types"
)

func TestNewProvider(t *testing.T) {
	c, err := lighthttp.New("chain-test", "192.168.0.1:26657")
	require.NoError(t, err)
	require.Equal(t, c.ID(), "http{http://192.168.0.1:26657}")

	c, err = lighthttp.New("chain-test", "http://153.200.0.1:26657")
	require.NoError(t, err)
	require.Equal(t, c.ID(), "http{http://153.200.0.1:26657}")

	c, err = lighthttp.New("chain-test", "153.200.0.1")
	require.NoError(t, err)
	require.Equal(t, c.ID(), "http{http://153.200.0.1}")
}

func TestProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)

	// start a Bhojpur State node in the background to test against
	app := kvstore.NewApplication()
	app.RetainBlocks = 9
	_, closer, err := rpctest.StartBhojpurState(ctx, cfg, app)
	require.NoError(t, err)

	rpcAddr := cfg.RPC.ListenAddress
	genDoc, err := types.GenesisDocFromFile(cfg.GenesisFile())
	require.NoError(t, err)

	chainID := genDoc.ChainID
	t.Log("chainID:", chainID)

	c, err := rpchttp.New(rpcAddr)
	require.NoError(t, err)

	p := lighthttp.NewWithClient(chainID, c)
	require.NoError(t, err)
	require.NotNil(t, p)

	// let it produce some blocks
	err = rpcclient.WaitForHeight(ctx, c, 10, nil)
	require.NoError(t, err)

	// let's get the highest block
	lb, err := p.LightBlock(ctx, 0)
	require.NoError(t, err)
	assert.True(t, lb.Height < 9001, "height=%d", lb.Height)

	// let's check this is valid somehow
	assert.Nil(t, lb.ValidateBasic(chainID))

	// historical queries now work :)
	lower := lb.Height - 3
	lb, err = p.LightBlock(ctx, lower)
	require.NoError(t, err)
	assert.Equal(t, lower, lb.Height)

	// fetching missing heights (both future and pruned) should return appropriate errors
	lb, err = p.LightBlock(ctx, 9001)
	require.Error(t, err)
	require.Nil(t, lb)
	assert.ErrorIs(t, err, provider.ErrHeightTooHigh)

	lb, err = p.LightBlock(ctx, 1)
	require.Error(t, err)
	require.Nil(t, lb)
	assert.ErrorIs(t, err, provider.ErrLightBlockNotFound)

	// if the provider is unable to provide four more blocks then we should return
	// an unreliable peer error
	for i := 0; i < 4; i++ {
		_, err = p.LightBlock(ctx, 1)
	}
	assert.IsType(t, provider.ErrUnreliableProvider{}, err)

	// shut down Bhojpur State node
	require.NoError(t, closer(ctx))
	cancel()

	time.Sleep(10 * time.Second)
	lb, err = p.LightBlock(ctx, lower+2)
	// Either the connection should be refused, or the context canceled.
	require.Error(t, err)
	require.Nil(t, lb)
	if !errors.Is(err, provider.ErrConnectionClosed) && !errors.Is(err, context.Canceled) {
		assert.Fail(t, "Incorrect error", "wanted connection closed or context canceled, got %v", err)
	}
}
