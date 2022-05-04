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
	"context"
	"time"

	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/service"
)

const (
	defaultMaxDialRetries        = 10
	defaultRetryWaitMilliseconds = 100
)

// SignerServiceEndpointOption sets an optional parameter on the SignerDialerEndpoint.
type SignerServiceEndpointOption func(*SignerDialerEndpoint)

// SignerDialerEndpointTimeoutReadWrite sets the read and write timeout for
// connections from client processes.
func SignerDialerEndpointTimeoutReadWrite(timeout time.Duration) SignerServiceEndpointOption {
	return func(ss *SignerDialerEndpoint) { ss.timeoutReadWrite = timeout }
}

// SignerDialerEndpointConnRetries sets the amount of attempted retries to
// acceptNewConnection.
func SignerDialerEndpointConnRetries(retries int) SignerServiceEndpointOption {
	return func(ss *SignerDialerEndpoint) { ss.maxConnRetries = retries }
}

// SignerDialerEndpointRetryWaitInterval sets the retry wait interval to a
// custom value.
func SignerDialerEndpointRetryWaitInterval(interval time.Duration) SignerServiceEndpointOption {
	return func(ss *SignerDialerEndpoint) { ss.retryWait = interval }
}

// SignerDialerEndpoint dials using its dialer and responds to any signature
// requests using its privVal.
type SignerDialerEndpoint struct {
	signerEndpoint

	dialer SocketDialer

	retryWait      time.Duration
	maxConnRetries int
}

// NewSignerDialerEndpoint returns a SignerDialerEndpoint that will dial using the given
// dialer and respond to any signature requests over the connection
// using the given privVal.
func NewSignerDialerEndpoint(
	logger log.Logger,
	dialer SocketDialer,
	options ...SignerServiceEndpointOption,
) *SignerDialerEndpoint {

	sd := &SignerDialerEndpoint{
		dialer:         dialer,
		retryWait:      defaultRetryWaitMilliseconds * time.Millisecond,
		maxConnRetries: defaultMaxDialRetries,
	}
	sd.signerEndpoint.logger = logger

	sd.BaseService = *service.NewBaseService(logger, "SignerDialerEndpoint", sd)
	sd.signerEndpoint.timeoutReadWrite = defaultTimeoutReadWriteSeconds * time.Second

	for _, optionFunc := range options {
		optionFunc(sd)
	}

	return sd
}

func (sd *SignerDialerEndpoint) OnStart(context.Context) error { return nil }
func (sd *SignerDialerEndpoint) OnStop()                       {}

func (sd *SignerDialerEndpoint) ensureConnection(ctx context.Context) error {
	if sd.IsConnected() {
		return nil
	}

	timer := time.NewTimer(0)
	defer timer.Stop()
	retries := 0
	for retries < sd.maxConnRetries {
		if err := ctx.Err(); err != nil {
			return err
		}
		conn, err := sd.dialer()

		if err != nil {
			retries++
			sd.logger.Debug("SignerDialer: Reconnection failed", "retries", retries, "max", sd.maxConnRetries, "err", err)

			// Wait between retries
			timer.Reset(sd.retryWait)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
			}
		} else {
			sd.SetConnection(conn)
			sd.logger.Debug("SignerDialer: Connection Ready")
			return nil
		}
	}

	sd.logger.Debug("SignerDialer: Max retries exceeded", "retries", retries, "max", sd.maxConnRetries)

	return ErrNoConnection
}
