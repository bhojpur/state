package debug

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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/cli"
	"github.com/bhojpur/state/pkg/libs/log"
	rpchttp "github.com/bhojpur/state/pkg/rpc/client/http"
)

func getDumpCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump [output-directory]",
		Short: "Continuously poll a Bhojpur State process and dump debugging data into a single location",
		Long: `Continuously poll a Bhojpur State process and dump debugging data into a single
location at a specified frequency. At each frequency interval, an archived and compressed
file will contain node debugging information including the goroutine and heap profiles
if enabled.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outDir := args[0]
			if outDir == "" {
				return errors.New("invalid output directory")
			}
			frequency, err := cmd.Flags().GetUint(flagFrequency)
			if err != nil {
				return fmt.Errorf("flag %q not defined: %w", flagFrequency, err)
			}

			if frequency == 0 {
				return errors.New("frequency must be positive")
			}

			nodeRPCAddr, err := cmd.Flags().GetString(flagNodeRPCAddr)
			if err != nil {
				return fmt.Errorf("flag %q not defined: %w", flagNodeRPCAddr, err)
			}

			profAddr, err := cmd.Flags().GetString(flagProfAddr)
			if err != nil {
				return fmt.Errorf("flag %q not defined: %w", flagProfAddr, err)
			}

			if _, err := os.Stat(outDir); os.IsNotExist(err) {
				if err := os.Mkdir(outDir, os.ModePerm); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			rpc, err := rpchttp.New(nodeRPCAddr)
			if err != nil {
				return fmt.Errorf("failed to create new http client: %w", err)
			}

			ctx := cmd.Context()

			home := viper.GetString(cli.HomeFlag)
			conf := config.DefaultConfig()
			conf = conf.SetRoot(home)
			config.EnsureRoot(conf.RootDir)

			dumpArgs := dumpDebugDataArgs{
				conf:     conf,
				outDir:   outDir,
				profAddr: profAddr,
			}
			dumpDebugData(ctx, logger, rpc, dumpArgs)

			ticker := time.NewTicker(time.Duration(frequency) * time.Second)
			for range ticker.C {
				dumpDebugData(ctx, logger, rpc, dumpArgs)
			}

			return nil
		},
	}
	cmd.Flags().Uint(
		flagFrequency,
		30,
		"the frequency (seconds) in which to poll, aggregate and dump Bhojpur State debug data",
	)

	cmd.Flags().String(
		flagProfAddr,
		"",
		"the profiling server address (<host>:<port>)",
	)

	return cmd

}

type dumpDebugDataArgs struct {
	conf     *config.Config
	outDir   string
	profAddr string
}

func dumpDebugData(ctx context.Context, logger log.Logger, rpc *rpchttp.HTTP, args dumpDebugDataArgs) {
	start := time.Now().UTC()

	tmpDir, err := os.MkdirTemp(args.outDir, "bhojpur_debug_tmp")
	if err != nil {
		logger.Error("failed to create temporary directory", "dir", tmpDir, "error", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	logger.Info("getting node status...")
	if err := dumpStatus(ctx, rpc, tmpDir, "status.json"); err != nil {
		logger.Error("failed to dump node status", "error", err)
		return
	}

	logger.Info("getting node network info...")
	if err := dumpNetInfo(ctx, rpc, tmpDir, "net_info.json"); err != nil {
		logger.Error("failed to dump node network info", "error", err)
		return
	}

	logger.Info("getting node consensus state...")
	if err := dumpConsensusState(ctx, rpc, tmpDir, "consensus_state.json"); err != nil {
		logger.Error("failed to dump node consensus state", "error", err)
		return
	}

	logger.Info("copying node WAL...")
	if err := copyWAL(args.conf, tmpDir); err != nil {
		logger.Error("failed to copy node WAL", "error", err)
		return
	}

	if args.profAddr != "" {
		logger.Info("getting node goroutine profile...")
		if err := dumpProfile(tmpDir, args.profAddr, "goroutine", 2); err != nil {
			logger.Error("failed to dump goroutine profile", "error", err)
			return
		}

		logger.Info("getting node heap profile...")
		if err := dumpProfile(tmpDir, args.profAddr, "heap", 2); err != nil {
			logger.Error("failed to dump heap profile", "error", err)
			return
		}
	}

	outFile := filepath.Join(args.outDir, fmt.Sprintf("%s.zip", start.Format(time.RFC3339)))
	if err := zipDir(tmpDir, outFile); err != nil {
		logger.Error("failed to create and compress archive", "file", outFile, "error", err)
	}
}
