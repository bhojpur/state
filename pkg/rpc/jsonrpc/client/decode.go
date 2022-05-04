package client

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
	"fmt"

	rpctypes "github.com/bhojpur/state/pkg/rpc/jsonrpc/types"
)

func unmarshalResponseBytes(responseBytes []byte, expectedID string, result interface{}) error {
	// Read response.  If rpc/core/types is imported, the result will unmarshal
	// into the correct type.
	var response rpctypes.RPCResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return response.Error
	}

	if got := response.ID(); got != expectedID {
		return fmt.Errorf("got response ID %q, wanted %q", got, expectedID)
	}

	// Unmarshal the RawMessage into the result.
	if err := json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("error unmarshaling result: %w", err)
	}
	return nil
}

func unmarshalResponseBytesArray(responseBytes []byte, expectedIDs []string, results []interface{}) error {
	var responses []rpctypes.RPCResponse
	if err := json.Unmarshal(responseBytes, &responses); err != nil {
		return fmt.Errorf("unmarshaling responses: %w", err)
	} else if len(responses) != len(results) {
		return fmt.Errorf("got %d results, wanted %d", len(responses), len(results))
	}

	// Intersect IDs from responses with expectedIDs.
	ids := make([]string, len(responses))
	for i, resp := range responses {
		ids[i] = resp.ID()
	}
	if err := validateResponseIDs(ids, expectedIDs); err != nil {
		return fmt.Errorf("wrong IDs: %w", err)
	}

	for i, resp := range responses {
		if err := json.Unmarshal(resp.Result, results[i]); err != nil {
			return fmt.Errorf("unmarshaling result %d: %w", i, err)
		}
	}
	return nil
}

func validateResponseIDs(ids, expectedIDs []string) error {
	m := make(map[string]struct{}, len(expectedIDs))
	for _, id := range expectedIDs {
		m[id] = struct{}{}
	}

	for i, id := range ids {
		if _, ok := m[id]; !ok {
			return fmt.Errorf("unexpected response ID %d: %q", i, id)
		}
	}
	return nil
}
