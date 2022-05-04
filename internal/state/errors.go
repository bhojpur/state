package state

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

import "fmt"

type (
	ErrInvalidBlock error
	ErrProxyAppConn error

	ErrUnknownBlock struct {
		Height int64
	}

	ErrBlockHashMismatch struct {
		CoreHash []byte
		AppHash  []byte
		Height   int64
	}

	ErrAppBlockHeightTooHigh struct {
		CoreHeight int64
		AppHeight  int64
	}

	ErrAppBlockHeightTooLow struct {
		AppHeight int64
		StoreBase int64
	}

	ErrLastStateMismatch struct {
		Height int64
		Core   []byte
		App    []byte
	}

	ErrStateMismatch struct {
		Got      *State
		Expected *State
	}

	ErrNoValSetForHeight struct {
		Height int64
		Err    error
	}

	ErrNoConsensusParamsForHeight struct {
		Height int64
	}

	ErrNoABCIResponsesForHeight struct {
		Height int64
	}
)

func (e ErrUnknownBlock) Error() string {
	return fmt.Sprintf("could not find block #%d", e.Height)
}

func (e ErrBlockHashMismatch) Error() string {
	return fmt.Sprintf(
		"app block hash (%X) does not match core block hash (%X) for height %d",
		e.AppHash,
		e.CoreHash,
		e.Height,
	)
}

func (e ErrAppBlockHeightTooHigh) Error() string {
	return fmt.Sprintf("app block height (%d) is higher than core (%d)", e.AppHeight, e.CoreHeight)
}

func (e ErrAppBlockHeightTooLow) Error() string {
	return fmt.Sprintf("app block height (%d) is too far below block store base (%d)", e.AppHeight, e.StoreBase)
}

func (e ErrLastStateMismatch) Error() string {
	return fmt.Sprintf(
		"latest Bhojpur State block (%d) LastAppHash (%X) does not match app's AppHash (%X)",
		e.Height,
		e.Core,
		e.App,
	)
}

func (e ErrStateMismatch) Error() string {
	return fmt.Sprintf(
		"state after replay does not match saved state. Got ----\n%v\nExpected ----\n%v\n",
		e.Got,
		e.Expected,
	)
}

func (e ErrNoValSetForHeight) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("could not find validator set for height #%d", e.Height)
	}
	return fmt.Sprintf("could not find validator set for height #%d: %s", e.Height, e.Err.Error())
}

func (e ErrNoValSetForHeight) Unwrap() error { return e.Err }

func (e ErrNoConsensusParamsForHeight) Error() string {
	return fmt.Sprintf("could not find consensus params for height #%d", e.Height)
}

func (e ErrNoABCIResponsesForHeight) Error() string {
	return fmt.Sprintf("could not find results for height #%d", e.Height)
}
