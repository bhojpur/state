package indexer_test

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
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/adlio/schema"
	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/eventbus"
	"github.com/bhojpur/state/internal/state/indexer"
	"github.com/bhojpur/state/internal/state/indexer/sink/kv"
	"github.com/bhojpur/state/internal/state/indexer/sink/psql"
	abcipb "github.com/bhojpur/state/pkg/abci/types"
	liblog "github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/types"

	// Register the Postgre database driver.
	_ "github.com/lib/pq"
)

var psqldb *sql.DB
var resource *dockertest.Resource
var pSink indexer.EventSink

var (
	user     = "postgres"
	password = "secret"
	port     = "5432"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dbName   = "postgres"
)

func TestIndexerServiceIndexesBlocks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := liblog.NewNopLogger()
	// event bus
	eventBus := eventbus.NewDefault(logger)
	err := eventBus.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(eventBus.Wait)

	assert.False(t, indexer.KVSinkEnabled([]indexer.EventSink{}))
	assert.False(t, indexer.IndexingEnabled([]indexer.EventSink{}))

	// event sink setup
	pool := setupDB(t)

	store := dbm.NewMemDB()
	eventSinks := []indexer.EventSink{kv.NewEventSink(store), pSink}
	assert.True(t, indexer.KVSinkEnabled(eventSinks))
	assert.True(t, indexer.IndexingEnabled(eventSinks))

	service := indexer.NewService(indexer.ServiceArgs{
		Logger:   logger,
		Sinks:    eventSinks,
		EventBus: eventBus,
	})
	require.NoError(t, service.Start(ctx))
	t.Cleanup(service.Wait)

	// publish block with txs
	err = eventBus.PublishEventNewBlockHeader(types.EventDataNewBlockHeader{
		Header: types.Header{Height: 1},
		NumTxs: int64(2),
	})
	require.NoError(t, err)
	txResult1 := &abcipb.TxResult{
		Height: 1,
		Index:  uint32(0),
		Tx:     types.Tx("foo"),
		Result: abcipb.ExecTxResult{Code: 0},
	}
	err = eventBus.PublishEventTx(types.EventDataTx{TxResult: *txResult1})
	require.NoError(t, err)
	txResult2 := &abcipb.TxResult{
		Height: 1,
		Index:  uint32(1),
		Tx:     types.Tx("bar"),
		Result: abcipb.ExecTxResult{Code: 0},
	}
	err = eventBus.PublishEventTx(types.EventDataTx{TxResult: *txResult2})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	res, err := eventSinks[0].GetTxByHash(types.Tx("foo").Hash())
	require.NoError(t, err)
	require.Equal(t, txResult1, res)

	ok, err := eventSinks[0].HasBlock(1)
	require.NoError(t, err)
	require.True(t, ok)

	res, err = eventSinks[0].GetTxByHash(types.Tx("bar").Hash())
	require.NoError(t, err)
	require.Equal(t, txResult2, res)

	assert.Nil(t, teardown(t, pool))
}

func readSchema() ([]*schema.Migration, error) {
	filename := "./sink/psql/schema.sql"
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read sql file from '%s': %w", filename, err)
	}

	mg := &schema.Migration{}
	mg.ID = time.Now().Local().String() + " db schema"
	mg.Script = string(contents)
	return append([]*schema.Migration{}, mg), nil
}

func resetDB(t *testing.T) {
	q := "DROP TABLE IF EXISTS block_events,tx_events,tx_results"
	_, err := psqldb.Exec(q)
	assert.NoError(t, err)

	q = "DROP TYPE IF EXISTS block_event_type"
	_, err = psqldb.Exec(q)
	assert.NoError(t, err)
}

func setupDB(t *testing.T) *dockertest.Pool {
	t.Helper()
	pool, err := dockertest.NewPool(os.Getenv("DOCKER_URL"))
	assert.NoError(t, err)
	if _, err := pool.Client.Info(); err != nil {
		t.Skipf("WARNING: Docker is not available: %v [skipping this test]", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{port},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	assert.NoError(t, err)

	// Set the container to expire in a minute to avoid orphaned containers
	// hanging around
	_ = resource.Expire(60)

	conn := fmt.Sprintf(dsn, user, password, resource.GetPort(port+"/tcp"), dbName)

	assert.NoError(t, pool.Retry(func() error {
		sink, err := psql.NewEventSink(conn, "test-chainID")
		if err != nil {
			return err
		}

		pSink = sink
		psqldb = sink.DB()
		return psqldb.Ping()
	}))

	resetDB(t)

	sm, err := readSchema()
	assert.NoError(t, err)

	migrator := schema.NewMigrator()
	err = migrator.Apply(psqldb, sm)
	assert.NoError(t, err)

	return pool
}

func teardown(t *testing.T, pool *dockertest.Pool) error {
	t.Helper()
	// When you're done, kill and remove the container
	assert.Nil(t, pool.Purge(resource))
	return psqldb.Close()
}
