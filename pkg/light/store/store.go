package store

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

import "github.com/bhojpur/state/pkg/types"

// Store is anything that can persistently store headers.
type Store interface {
	// SaveSignedHeaderAndValidatorSet saves a SignedHeader (h: sh.Height) and a
	// ValidatorSet (h: sh.Height).
	//
	// height must be > 0.
	SaveLightBlock(lb *types.LightBlock) error

	// DeleteSignedHeaderAndValidatorSet deletes SignedHeader (h: height) and
	// ValidatorSet (h: height).
	//
	// height must be > 0.
	DeleteLightBlock(height int64) error

	// LightBlock returns the LightBlock that corresponds to the given
	// height.
	//
	// height must be > 0.
	//
	// If LightBlock is not found, ErrLightBlockNotFound is returned.
	LightBlock(height int64) (*types.LightBlock, error)

	// LastLightBlockHeight returns the last (newest) LightBlock height.
	//
	// If the store is empty, -1 and nil error are returned.
	LastLightBlockHeight() (int64, error)

	// FirstLightBlockHeight returns the first (oldest) LightBlock height.
	//
	// If the store is empty, -1 and nil error are returned.
	FirstLightBlockHeight() (int64, error)

	// LightBlockBefore returns the LightBlock before a certain height.
	//
	// height must be > 0 && <= LastLightBlockHeight.
	LightBlockBefore(height int64) (*types.LightBlock, error)

	// Prune removes headers & the associated validator sets when Store reaches a
	// defined size (number of header & validator set pairs).
	Prune(size uint16) error

	// Size returns a number of currently existing header & validator set pairs.
	Size() uint16
}
