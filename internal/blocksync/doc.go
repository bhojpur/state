package blocksync

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

/*
It implements two versions of a reactor Service that are responsible for block
propagation and gossip between peers.

In order for a full node to successfully participate in consensus, it must have
the latest view of state. The blocksync protocol is a mechanism in which peers
may exchange and gossip entire blocks with one another, in a request/response
type model, until they've successfully synced to the latest head block. Once
succussfully synced, the full node can switch to an active role in consensus and
will no longer blocksync and thus no longer run the blocksync process.

Note, the blocksync reactor Service gossips entire block and relevant data such
that each receiving peer may construct the entire view of the blocksync state.

There is currently only one version of the blocksync reactor Service
that is battle-tested, but whose test coverage is lacking and is not
formally verified.

The v0 blocksync reactor Service has one P2P channel, BlockchainChannel. This
channel is responsible for handling messages that both request blocks and respond
to block requests from peers. For every block request from a peer, the reactor
will execute respondToPeer which will fetch the block from the node's state store
and respond to the peer. For every block response, the node will add the block
to its pool via AddBlock.

Internally, v0 runs a poolRoutine that constantly checks for what blocks it needs
and requests them. The poolRoutine is also responsible for taking blocks from the
pool, saving and executing each block.
*/
