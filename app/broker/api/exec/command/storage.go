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

package command

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var log = logger.GetLogger("Exec", "Command")

// databaseConfig represents the configuration of all databases by storage node.
type databaseConfig struct {
	nodeID    models.NodeID
	databases map[string]models.DatabaseConfig
}

// storageCommandFn represents storage command function define.
type storageCommandFn = func(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Storage) (interface{}, error)

// storageCommands registers all storage related commands.
var storageCommands = map[stmtpkg.StorageOpType]storageCommandFn{
	stmtpkg.StorageOpRecover: recoverStorage,
}

// StorageCommand executes lin query language for storage related.
func StorageCommand(ctx context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	storageStmt := stmt.(*stmtpkg.Storage)
	if commandFn, ok := storageCommands[storageStmt.Type]; ok {
		return commandFn(ctx, deps, storageStmt)
	}
	return nil, nil
}

// recoverStorage recovers all database config/shard assignment by given storage.
func recoverStorage(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Storage) (interface{}, error) {
	storage := deps.StateMgr.GetStorage()
	nodes := storage.LiveNodes
	size := len(nodes)
	result := make([]databaseConfig, size)
	var wait sync.WaitGroup
	wait.Add(size)
	idx := 0
	for nodeID := range nodes {
		i := idx
		node := nodes[nodeID]
		idx++
		go func() {
			defer wait.Done()

			address := node.HTTPAddress()
			databases := make(map[string]models.DatabaseConfig)
			_, err := resty.New().R().
				SetHeader("Accept", "application/json").
				SetResult(&databases).
				Get(address + constants.APIVersion1CliPath + "/state/metadata/local/database/config")
			if err != nil {
				log.Error("get database config from alive node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = databaseConfig{
				nodeID:    node.ID,
				databases: databases,
			}
		}()
	}
	wait.Wait()

	databases := make(map[string]*models.ShardAssignment)
	databaseSchema := make(map[string]*models.Database)
	for _, cfg := range result {
		for databaseName, databaseCfg := range cfg.databases {
			shardAssignment, ok := databases[databaseName]
			if !ok {
				shardAssignment = models.NewShardAssignment(databaseName)
				databases[databaseName] = shardAssignment
				databaseSchema[databaseName] = &models.Database{
					Name:   databaseName,
					Option: databaseCfg.Option,
				}
			}
			for _, shardID := range databaseCfg.ShardIDs {
				shardAssignment.AddReplica(shardID, cfg.nodeID)
			}
		}
	}

	var databaseNames []string
	for databaseName, shardAssignment := range databases {
		log.Info("recover shard assign",
			logger.String("database", databaseName),
			logger.Any("shardAssign", shardAssignment))
		if err := deps.Repo.Put(ctx, constants.GetShardAssignPath(databaseName), encoding.JSONMarshal(shardAssignment)); err != nil {
			return nil, err
		}
		log.Info("recover database schema", logger.String("config", stmt.Value))
		schema := databaseSchema[databaseName]
		schema.NumOfShard = len(shardAssignment.Shards)
		schema.ReplicaFactor = shardAssignment.GetReplicaFactor()
		if err := deps.Repo.Put(ctx, constants.GetDatabaseConfigPath(databaseName), encoding.JSONMarshal(schema)); err != nil {
			return nil, err
		}
		databaseNames = append(databaseNames, databaseName)
	}

	return &databaseNames, nil
}
