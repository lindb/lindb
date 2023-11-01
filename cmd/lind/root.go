// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"

	"github.com/lindb/common/pkg/ltoml"
	"github.com/spf13/cobra"

	"github.com/lindb/lindb/app/root"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
)

const (
	rootCfgName        = "root.toml"
	rootLogFileName    = "lind-root.log"
	defaultRootCfgFile = "./" + rootCfgName
)

// newRootCmd returns a new root-cmd.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Run as a root compute node with multi idc/regions mode enabled",
	}
	runRootCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("root config file path, default is %s", defaultRootCfgFile))
	runRootCmd.PersistentFlags().BoolVar(&pprof, "pprof", false,
		"profiling Go programs with pprof")
	rootCmd.AddCommand(
		runRootCmd,
		initializeRootConfigCmd,
	)
	return rootCmd
}

var runRootCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the root",
	RunE:  serveRoot,
}

// initialize config for root
var initializeRootConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default root-config",
	RunE: func(_ *cobra.Command, _ []string) error {
		path := cfg
		if path == "" {
			path = defaultRootCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		return ltoml.WriteConfig(path, config.NewDefaultRootTOML())
	},
}

// serveRoot runs the root.
func serveRoot(_ *cobra.Command, _ []string) error {
	ctx := newCtxWithSignals()

	rootCfg := config.Root{}
	if err := config.LoadAndSetRootConfig(cfg, defaultRootCfgFile, &rootCfg); err != nil {
		return err
	}

	if err := logger.InitLogger(rootCfg.Logging, rootLogFileName); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}
	if err := logger.InitAccessLogger(rootCfg.Logging, logger.AccessLogFileName); err != nil {
		return fmt.Errorf("init http access logger error: %s", err)
	}

	// start root server
	rootRuntime := root.NewRootRuntime(config.Version, &rootCfg)
	return run(ctx, rootRuntime, func() error {
		newRootCfg := config.Root{}
		return config.LoadAndSetRootConfig(cfg, defaultRootCfgFile, &newRootCfg)
	})
}
