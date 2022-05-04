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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/cli"
	"github.com/bhojpur/state/pkg/libs/log"
	rpchttp "github.com/bhojpur/state/pkg/rpc/client/http"
)

func getKillCmd(logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill [pid] [compressed-output-file]",
		Short: "Kill a Bhojpur State process while aggregating and packaging debugging data",
		Long: `Kill a Bhojpur State process while also aggregating Bhojpur State process data
such as the latest node state, including consensus and networking state,
go-routine state, and the node's WAL and config information. This aggregated data
is packaged into a compressed archive.

Example:
$ statectl debug kill 34255 /path/to/tm-debug.zip`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			pid, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			outFile := args[1]
			if outFile == "" {
				return errors.New("invalid output file")
			}
			nodeRPCAddr, err := cmd.Flags().GetString(flagNodeRPCAddr)
			if err != nil {
				return fmt.Errorf("flag %q not defined: %w", flagNodeRPCAddr, err)
			}

			rpc, err := rpchttp.New(nodeRPCAddr)
			if err != nil {
				return fmt.Errorf("failed to create new http client: %w", err)
			}

			home := viper.GetString(cli.HomeFlag)
			conf := config.DefaultConfig()
			conf = conf.SetRoot(home)
			config.EnsureRoot(conf.RootDir)

			// Create a temporary directory which will contain all the state dumps and
			// relevant files and directories that will be compressed into a file.
			tmpDir, err := os.MkdirTemp(os.TempDir(), "bhojpur_debug_tmp")
			if err != nil {
				return fmt.Errorf("failed to create temporary directory: %w", err)
			}
			defer os.RemoveAll(tmpDir)

			logger.Info("getting node status...")
			if err := dumpStatus(ctx, rpc, tmpDir, "status.json"); err != nil {
				return err
			}

			logger.Info("getting node network info...")
			if err := dumpNetInfo(ctx, rpc, tmpDir, "net_info.json"); err != nil {
				return err
			}

			logger.Info("getting node consensus state...")
			if err := dumpConsensusState(ctx, rpc, tmpDir, "consensus_state.json"); err != nil {
				return err
			}

			logger.Info("copying node WAL...")
			if err := copyWAL(conf, tmpDir); err != nil {
				if !os.IsNotExist(err) {
					return err
				}

				logger.Info("node WAL does not exist; continuing...")
			}

			logger.Info("copying node configuration...")
			if err := copyConfig(home, tmpDir); err != nil {
				return err
			}

			logger.Info("killing Bhojpur State process")
			if err := killProc(int(pid), tmpDir); err != nil {
				return err
			}

			logger.Info("archiving and compressing debug directory...")
			return zipDir(tmpDir, outFile)
		},
	}

	return cmd
}

// killProc attempts to kill the Bhojpur State process with a given PID with an
// ABORT signal which should result in a goroutine stacktrace. The PID's STDERR
// is tailed and piped to a file under the directory dir. An error is returned
// if the output file cannot be created or the tail command cannot be started.
// An error is not returned if any subsequent syscall fails.
func killProc(pid int, dir string) error {
	// pipe STDERR output from tailing the Bhojpur State process to a file
	//
	// NOTE: This will only work on UNIX systems.
	cmd := exec.Command("tail", "-f", fmt.Sprintf("/proc/%d/fd/2", pid)) // nolint: gosec

	outFile, err := os.Create(filepath.Join(dir, "stacktrace.out"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = outFile

	if err := cmd.Start(); err != nil {
		return err
	}

	// kill the underlying Bhojpur State process and subsequent tailing process
	go func() {
		// Killing the Bhojpur State process with the '-ABRT|-6' signal will result in
		// a goroutine stacktrace.
		p, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to find PID to kill Bhojpur State process: %s", err)
		} else if err = p.Signal(syscall.SIGABRT); err != nil {
			fmt.Fprintf(os.Stderr, "failed to kill Bhojpur State process: %s", err)
		}

		// allow some time to allow the Bhojpur State process to be killed
		//
		// TODO: We should 'wait' for a kill to succeed (e.g. poll for PID until it
		// cannot be found). Regardless, this should be ample time.
		time.Sleep(5 * time.Second)

		if err := cmd.Process.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to kill Bhojpur State process output redirection: %s", err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		// only return an error not invoked by a manual kill
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}

	return nil
}
