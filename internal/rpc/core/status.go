package core

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
	"fmt"
	"time"

	libytes "github.com/bhojpur/state/pkg/libs/bytes"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	"github.com/bhojpur/state/pkg/types"
)

// Status returns Bhojpur State status including node info, pubkey, latest block
// hash, app hash, block height, current max peer block height, and time.
func (env *Environment) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	var (
		earliestBlockHeight   int64
		earliestBlockHash     libytes.HexBytes
		earliestAppHash       libytes.HexBytes
		earliestBlockTimeNano int64
	)

	if earliestBlockMeta := env.BlockStore.LoadBaseMeta(); earliestBlockMeta != nil {
		earliestBlockHeight = earliestBlockMeta.Header.Height
		earliestAppHash = earliestBlockMeta.Header.AppHash
		earliestBlockHash = earliestBlockMeta.BlockID.Hash
		earliestBlockTimeNano = earliestBlockMeta.Header.Time.UnixNano()
	}

	var (
		latestBlockHash     libytes.HexBytes
		latestAppHash       libytes.HexBytes
		latestBlockTimeNano int64

		latestHeight = env.BlockStore.Height()
	)

	if latestHeight != 0 {
		if latestBlockMeta := env.BlockStore.LoadBlockMeta(latestHeight); latestBlockMeta != nil {
			latestBlockHash = latestBlockMeta.BlockID.Hash
			latestAppHash = latestBlockMeta.Header.AppHash
			latestBlockTimeNano = latestBlockMeta.Header.Time.UnixNano()
		}
	}

	// Return the very last voting power, not the voting power of this validator
	// during the last block.
	var votingPower int64
	if val := env.validatorAtHeight(env.latestUncommittedHeight()); val != nil {
		votingPower = val.VotingPower
	}
	validatorInfo := coretypes.ValidatorInfo{}
	if env.PubKey != nil {
		validatorInfo = coretypes.ValidatorInfo{
			Address:     env.PubKey.Address(),
			PubKey:      env.PubKey,
			VotingPower: votingPower,
		}
	}

	var applicationInfo coretypes.ApplicationInfo
	if abciInfo, err := env.ABCIInfo(ctx); err == nil {
		applicationInfo.Version = fmt.Sprint(abciInfo.Response.AppVersion)
	}

	result := &coretypes.ResultStatus{
		NodeInfo:        env.NodeInfo,
		ApplicationInfo: applicationInfo,
		SyncInfo: coretypes.SyncInfo{
			LatestBlockHash:     latestBlockHash,
			LatestAppHash:       latestAppHash,
			LatestBlockHeight:   latestHeight,
			LatestBlockTime:     time.Unix(0, latestBlockTimeNano),
			EarliestBlockHash:   earliestBlockHash,
			EarliestAppHash:     earliestAppHash,
			EarliestBlockHeight: earliestBlockHeight,
			EarliestBlockTime:   time.Unix(0, earliestBlockTimeNano),
			// this should start as true, if consensus
			// hasn't started yet, and then flip to false
			// (or true,) depending on what's actually
			// happening.
			CatchingUp: true,
		},
		ValidatorInfo: validatorInfo,
	}

	if env.ConsensusReactor != nil {
		result.SyncInfo.CatchingUp = env.ConsensusReactor.WaitSync()
	}

	if env.BlockSyncReactor != nil {
		result.SyncInfo.MaxPeerBlockHeight = env.BlockSyncReactor.GetMaxPeerBlockHeight()
		result.SyncInfo.TotalSyncedTime = env.BlockSyncReactor.GetTotalSyncedTime()
		result.SyncInfo.RemainingTime = env.BlockSyncReactor.GetRemainingSyncTime()
	}

	if env.StateSyncMetricer != nil {
		result.SyncInfo.TotalSnapshots = env.StateSyncMetricer.TotalSnapshots()
		result.SyncInfo.ChunkProcessAvgTime = env.StateSyncMetricer.ChunkProcessAvgTime()
		result.SyncInfo.SnapshotHeight = env.StateSyncMetricer.SnapshotHeight()
		result.SyncInfo.SnapshotChunksCount = env.StateSyncMetricer.SnapshotChunksCount()
		result.SyncInfo.SnapshotChunksTotal = env.StateSyncMetricer.SnapshotChunksTotal()
		result.SyncInfo.BackFilledBlocks = env.StateSyncMetricer.BackFilledBlocks()
		result.SyncInfo.BackFillBlocksTotal = env.StateSyncMetricer.BackFillBlocksTotal()
	}

	return result, nil
}

func (env *Environment) validatorAtHeight(h int64) *types.Validator {
	valsWithH, err := env.StateStore.LoadValidators(h)
	if err != nil {
		return nil
	}
	if env.ConsensusState == nil {
		return nil
	}
	if env.PubKey == nil {
		return nil
	}
	privValAddress := env.PubKey.Address()

	// If we're still at height h, search in the current validator set.
	lastBlockHeight, vals := env.ConsensusState.GetValidators()
	if lastBlockHeight == h {
		for _, val := range vals {
			if bytes.Equal(val.Address, privValAddress) {
				return val
			}
		}
	}

	_, val := valsWithH.GetByAddress(privValAddress)
	return val
}
