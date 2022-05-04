package commands

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
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/spf13/cobra"

	"github.com/bhojpur/state/internal/libs/progressbar"
	"github.com/bhojpur/state/internal/state"
	"github.com/bhojpur/state/internal/state/indexer"
	"github.com/bhojpur/state/internal/state/indexer/sink/kv"
	"github.com/bhojpur/state/internal/state/indexer/sink/psql"
	"github.com/bhojpur/state/internal/store"
	abcipb "github.com/bhojpur/state/pkg/abci/types"
	cfgsvc "github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/os"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	"github.com/bhojpur/state/pkg/types"
)

const (
	reindexFailed = "event re-index failed: "
)

// MakeReindexEventCommand constructs a command to re-index events in a block height interval.
func MakeReindexEventCommand(conf *cfgsvc.Config, logger log.Logger) *cobra.Command {
	var (
		startHeight int64
		endHeight   int64
	)

	cmd := &cobra.Command{
		Use:   "reindex-event",
		Short: "reindex events to the event store backends",
		Long: `
reindex-event is an offline tooling to re-index block and tx events to the eventsinks,
you can run this command when the event store backend dropped/disconnected or you want to
replace the backend. The default start-height is 0, meaning the tooling will start
reindex from the base block height(inclusive); and the default end-height is 0, meaning
the tooling will reindex until the latest block height(inclusive). User can omit
either or both arguments.
	`,
		Example: `
	statectl reindex-event
	statectl reindex-event --start-height 2
	statectl reindex-event --end-height 10
	statectl reindex-event --start-height 2 --end-height 10
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, ss, err := loadStateAndBlockStore(conf)
			if err != nil {
				return fmt.Errorf("%s: %w", reindexFailed, err)
			}

			cvhArgs := checkValidHeightArgs{
				startHeight: startHeight,
				endHeight:   endHeight,
			}
			if err := checkValidHeight(bs, cvhArgs); err != nil {
				return fmt.Errorf("%s: %w", reindexFailed, err)
			}

			es, err := loadEventSinks(conf)
			if err != nil {
				return fmt.Errorf("%s: %w", reindexFailed, err)
			}

			riArgs := eventReIndexArgs{
				startHeight: startHeight,
				endHeight:   endHeight,
				sinks:       es,
				blockStore:  bs,
				stateStore:  ss,
			}
			if err := eventReIndex(cmd, riArgs); err != nil {
				return fmt.Errorf("%s: %w", reindexFailed, err)
			}

			logger.Info("event re-index finished")
			return nil
		},
	}

	cmd.Flags().Int64Var(&startHeight, "start-height", 0, "the block height would like to start for re-index")
	cmd.Flags().Int64Var(&endHeight, "end-height", 0, "the block height would like to finish for re-index")
	return cmd
}

func loadEventSinks(cfg *cfgsvc.Config) ([]indexer.EventSink, error) {
	// Check duplicated sinks.
	sinks := map[string]bool{}
	for _, s := range cfg.TxIndex.Indexer {
		sl := strings.ToLower(s)
		if sinks[sl] {
			return nil, errors.New("found duplicated sinks, please check the tx-index section in the config.toml")
		}
		sinks[sl] = true
	}

	eventSinks := []indexer.EventSink{}

	for k := range sinks {
		switch k {
		case string(indexer.NULL):
			return nil, errors.New("found null event sink, please check the tx-index section in the config.toml")
		case string(indexer.KV):
			store, err := cfgsvc.DefaultDBProvider(&cfgsvc.DBContext{ID: "tx_index", Config: cfg})
			if err != nil {
				return nil, err
			}
			eventSinks = append(eventSinks, kv.NewEventSink(store))
		case string(indexer.PSQL):
			conn := cfg.TxIndex.PsqlConn
			if conn == "" {
				return nil, errors.New("the psql connection settings cannot be empty")
			}
			es, err := psql.NewEventSink(conn, cfg.ChainID())
			if err != nil {
				return nil, err
			}
			eventSinks = append(eventSinks, es)
		default:
			return nil, errors.New("unsupported event sink type")
		}
	}

	if len(eventSinks) == 0 {
		return nil, errors.New("no proper event sink can do event re-indexing," +
			" please check the tx-index section in the config.toml")
	}

	if !indexer.IndexingEnabled(eventSinks) {
		return nil, fmt.Errorf("no event sink has been enabled")
	}

	return eventSinks, nil
}

func loadStateAndBlockStore(cfg *cfgsvc.Config) (*store.BlockStore, state.Store, error) {
	dbType := dbm.BackendType(cfg.DBBackend)

	if !os.FileExists(filepath.Join(cfg.DBDir(), "blockstore.db")) {
		return nil, nil, fmt.Errorf("no blockstore found in %v", cfg.DBDir())
	}

	// Get BlockStore
	blockStoreDB, err := dbm.NewDB("blockstore", dbType, cfg.DBDir())
	if err != nil {
		return nil, nil, err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	if !os.FileExists(filepath.Join(cfg.DBDir(), "state.db")) {
		return nil, nil, fmt.Errorf("no blockstore found in %v", cfg.DBDir())
	}

	// Get StateStore
	stateDB, err := dbm.NewDB("state", dbType, cfg.DBDir())
	if err != nil {
		return nil, nil, err
	}
	stateStore := state.NewStore(stateDB)

	return blockStore, stateStore, nil
}

type eventReIndexArgs struct {
	startHeight int64
	endHeight   int64
	sinks       []indexer.EventSink
	blockStore  state.BlockStore
	stateStore  state.Store
}

func eventReIndex(cmd *cobra.Command, args eventReIndexArgs) error {
	var bar progressbar.Bar
	bar.NewOption(args.startHeight-1, args.endHeight)

	fmt.Println("start re-indexing events:")
	defer bar.Finish()
	for i := args.startHeight; i <= args.endHeight; i++ {
		select {
		case <-cmd.Context().Done():
			return fmt.Errorf("event re-index terminated at height %d: %w", i, cmd.Context().Err())
		default:
			b := args.blockStore.LoadBlock(i)
			if b == nil {
				return fmt.Errorf("not able to load block at height %d from the blockstore", i)
			}

			r, err := args.stateStore.LoadABCIResponses(i)
			if err != nil {
				return fmt.Errorf("not able to load ABCI Response at height %d from the statestore", i)
			}

			e := types.EventDataNewBlockHeader{
				Header:              b.Header,
				NumTxs:              int64(len(b.Txs)),
				ResultFinalizeBlock: *r.FinalizeBlock,
			}

			var batch *indexer.Batch
			if e.NumTxs > 0 {
				batch = indexer.NewBatch(e.NumTxs)

				for i := range b.Data.Txs {
					tr := abcipb.TxResult{
						Height: b.Height,
						Index:  uint32(i),
						Tx:     b.Data.Txs[i],
						Result: (r.FinalizeBlock.TxResults[i]),
					}

					_ = batch.Add(&tr)
				}
			}

			for _, sink := range args.sinks {
				if err := sink.IndexBlockEvents(e); err != nil {
					return fmt.Errorf("block event re-index at height %d failed: %w", i, err)
				}

				if batch != nil {
					if err := sink.IndexTxEvents(batch.Ops); err != nil {
						return fmt.Errorf("tx event re-index at height %d failed: %w", i, err)
					}
				}
			}
		}

		bar.Play(i)
	}

	return nil
}

type checkValidHeightArgs struct {
	startHeight int64
	endHeight   int64
}

func checkValidHeight(bs state.BlockStore, args checkValidHeightArgs) error {
	base := bs.Base()

	if args.startHeight == 0 {
		args.startHeight = base
		fmt.Printf("set the start block height to the base height of the blockstore %d \n", base)
	}

	if args.startHeight < base {
		return fmt.Errorf("%s (requested start height: %d, base height: %d)",
			coretypes.ErrHeightNotAvailable, args.startHeight, base)
	}

	height := bs.Height()

	if args.startHeight > height {
		return fmt.Errorf(
			"%s (requested start height: %d, store height: %d)", coretypes.ErrHeightNotAvailable, args.startHeight, height)
	}

	if args.endHeight == 0 || args.endHeight > height {
		args.endHeight = height
		fmt.Printf("set the end block height to the latest height of the blockstore %d \n", height)
	}

	if args.endHeight < base {
		return fmt.Errorf(
			"%s (requested end height: %d, base height: %d)", coretypes.ErrHeightNotAvailable, args.endHeight, base)
	}

	if args.endHeight < args.startHeight {
		return fmt.Errorf(
			"%s (requested the end height: %d is less than the start height: %d)",
			coretypes.ErrInvalidRequest, args.startHeight, args.endHeight)
	}

	return nil
}
