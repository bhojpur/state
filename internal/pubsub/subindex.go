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
	"github.com/bhojpur/state/internal/pubsub/query"
	abci "github.com/bhojpur/state/pkg/abci/types"
	"github.com/bhojpur/state/pkg/types"
)

// An item to be published to subscribers.
type item struct {
	Data   types.EventData
	Events []abci.Event
}

// A subInfo value records a single subscription.
type subInfo struct {
	clientID string        // chosen by the client
	query    *query.Query  // chosen by the client
	subID    string        // assigned at registration
	sub      *Subscription // receives published events
}

// A subInfoSet is an unordered set of subscription info records.
type subInfoSet map[*subInfo]struct{}

func (s subInfoSet) contains(si *subInfo) bool { _, ok := s[si]; return ok }
func (s subInfoSet) add(si *subInfo)           { s[si] = struct{}{} }
func (s subInfoSet) remove(si *subInfo)        { delete(s, si) }

// withQuery returns the subset of s whose query string matches qs.
func (s subInfoSet) withQuery(qs string) subInfoSet {
	out := make(subInfoSet)
	for si := range s {
		if si.query.String() == qs {
			out.add(si)
		}
	}
	return out
}

// A subIndex is an indexed collection of subscription info records.
// The index is not safe for concurrent use without external synchronization.
type subIndex struct {
	all      subInfoSet            // all subscriptions
	byClient map[string]subInfoSet // per-client subscriptions
	byQuery  map[string]subInfoSet // per-query subscriptions

	// TODO: We allow indexing by query to support existing use by
	// the RPC service methods for event streaming. Fix up those methods not to
	// require this, and then remove indexing by query.
}

// newSubIndex constructs a new, empty subscription index.
func newSubIndex() *subIndex {
	return &subIndex{
		all:      make(subInfoSet),
		byClient: make(map[string]subInfoSet),
		byQuery:  make(map[string]subInfoSet),
	}
}

// findClients returns the set of subscriptions for the given client ID, or nil.
func (idx *subIndex) findClientID(id string) subInfoSet { return idx.byClient[id] }

// findQuery returns the set of subscriptions on the given query string, or nil.
func (idx *subIndex) findQuery(qs string) subInfoSet { return idx.byQuery[qs] }

// contains reports whether idx contains any subscription matching the given
// client ID and query pair.
func (idx *subIndex) contains(clientID, query string) bool {
	csubs, qsubs := idx.byClient[clientID], idx.byQuery[query]
	if len(csubs) == 0 || len(qsubs) == 0 {
		return false
	}
	for si := range csubs {
		if qsubs.contains(si) {
			return true
		}
	}
	return false
}

// add adds si to the index, replacing any previous entry with the same terms.
// It is the caller's responsibility to check for duplicates before adding.
// See also the contains method.
func (idx *subIndex) add(si *subInfo) {
	idx.all.add(si)
	if m := idx.byClient[si.clientID]; m == nil {
		idx.byClient[si.clientID] = subInfoSet{si: struct{}{}}
	} else {
		m.add(si)
	}
	qs := si.query.String()
	if m := idx.byQuery[qs]; m == nil {
		idx.byQuery[qs] = subInfoSet{si: struct{}{}}
	} else {
		m.add(si)
	}
}

// removeAll removes all the elements of s from the index.
func (idx *subIndex) removeAll(s subInfoSet) {
	for si := range s {
		idx.all.remove(si)
		idx.byClient[si.clientID].remove(si)
		if len(idx.byClient[si.clientID]) == 0 {
			delete(idx.byClient, si.clientID)
		}
		if si.query != nil {
			qs := si.query.String()
			idx.byQuery[qs].remove(si)
			if len(idx.byQuery[qs]) == 0 {
				delete(idx.byQuery, qs)
			}
		}
	}
}
