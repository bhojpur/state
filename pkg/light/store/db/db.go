package db

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
	"encoding/binary"
	"fmt"
	"sync"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/google/orderedcode"

	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/light/store"
	"github.com/bhojpur/state/pkg/types"
)

const (
	prefixLightBlock = int64(11)
	prefixSize       = int64(12)
)

type dbs struct {
	db dbm.DB

	mtx  sync.RWMutex
	size uint16
}

// New returns a Store that wraps any DB
// If you want to share one DB across many light clients consider using PrefixDB
func New(db dbm.DB) store.Store {

	lightStore := &dbs{db: db}

	// retrieve the size of the db
	size := uint16(0)
	bz, err := lightStore.db.Get(lightStore.sizeKey())
	if err == nil && len(bz) > 0 {
		size = unmarshalSize(bz)
	}
	lightStore.size = size

	return lightStore
}

// SaveLightBlock persists LightBlock to the db.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) SaveLightBlock(lb *types.LightBlock) error {
	if lb.Height <= 0 {
		panic("negative or zero height")
	}

	lbpb, err := lb.ToProto()
	if err != nil {
		return fmt.Errorf("unable to convert light block to protobuf: %w", err)
	}

	lbBz, err := lbpb.Marshal()
	if err != nil {
		return fmt.Errorf("marshaling LightBlock: %w", err)
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	b := s.db.NewBatch()
	defer b.Close()
	if err = b.Set(s.lbKey(lb.Height), lbBz); err != nil {
		return err
	}
	if err = b.Set(s.sizeKey(), marshalSize(s.size+1)); err != nil {
		return err
	}
	if err = b.WriteSync(); err != nil {
		return err
	}
	s.size++

	return nil
}

// DeleteLightBlockAndValidatorSet deletes the LightBlock from
// the db.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) DeleteLightBlock(height int64) error {
	if height <= 0 {
		panic("negative or zero height")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	b := s.db.NewBatch()
	defer b.Close()
	if err := b.Delete(s.lbKey(height)); err != nil {
		return err
	}
	if err := b.Set(s.sizeKey(), marshalSize(s.size-1)); err != nil {
		return err
	}
	if err := b.WriteSync(); err != nil {
		return err
	}
	s.size--

	return nil
}

// LightBlock retrieves the LightBlock at the given height.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) LightBlock(height int64) (*types.LightBlock, error) {
	if height <= 0 {
		panic("negative or zero height")
	}

	bz, err := s.db.Get(s.lbKey(height))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil, store.ErrLightBlockNotFound
	}

	var lbpb v1.LightBlock
	err = lbpb.Unmarshal(bz)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	lightBlock, err := types.LightBlockFromProto(&lbpb)
	if err != nil {
		return nil, fmt.Errorf("proto conversion error: %w", err)
	}

	return lightBlock, err
}

// LastLightBlockHeight returns the last LightBlock height stored.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) LastLightBlockHeight() (int64, error) {
	itr, err := s.db.ReverseIterator(
		s.lbKey(1),
		append(s.lbKey(1<<63-1), byte(0x00)),
	)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	if itr.Valid() {
		return s.decodeLbKey(itr.Key())
	}

	return -1, itr.Error()
}

// FirstLightBlockHeight returns the first LightBlock height stored.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) FirstLightBlockHeight() (int64, error) {
	itr, err := s.db.Iterator(
		s.lbKey(1),
		append(s.lbKey(1<<63-1), byte(0x00)),
	)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	if itr.Valid() {
		return s.decodeLbKey(itr.Key())
	}

	return -1, itr.Error()
}

// LightBlockBefore iterates over light blocks until it finds a block before
// the given height. It returns ErrLightBlockNotFound if no such block exists.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) LightBlockBefore(height int64) (*types.LightBlock, error) {
	if height <= 0 {
		panic("negative or zero height")
	}

	itr, err := s.db.ReverseIterator(
		s.lbKey(1),
		s.lbKey(height),
	)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	if itr.Valid() {
		var lbpb v1.LightBlock
		err = lbpb.Unmarshal(itr.Value())
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		lightBlock, err := types.LightBlockFromProto(&lbpb)
		if err != nil {
			return nil, fmt.Errorf("proto conversion error: %w", err)
		}
		return lightBlock, nil
	}
	if err = itr.Error(); err != nil {
		return nil, err
	}

	return nil, store.ErrLightBlockNotFound
}

// Prune prunes header & validator set pairs until there are only size pairs
// left.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) Prune(size uint16) error {
	// 1) Check how many we need to prune.
	s.mtx.Lock()
	defer s.mtx.Unlock()
	sSize := s.size

	if sSize <= size { // nothing to prune
		return nil
	}
	numToPrune := sSize - size

	b := s.db.NewBatch()
	defer b.Close()

	// 2) use an iterator to batch together all the blocks that need to be deleted
	if err := s.batchDelete(b, numToPrune); err != nil {
		return err
	}

	// 3) // update size
	s.size = size
	if err := b.Set(s.sizeKey(), marshalSize(size)); err != nil {
		return fmt.Errorf("failed to persist size: %w", err)
	}

	// 4) write batch deletion to disk
	return b.WriteSync()
}

// Size returns the number of header & validator set pairs.
//
// Safe for concurrent use by multiple goroutines.
func (s *dbs) Size() uint16 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.size
}

func (s *dbs) batchDelete(batch dbm.Batch, numToPrune uint16) error {
	itr, err := s.db.Iterator(
		s.lbKey(1),
		append(s.lbKey(1<<63-1), byte(0x00)),
	)
	if err != nil {
		return err
	}
	defer itr.Close()

	for itr.Valid() && numToPrune > 0 {
		if err = batch.Delete(itr.Key()); err != nil {
			return err
		}
		itr.Next()
		numToPrune--
	}

	return itr.Error()
}

func (s *dbs) sizeKey() []byte {
	key, err := orderedcode.Append(nil, prefixSize)
	if err != nil {
		panic(err)
	}
	return key
}

func (s *dbs) lbKey(height int64) []byte {
	key, err := orderedcode.Append(nil, prefixLightBlock, height)
	if err != nil {
		panic(err)
	}
	return key
}

func (s *dbs) decodeLbKey(key []byte) (height int64, err error) {
	var lightBlockPrefix int64
	remaining, err := orderedcode.Parse(string(key), &lightBlockPrefix, &height)
	if err != nil {
		err = fmt.Errorf("failed to parse light block key: %w", err)
	}
	if len(remaining) != 0 {
		err = fmt.Errorf("expected no remainder when parsing light block key but got: %s", remaining)
	}
	if lightBlockPrefix != prefixLightBlock {
		err = fmt.Errorf("expected light block prefix but got: %d", lightBlockPrefix)
	}
	return
}

func marshalSize(size uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, size)
	return bs
}

func unmarshalSize(bz []byte) uint16 {
	return binary.LittleEndian.Uint16(bz)
}
