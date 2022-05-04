package clist_test

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
	"pgregory.net/rapid"

	"github.com/bhojpur/state/internal/libs/clist"
)

func TestCListProperties(t *testing.T) {
	rapid.Check(t, rapid.Run(&clistModel{}))
}

// clistModel is used by the rapid state machine testing framework.
// clistModel contains both the clist that is being tested and a slice of *clist.CElements
// that will be used to model the expected clist behavior.
type clistModel struct {
	clist *clist.CList

	model []*clist.CElement
}

// Init is a method used by the rapid state machine testing library.
// Init is called when the test starts to initialize the data that will be used
// in the state machine test.
func (m *clistModel) Init(t *rapid.T) {
	m.clist = clist.New()
	m.model = []*clist.CElement{}
}

// PushBack defines an action that will be randomly selected across by the rapid state
// machines testing library. Every call to PushBack calls PushBack on the clist and
// performs a similar action on the model data.
func (m *clistModel) PushBack(t *rapid.T) {
	value := rapid.String().Draw(t, "value").(string)
	el := m.clist.PushBack(value)
	m.model = append(m.model, el)
}

// Remove defines an action that will be randomly selected across by the rapid state
// machine testing library. Every call to Remove selects an element from the model
// and calls Remove on the CList with that element. The same element is removed from
// the model to keep the objects in sync.
func (m *clistModel) Remove(t *rapid.T) {
	if len(m.model) == 0 {
		return
	}
	ix := rapid.IntRange(0, len(m.model)-1).Draw(t, "index").(int)
	value := m.model[ix]
	m.model = append(m.model[:ix], m.model[ix+1:]...)
	m.clist.Remove(value)
}

// Check is a method required by the rapid state machine testing library.
// Check is run after each action and is used to verify that the state of the object,
// in this case a clist.CList matches the state of the objec.
func (m *clistModel) Check(t *rapid.T) {
	require.Equal(t, len(m.model), m.clist.Len())
	if len(m.model) == 0 {
		return
	}
	require.Equal(t, m.model[0], m.clist.Front())
	require.Equal(t, m.model[len(m.model)-1], m.clist.Back())

	iter := m.clist.Front()
	for _, val := range m.model {
		require.Equal(t, val, iter)
		iter = iter.Next()
	}
}
