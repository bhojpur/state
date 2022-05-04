package provider

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

	"github.com/bhojpur/state/pkg/types"
)

//go:generate ../../../scripts/mockery_generate.sh Provider

// Provider provides information for the light client to sync (verification
// happens in the client).
type Provider interface {
	// LightBlock returns the LightBlock that corresponds to the given
	// height.
	//
	// 0 - the latest.
	// height must be >= 0.
	//
	// If the provider fails to fetch the LightBlock due to the IO or other
	// issues, an error will be returned.
	// If there's no LightBlock for the given height, ErrLightBlockNotFound
	// error is returned.
	LightBlock(ctx context.Context, height int64) (*types.LightBlock, error)

	// ReportEvidence reports an evidence of misbehavior.
	ReportEvidence(context.Context, types.Evidence) error

	// Returns the ID of a provider. For RPC providers it returns the IP address of the client
	// For p2p providers it returns a combination of NodeID and IP address
	ID() string
}
