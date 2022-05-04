package light_test

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
	"time"

	"github.com/stretchr/testify/assert"

	libmath "github.com/bhojpur/state/pkg/libs/math"
	"github.com/bhojpur/state/pkg/light"
	"github.com/bhojpur/state/pkg/types"
)

const (
	maxClockDrift = 10 * time.Second
)

func TestVerifyAdjacentHeaders(t *testing.T) {
	const (
		chainID    = "TestVerifyAdjacentHeaders"
		lastHeight = 1
		nextHeight = 2
	)

	var (
		keys = genPrivKeys(4)
		// 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
		vals     = keys.ToValidators(20, 10)
		bTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		header   = keys.GenSignedHeader(t, chainID, lastHeight, bTime, nil, vals, vals,
			hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys))
	)

	testCases := []struct {
		newHeader      *types.SignedHeader
		newVals        *types.ValidatorSet
		trustingPeriod time.Duration
		now            time.Time
		expErr         error
		expErrText     string
	}{
		// same header -> no error
		0: {
			header,
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"headers must be adjacent in height",
		},
		// different chainID -> error
		1: {
			keys.GenSignedHeader(t, "different-chainID", nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"header belongs to another chain",
		},
		// new header's time is before old header's time -> error
		2: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(-1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			4 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"to be after old header time",
		},
		// new header's time is from the future -> error
		3: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(3*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"new header has a time from the future",
		},
		// new header's time is from the future, but it's acceptable (< maxClockDrift) -> no error
		4: {
			keys.GenSignedHeader(t, chainID, nextHeight,
				bTime.Add(2*time.Hour).Add(maxClockDrift).Add(-1*time.Millisecond), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 3/3 signed -> no error
		5: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 2/3 signed -> no error
		6: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 1, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 1/3 signed -> error
		7: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), len(keys)-1, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			light.ErrInvalidHeader{Reason: types.ErrNotEnoughVotingPowerSigned{Got: 50, Needed: 93}},
			"",
		},
		// vals does not match with what we have -> error
		8: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, keys.ToValidators(10, 1), vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			keys.ToValidators(10, 1),
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"to match those from new header",
		},
		// vals are inconsistent with newHeader -> error
		9: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			keys.ToValidators(10, 1),
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"to match those that were supplied",
		},
		// old header has expired -> error
		10: {
			keys.GenSignedHeader(t, chainID, nextHeight, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			keys.ToValidators(10, 1),
			1 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"old header has expired",
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			err := light.VerifyAdjacent(header, tc.newHeader, tc.newVals, tc.trustingPeriod, tc.now, maxClockDrift)
			switch {
			case tc.expErr != nil && assert.Error(t, err):
				assert.Equal(t, tc.expErr, err)
			case tc.expErrText != "":
				assert.Contains(t, err.Error(), tc.expErrText)
			default:
				assert.NoError(t, err)
			}
		})
	}

}

func TestVerifyNonAdjacentHeaders(t *testing.T) {
	const (
		chainID    = "TestVerifyNonAdjacentHeaders"
		lastHeight = 1
	)

	var (
		keys = genPrivKeys(4)
		// 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
		vals     = keys.ToValidators(20, 10)
		bTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		header   = keys.GenSignedHeader(t, chainID, lastHeight, bTime, nil, vals, vals,
			hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys))

		// 30, 40, 50
		twoThirds     = keys[1:]
		twoThirdsVals = twoThirds.ToValidators(30, 10)

		// 50
		oneThird     = keys[len(keys)-1:]
		oneThirdVals = oneThird.ToValidators(50, 10)

		// 20
		lessThanOneThird     = keys[0:1]
		lessThanOneThirdVals = lessThanOneThird.ToValidators(20, 10)
	)

	testCases := []struct {
		newHeader      *types.SignedHeader
		newVals        *types.ValidatorSet
		trustingPeriod time.Duration
		now            time.Time
		expErr         error
		expErrText     string
	}{
		// 3/3 new vals signed, 3/3 old vals present -> no error
		0: {
			keys.GenSignedHeader(t, chainID, 3, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 2/3 new vals signed, 3/3 old vals present -> no error
		1: {
			keys.GenSignedHeader(t, chainID, 4, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 1, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 1/3 new vals signed, 3/3 old vals present -> error
		2: {
			keys.GenSignedHeader(t, chainID, 5, bTime.Add(1*time.Hour), nil, vals, vals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), len(keys)-1, len(keys)),
			vals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			light.ErrInvalidHeader{types.ErrNotEnoughVotingPowerSigned{Got: 50, Needed: 93}},
			"",
		},
		// 3/3 new vals signed, 2/3 old vals present -> no error
		3: {
			twoThirds.GenSignedHeader(t, chainID, 5, bTime.Add(1*time.Hour), nil, twoThirdsVals, twoThirdsVals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(twoThirds)),
			twoThirdsVals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 3/3 new vals signed, 1/3 old vals present -> no error
		4: {
			oneThird.GenSignedHeader(t, chainID, 5, bTime.Add(1*time.Hour), nil, oneThirdVals, oneThirdVals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(oneThird)),
			oneThirdVals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			nil,
			"",
		},
		// 3/3 new vals signed, less than 1/3 old vals present -> error
		5: {
			lessThanOneThird.GenSignedHeader(t, chainID, 5, bTime.Add(1*time.Hour), nil, lessThanOneThirdVals, lessThanOneThirdVals,
				hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(lessThanOneThird)),
			lessThanOneThirdVals,
			3 * time.Hour,
			bTime.Add(2 * time.Hour),
			light.ErrNewValSetCantBeTrusted{types.ErrNotEnoughVotingPowerSigned{Got: 20, Needed: 46}},
			"",
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			err := light.VerifyNonAdjacent(header, vals, tc.newHeader, tc.newVals, tc.trustingPeriod,
				tc.now, maxClockDrift,
				light.DefaultTrustLevel)

			switch {
			case tc.expErr != nil && assert.Error(t, err):
				assert.Equal(t, tc.expErr, err)
			case tc.expErrText != "":
				assert.Contains(t, err.Error(), tc.expErrText)
			default:
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyReturnsErrorIfTrustLevelIsInvalid(t *testing.T) {
	const (
		chainID    = "TestVerifyReturnsErrorIfTrustLevelIsInvalid"
		lastHeight = 1
	)

	var (
		keys = genPrivKeys(4)
		// 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
		vals     = keys.ToValidators(20, 10)
		bTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		header   = keys.GenSignedHeader(t, chainID, lastHeight, bTime, nil, vals, vals,
			hash("app_hash"), hash("cons_hash"), hash("results_hash"), 0, len(keys))
	)

	err := light.Verify(header, vals, header, vals, 2*time.Hour, time.Now(), maxClockDrift,
		libmath.Fraction{Numerator: 2, Denominator: 1})
	assert.Error(t, err)
}

func TestValidateTrustLevel(t *testing.T) {
	testCases := []struct {
		lvl   libmath.Fraction
		valid bool
	}{
		// valid
		0: {libmath.Fraction{Numerator: 1, Denominator: 3}, true},
		1: {libmath.Fraction{Numerator: 2, Denominator: 3}, true},
		2: {libmath.Fraction{Numerator: 4, Denominator: 5}, true},
		3: {libmath.Fraction{Numerator: 99, Denominator: 100}, true},

		// invalid
		4: {libmath.Fraction{Numerator: 3, Denominator: 3}, false},
		5: {libmath.Fraction{Numerator: 6, Denominator: 5}, false},
		6: {libmath.Fraction{Numerator: 3, Denominator: 10}, false},
		7: {libmath.Fraction{Numerator: 0, Denominator: 1}, false},
		8: {libmath.Fraction{Numerator: 0, Denominator: 0}, false},
		9: {libmath.Fraction{Numerator: 1, Denominator: 0}, false},
	}

	for idx, tc := range testCases {
		err := light.ValidateTrustLevel(tc.lvl)
		if !tc.valid {
			assert.Error(t, err, idx)
		} else {
			assert.NoError(t, err, idx)
		}
	}
}
