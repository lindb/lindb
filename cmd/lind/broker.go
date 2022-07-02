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

	"github.com/spf13/cobra"

	"github.com/lindb/lindb/app/broker"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
)

const (
	brokerCfgName        = "broker.toml"
	brokerLogFileName    = "lind-broker.log"
	defaultBrokerCfgFile = "./" + brokerCfgName
)

// newBrokerCmd returns a new broker-cmd
func newBrokerCmd() *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:   "broker",
		Short: "Run as a compute node with cluster mode enabled",
	}
	runBrokerCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("broker config file path, default is %s", defaultBrokerCfgFile))
	runBrokerCmd.PersistentFlags().BoolVar(&doc, "doc", false,
		"enable swagger api doc")
	runBrokerCmd.PersistentFlags().BoolVar(&pprof, "pprof", false,
		"profiling Go programs with pprof")
	brokerCmd.AddCommand(
		runBrokerCmd,
		initializeBrokerConfigCmd,
	)
	return brokerCmd
}

var runBrokerCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the broker",
	RunE:  serveBroker,
}

// initialize config for broker
var initializeBrokerConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default broker-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if path == "" {
			path = defaultBrokerCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		return ltoml.WriteConfig(path, config.NewDefaultBrokerTOML())
	},
}

// serveBroker runs the broker
func serveBroker(_ *cobra.Command, _ []string) error {
	ctx := newCtxWithSignals()

	brokerCfg := config.Broker{}
	if err := config.LoadAndSetBrokerConfig(cfg, defaultBrokerCfgFile, &brokerCfg); err != nil {
		return err
	}

	if err := logger.InitLogger(brokerCfg.Logging, brokerLogFileName); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// start broker server
	brokerRuntime := broker.NewBrokerRuntime(config.Version, &brokerCfg, true)
	return run(ctx, brokerRuntime, func() error {
		newBrokerCfg := config.Broker{}
		return config.LoadAndSetBrokerConfig(cfg, defaultBrokerCfgFile, &newBrokerCfg)
	})
}
