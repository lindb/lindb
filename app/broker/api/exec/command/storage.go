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
	"time"

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

var log = logger.GetLogger("exec", "Command")

// storageCommandFn represents storage command function define.
type storageCommandFn = func(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Storage) (interface{}, error)

// storageCommands registers all storage related commands.
var storageCommands = map[stmtpkg.StorageOpType]storageCommandFn{
	stmtpkg.StorageOpShow:   listStorages,
	stmtpkg.StorageOpCreate: createStorage,
}

// StorageCommand executes lin query language for storage related.
func StorageCommand(ctx context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	storageStmt := stmt.(*stmtpkg.Storage)
	commandFn, ok := storageCommands[storageStmt.Type]
	if ok {
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
			_, ok := stateMgr.GetStorage(storage.Config.Namespace)
			if ok {
				storage.Status = models.StorageStatusReady
			} else {
				storage.Status = models.StorageStatusInitialize
				// TODO check storage un-health
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
