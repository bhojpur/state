package eventbus

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
	"fmt"
	"strings"

	pbsb "github.com/bhojpur/state/internal/pubsub"
	tmquery "github.com/bhojpur/state/internal/pubsub/query"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/service"
	"github.com/bhojpur/state/pkg/types"
)

// Subscription is a proxy interface for a pubsub Subscription.
type Subscription interface {
	ID() string
	Next(context.Context) (pbsb.Message, error)
}

// EventBus is a common bus for all events going through the system.
// It is a type-aware wrapper around an underlying pubsub server.
// All events should be published via the bus.
type EventBus struct {
	service.BaseService
	pubsub *pbsb.Server
}

// NewDefault returns a new event bus with default options.
func NewDefault(l log.Logger) *EventBus {
	logger := l.With("module", "eventbus")
	pubsub := pbsb.NewServer(l, pbsb.BufferCapacity(0))
	b := &EventBus{pubsub: pubsub}
	b.BaseService = *service.NewBaseService(logger, "EventBus", b)
	return b
}

func (b *EventBus) OnStart(ctx context.Context) error {
	return b.pubsub.Start(ctx)
}

func (b *EventBus) OnStop() {}

func (b *EventBus) NumClients() int {
	return b.pubsub.NumClients()
}

func (b *EventBus) NumClientSubscriptions(clientID string) int {
	return b.pubsub.NumClientSubscriptions(clientID)
}

func (b *EventBus) SubscribeWithArgs(ctx context.Context, args pbsb.SubscribeArgs) (Subscription, error) {
	return b.pubsub.SubscribeWithArgs(ctx, args)
}

func (b *EventBus) Unsubscribe(ctx context.Context, args pbsb.UnsubscribeArgs) error {
	return b.pubsub.Unsubscribe(ctx, args)
}

func (b *EventBus) UnsubscribeAll(ctx context.Context, subscriber string) error {
	return b.pubsub.UnsubscribeAll(ctx, subscriber)
}

func (b *EventBus) Observe(ctx context.Context, observe func(pbsb.Message) error, queries ...*tmquery.Query) error {
	return b.pubsub.Observe(ctx, observe, queries...)
}

func (b *EventBus) Publish(eventValue string, eventData types.EventData) error {
	tokens := strings.Split(types.EventTypeKey, ".")
	event := abci.Event{
		Type: tokens[0],
		Attributes: []abci.EventAttribute{
			{
				Key:   tokens[1],
				Value: eventValue,
			},
		},
	}

	return b.pubsub.PublishWithEvents(eventData, []abci.Event{event})
}

func (b *EventBus) PublishEventNewBlock(data types.EventDataNewBlock) error {
	events := data.ResultFinalizeBlock.Events

	// add Bhojpur State-reserved new block event
	events = append(events, types.EventNewBlock)

	return b.pubsub.PublishWithEvents(data, events)
}

func (b *EventBus) PublishEventNewBlockHeader(data types.EventDataNewBlockHeader) error {
	// no explicit deadline for publishing events

	events := data.ResultFinalizeBlock.Events

	// add Bhojpur State-reserved new block header event
	events = append(events, types.EventNewBlockHeader)

	return b.pubsub.PublishWithEvents(data, events)
}

func (b *EventBus) PublishEventNewEvidence(evidence types.EventDataNewEvidence) error {
	return b.Publish(types.EventNewEvidenceValue, evidence)
}

func (b *EventBus) PublishEventVote(data types.EventDataVote) error {
	return b.Publish(types.EventVoteValue, data)
}

func (b *EventBus) PublishEventValidBlock(data types.EventDataRoundState) error {
	return b.Publish(types.EventValidBlockValue, data)
}

func (b *EventBus) PublishEventBlockSyncStatus(data types.EventDataBlockSyncStatus) error {
	return b.Publish(types.EventBlockSyncStatusValue, data)
}

func (b *EventBus) PublishEventStateSyncStatus(data types.EventDataStateSyncStatus) error {
	return b.Publish(types.EventStateSyncStatusValue, data)
}

// PublishEventTx publishes tx event with events from Result. Note it will add
// predefined keys (EventTypeKey, TxHashKey). Existing events with the same keys
// will be overwritten.
func (b *EventBus) PublishEventTx(data types.EventDataTx) error {
	events := data.Result.Events

	// add Bhojpur State-reserved events
	events = append(events, types.EventTx)

	tokens := strings.Split(types.TxHashKey, ".")
	events = append(events, abci.Event{
		Type: tokens[0],
		Attributes: []abci.EventAttribute{
			{
				Key:   tokens[1],
				Value: fmt.Sprintf("%X", types.Tx(data.Tx).Hash()),
			},
		},
	})

	tokens = strings.Split(types.TxHeightKey, ".")
	events = append(events, abci.Event{
		Type: tokens[0],
		Attributes: []abci.EventAttribute{
			{
				Key:   tokens[1],
				Value: fmt.Sprintf("%d", data.Height),
			},
		},
	})

	return b.pubsub.PublishWithEvents(data, events)
}

func (b *EventBus) PublishEventNewRoundStep(data types.EventDataRoundState) error {
	return b.Publish(types.EventNewRoundStepValue, data)
}

func (b *EventBus) PublishEventTimeoutPropose(data types.EventDataRoundState) error {
	return b.Publish(types.EventTimeoutProposeValue, data)
}

func (b *EventBus) PublishEventTimeoutWait(data types.EventDataRoundState) error {
	return b.Publish(types.EventTimeoutWaitValue, data)
}

func (b *EventBus) PublishEventNewRound(data types.EventDataNewRound) error {
	return b.Publish(types.EventNewRoundValue, data)
}

func (b *EventBus) PublishEventCompleteProposal(data types.EventDataCompleteProposal) error {
	return b.Publish(types.EventCompleteProposalValue, data)
}

func (b *EventBus) PublishEventPolka(data types.EventDataRoundState) error {
	return b.Publish(types.EventPolkaValue, data)
}

func (b *EventBus) PublishEventRelock(data types.EventDataRoundState) error {
	return b.Publish(types.EventRelockValue, data)
}

func (b *EventBus) PublishEventLock(data types.EventDataRoundState) error {
	return b.Publish(types.EventLockValue, data)
}

func (b *EventBus) PublishEventValidatorSetUpdates(data types.EventDataValidatorSetUpdates) error {
	return b.Publish(types.EventValidatorSetUpdatesValue, data)
}

func (b *EventBus) PublishEventEvidenceValidated(evidence types.EventDataEvidenceValidated) error {
	return b.Publish(types.EventEvidenceValidatedValue, evidence)
}
