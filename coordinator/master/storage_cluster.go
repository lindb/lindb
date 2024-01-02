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

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_cluster.go -destination=./storage_cluster_mock.go -package=master

// StorageCluster represents storage cluster controller,
// 1) discovery active node list in storage cluster
// 2) save shard assignment
// 3) generate coordinator task
type StorageCluster interface {
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
	// SetDatabaseLimits sets the database's limits.
	SetDatabaseLimits(database string, limits []byte) error
	// DropDatabaseAssignment drops database assignment from storage state repo.
	DropDatabaseAssignment(databaseName string) error
}

// storageCluster implements StorageCluster controller, master will maintain multi storage cluster.
type storageCluster struct {
	ctx  context.Context
	repo state.Repository

	state *models.StorageState

	logger logger.Logger
}

// newStorageCluster creates storage cluster controller, init active node list if exist node, must return a storage cluster instance.
func newStorageCluster(ctx context.Context,
	repo state.Repository) StorageCluster {
	log := logger.GetLogger("Master", "Storage")
	cluster := &storageCluster{
		ctx:    ctx,
		repo:   repo,
		state:  models.NewStorageState(),
		logger: log,
	}

	log.Info("init storage cluster success")
	return cluster
}

// GetState returns the current state of storage cluster.
func (c *storageCluster) GetState() *models.StorageState {
	return c.state
}

// GetLiveNodes returns the current live nodes of storage cluster.
func (c *storageCluster) GetLiveNodes() (rs []models.StatefulNode, err error) {
	// TODO: add timeout ctx
	kvs, err := c.repo.List(c.ctx, constants.StorageLiveNodesPath)
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

// FlushDatabase submits the coordinator task for flushing memory database by name
func (c *storageCluster) FlushDatabase(_ string) error {
	// FIXME: need impl
	panic("need impl")
}

// SetDatabaseLimits sets the database's limits.
func (c *storageCluster) SetDatabaseLimits(database string, limits []byte) error {
	if err := c.repo.Put(c.ctx, constants.GetDatabaseLimitPath(database), limits); err != nil {
		return err
	}
	c.logger.Info("set database's limits successfully",
		logger.String("database", database))
	return nil
}

// SaveDatabaseAssignment saves database assignment in storage state repo.
func (c *storageCluster) SaveDatabaseAssignment(
	shardAssign *models.ShardAssignment,
	databaseOption *option.DatabaseOption,
) error {
	// TODO: timeout ctx
	data := encoding.JSONMarshal(shardAssign)
	if err := c.repo.Put(c.ctx, constants.GetShardAssignPath(shardAssign.Name), data); err != nil {
		return err
	}
	c.logger.Info("save database assignment successfully",
		logger.String("database", shardAssign.Name))
	return nil
}

// DropDatabaseAssignment drops database assignment from storage state repo.
func (c *storageCluster) DropDatabaseAssignment(databaseName string) error {
	if err := c.repo.Delete(c.ctx, constants.GetShardAssignPath(databaseName)); err != nil {
		return err
	}
	c.logger.Info("drop database assignment successfully",
		logger.String("database", databaseName))
	return nil
}
