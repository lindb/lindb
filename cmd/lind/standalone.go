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

package lind

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/standalone"
)

const (
	standaloneCfgName     = "standalone.toml"
	standaloneLogFileName = "lind-standalone.log"
	// DefaultStandaloneCfgFile defines default config file path for standalone mode
	defaultStandaloneCfgFile = "./" + standaloneCfgName
)

// newStandaloneCmd returns a new standalone-cmd
func newStandaloneCmd() *cobra.Command {
	standaloneCmd := &cobra.Command{
		Use:   "standalone",
		Short: "Run as a standalone node with embed broker, storage, etcd)",
	}

	standaloneCmd.AddCommand(
		runStandaloneCmd,
		initializeStandaloneConfigCmd,
	)

	runStandaloneCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"profiling Go programs with pprof")
	runStandaloneCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("config file path for standalone mode, default is %s", defaultStandaloneCfgFile))

	return standaloneCmd
}

var runStandaloneCmd = &cobra.Command{
	Use:   "run",
	Short: "run as standalone mode",
	RunE:  serveStandalone,
}

// initializeStandaloneConfigCmd initializes config for standalone mode
var initializeStandaloneConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default standalone-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if len(path) == 0 {
			path = defaultStandaloneCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		return ltoml.WriteConfig(path, config.NewDefaultStandaloneTOML())
	},
}

// serveStandalone runs the cluster as standalone mode
func serveStandalone(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	standaloneCfg := config.Standalone{}
	if err := ltoml.LoadConfig(cfg, defaultStandaloneCfgFile, &standaloneCfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	if err := logger.InitLogger(standaloneCfg.Logging, standaloneLogFileName); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// run cluster as standalone mode
	runtime := standalone.NewStandaloneRuntime(getVersion(), &standaloneCfg)
	if err := run(ctx, runtime); err != nil {
		return err
	}
	return nil
}
