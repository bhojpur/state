package pubsub

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

	"github.com/google/uuid"

	"github.com/bhojpur/state/internal/libs/queue"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

var (
	// ErrUnsubscribed is returned by Next when the client has unsubscribed.
	ErrUnsubscribed = errors.New("subscription removed by client")

	// ErrTerminated is returned by Next when the subscription was terminated by
	// the publisher.
	ErrTerminated = errors.New("subscription terminated by publisher")
)

// A Subscription represents a client subscription for a particular query.
type Subscription struct {
	id      string
	queue   *queue.Queue // open until the subscription ends
	stopErr error        // after queue is closed, the reason why
}

// newSubscription returns a new subscription with the given queue capacity.
func newSubscription(quota, limit int) (*Subscription, error) {
	queue, err := queue.New(queue.Options{
		SoftQuota: quota,
		HardLimit: limit,
	})
	if err != nil {
		return nil, err
	}
	return &Subscription{
		id:    uuid.NewString(),
		queue: queue,
	}, nil
}

// Next blocks until a message is available, ctx ends, or the subscription
// ends.  Next returns ErrUnsubscribed if s was unsubscribed, ErrTerminated if
// s was terminated by the publisher, or a context error if ctx ended without a
// message being available.
func (s *Subscription) Next(ctx context.Context) (Message, error) {
	next, err := s.queue.Wait(ctx)
	if errors.Is(err, queue.ErrQueueClosed) {
		return Message{}, s.stopErr
	} else if err != nil {
		return Message{}, err
	}
	return next.(Message), nil
}

// ID returns the unique subscription identifier for s.
func (s *Subscription) ID() string { return s.id }

// publish transmits msg to the subscriber. It reports a queue error if the
// queue cannot accept any further messages.
func (s *Subscription) publish(msg Message) error { return s.queue.Add(msg) }

// stop terminates the subscription with the given error reason.
func (s *Subscription) stop(err error) {
	if err == nil {
		panic("nil stop error")
	}
	s.stopErr = err
	s.queue.Close()
}

// Message glues data and events together.
type Message struct {
	subID  string
	data   types.EventData
	events []abci.Event
}

// SubscriptionID returns the unique identifier for the subscription
// that produced this message.
func (msg Message) SubscriptionID() string { return msg.subID }

// Data returns an original data published.
func (msg Message) Data() types.EventData { return msg.data }

// Events returns events, which matched the client's query.
func (msg Message) Events() []abci.Event { return msg.events }
