package commands_test

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

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/cmd/client/commands"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/rpc/client/local"
	rpctest "github.com/bhojpur/state/pkg/rpc/test"
	e2e "github.com/bhojpur/state/test/e2e/app"
)

func TestRollbackIntegration(t *testing.T) {
	var height int64
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)
	cfg.BaseConfig.DBBackend = "goleveldb"

	app, err := e2e.NewApplication(e2e.DefaultConfig(dir))
	require.NoError(t, err)

	t.Run("First run", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		require.NoError(t, err)
		node, _, err := rpctest.StartBhojpurState(ctx, cfg, app, rpctest.SuppressStdout)
		require.NoError(t, err)
		require.True(t, node.IsRunning())

		time.Sleep(3 * time.Second)
		cancel()
		node.Wait()

		require.False(t, node.IsRunning())
	})
	t.Run("Rollback", func(t *testing.T) {
		time.Sleep(time.Second)
		require.NoError(t, app.Rollback())
		height, _, err = commands.RollbackState(cfg)
		require.NoError(t, err, "%d", height)
	})
	t.Run("Restart", func(t *testing.T) {
		require.True(t, height > 0, "%d", height)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		node2, _, err2 := rpctest.StartBhojpurState(ctx, cfg, app, rpctest.SuppressStdout)
		require.NoError(t, err2)
		t.Cleanup(node2.Wait)

		logger := log.NewNopLogger()

		client, err := local.New(logger, node2.(local.NodeService))
		require.NoError(t, err)

		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("failed to make progress after 20 seconds. Min height: %d", height)
			case <-ticker.C:
				status, err := client.Status(ctx)
				require.NoError(t, err)

				if status.SyncInfo.LatestBlockHeight > height {
					return
				}
			}
		}
	})

}
