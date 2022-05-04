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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bhojpur/state/pkg/privval"
	"github.com/bhojpur/state/pkg/types"
)

// GenValidatorCmd allows the generation of a keypair for a
// validator.
func MakeGenValidatorCommand() *cobra.Command {
	var keyType string
	cmd := &cobra.Command{
		Use:   "gen-validator",
		Short: "Generate new Bhojpur State validator keypair",
		RunE: func(cmd *cobra.Command, args []string) error {
			pv, err := privval.GenFilePV("", "", keyType)
			if err != nil {
				return err
			}

			jsbz, err := json.Marshal(pv)
			if err != nil {
				return fmt.Errorf("validator -> json: %w", err)
			}

			fmt.Printf("%v\n", string(jsbz))

			return nil
		},
	}

	cmd.Flags().StringVar(&keyType, "key", types.ABCIPubKeyTypeEd25519,
		"Key type to generate privval file with. Options: ed25519, secp256k1")

	return cmd
}
