package p2p

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
	"strings"
	"testing"
	"time"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/crypto/ed25519"
	"github.com/bhojpur/state/pkg/types"
)

func TestPeerScoring(t *testing.T) {
	// coppied from p2p_test shared variables
	selfKey := ed25519.GenPrivKeyFromSecret([]byte{0xf9, 0x1b, 0x08, 0xaa, 0x38, 0xee, 0x34, 0xdd})
	selfID := types.NodeIDFromPubKey(selfKey.PubKey())

	// create a mock peer manager
	db := dbm.NewMemDB()
	peerManager, err := NewPeerManager(selfID, db, PeerManagerOptions{})
	require.NoError(t, err)

	// create a fake node
	id := types.NodeID(strings.Repeat("a1", 20))
	added, err := peerManager.Add(NodeAddress{NodeID: id, Protocol: "memory"})
	require.NoError(t, err)
	require.True(t, added)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Synchronous", func(t *testing.T) {
		// update the manager and make sure it's correct
		require.EqualValues(t, 0, peerManager.Scores()[id])

		// add a bunch of good status updates and watch things increase.
		for i := 1; i < 10; i++ {
			peerManager.processPeerEvent(ctx, PeerUpdate{
				NodeID: id,
				Status: PeerStatusGood,
			})
			require.EqualValues(t, i, peerManager.Scores()[id])
		}

		// watch the corresponding decreases respond to update
		for i := 10; i == 0; i-- {
			peerManager.processPeerEvent(ctx, PeerUpdate{
				NodeID: id,
				Status: PeerStatusBad,
			})
			require.EqualValues(t, i, peerManager.Scores()[id])
		}
	})
	t.Run("AsynchronousIncrement", func(t *testing.T) {
		start := peerManager.Scores()[id]
		pu := peerManager.Subscribe(ctx)
		pu.SendUpdate(ctx, PeerUpdate{
			NodeID: id,
			Status: PeerStatusGood,
		})
		require.Eventually(t,
			func() bool { return start+1 == peerManager.Scores()[id] },
			time.Second,
			time.Millisecond,
			"startAt=%d score=%d", start, peerManager.Scores()[id])
	})
	t.Run("AsynchronousDecrement", func(t *testing.T) {
		start := peerManager.Scores()[id]
		pu := peerManager.Subscribe(ctx)
		pu.SendUpdate(ctx, PeerUpdate{
			NodeID: id,
			Status: PeerStatusBad,
		})
		require.Eventually(t,
			func() bool { return start-1 == peerManager.Scores()[id] },
			time.Second,
			time.Millisecond,
			"startAt=%d score=%d", start, peerManager.Scores()[id])
	})
	t.Run("TestNonPersistantPeerUpperBound", func(t *testing.T) {
		start := int64(peerManager.Scores()[id] + 1)

		for i := start; i <= int64(PeerScorePersistent); i++ {
			peerManager.processPeerEvent(ctx, PeerUpdate{
				NodeID: id,
				Status: PeerStatusGood,
			})

			if i == int64(PeerScorePersistent) {
				require.EqualValues(t, MaxPeerScoreNotPersistent, peerManager.Scores()[id])
			} else {
				require.EqualValues(t, i, peerManager.Scores()[id])
			}
		}
	})
}
