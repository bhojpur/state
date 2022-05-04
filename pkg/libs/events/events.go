package events

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

// It provides Pub-Sub with event caching

import (
	"sync"
)

// EventData is a generic event data can be typed and registered with
// bhojpur/ems via concrete implementation of this interface.
type EventData interface{}

// Eventable is the interface reactors and other modules must export to become
// eventable.
type Eventable interface {
	SetEventSwitch(evsw EventSwitch)
}

// Fireable is the interface that wraps the FireEvent method.
//
// FireEvent fires an event with the given name and data.
type Fireable interface {
	FireEvent(eventValue string, data EventData)
}

// EventSwitch is the interface for synchronous pubsub, where listeners
// subscribe to certain events and, when an event is fired (see Fireable),
// notified via a callback function.
//
// Listeners are added by calling AddListenerForEvent function.
// They can be removed by calling either RemoveListenerForEvent or
// RemoveListener (for all events).
type EventSwitch interface {
	Fireable
	AddListenerForEvent(listenerID, eventValue string, cb EventCallback) error
}

type eventSwitch struct {
	mtx        sync.RWMutex
	eventCells map[string]*eventCell
}

func NewEventSwitch() EventSwitch {
	evsw := &eventSwitch{
		eventCells: make(map[string]*eventCell),
	}
	return evsw
}

func (evsw *eventSwitch) AddListenerForEvent(listenerID, eventValue string, cb EventCallback) error {
	// Get/Create eventCell and listener.
	evsw.mtx.Lock()

	eventCell := evsw.eventCells[eventValue]
	if eventCell == nil {
		eventCell = newEventCell()
		evsw.eventCells[eventValue] = eventCell
	}
	evsw.mtx.Unlock()

	eventCell.addListener(listenerID, cb)
	return nil
}

func (evsw *eventSwitch) FireEvent(event string, data EventData) {
	// Get the eventCell
	evsw.mtx.RLock()
	eventCell := evsw.eventCells[event]
	evsw.mtx.RUnlock()

	if eventCell == nil {
		return
	}

	// Fire event for all listeners in eventCell
	eventCell.fireEvent(data)
}

type EventCallback func(data EventData) error

// eventCell handles keeping track of listener callbacks for a given event.
type eventCell struct {
	mtx       sync.RWMutex
	listeners map[string]EventCallback
}

func newEventCell() *eventCell {
	return &eventCell{
		listeners: make(map[string]EventCallback),
	}
}

func (cell *eventCell) addListener(listenerID string, cb EventCallback) {
	cell.mtx.Lock()
	defer cell.mtx.Unlock()
	cell.listeners[listenerID] = cb
}

func (cell *eventCell) fireEvent(data EventData) {
	cell.mtx.RLock()
	eventCallbacks := make([]EventCallback, 0, len(cell.listeners))
	for _, cb := range cell.listeners {
		eventCallbacks = append(eventCallbacks, cb)
	}
	cell.mtx.RUnlock()

	for _, cb := range eventCallbacks {
		if err := cb(data); err != nil {
			// should we log or abort here?
			continue
		}
	}
}
