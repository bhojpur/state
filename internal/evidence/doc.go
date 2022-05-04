package evidence

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
It handles all evidence storage and gossiping from detection to block proposal.
For the different types of evidence refer to the `evidence.go` file in the types package

Gossiping

The core functionality begins with the evidence reactor (see reactor.
go) which operates both the sending and receiving of evidence.

The `Receive` function takes a list of evidence and does the following:

1. Checks that it does not already have the evidence stored

2. Verifies the evidence against the node's state (see state/validation.go#VerifyEvidence)

3. Stores the evidence to a db and a concurrent list

The gossiping of evidence is initiated when a peer is added which starts a go routine to broadcast currently
uncommitted evidence at intervals of 60 seconds (set by the by broadcastEvidenceIntervalS).
It uses a concurrent list to store the evidence and before sending verifies that each evidence is still valid in the
sense that it has not exceeded the max evidence age and height (see types/params.go#EvidenceParams).

There are two buckets that evidence can be stored in: Pending & Committed.

1. Pending is awaiting to be committed (evidence is usually broadcasted then)

2. Committed is for those already on the block and is to ensure that evidence isn't submitted twice

All evidence is proto encoded to disk.

Proposing

When a new block is being proposed (in state/execution.go#CreateProposalBlock),
`PendingEvidence(maxBytes)` is called to send up to the maxBytes of uncommitted evidence, from the evidence store,
prioritized in order of age. All evidence is checked for expiration.

When a node receives evidence in a block it will use the evidence module as a cache first to see if it has
already verified the evidence before trying to verify it again.

Once the proposed evidence is submitted,
the evidence is marked as committed and is moved from the broadcasted set to the committed set.
As a result it is also removed from the concurrent list so that it is no longer gossiped.

Minor Functionality

As all evidence (including POLC's) are bounded by an expiration date, those that exceed this are no longer needed
and hence pruned. Currently, only committed evidence in which a marker to the height that the evidence was committed
and hence very small is saved. All updates are made from the `Update(block, state)` function which should be called
when a new block is committed.

*/
