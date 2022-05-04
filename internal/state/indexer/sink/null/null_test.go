package null

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/state/internal/state/indexer"
	"github.com/bhojpur/state/pkg/types"
)

func TestNullEventSink(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nullIndexer := NewEventSink()

	assert.Nil(t, nullIndexer.IndexTxEvents(nil))
	assert.Nil(t, nullIndexer.IndexBlockEvents(types.EventDataNewBlockHeader{}))
	val1, err1 := nullIndexer.SearchBlockEvents(ctx, nil)
	assert.Nil(t, val1)
	assert.NoError(t, err1)
	val2, err2 := nullIndexer.SearchTxEvents(ctx, nil)
	assert.Nil(t, val2)
	assert.NoError(t, err2)
	val3, err3 := nullIndexer.GetTxByHash(nil)
	assert.Nil(t, val3)
	assert.NoError(t, err3)
	val4, err4 := nullIndexer.HasBlock(0)
	assert.False(t, val4)
	assert.NoError(t, err4)
}

func TestType(t *testing.T) {
	nullIndexer := NewEventSink()
	assert.Equal(t, indexer.NULL, nullIndexer.Type())
}

func TestStop(t *testing.T) {
	nullIndexer := NewEventSink()
	assert.Nil(t, nullIndexer.Stop())
}
