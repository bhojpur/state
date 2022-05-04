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

	cfg "github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/scripts/keymigrate"
	"github.com/bhojpur/state/scripts/scmigrate"
)

func MakeKeyMigrateCommand(conf *cfg.Config, logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-migrate",
		Short: "Run Database key migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			contexts := []string{
				// this is ordered to put the
				// (presumably) biggest/most important
				// subsets first.
				"blockstore",
				"state",
				"peerstore",
				"tx_index",
				"evidence",
				"light",
			}

			for idx, dbctx := range contexts {
				logger.Info("beginning a key migration",
					"dbctx", dbctx,
					"num", idx+1,
					"total", len(contexts),
				)

				db, err := cfg.DefaultDBProvider(&cfg.DBContext{
					ID:     dbctx,
					Config: conf,
				})

				if err != nil {
					return fmt.Errorf("constructing database handle: %w", err)
				}

				if err = keymigrate.Migrate(ctx, db); err != nil {
					return fmt.Errorf("running migration for context %q: %w",
						dbctx, err)
				}

				if dbctx == "blockstore" {
					if err := scmigrate.Migrate(ctx, db); err != nil {
						return fmt.Errorf("running seen commit migration: %w", err)

					}
				}
			}

			logger.Info("completed database migration successfully")

			return nil
		},
	}

	// allow database info to be overridden via cli
	addDBFlags(cmd, conf)

	return cmd
}
