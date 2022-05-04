package types

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verify that the event data types satisfy their shared interface.
var (
	_ EventData = EventDataBlockSyncStatus{}
	_ EventData = EventDataCompleteProposal{}
	_ EventData = EventDataNewBlock{}
	_ EventData = EventDataNewBlockHeader{}
	_ EventData = EventDataNewEvidence{}
	_ EventData = EventDataNewRound{}
	_ EventData = EventDataRoundState{}
	_ EventData = EventDataStateSyncStatus{}
	_ EventData = EventDataTx{}
	_ EventData = EventDataValidatorSetUpdates{}
	_ EventData = EventDataVote{}
	_ EventData = EventDataString("")
)

func TestQueryTxFor(t *testing.T) {
	tx := Tx("foo")
	assert.Equal(t,
		fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%X'", tx.Hash()),
		EventQueryTxFor(tx).String(),
	)
}

func TestQueryForEvent(t *testing.T) {
	assert.Equal(t,
		"tm.event = 'NewBlock'",
		QueryForEvent(EventNewBlockValue).String(),
	)
	assert.Equal(t,
		"tm.event = 'NewEvidence'",
		QueryForEvent(EventNewEvidenceValue).String(),
	)
}
