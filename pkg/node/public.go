package node

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

// It provides a high level wrapper around Bhojpur State services.

import (
	"context"
	"fmt"

	abciclient "github.com/bhojpur/state/pkg/abci/client"
	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/libs/service"
	"github.com/bhojpur/state/pkg/privval"
	"github.com/bhojpur/state/pkg/types"
)

// NewDefault constructs a Bhojpur State node service for use in go
// process that host their own process-local Bhojpur State node. This is
// equivalent to running Bhojpur State in it's own process communicating
// to an external ABCI application.
func NewDefault(
	ctx context.Context,
	conf *config.Config,
	logger log.Logger,
) (service.Service, error) {
	return newDefaultNode(ctx, conf, logger)
}

// New constructs a Bhojpur State node. The ClientCreator makes it
// possible to construct an ABCI application that runs in the same
// process as the Bhojpur State node.  The final option is a pointer to a
// Genesis document: if the value is nil, the genesis document is read
// from the file specified in the config, and otherwise the node uses
// value of the final argument.
func New(
	ctx context.Context,
	conf *config.Config,
	logger log.Logger,
	cf abciclient.Client,
	gen *types.GenesisDoc,
) (service.Service, error) {
	nodeKey, err := types.LoadOrGenNodeKey(conf.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load or gen node key %s: %w", conf.NodeKeyFile(), err)
	}

	var genProvider genesisDocProvider
	switch gen {
	case nil:
		genProvider = defaultGenesisDocProviderFunc(conf)
	default:
		genProvider = func() (*types.GenesisDoc, error) { return gen, nil }
	}

	switch conf.Mode {
	case config.ModeFull, config.ModeValidator:
		pval, err := privval.LoadOrGenFilePV(conf.PrivValidator.KeyFile(), conf.PrivValidator.StateFile())
		if err != nil {
			return nil, err
		}

		return makeNode(
			ctx,
			conf,
			pval,
			nodeKey,
			cf,
			genProvider,
			config.DefaultDBProvider,
			logger)
	case config.ModeSeed:
		return makeSeedNode(logger, conf, config.DefaultDBProvider, nodeKey, genProvider)
	default:
		return nil, fmt.Errorf("%q is not a valid mode", conf.Mode)
	}
}
