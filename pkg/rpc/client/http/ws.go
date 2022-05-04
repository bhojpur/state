package http

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
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bhojpur/state/internal/pubsub"
	"github.com/bhojpur/state/pkg/libs/log"
	rpcclient "github.com/bhojpur/state/pkg/rpc/client"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
	jsonrpcclient "github.com/bhojpur/state/pkg/rpc/jsonrpc/client"
)

// wsEvents is a wrapper around WSClient, which implements SubscriptionClient.
type wsEvents struct {
	Logger log.Logger
	ws     *jsonrpcclient.WSClient

	mtx           sync.RWMutex
	subscriptions map[string]*wsSubscription
}

type wsSubscription struct {
	res   chan coretypes.ResultEvent
	id    string
	query string
}

var _ rpcclient.SubscriptionClient = (*wsEvents)(nil)

func newWsEvents(remote string) (*wsEvents, error) {
	w := &wsEvents{
		Logger:        log.NewNopLogger(),
		subscriptions: make(map[string]*wsSubscription),
	}

	var err error
	w.ws, err = jsonrpcclient.NewWS(strings.TrimSuffix(remote, "/"), "/websocket")
	if err != nil {
		return nil, fmt.Errorf("can't create WS client: %w", err)
	}
	w.ws.OnReconnect(func() {
		// resubscribe immediately
		w.redoSubscriptionsAfter(0 * time.Second)
	})
	w.ws.Logger = w.Logger

	return w, nil
}

// Start starts the websocket client and the event loop.
func (w *wsEvents) Start(ctx context.Context) error {
	if err := w.ws.Start(ctx); err != nil {
		return err
	}
	go w.eventListener(ctx)
	return nil
}

// Stop shuts down the websocket client.
func (w *wsEvents) Stop() error { return w.ws.Stop() }

// Subscribe implements SubscriptionClient by using WSClient to subscribe given
// subscriber to query. By default, it returns a channel with cap=1. Error is
// returned if it fails to subscribe.
//
// When reading from the channel, keep in mind there's a single events loop, so
// if you don't read events for this subscription fast enough, other
// subscriptions will slow down in effect.
//
// The channel is never closed to prevent clients from seeing an erroneous
// event.
//
// It returns an error if wsEvents is not running.
func (w *wsEvents) Subscribe(ctx context.Context, subscriber, query string,
	outCapacity ...int) (out <-chan coretypes.ResultEvent, err error) {
	if err := w.ws.Subscribe(ctx, query); err != nil {
		return nil, err
	}

	outCap := 1
	if len(outCapacity) > 0 {
		outCap = outCapacity[0]
	}

	outc := make(chan coretypes.ResultEvent, outCap)
	w.mtx.Lock()
	defer w.mtx.Unlock()
	// subscriber param is ignored because Bhojpur State will override it with
	// remote IP anyway.
	w.subscriptions[query] = &wsSubscription{res: outc, query: query}

	return outc, nil
}

// Unsubscribe implements SubscriptionClient by using WSClient to unsubscribe
// given subscriber from query.
//
// It returns an error if wsEvents is not running.
func (w *wsEvents) Unsubscribe(ctx context.Context, subscriber, query string) error {
	if err := w.ws.Unsubscribe(ctx, query); err != nil {
		return err
	}

	w.mtx.Lock()
	defer w.mtx.Unlock()
	info, ok := w.subscriptions[query]
	if ok {
		if info.id != "" {
			delete(w.subscriptions, info.id)
		}
		delete(w.subscriptions, info.query)
	}

	return nil
}

// UnsubscribeAll implements SubscriptionClient by using WSClient to
// unsubscribe given subscriber from all the queries.
//
// It returns an error if wsEvents is not running.
func (w *wsEvents) UnsubscribeAll(ctx context.Context, subscriber string) error {
	if err := w.ws.UnsubscribeAll(ctx); err != nil {
		return err
	}

	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.subscriptions = make(map[string]*wsSubscription)

	return nil
}

// After being reconnected, it is necessary to redo subscription to server
// otherwise no data will be automatically received.
func (w *wsEvents) redoSubscriptionsAfter(d time.Duration) {
	time.Sleep(d)

	ctx := context.TODO()

	w.mtx.Lock()
	defer w.mtx.Unlock()

	for q, info := range w.subscriptions {
		if q != "" && q == info.id {
			continue
		}
		err := w.ws.Subscribe(ctx, q)
		if err != nil {
			w.Logger.Error("failed to resubscribe", "query", q, "err", err)
			delete(w.subscriptions, q)
		}
	}
}

func isErrAlreadySubscribed(err error) bool {
	return strings.Contains(err.Error(), pubsub.ErrAlreadySubscribed.Error())
}

func (w *wsEvents) eventListener(ctx context.Context) {
	for {
		select {
		case resp, ok := <-w.ws.ResponsesCh:
			if !ok {
				return
			}

			if resp.Error != nil {
				w.Logger.Error("WS error", "err", resp.Error.Error())
				// Error can be ErrAlreadySubscribed or max client (subscriptions per
				// client) reached or Bhojpur State exited.
				// We can ignore ErrAlreadySubscribed, but need to retry in other
				// cases.
				if !isErrAlreadySubscribed(resp.Error) {
					// Resubscribe after 1 second to give Bhojpur State time to restart (if
					// crashed).
					w.redoSubscriptionsAfter(1 * time.Second)
				}
				continue
			}

			result := new(coretypes.ResultEvent)
			err := json.Unmarshal(resp.Result, result)
			if err != nil {
				w.Logger.Error("failed to unmarshal response", "err", err)
				continue
			}

			w.mtx.RLock()
			out, ok := w.subscriptions[result.Query]
			if ok {
				if _, idOk := w.subscriptions[result.SubscriptionID]; !idOk {
					out.id = result.SubscriptionID
					w.subscriptions[result.SubscriptionID] = out
				}
			}

			w.mtx.RUnlock()
			if ok {
				select {
				case out.res <- *result:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}