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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/bhojpur/state/pkg/abci/types"
	librand "github.com/bhojpur/state/pkg/libs/rand"
	"github.com/bhojpur/state/pkg/rpc/client"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	"github.com/bhojpur/state/pkg/types"
)

const waitForEventTimeout = 2 * time.Second

// MakeTxKV returns a text transaction, allong with expected key, value pair
func MakeTxKV() ([]byte, []byte, []byte) {
	k := []byte(librand.Str(8))
	v := []byte(librand.Str(8))
	return k, v, append(k, append([]byte("="), v...)...)
}

func testTxEventsSent(ctx context.Context, t *testing.T, broadcastMethod string, c client.Client) {
	t.Helper()
	// make the tx
	_, _, tx := MakeTxKV()

	// send
	done := make(chan struct{})
	go func() {
		defer close(done)
		var (
			txres *coretypes.ResultBroadcastTx
			err   error
		)
		switch broadcastMethod {
		case "async":
			txres, err = c.BroadcastTxAsync(ctx, tx)
		case "sync":
			txres, err = c.BroadcastTxSync(ctx, tx)
		default:
			require.FailNowf(t, "Unknown broadcastMethod %s", broadcastMethod)
		}
		if assert.NoError(t, err) {
			assert.Equal(t, txres.Code, abci.CodeTypeOK)
		}
	}()

	// and wait for confirmation
	ectx, cancel := context.WithTimeout(ctx, waitForEventTimeout)
	defer cancel()

	// Wait for the transaction we sent to be confirmed.
	query := fmt.Sprintf(`tm.event = '%s' AND tx.hash = '%X'`,
		types.EventTxValue, types.Tx(tx).Hash())
	evt, err := client.WaitForOneEvent(ectx, c, query)
	require.NoError(t, err)

	// and make sure it has the proper info
	txe, ok := evt.(types.EventDataTx)
	require.True(t, ok)

	// make sure this is the proper tx
	require.EqualValues(t, tx, txe.Tx)
	require.True(t, txe.Result.IsOK())
	<-done
}
