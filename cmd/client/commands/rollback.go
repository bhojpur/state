package commands

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

	"github.com/spf13/cobra"

	"github.com/bhojpur/state/internal/state"
	"github.com/bhojpur/state/pkg/config"
)

func MakeRollbackStateCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "rollback Bhojpur State machine state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when Bhojpur State has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application should also roll back to height n - 1. No blocks are removed, so upon
restarting Bhojpur State the transactions in block n will be re-executed against the
application.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			height, hash, err := RollbackState(conf)
			if err != nil {
				return fmt.Errorf("failed to rollback state: %w", err)
			}

			fmt.Printf("Rolled back state to height %d and hash %X", height, hash)
			return nil
		},
	}

}

// RollbackState takes the state at the current height n and overwrites it with the state
// at height n - 1. Note state here refers to Bhojpur State state not application state.
// Returns the latest state height and app hash alongside an error if there was one.
func RollbackState(config *config.Config) (int64, []byte, error) {
	// use the parsed config to load the block and state store
	blockStore, stateStore, err := loadStateAndBlockStore(config)
	if err != nil {
		return -1, nil, err
	}
	defer func() {
		_ = blockStore.Close()
		_ = stateStore.Close()
	}()

	// rollback the last state
	return state.Rollback(blockStore, stateStore)
}
