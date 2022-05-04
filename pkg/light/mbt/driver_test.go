package mbt

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
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/pkg/light"
	"github.com/bhojpur/state/pkg/types"
)

const jsonDir = "./json"

func TestVerify(t *testing.T) {
	filenames := jsonFilenames(t)

	for _, filename := range filenames {
		filename := filename
		t.Run(filename, func(t *testing.T) {

			jsonBlob, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}

			var tc testCase
			err = json.Unmarshal(jsonBlob, &tc)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(tc.Description)

			var (
				trustedSignedHeader = tc.Initial.SignedHeader
				trustedNextVals     = tc.Initial.NextValidatorSet
				trustingPeriod      = time.Duration(tc.Initial.TrustingPeriod) * time.Nanosecond
			)

			for _, input := range tc.Input {
				var (
					newSignedHeader = input.LightBlock.SignedHeader
					newVals         = input.LightBlock.ValidatorSet
				)

				err = light.Verify(
					&trustedSignedHeader,
					&trustedNextVals,
					newSignedHeader,
					newVals,
					trustingPeriod,
					input.Now,
					1*time.Second,
					light.DefaultTrustLevel,
				)

				t.Logf("%d -> %d", trustedSignedHeader.Height, newSignedHeader.Height)

				switch input.Verdict {
				case "SUCCESS":
					require.NoError(t, err)
				case "NOT_ENOUGH_TRUST":
					require.IsType(t, light.ErrNewValSetCantBeTrusted{}, err)
				case "INVALID":
					switch err.(type) {
					case light.ErrOldHeaderExpired:
					case light.ErrInvalidHeader:
					default:
						t.Fatalf("expected either ErrInvalidHeader or ErrOldHeaderExpired, but got %v", err)
					}
				default:
					t.Fatalf("unexpected verdict: %q", input.Verdict)
				}

				if err == nil { // advance
					trustedSignedHeader = *newSignedHeader
					trustedNextVals = *input.LightBlock.NextValidatorSet
				}
			}
		})
	}
}

// jsonFilenames returns a list of files in jsonDir directory
func jsonFilenames(t *testing.T) []string {
	matches, err := filepath.Glob(filepath.Join(jsonDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	return matches
}

type testCase struct {
	Description string      `json:"description"`
	Initial     initialData `json:"initial"`
	Input       []inputData `json:"input"`
}

type initialData struct {
	SignedHeader     types.SignedHeader `json:"signed_header"`
	NextValidatorSet types.ValidatorSet `json:"next_validator_set"`
	TrustingPeriod   uint64             `json:"trusting_period,string"`
	Now              time.Time          `json:"now"`
}

type inputData struct {
	LightBlock lightBlockWithNextValidatorSet `json:"block"`
	Now        time.Time                      `json:"now"`
	Verdict    string                         `json:"verdict"`
}

// In state-rs, NextValidatorSet is used to verify new blocks (opposite to
// Go Bhojpur State).
type lightBlockWithNextValidatorSet struct {
	*types.SignedHeader `json:"signed_header"`
	ValidatorSet        *types.ValidatorSet `json:"validator_set"`
	NextValidatorSet    *types.ValidatorSet `json:"next_validator_set"`
}
