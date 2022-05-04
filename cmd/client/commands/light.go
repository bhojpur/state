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
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	dbm "github.com/bhojpur/state/pkg/database"
	"github.com/spf13/cobra"

	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	libmath "github.com/bhojpur/state/pkg/libs/math"
	"github.com/bhojpur/state/pkg/light"
	lproxy "github.com/bhojpur/state/pkg/light/proxy"
	lrpc "github.com/bhojpur/state/pkg/light/rpc"
	dbs "github.com/bhojpur/state/pkg/light/store/db"
	rpcserver "github.com/bhojpur/state/pkg/rpc/jsonrpc/server"
)

// LightCmd constructs the base command called when invoked without any subcommands.
func MakeLightCommand(conf *config.Config, logger log.Logger) *cobra.Command {
	var (
		listenAddr         string
		primaryAddr        string
		witnessAddrsJoined string
		chainID            string
		dir                string
		maxOpenConnections int

		sequential     bool
		trustingPeriod time.Duration
		trustedHeight  int64
		trustedHash    []byte
		trustLevelStr  string

		logLevel  string
		logFormat string

		primaryKey   = []byte("primary")
		witnessesKey = []byte("witnesses")
	)

	checkForExistingProviders := func(db dbm.DB) (string, []string, error) {
		primaryBytes, err := db.Get(primaryKey)
		if err != nil {
			return "", []string{""}, err
		}
		witnessesBytes, err := db.Get(witnessesKey)
		if err != nil {
			return "", []string{""}, err
		}
		witnessesAddrs := strings.Split(string(witnessesBytes), ",")
		return string(primaryBytes), witnessesAddrs, nil
	}

	saveProviders := func(db dbm.DB, primaryAddr, witnessesAddrs string) error {
		err := db.Set(primaryKey, []byte(primaryAddr))
		if err != nil {
			return fmt.Errorf("failed to save primary provider: %w", err)
		}
		err = db.Set(witnessesKey, []byte(witnessesAddrs))
		if err != nil {
			return fmt.Errorf("failed to save witness providers: %w", err)
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:   "light [chainID]",
		Short: "Run a light client proxy server, verifying Bhojpur State rpc",
		Long: `Run a light client proxy server, verifying Bhojpur State rpc.

All calls that can be tracked back to a block header by a proof
will be verified before passing them back to the caller. Other than
that, it will present the same interface as a full Bhojpur State node.

Furthermore to the chainID, a fresh instance of a light client will
need a primary RPC address and a trusted hash and height. It is also highly
recommended to provide additional witness RPC addresses, especially if
not using sequential verification.

To restart the node, thereafter only the chainID is required.

When /abci_query is called, the Merkle key path format is:

	/{store name}/{key}

Please verify with your application that this Merkle key format is used (true
for applications built w/ Bhojpur Chain SDK).
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID = args[0]
			logger.Info("Creating client...", "chainID", chainID)

			var witnessesAddrs []string
			if witnessAddrsJoined != "" {
				witnessesAddrs = strings.Split(witnessAddrsJoined, ",")
			}

			lightDB, err := dbm.NewGoLevelDB("light-client-db", dir)
			if err != nil {
				return fmt.Errorf("can't create a database: %w", err)
			}
			// create a prefixed db on the chainID
			db := dbm.NewPrefixDB(lightDB, []byte(chainID))

			if primaryAddr == "" { // check to see if we can start from an existing state
				var err error
				primaryAddr, witnessesAddrs, err = checkForExistingProviders(db)
				if err != nil {
					return fmt.Errorf("failed to retrieve primary or witness from db: %w", err)
				}
				if primaryAddr == "" {
					return errors.New("no primary address was provided nor found. Please provide a primary (using -p)." +
						" Run the command: statectl light --help for more information")
				}
			} else {
				err := saveProviders(db, primaryAddr, witnessAddrsJoined)
				if err != nil {
					logger.Error("Unable to save primary and or witness addresses", "err", err)
				}
			}

			if len(witnessesAddrs) < 1 && !sequential {
				logger.Info("In skipping verification mode it is highly recommended to provide at least one witness")
			}

			trustLevel, err := libmath.ParseFraction(trustLevelStr)
			if err != nil {
				return fmt.Errorf("can't parse trust level: %w", err)
			}

			options := []light.Option{light.Logger(logger)}

			vo := light.SkippingVerification(trustLevel)
			if sequential {
				vo = light.SequentialVerification()
			}
			options = append(options, vo)

			// Initiate the light client. If the trusted store already has blocks in it, this
			// will be used else we use the trusted options.
			c, err := light.NewHTTPClient(
				cmd.Context(),
				chainID,
				light.TrustOptions{
					Period: trustingPeriod,
					Height: trustedHeight,
					Hash:   trustedHash,
				},
				primaryAddr,
				witnessesAddrs,
				dbs.New(db),
				options...,
			)
			if err != nil {
				return err
			}

			cfg := rpcserver.DefaultConfig()
			cfg.MaxBodyBytes = conf.RPC.MaxBodyBytes
			cfg.MaxHeaderBytes = conf.RPC.MaxHeaderBytes
			cfg.MaxOpenConnections = maxOpenConnections
			// If necessary adjust global WriteTimeout to ensure it's greater than
			// TimeoutBroadcastTxCommit.
			if cfg.WriteTimeout <= conf.RPC.TimeoutBroadcastTxCommit {
				cfg.WriteTimeout = conf.RPC.TimeoutBroadcastTxCommit + 1*time.Second
			}

			p, err := lproxy.NewProxy(c, listenAddr, primaryAddr, cfg, logger, lrpc.KeyPathFn(lrpc.DefaultMerkleKeyPathFn()))
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			go func() {
				<-ctx.Done()
				p.Listener.Close()
			}()

			logger.Info("Starting proxy...", "laddr", listenAddr)
			if err := p.ListenAndServe(ctx); err != http.ErrServerClosed {
				// Error starting or closing listener:
				logger.Error("proxy ListenAndServe", "err", err)
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
		Example: `light hub-3 -p http://25.8.19.73:26657 -w http://public-seed-node.bhojpur.net:26657
	--height 962118 --hash 28B97BE9F6DE51AC69F70E0B7BFD7E5C9CD1A595B7DC31AFF27C50D4948020CD`,
	}

	cmd.Flags().StringVar(&listenAddr, "laddr", "tcp://localhost:8888",
		"serve the proxy on the given address")
	cmd.Flags().StringVarP(&primaryAddr, "primary", "p", "",
		"connect to a Bhojpur State node at this address")
	cmd.Flags().StringVarP(&witnessAddrsJoined, "witnesses", "w", "",
		"Bhojpur State nodes to cross-check the primary node, comma-separated")
	cmd.Flags().StringVarP(&dir, "dir", "d", os.ExpandEnv(filepath.Join("$HOME", ".bhojpur-light")),
		"specify the directory")
	cmd.Flags().IntVar(
		&maxOpenConnections,
		"max-open-connections",
		900,
		"maximum number of simultaneous connections (including WebSocket).")
	cmd.Flags().DurationVar(&trustingPeriod, "trusting-period", 168*time.Hour,
		"trusting period that headers can be verified within. Should be significantly less than the unbonding period")
	cmd.Flags().Int64Var(&trustedHeight, "height", 1, "Trusted header's height")
	cmd.Flags().BytesHexVar(&trustedHash, "hash", []byte{}, "Trusted header's hash")
	cmd.Flags().StringVar(&logLevel, "log-level", log.LogLevelInfo, "The logging level (debug|info|warn|error|fatal)")
	cmd.Flags().StringVar(&logFormat, "log-format", log.LogFormatPlain, "The logging format (text|json)")
	cmd.Flags().StringVar(&trustLevelStr, "trust-level", "1/3",
		"trust level. Must be between 1/3 and 3/3",
	)
	cmd.Flags().BoolVar(&sequential, "sequential", false,
		"sequential verification. Verify all headers sequentially as opposed to using skipping verification",
	)

	return cmd

}
