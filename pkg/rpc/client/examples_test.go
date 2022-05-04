package client_test

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
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/example/kvstore"
	rpchttp "github.com/bhojpur/state/pkg/rpc/client/http"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	rpctest "github.com/bhojpur/state/pkg/rpc/test"
)

func TestHTTPSimple(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a Bhojpur State node (and kvstore) in the background to test against
	app := kvstore.NewApplication()
	conf, err := rpctest.CreateConfig(t, "ExampleHTTP_simple")
	require.NoError(t, err)

	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	if err != nil {
		log.Fatal(err) //nolint:gocritic
	}
	defer func() { _ = closer(ctx) }()

	// Create our RPC client
	rpcAddr := conf.RPC.ListenAddress
	c, err := rpchttp.New(rpcAddr)
	require.NoError(t, err)

	// Create a transaction
	k := []byte("name")
	v := []byte("satoshi")
	tx := append(k, append([]byte("="), v...)...)

	// Broadcast the transaction and wait for it to commit (rather use
	// c.BroadcastTxSync though in production).
	bres, err := c.BroadcastTxCommit(ctx, tx)
	require.NoError(t, err)
	if err != nil {
		log.Fatal(err)
	}
	if bres.CheckTx.IsErr() || bres.TxResult.IsErr() {
		log.Fatal("BroadcastTxCommit transaction failed")
	}

	// Now try to fetch the value for the key
	qres, err := c.ABCIQuery(ctx, "/key", k)
	require.NoError(t, err)
	require.False(t, qres.Response.IsErr(), "ABCIQuery failed")
	require.True(t, bytes.Equal(qres.Response.Key, k),
		"returned key does not match queried key")
	require.True(t, bytes.Equal(qres.Response.Value, v),
		"returned value does not match sent value [%s]", string(v))

	assert.Equal(t, "name=satoshi", string(tx), "sent tx")
	assert.Equal(t, "name", string(qres.Response.Key), "queried for")
	assert.Equal(t, "satoshi", string(qres.Response.Value), "got value")
}

func TestHTTPBatching(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a Bhojpur State node (and kvstore) in the background to test against
	app := kvstore.NewApplication()
	conf, err := rpctest.CreateConfig(t, "ExampleHTTP_batching")
	require.NoError(t, err)

	_, closer, err := rpctest.StartBhojpurState(ctx, conf, app, rpctest.SuppressStdout)
	if err != nil {
		log.Fatal(err) //nolint:gocritic
	}
	defer func() { _ = closer(ctx) }()

	rpcAddr := conf.RPC.ListenAddress
	c, err := rpchttp.NewWithClient(rpcAddr, http.DefaultClient)
	require.NoError(t, err)

	// Create our two transactions
	k1 := []byte("firstName")
	v1 := []byte("satoshi")
	tx1 := append(k1, append([]byte("="), v1...)...)

	k2 := []byte("lastName")
	v2 := []byte("nakamoto")
	tx2 := append(k2, append([]byte("="), v2...)...)

	txs := [][]byte{tx1, tx2}

	// Create a new batch
	batch := c.NewBatch()

	// Queue up our transactions
	for _, tx := range txs {
		// Broadcast the transaction and wait for it to commit (rather use
		// c.BroadcastTxSync though in production).
		_, err := batch.BroadcastTxSync(ctx, tx)
		require.NoError(t, err)
	}

	// Send the batch of 2 transactions
	_, err = batch.Send(ctx)
	require.NoError(t, err)

	// wait for the transaction to land, we could poll more for
	// the transactions to land definitively.
	require.Eventually(t,
		func() bool {
			// Now let's query for the original results as a batch
			exists := 0
			for _, key := range [][]byte{k1, k2} {
				_, err := batch.ABCIQuery(ctx, "/key", key)
				if err == nil {
					exists++

				}
			}
			return exists == 2
		},
		10*time.Second,
		time.Second,
	)

	// Send the 2 queries and keep the results
	results, err := batch.Send(ctx)
	require.NoError(t, err)

	require.Len(t, results, 2)
	// Each result in the returned list is the deserialized result of each
	// respective ABCIQuery response
	for _, result := range results {
		qr, ok := result.(*coretypes.ResultABCIQuery)
		require.True(t, ok, "invalid result type from ABCIQuery request")

		switch string(qr.Response.Key) {
		case "firstName":
			require.Equal(t, "satoshi", string(qr.Response.Value))
		case "lastName":
			require.Equal(t, "nakamoto", string(qr.Response.Value))
		default:
			t.Fatalf("encountered unknown key %q", string(qr.Response.Key))
		}
	}
}
