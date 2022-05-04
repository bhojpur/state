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
	"context"

	"github.com/bhojpur/state/pkg/light/provider"
	"github.com/bhojpur/state/pkg/light/provider/http"
	"github.com/bhojpur/state/pkg/light/store"
)

// NewHTTPClient initiates an instance of a light client using HTTP addresses
// for both the primary provider and witnesses of the light client. A trusted
// header and hash must be passed to initialize the client.
//
// See all Option(s) for the additional configuration.
// See NewClient.
func NewHTTPClient(
	ctx context.Context,
	chainID string,
	trustOptions TrustOptions,
	primaryAddress string,
	witnessesAddresses []string,
	trustedStore store.Store,
	options ...Option) (*Client, error) {

	providers, err := providersFromAddresses(append(witnessesAddresses, primaryAddress), chainID)
	if err != nil {
		return nil, err
	}

	return NewClient(
		ctx,
		chainID,
		trustOptions,
		providers[len(providers)-1],
		providers[:len(providers)-1],
		trustedStore,
		options...)
}

func providersFromAddresses(addrs []string, chainID string) ([]provider.Provider, error) {
	providers := make([]provider.Provider, len(addrs))
	for idx, address := range addrs {
		p, err := http.New(chainID, address)
		if err != nil {
			return nil, err
		}
		providers[idx] = p
	}
	return providers, nil
}
