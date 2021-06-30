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

package cli

import (
	"fmt"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"

	"github.com/spf13/cobra"
)

var createDatabaseCmd = &cobra.Command{
	Use:   "database-create",
	Short: "Creates the database for LinDB",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storageCfg := config.Storage{}
		if err := ltoml.LoadConfig(cliConfigFile, "", &storageCfg); err != nil {
			return fmt.Errorf("decode config file error: %s", err)
		}
		numOfShards, err := readInt("number of shards? :")
		if err != nil {
			return err
		}
		replicaFactor, err := readInt("replica factor? :")
		if err != nil {
			return err
		}
		interval, err := readString("write interval? :")
		if err != nil {
			return err
		}
		initializer := bootstrap.NewClusterInitializer(cliBrokerEndpoint)
		if err := initializer.InitStorageCluster(config.StorageCluster{
			Name:   args[0],
			Config: storageCfg.StorageBase.Coordinator},
		); err != nil {
			return err
		}
		if err := initializer.InitInternalDatabase(models.Database{
			Name:          args[0],
			Cluster:       args[0],
			NumOfShard:    numOfShards,
			ReplicaFactor: replicaFactor,
			Option: option.DatabaseOption{
				Interval: interval,
			},
		}); err != nil {
			return err
		}
		return nil
	},
}

var listDatabaseCmd = &cobra.Command{
	Use:   "database-list",
	Short: "Lists all databases of LinDB",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

var getDatabaseCmd = &cobra.Command{
	Use:   "database-get",
	Short: "Gets Detailed information of a database",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

var deleteDatabaseCmd = &cobra.Command{
	Use:   "database-delete",
	Short: "Deletes a database",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}
