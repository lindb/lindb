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
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/validate"
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
	stmtpkg.StorageOpShow:    listStorages,
	stmtpkg.StorageOpCreate:  createStorage,
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

// List lists all storage clusters
func listStorages(ctx context.Context, deps *depspkg.HTTPDeps, _ *stmtpkg.Storage) (interface{}, error) {
	data, err := deps.Repo.List(ctx, constants.StorageConfigPath)
	if err != nil {
		return nil, err
	}
	stateMgr := deps.StateMgr
	var storages models.Storages
	for _, val := range data {
		storage := models.Storage{}
		err = encoding.JSONUnmarshal(val.Value, &storage)
		if err != nil {
			log.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			if _, ok := stateMgr.GetStorage(storage.Config.Namespace); ok {
				storage.Status = models.ClusterStatusReady
			} else {
				storage.Status = models.ClusterStatusInitialize
				// TODO: check storage un-health
			}
			storages = append(storages, storage)
		}
	}

	if err != nil {
		return nil, err
	}
	return storages, nil
}

// createStorage creates config of storage cluster.
func createStorage(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Storage) (interface{}, error) {
	data := []byte(stmt.Value)
	storage := &config.StorageCluster{}
	err := encoding.JSONUnmarshal(data, storage)
	if err != nil {
		return nil, err
	}
	err = validate.Validator.Struct(storage)
	if err != nil {
		return nil, err
	}
	// copy config for testing
	cfg := &config.RepoState{}
	_ = encoding.JSONUnmarshal(encoding.JSONMarshal(storage.Config), cfg)
	cfg.Timeout = ltoml.Duration(time.Second)
	cfg.DialTimeout = ltoml.Duration(time.Second)
	// check storage repo config if valid
	repo, err := deps.RepoFactory.CreateStorageRepo(cfg)
	if err != nil {
		return nil, err
	}
	err = repo.Close()
	if err != nil {
		return nil, err
	}
	// re-marshal storage config, keep same structure with repo.
	data = encoding.JSONMarshal(storage)
	log.Info("Creating storage cluster", logger.String("config", stmt.Value))
	ok, err := deps.Repo.PutWithTX(ctx, constants.GetStorageClusterConfigPath(storage.Config.Namespace), data, func(oldVal []byte) error {
		if bytes.Equal(data, oldVal) {
			log.Info("storage cluster exist", logger.String("config", string(oldVal)))
			return state.ErrNotExist
		}
		return nil
	})
	if errors.Is(state.ErrNotExist, err) {
		rs := "Storage is exist"
		return &rs, nil
	}
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create storage failure")
	}
	rs := "Create storage ok"
	return &rs, nil
}

// recoverStorage recovers all database config/shard assignment by given storage.
func recoverStorage(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Storage) (interface{}, error) {
	storage, ok := deps.StateMgr.GetStorage(stmt.Value)
	if !ok {
		return nil, fmt.Errorf("storage not found")
	}
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

	storageName := storage.Name

	databases := make(map[string]*models.ShardAssignment)
	databaseSchema := make(map[string]*models.Database)
	for _, cfg := range result {
		for databaseName, databaseCfg := range cfg.databases {
			shardAssignment, ok := databases[databaseName]
			if !ok {
				shardAssignment = models.NewShardAssignment(databaseName)
				databases[databaseName] = shardAssignment
				databaseSchema[databaseName] = &models.Database{
					Name:    databaseName,
					Storage: storageName,
					Option:  databaseCfg.Option,
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
		if err := deps.Repo.Put(ctx, constants.GetDatabaseAssignPath(databaseName), encoding.JSONMarshal(shardAssignment)); err != nil {
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
