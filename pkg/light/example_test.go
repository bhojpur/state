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

	"github.com/bhojpur/state/example/kvstore"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/light"
	httpp "github.com/bhojpur/state/pkg/light/provider/http"
	dbs "github.com/bhojpur/state/pkg/light/store/db"
	rpctest "github.com/bhojpur/state/pkg/rpc/test"
)

// Manually getting light blocks and verifying them.
func TestExampleClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := rpctest.CreateConfig(t, "ExampleClient_VerifyLightBlockAtHeight")
	if err != nil {
		t.Fatal(err)
	}

	logger, err := log.NewDefaultLogger(log.LogFormatPlain, log.LogLevelInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Start a test application
	app := kvstore.NewApplication()

	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = closer(ctx) }()

	dbDir := t.TempDir()
	chainID := conf.ChainID()

	primary, err := httpp.New(chainID, conf.RPC.ListenAddress)
	if err != nil {
		t.Fatal(err)
	}

	// give Bhojpur State time to generate some blocks
	time.Sleep(5 * time.Second)

	block, err := primary.LightBlock(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Cleanup(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait for a few more blocks to be produced
	time.Sleep(2 * time.Second)

	// veify the block at height 3
	_, err = c.VerifyLightBlockAtHeight(ctx, 3, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	// retrieve light block at height 3
	_, err = c.TrustedLightBlock(3)
	if err != nil {
		t.Fatal(err)
	}

	// update to the latest height
	lb, err := c.Update(ctx, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("verified light block", "light-block", lb)
}
