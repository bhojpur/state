package light_test

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
	"testing"
	"time"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/example/kvstore"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/light"
	"github.com/bhojpur/state/pkg/light/provider"
	httpp "github.com/bhojpur/state/pkg/light/provider/http"
	dbs "github.com/bhojpur/state/pkg/light/store/db"
	rpctest "github.com/bhojpur/state/pkg/rpc/test"
	"github.com/bhojpur/state/pkg/types"
)

// NOTE: these are ports of the tests from example_test.go but
// rewritten as more conventional tests.

// Automatically getting new headers and verifying them.
func TestClientIntegration_Update(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)

	logger := log.NewNopLogger()

	// Start a test application
	app := kvstore.NewApplication()
	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	require.NoError(t, err)
	defer func() { require.NoError(t, closer(ctx)) }()

	// give Bhojpur State time to generate some blocks
	time.Sleep(5 * time.Second)

	dbDir := t.TempDir()
	chainID := conf.ChainID()

	primary, err := httpp.New(chainID, conf.RPC.ListenAddress)
	require.NoError(t, err)

	// give Bhojpur State time to generate some blocks
	block, err := waitForBlock(ctx, primary, 2)
	require.NoError(t, err)

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	require.NoError(t, err)

	c, err := light.NewClient(
		ctx,
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: 2,
			Hash:   block.Hash(),
		},
		primary,
		nil,
		dbs.New(db),
		light.Logger(logger),
	)
	require.NoError(t, err)

	defer func() { require.NoError(t, c.Cleanup()) }()

	// ensure Bhojpur State is at height 3 or higher
	_, err = waitForBlock(ctx, primary, 3)
	require.NoError(t, err)

	h, err := c.Update(ctx, time.Now())
	require.NoError(t, err)
	require.NotNil(t, h)

	require.True(t, h.Height > 2)
}

// Manually getting light blocks and verifying them.
func TestClientIntegration_VerifyLightBlockAtHeight(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)

	logger := log.NewNopLogger()

	// Start a test application
	app := kvstore.NewApplication()

	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	require.NoError(t, err)
	defer func() { require.NoError(t, closer(ctx)) }()

	dbDir := t.TempDir()
	chainID := conf.ChainID()

	primary, err := httpp.New(chainID, conf.RPC.ListenAddress)
	require.NoError(t, err)

	// give Bhojpur State time to generate some blocks
	block, err := waitForBlock(ctx, primary, 2)
	require.NoError(t, err)

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	require.NoError(t, err)

	c, err := light.NewClient(ctx,
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: 2,
			Hash:   block.Hash(),
		},
		primary,
		nil,
		dbs.New(db),
		light.Logger(logger),
	)
	require.NoError(t, err)

	defer func() { require.NoError(t, c.Cleanup()) }()

	// ensure Bhojpur State is at height 3 or higher
	_, err = waitForBlock(ctx, primary, 3)
	require.NoError(t, err)

	_, err = c.VerifyLightBlockAtHeight(ctx, 3, time.Now())
	require.NoError(t, err)

	h, err := c.TrustedLightBlock(3)
	require.NoError(t, err)

	require.EqualValues(t, 3, h.Height)
}

func waitForBlock(ctx context.Context, p provider.Provider, height int64) (*types.LightBlock, error) {
	for {
		block, err := p.LightBlock(ctx, height)
		switch err {
		case nil:
			return block, nil
		// node isn't running yet, wait 1 second and repeat
		case provider.ErrNoResponse, provider.ErrHeightTooHigh:
			timer := time.NewTimer(1 * time.Second)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-timer.C:
			}
		default:
			return nil, err
		}
	}
}

func TestClientStatusRPC(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)

	// Start a test application
	app := kvstore.NewApplication()

	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	require.NoError(t, err)
	defer func() { require.NoError(t, closer(ctx)) }()

	dbDir := t.TempDir()
	chainID := conf.ChainID()

	primary, err := httpp.New(chainID, conf.RPC.ListenAddress)
	require.NoError(t, err)

	// give Bhojpur State time to generate some blocks
	block, err := waitForBlock(ctx, primary, 2)
	require.NoError(t, err)

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	require.NoError(t, err)

	// In order to not create a full testnet we create the light client with no witnesses
	// and only verify the primary IP address.
	witnesses := []provider.Provider{}

	c, err := light.NewClient(ctx,
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: 2,
			Hash:   block.Hash(),
		},
		primary,
		witnesses,
		dbs.New(db),
		light.Logger(log.NewNopLogger()),
	)
	require.NoError(t, err)

	defer func() { require.NoError(t, c.Cleanup()) }()

	lightStatus := c.Status(ctx)

	// Verify primary IP
	require.True(t, lightStatus.PrimaryID == primary.ID())

	// Verify that number of peers is equal to number of witnesses  (+ 1 if the primary is not a witness)
	require.Equal(t, len(witnesses)+1*primaryNotInWitnessList(witnesses, primary), lightStatus.NumPeers)

	// Verify that the last trusted hash returned matches the stored hash of the trusted
	// block at the last trusted height.
	blockAtTrustedHeight, err := c.TrustedLightBlock(lightStatus.LastTrustedHeight)
	require.NoError(t, err)

	require.EqualValues(t, lightStatus.LastTrustedHash, blockAtTrustedHeight.Hash())

}

// If the primary is not in the witness list, we will return 1
// Otherwise, return 0
func primaryNotInWitnessList(witnesses []provider.Provider, primary provider.Provider) int {
	for _, el := range witnesses {
		if el == primary {
			return 0
		}
	}
	return 1
}
