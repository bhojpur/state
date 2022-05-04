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
	"context"
	"errors"
	"fmt"

	"github.com/bhojpur/state/pkg/rpc/coretypes"
)

// NetInfo returns network info.
func (env *Environment) NetInfo(ctx context.Context) (*coretypes.ResultNetInfo, error) {
	peerList := env.PeerManager.Peers()

	peers := make([]coretypes.Peer, 0, len(peerList))
	for _, peer := range peerList {
		addrs := env.PeerManager.Addresses(peer)
		if len(addrs) == 0 {
			continue
		}

		peers = append(peers, coretypes.Peer{
			ID:  peer,
			URL: addrs[0].String(),
		})
	}

	return &coretypes.ResultNetInfo{
		Listening: env.IsListening,
		Listeners: env.Listeners,
		NPeers:    len(peers),
		Peers:     peers,
	}, nil
}

// Genesis returns genesis file.
func (env *Environment) Genesis(ctx context.Context) (*coretypes.ResultGenesis, error) {
	if len(env.genChunks) > 1 {
		return nil, errors.New("genesis response is large, please use the genesis_chunked API instead")
	}

	return &coretypes.ResultGenesis{Genesis: env.GenDoc}, nil
}

func (env *Environment) GenesisChunked(ctx context.Context, req *coretypes.RequestGenesisChunked) (*coretypes.ResultGenesisChunk, error) {
	if env.genChunks == nil {
		return nil, fmt.Errorf("service configuration error, genesis chunks are not initialized")
	}

	if len(env.genChunks) == 0 {
		return nil, fmt.Errorf("service configuration error, there are no chunks")
	}

	id := int(req.Chunk)

	if id > len(env.genChunks)-1 {
		return nil, fmt.Errorf("there are %d chunks, %d is invalid", len(env.genChunks)-1, id)
	}

	return &coretypes.ResultGenesisChunk{
		TotalChunks: len(env.genChunks),
		ChunkNumber: id,
		Data:        env.genChunks[id],
	}, nil
}
