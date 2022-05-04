package database

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
)

func mockDBWithStuff(t *testing.T) DB {
	db := NewMemDB()
	// Under "key" prefix
	require.NoError(t, db.Set(bz("key"), bz("value")))
	require.NoError(t, db.Set(bz("key1"), bz("value1")))
	require.NoError(t, db.Set(bz("key2"), bz("value2")))
	require.NoError(t, db.Set(bz("key3"), bz("value3")))
	require.NoError(t, db.Set(bz("something"), bz("else")))
	require.NoError(t, db.Set(bz("k"), bz("val")))
	require.NoError(t, db.Set(bz("ke"), bz("valu")))
	require.NoError(t, db.Set(bz("kee"), bz("valuu")))
	return db
}

func TestPrefixDBSimple(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	checkValue(t, pdb, bz("key"), nil)
	checkValue(t, pdb, bz("key1"), nil)
	checkValue(t, pdb, bz("1"), bz("value1"))
	checkValue(t, pdb, bz("key2"), nil)
	checkValue(t, pdb, bz("2"), bz("value2"))
	checkValue(t, pdb, bz("key3"), nil)
	checkValue(t, pdb, bz("3"), bz("value3"))
	checkValue(t, pdb, bz("something"), nil)
	checkValue(t, pdb, bz("k"), nil)
	checkValue(t, pdb, bz("ke"), nil)
	checkValue(t, pdb, bz("kee"), nil)
}

func TestPrefixDBIterator1(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	itr, err := pdb.Iterator(nil, nil)
	require.NoError(t, err)
	checkDomain(t, itr, nil, nil)
	checkItem(t, itr, bz("1"), bz("value1"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("2"), bz("value2"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("3"), bz("value3"))
	checkNext(t, itr, false)
	checkInvalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator1(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	itr, err := pdb.ReverseIterator(nil, nil)
	require.NoError(t, err)
	checkDomain(t, itr, nil, nil)
	checkItem(t, itr, bz("3"), bz("value3"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("2"), bz("value2"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("1"), bz("value1"))
	checkNext(t, itr, false)
	checkInvalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator5(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	itr, err := pdb.ReverseIterator(bz("1"), nil)
	require.NoError(t, err)
	checkDomain(t, itr, bz("1"), nil)
	checkItem(t, itr, bz("3"), bz("value3"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("2"), bz("value2"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("1"), bz("value1"))
	checkNext(t, itr, false)
	checkInvalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator6(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	itr, err := pdb.ReverseIterator(bz("2"), nil)
	require.NoError(t, err)
	checkDomain(t, itr, bz("2"), nil)
	checkItem(t, itr, bz("3"), bz("value3"))
	checkNext(t, itr, true)
	checkItem(t, itr, bz("2"), bz("value2"))
	checkNext(t, itr, false)
	checkInvalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator7(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := NewPrefixDB(db, bz("key"))

	itr, err := pdb.ReverseIterator(nil, bz("2"))
	require.NoError(t, err)
	checkDomain(t, itr, nil, bz("2"))
	checkItem(t, itr, bz("1"), bz("value1"))
	checkNext(t, itr, false)
	checkInvalid(t, itr)
	itr.Close()
}
