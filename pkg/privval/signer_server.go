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
	"io"
	"sync"

	privvalproto "github.com/bhojpur/state/pkg/api/v1/privval"
	"github.com/bhojpur/state/pkg/libs/service"
	"github.com/bhojpur/state/pkg/types"
)

// ValidationRequestHandlerFunc handles different remoteSigner requests
type ValidationRequestHandlerFunc func(
	ctx context.Context,
	privVal types.PrivValidator,
	requestMessage privvalproto.Message,
	chainID string) (privvalproto.Message, error)

type SignerServer struct {
	service.BaseService

	endpoint *SignerDialerEndpoint
	chainID  string
	privVal  types.PrivValidator

	handlerMtx               sync.Mutex
	validationRequestHandler ValidationRequestHandlerFunc
}

func NewSignerServer(endpoint *SignerDialerEndpoint, chainID string, privVal types.PrivValidator) *SignerServer {
	ss := &SignerServer{
		endpoint:                 endpoint,
		chainID:                  chainID,
		privVal:                  privVal,
		validationRequestHandler: DefaultValidationRequestHandler,
	}

	ss.BaseService = *service.NewBaseService(endpoint.logger, "SignerServer", ss)

	return ss
}

// OnStart implements service.Service.
func (ss *SignerServer) OnStart(ctx context.Context) error {
	go ss.serviceLoop(ctx)
	return nil
}

// OnStop implements service.Service.
func (ss *SignerServer) OnStop() {
	ss.endpoint.logger.Debug("SignerServer: OnStop calling Close")
	_ = ss.endpoint.Close()
}

// SetRequestHandler override the default function that is used to service requests
func (ss *SignerServer) SetRequestHandler(validationRequestHandler ValidationRequestHandlerFunc) {
	ss.handlerMtx.Lock()
	defer ss.handlerMtx.Unlock()
	ss.validationRequestHandler = validationRequestHandler
}

func (ss *SignerServer) servicePendingRequest(ctx context.Context) {
	if !ss.IsRunning() {
		return // Ignore error from closing.
	}

	req, err := ss.endpoint.ReadMessage()
	if err != nil {
		if err != io.EOF {
			ss.endpoint.logger.Error("SignerServer: HandleMessage", "err", err)
		}
		return
	}

	var res privvalproto.Message
	{
		// limit the scope of the lock
		ss.handlerMtx.Lock()
		defer ss.handlerMtx.Unlock()
		res, err = ss.validationRequestHandler(ctx, ss.privVal, req, ss.chainID) // todo
		if err != nil {
			// only log the error; we'll reply with an error in res
			ss.endpoint.logger.Error("SignerServer: handleMessage", "err", err)
		}
	}

	err = ss.endpoint.WriteMessage(res)
	if err != nil {
		ss.endpoint.logger.Error("SignerServer: writeMessage", "err", err)
	}
}

func (ss *SignerServer) serviceLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := ss.endpoint.ensureConnection(ctx); err != nil {
				return
			}
			ss.servicePendingRequest(ctx)
		}
	}
}
