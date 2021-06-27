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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/storage"

	"github.com/spf13/cobra"
)

const (
	storageCfgName        = "storage.toml"
	storageLogFileName    = "lind-storage.log"
	defaultStorageCfgFile = "./" + storageCfgName
)

var runStorageCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the storage",
	RunE:  serveStorage,
}

// newStorageCmd returns a new storage-cmd
func newStorageCmd() *cobra.Command {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "Run as a storage node with cluster mode enabled",
	}
	runStorageCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"profiling Go programs with pprof")
	runStorageCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("storage config file path, default is %s", defaultStorageCfgFile))

	storageCmd.AddCommand(
		runStorageCmd,
		initializeStorageConfigCmd,
	)
	return storageCmd
}

var initializeStorageConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default storage-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if len(path) == 0 {
			path = defaultStorageCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		return ltoml.WriteConfig(path, config.NewDefaultStorageTOML())
	},
}

func serveStorage(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	storageCfg := config.Storage{}
	if err := ltoml.LoadConfig(cfg, defaultStorageCfgFile, &storageCfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	if err := logger.InitLogger(storageCfg.Logging, storageLogFileName); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// start storage server
	storageRuntime := storage.NewStorageRuntime(getVersion(), &storageCfg)
	if err := run(ctx, storageRuntime); err != nil {
		return err
	}
	return nil
}
