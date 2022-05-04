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
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/cli"
	"github.com/bhojpur/state/pkg/libs/log"
)

const ctxTimeout = 4 * time.Second

// ParseConfig retrieves the default environment configuration,
// sets up the Bhojpur State root and ensures that the root exists
func ParseConfig(conf *config.Config) (*config.Config, error) {
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	conf.SetRoot(conf.RootDir)

	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %w", err)
	}
	return conf, nil
}

// RootCommand constructs the root command-line entry point for Bhojpur State core.
func RootCommand(conf *config.Config, logger log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statectl",
		Short: "Bhojpur State replication tool for applications in any programming languages",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == VersionCmd.Name() {
				return nil
			}

			if err := cli.BindFlagsLoadViper(cmd, args); err != nil {
				return err
			}

			pconf, err := ParseConfig(conf)
			if err != nil {
				return err
			}
			*conf = *pconf
			config.EnsureRoot(conf.RootDir)
			if err := log.OverrideWithNewLogger(logger, conf.LogFormat, conf.LogLevel); err != nil {
				return err
			}
			if warning := pconf.DeprecatedFieldWarning(); warning != nil {
				logger.Info("WARNING", "deprecated field warning", warning)
			}

			return nil
		},
	}
	cmd.PersistentFlags().StringP(cli.HomeFlag, "", os.ExpandEnv(filepath.Join("$HOME", config.DefaultStateDir)), "directory for config and data")
	cmd.PersistentFlags().Bool(cli.TraceFlag, false, "print out full stack trace on errors")
	cmd.PersistentFlags().String("log-level", conf.LogLevel, "log level")
	cobra.OnInitialize(func() { cli.InitEnv("STATE") })
	return cmd
}
