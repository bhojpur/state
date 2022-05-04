package light

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
	"time"

	"github.com/bhojpur/state/pkg/crypto"
)

// TrustOptions are the trust parameters needed when a new light client
// connects to the network or when an existing light client that has been
// offline for longer than the trusting period connects to the network.
//
// The expectation is the user will get this information from a trusted source
// like a validator, a friend, or a secure website. A more user friendly
// solution with trust tradeoffs is that we establish an https based protocol
// with a default end point that populates this information. Also an on-chain
// registry of roots-of-trust (e.g. on the Bhojpur.NET Hub) seems likely in the
// future.
type TrustOptions struct {
	// tp: trusting period.
	//
	// Should be significantly less than the unbonding period (e.g. unbonding
	// period = 3 weeks, trusting period = 2 weeks).
	//
	// More specifically, trusting period + time needed to check headers + time
	// needed to report and punish misbehavior should be less than the unbonding
	// period.
	Period time.Duration

	// Header's Height and Hash must both be provided to force the trusting of a
	// particular header.
	Height int64
	Hash   []byte
}

// ValidateBasic performs basic validation.
func (opts TrustOptions) ValidateBasic() error {
	if opts.Period <= 0 {
		return errors.New("negative or zero period")
	}
	if opts.Height <= 0 {
		return errors.New("negative or zero height")
	}
	if len(opts.Hash) != crypto.HashSize {
		return fmt.Errorf("expected hash size to be %d bytes, got %d bytes",
			crypto.HashSize,
			len(opts.Hash),
		)
	}
	return nil
}
