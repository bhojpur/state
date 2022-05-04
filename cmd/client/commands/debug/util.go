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
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/state/pkg/config"
	rpchttp "github.com/bhojpur/state/pkg/rpc/client/http"
)

// dumpStatus gets node status state dump from the Bhojpur State RPC and writes it
// to file. It returns an error upon failure.
func dumpStatus(ctx context.Context, rpc *rpchttp.HTTP, dir, filename string) error {
	status, err := rpc.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get node status: %w", err)
	}

	return writeStateJSONToFile(status, dir, filename)
}

// dumpNetInfo gets network information state dump from the Bhojpur State RPC and
// writes it to file. It returns an error upon failure.
func dumpNetInfo(ctx context.Context, rpc *rpchttp.HTTP, dir, filename string) error {
	netInfo, err := rpc.NetInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get node network information: %w", err)
	}

	return writeStateJSONToFile(netInfo, dir, filename)
}

// dumpConsensusState gets consensus state dump from the Bhojpur State RPC and
// writes it to file. It returns an error upon failure.
func dumpConsensusState(ctx context.Context, rpc *rpchttp.HTTP, dir, filename string) error {
	consDump, err := rpc.DumpConsensusState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get node consensus dump: %w", err)
	}

	return writeStateJSONToFile(consDump, dir, filename)
}

// copyWAL copies the Bhojpur State node's WAL file. It returns an error if the
// WAL file cannot be read or copied.
func copyWAL(conf *config.Config, dir string) error {
	walPath := conf.Consensus.WalFile()
	walFile := filepath.Base(walPath)

	return copyFile(walPath, filepath.Join(dir, walFile))
}

// copyConfig copies the Bhojpur State node's config file. It returns an error if
// the config file cannot be read or copied.
func copyConfig(home, dir string) error {
	configFile := "config.toml"
	configPath := filepath.Join(home, "config", configFile)

	return copyFile(configPath, filepath.Join(dir, configFile))
}

func dumpProfile(dir, addr, profile string, debug int) error {
	endpoint := fmt.Sprintf("%s/debug/pprof/%s?debug=%d", addr, profile, debug)

	resp, err := http.Get(endpoint) // nolint: gosec
	if err != nil {
		return fmt.Errorf("failed to query for %s profile: %w", profile, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read %s profile response body: %w", profile, err)
	}

	return os.WriteFile(path.Join(dir, fmt.Sprintf("%s.out", profile)), body, os.ModePerm)
}
