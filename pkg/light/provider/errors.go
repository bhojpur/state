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
	"errors"
	"fmt"
)

var (
	// ErrHeightTooHigh is returned when the height is higher than the last
	// block that the provider has. The light client will not remove the provider
	ErrHeightTooHigh = errors.New("height requested is too high")
	// ErrLightBlockNotFound is returned when a provider can't find the
	// requested header (i.e. it has been pruned).
	// The light client will not remove the provider
	ErrLightBlockNotFound = errors.New("light block not found")
	// ErrNoResponse is returned if the provider doesn't respond to the
	// request in a given time. The light client will not remove the provider
	ErrNoResponse = errors.New("client failed to respond")
	// ErrConnectionClosed is returned if the provider closes the connection.
	// In this case we remove the provider.
	ErrConnectionClosed = errors.New("client closed connection")
)

// ErrBadLightBlock is returned when a provider returns an invalid
// light block. The light client will remove the provider.
type ErrBadLightBlock struct {
	Reason error
}

func (e ErrBadLightBlock) Error() string {
	return fmt.Sprintf("client provided bad signed header: %v", e.Reason)
}

func (e ErrBadLightBlock) Unwrap() error { return e.Reason }

// ErrUnreliableProvider is a generic error that indicates that the provider isn't
// behaving in a reliable manner to the light client. The light client will
// remove the provider
type ErrUnreliableProvider struct {
	Reason error
}

func (e ErrUnreliableProvider) Error() string {
	return fmt.Sprintf("client deemed unreliable: %v", e.Reason)
}

func (e ErrUnreliableProvider) Unwrap() error { return e.Reason }
