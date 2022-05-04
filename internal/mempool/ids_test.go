package mempool

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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/types"
)

func TestMempoolIDsBasic(t *testing.T) {
	ids := NewMempoolIDs()

	peerID, err := types.NewNodeID("0011223344556677889900112233445566778899")
	require.NoError(t, err)
	require.EqualValues(t, 0, ids.GetForPeer(peerID))

	ids.ReserveForPeer(peerID)
	require.EqualValues(t, 1, ids.GetForPeer(peerID))

	ids.Reclaim(peerID)
	require.EqualValues(t, 0, ids.GetForPeer(peerID))

	ids.ReserveForPeer(peerID)
	require.EqualValues(t, 1, ids.GetForPeer(peerID))
}

func TestMempoolIDsPeerDupReserve(t *testing.T) {
	ids := NewMempoolIDs()

	peerID, err := types.NewNodeID("0011223344556677889900112233445566778899")
	require.NoError(t, err)
	require.EqualValues(t, 0, ids.GetForPeer(peerID))

	ids.ReserveForPeer(peerID)
	require.EqualValues(t, 1, ids.GetForPeer(peerID))

	ids.ReserveForPeer(peerID)
	require.EqualValues(t, 1, ids.GetForPeer(peerID))
}

func TestMempoolIDs2Peers(t *testing.T) {
	ids := NewMempoolIDs()

	peer1ID, _ := types.NewNodeID("0011223344556677889900112233445566778899")
	require.EqualValues(t, 0, ids.GetForPeer(peer1ID))

	ids.ReserveForPeer(peer1ID)
	require.EqualValues(t, 1, ids.GetForPeer(peer1ID))

	ids.Reclaim(peer1ID)
	require.EqualValues(t, 0, ids.GetForPeer(peer1ID))

	peer2ID, _ := types.NewNodeID("1011223344556677889900112233445566778899")

	ids.ReserveForPeer(peer2ID)
	require.EqualValues(t, 1, ids.GetForPeer(peer2ID))

	ids.ReserveForPeer(peer1ID)
	require.EqualValues(t, 2, ids.GetForPeer(peer1ID))
}

func TestMempoolIDsNextExistID(t *testing.T) {
	ids := NewMempoolIDs()

	peer1ID, _ := types.NewNodeID("0011223344556677889900112233445566778899")
	ids.ReserveForPeer(peer1ID)
	require.EqualValues(t, 1, ids.GetForPeer(peer1ID))

	peer2ID, _ := types.NewNodeID("1011223344556677889900112233445566778899")
	ids.ReserveForPeer(peer2ID)
	require.EqualValues(t, 2, ids.GetForPeer(peer2ID))

	peer3ID, _ := types.NewNodeID("2011223344556677889900112233445566778899")
	ids.ReserveForPeer(peer3ID)
	require.EqualValues(t, 3, ids.GetForPeer(peer3ID))

	ids.Reclaim(peer1ID)
	require.EqualValues(t, 0, ids.GetForPeer(peer1ID))

	ids.Reclaim(peer3ID)
	require.EqualValues(t, 0, ids.GetForPeer(peer3ID))

	ids.ReserveForPeer(peer1ID)
	require.EqualValues(t, 1, ids.GetForPeer(peer1ID))

	ids.ReserveForPeer(peer3ID)
	require.EqualValues(t, 3, ids.GetForPeer(peer3ID))
}
