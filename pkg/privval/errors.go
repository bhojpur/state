package privval

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

// EndpointTimeoutError occurs when endpoint times out.
type EndpointTimeoutError struct{}

// Implement the net.Error interface.
func (e EndpointTimeoutError) Error() string   { return "endpoint connection timed out" }
func (e EndpointTimeoutError) Timeout() bool   { return true }
func (e EndpointTimeoutError) Temporary() bool { return true }

// Socket errors.
var (
	ErrConnectionTimeout  = EndpointTimeoutError{}
	ErrNoConnection       = errors.New("endpoint is not connected")
	ErrReadTimeout        = errors.New("endpoint read timed out")
	ErrUnexpectedResponse = errors.New("empty response")
	ErrWriteTimeout       = errors.New("endpoint write timed out")
)

// RemoteSignerError allows (remote) validators to include meaningful error
// descriptions in their reply.
type RemoteSignerError struct {
	// TODO: create an enum of known errors
	Code        int
	Description string
}

func (e *RemoteSignerError) Error() string {
	return fmt.Sprintf("signerEndpoint returned error #%d: %s", e.Code, e.Description)
}
