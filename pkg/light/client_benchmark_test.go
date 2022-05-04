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
	"errors"
	"testing"
	"time"

	dbm "github.com/bhojpur/state/pkg/database"

	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/light"
	"github.com/bhojpur/state/pkg/light/provider"
	dbs "github.com/bhojpur/state/pkg/light/store/db"
	"github.com/bhojpur/state/pkg/types"
)

// NOTE: block is produced every minute. Make sure the verification time
// provided in the function call is correct for the size of the blockchain. The
// benchmarking may take some time hence it can be more useful to set the time
// or the amount of iterations use the flag -benchtime t -> i.e. -benchtime 5m
// or -benchtime 100x.
//
// Remember that none of these benchmarks account for network latency.
var ()

type providerBenchmarkImpl struct {
	currentHeight int64
	blocks        map[int64]*types.LightBlock
}

func newProviderBenchmarkImpl(headers map[int64]*types.SignedHeader,
	vals map[int64]*types.ValidatorSet) provider.Provider {
	impl := providerBenchmarkImpl{
		blocks: make(map[int64]*types.LightBlock, len(headers)),
	}
	for height, header := range headers {
		if height > impl.currentHeight {
			impl.currentHeight = height
		}
		impl.blocks[height] = &types.LightBlock{
			SignedHeader: header,
			ValidatorSet: vals[height],
		}
	}
	return &impl
}

func (impl *providerBenchmarkImpl) LightBlock(ctx context.Context, height int64) (*types.LightBlock, error) {
	if height == 0 {
		return impl.blocks[impl.currentHeight], nil
	}
	lb, ok := impl.blocks[height]
	if !ok {
		return nil, provider.ErrLightBlockNotFound
	}
	return lb, nil
}

func (impl *providerBenchmarkImpl) ReportEvidence(_ context.Context, _ types.Evidence) error {
	return errors.New("not implemented")
}

// provierBenchmarkImpl does not have an ID iteself.
// Thus we return a sample string
func (impl *providerBenchmarkImpl) ID() string { return "ip-not-defined.com" }

func BenchmarkSequence(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headers, vals, _ := genLightBlocksWithKeys(b, 1000, 100, 1, bTime)
	benchmarkFullNode := newProviderBenchmarkImpl(headers, vals)
	genesisBlock, _ := benchmarkFullNode.LightBlock(ctx, 1)

	logger := log.NewTestingLogger(b)

	c, err := light.NewClient(
		ctx,
		chainID,
		light.TrustOptions{
			Period: 24 * time.Hour,
			Height: 1,
			Hash:   genesisBlock.Hash(),
		},
		benchmarkFullNode,
		nil,
		dbs.New(dbm.NewMemDB()),
		light.Logger(logger),
		light.SequentialVerification(),
	)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err = c.VerifyLightBlockAtHeight(ctx, 1000, bTime.Add(1000*time.Minute))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBisection(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headers, vals, _ := genLightBlocksWithKeys(b, 1000, 100, 1, bTime)
	benchmarkFullNode := newProviderBenchmarkImpl(headers, vals)
	genesisBlock, _ := benchmarkFullNode.LightBlock(ctx, 1)

	logger := log.NewTestingLogger(b)

	c, err := light.NewClient(
		context.Background(),
		chainID,
		light.TrustOptions{
			Period: 24 * time.Hour,
			Height: 1,
			Hash:   genesisBlock.Hash(),
		},
		benchmarkFullNode,
		nil,
		dbs.New(dbm.NewMemDB()),
		light.Logger(logger),
	)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err = c.VerifyLightBlockAtHeight(ctx, 1000, bTime.Add(1000*time.Minute))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBackwards(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headers, vals, _ := genLightBlocksWithKeys(b, 1000, 100, 1, bTime)
	benchmarkFullNode := newProviderBenchmarkImpl(headers, vals)
	trustedBlock, _ := benchmarkFullNode.LightBlock(ctx, 0)

	logger := log.NewTestingLogger(b)

	c, err := light.NewClient(
		ctx,
		chainID,
		light.TrustOptions{
			Period: 24 * time.Hour,
			Height: trustedBlock.Height,
			Hash:   trustedBlock.Hash(),
		},
		benchmarkFullNode,
		nil,
		dbs.New(dbm.NewMemDB()),
		light.Logger(logger),
	)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err = c.VerifyLightBlockAtHeight(ctx, 1, bTime)
		if err != nil {
			b.Fatal(err)
		}
	}

}