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

	"github.com/lindb/lindb/app/streaming"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
)

const (
	streamingCfgName        = "streaming.toml"
	streamingLogFileName    = "lind-streaming.log"
	defaultStreamingCfgFile = currentDir + streamingCfgName
)

var runStreamingCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the streaming",
	RunE:  serveStreaming,
}

// newStreamingCmd returns a new streaming-cmd
func newStreamingCmd() *cobra.Command {
	streamingCmd := &cobra.Command{
		Use:   "streaming",
		Short: "Run as a streaming node with cluster mode enabled",
	}
	runStreamingCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("streaming config file path, default is %s", defaultStreamingCfgFile))
	runStreamingCmd.PersistentFlags().BoolVar(&doc, "doc", false,
		"enable swagger api doc")
	runStreamingCmd.PersistentFlags().BoolVar(&pprof, "pprof", false,
		"profiling Go programs with pprof")

	streamingCmd.AddCommand(
		runStreamingCmd,
		initializeStreamingConfigCmd,
	)
	return streamingCmd
}

var initializeStreamingConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default streaming-config",
	RunE: func(_ *cobra.Command, _ []string) error {
		path := cfg
		if path == "" {
			path = defaultStreamingCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		return ltoml.WriteConfig(path, config.NewDefaultStreamingTOML())
	},
}

func serveStreaming(_ *cobra.Command, _ []string) error {
	ctx := newCtxWithSignals()
	streamingCfg := config.Streaming{}

	if err := config.LoadAndSetStreamingConfig(cfg, defaultStreamingCfgFile, &streamingCfg); err != nil {
		return err
	}
	if err := logger.InitLogger(streamingCfg.Logging, streamingLogFileName); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// start streaming server
	streamingRuntime := streaming.NewStreamingRuntime(config.Version, &streamingCfg)
	return run(ctx, streamingRuntime, func() error {
		newStreamingCfg := config.Streaming{}
		return config.LoadAndSetStreamingConfig(cfg, defaultStreamingCfgFile, &newStreamingCfg)
	})
}
