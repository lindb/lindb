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

package master

import (
	"context"
	"encoding/json"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_cluster.go -destination=./storage_cluster_mock.go -package=master

// StorageCluster represents storage storageCluster controller,
// 1) discovery active node list in storageCluster
// 2) save shard assignment
// 3) generate coordinator task
type StorageCluster interface {
	// GetConfig returns storage configuration.
	GetConfig() *config.StorageCluster
	// Start starts the state machine for storage state change.
	Start() error
	// GetState returns the current state of storage cluster.
	GetState() *models.StorageState
	// GetLiveNodes returns the current live nodes of storage cluster.
	GetLiveNodes() ([]models.StatefulNode, error)
	// FlushDatabase submits the coordinator task for flushing memory database by name
	FlushDatabase(databaseName string) error
	// SaveDatabaseAssignment saves database assignment in storage state repo.
	SaveDatabaseAssignment(
		shardAssign *models.ShardAssignment,
		databaseOption *option.DatabaseOption,
	) error
	// DropDatabaseAssignment drops database assignment from storage state repo.
	DropDatabaseAssignment(databaseName string) error
	// GetRepo returns current storage storageCluster's state repo
	GetRepo() state.Repository
	// Close closes storageCluster controller
	Close()
}

// storageCluster implements StorageCluster controller, master will maintain multi storage storageCluster
type storageCluster struct {
	ctx         context.Context
	cfg         *config.StorageCluster
	storageRepo state.Repository
	stateMgr    StateManager

	state *models.StorageState
	sm    discovery.StateMachine

	logger *logger.Logger
}

// newStorageCluster creates storageCluster controller, init active node list if exist node, must return storageCluster
func newStorageCluster(ctx context.Context,
	cfg *config.StorageCluster,
	stateMgr StateManager,
	repoFactory state.RepositoryFactory) (cluster StorageCluster, err error) {
	var storageRepo state.Repository
	storageRepo, err = repoFactory.CreateStorageRepo(cfg.Config)
	if err != nil {
		return nil, err
	}

	log := logger.GetLogger("coordinator", "Storage")
	cluster = &storageCluster{
		ctx:         ctx,
		cfg:         cfg,
		storageRepo: storageRepo,
		stateMgr:    stateMgr,
		state:       models.NewStorageState(cfg.Config.Namespace),
		logger:      log,
	}

	log.Info("init storage cluster success", logger.String("storage", cfg.Config.Namespace))
	return cluster, nil
}

// GetConfig returns storage configuration.
func (c *storageCluster) GetConfig() *config.StorageCluster {
	return c.cfg
}

// Start starts the state machine for storage state change.
func (c *storageCluster) Start() error {
	sm, err := c.stateMgr.GetStateMachineFactory().
		createStorageNodeStateMachine(c.cfg.Config.Namespace, discovery.NewFactory(c.storageRepo))
	if err != nil {
		return err
	}
	c.sm = sm

	c.logger.Info("start storage cluster successfully", logger.String("storage", c.cfg.Config.Namespace))
	return nil
}

// GetState returns the current state of storage cluster.
func (c *storageCluster) GetState() *models.StorageState {
	return c.state
}

// GetLiveNodes returns the current live nodes of storage cluster.
func (c *storageCluster) GetLiveNodes() (rs []models.StatefulNode, err error) {
	// TODO add timeout ctx
	kvs, err := c.storageRepo.List(c.ctx, constants.LiveNodesPath)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		node := models.StatefulNode{}
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			return nil, err
		}
		rs = append(rs, node)
	}
	return rs, nil
}

// GetRepo returns current storage storageCluster's state repo
func (c *storageCluster) GetRepo() state.Repository {
	return c.storageRepo
}

// FlushDatabase submits the coordinator task for flushing memory database by name
func (c *storageCluster) FlushDatabase(_ string) error {
	panic("need impl")
}

// SaveDatabaseAssignment saves database assignment in storage state repo.
func (c *storageCluster) SaveDatabaseAssignment(
	shardAssign *models.ShardAssignment,
	databaseOption *option.DatabaseOption,
) error {
	// TODO timeout ctx
	data := encoding.JSONMarshal(&models.DatabaseAssignment{
		ShardAssignment: shardAssign,
		Option:          databaseOption,
	})
	if err := c.storageRepo.Put(c.ctx, constants.GetDatabaseAssignPath(shardAssign.Name), data); err != nil {
		return err
	}
	c.logger.Info("save database assignment successfully",
		logger.String("storage", c.cfg.Config.Namespace),
		logger.String("database", shardAssign.Name))
	return nil
}

// DropDatabaseAssignment drops database assignment from storage state repo.
func (c *storageCluster) DropDatabaseAssignment(databaseName string) error {
	if err := c.storageRepo.Delete(c.ctx, constants.GetDatabaseAssignPath(databaseName)); err != nil {
		return err
	}
	c.logger.Info("drop database assignment successfully",
		logger.String("storage", c.cfg.Config.Namespace),
		logger.String("database", databaseName))
	return nil
}

// Close stops watch, and cleanups storageCluster's metadata
func (c *storageCluster) Close() {
	c.logger.Info("close storage cluster state machine", logger.String("storage", c.cfg.Config.Namespace))
	if c.sm != nil {
		if err := c.sm.Close(); err != nil {
			c.logger.Error("close storage node state machine of storage cluster",
				logger.String("storage", c.cfg.Config.Namespace), logger.Error(err), logger.Stack())
		}
	}
	if err := c.storageRepo.Close(); err != nil {
		c.logger.Error("close state repo of storage cluster",
			logger.String("storage", c.cfg.Config.Namespace), logger.Error(err), logger.Stack())
	}
}
