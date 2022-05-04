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
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bhojpur/state/internal/jsontypes"
	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/libs/log"
	libnet "github.com/bhojpur/state/pkg/libs/net"
	bos "github.com/bhojpur/state/pkg/libs/os"
	"github.com/bhojpur/state/pkg/privval"
	privrpc "github.com/bhojpur/state/pkg/privval/grpc"
)

// MakeShowValidatorCommand constructs a command to show the validator info.
func MakeShowValidatorCommand(conf *config.Config, logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				pubKey crypto.PubKey
				err    error
				bctx   = cmd.Context()
			)
			//TODO: remove once gRPC is the only supported protocol
			protocol, _ := libnet.ProtocolAndAddress(conf.PrivValidator.ListenAddr)
			switch protocol {
			case "grpc":
				pvsc, err := privrpc.DialRemoteSigner(
					bctx,
					conf.PrivValidator,
					conf.ChainID(),
					logger,
					conf.Instrumentation.Prometheus,
				)
				if err != nil {
					return fmt.Errorf("can't connect to remote validator %w", err)
				}

				ctx, cancel := context.WithTimeout(bctx, ctxTimeout)
				defer cancel()

				pubKey, err = pvsc.GetPubKey(ctx)
				if err != nil {
					return fmt.Errorf("can't get pubkey: %w", err)
				}
			default:

				keyFilePath := conf.PrivValidator.KeyFile()
				if !bos.FileExists(keyFilePath) {
					return fmt.Errorf("private validator file %s does not exist", keyFilePath)
				}

				pv, err := privval.LoadFilePV(keyFilePath, conf.PrivValidator.StateFile())
				if err != nil {
					return err
				}

				ctx, cancel := context.WithTimeout(bctx, ctxTimeout)
				defer cancel()

				pubKey, err = pv.GetPubKey(ctx)
				if err != nil {
					return fmt.Errorf("can't get pubkey: %w", err)
				}
			}

			bz, err := jsontypes.Marshal(pubKey)
			if err != nil {
				return fmt.Errorf("failed to marshal private validator pubkey: %w", err)
			}

			fmt.Println(string(bz))
			return nil
		},
	}

}
