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
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bhojpur/state/internal/inspect"
	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
)

// InspectCmd constructs the command to start an inspect server.
func MakeInspectCommand(conf *config.Config, logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Run an inspect server for investigating Bhojpur State",
		Long: `
	inspect runs a subset of Bhojpur State's RPC endpoints that are useful for debugging
	issues with Bhojpur State.

	When the Bhojpur State consensus engine detects inconsistent state, it will crash the
	statectl process. Bhojpur State will not start up while in this inconsistent state. 
	The inspect command can be used to query the block and state store using Bhojpur State
	RPC calls to debug issues of inconsistent state.
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT)
			defer cancel()

			ins, err := inspect.NewFromConfig(logger, conf)
			if err != nil {
				return err
			}

			logger.Info("starting inspect server")
			if err := ins.Run(ctx); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("rpc.laddr",
		conf.RPC.ListenAddress, "RPC listenener address. Port required")
	cmd.Flags().String("db-backend",
		conf.DBBackend, "database backend: goleveldb | cleveldb | boltdb | rocksdb | badgerdb")
	cmd.Flags().String("db-dir", conf.DBPath, "database directory")

	return cmd
}
