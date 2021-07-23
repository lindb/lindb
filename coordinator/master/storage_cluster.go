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
	"fmt"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_cluster.go -destination=./storage_cluster_mock.go -package=master

// StorageCluster represents storage storageCluster controller,
// 1) discovery active node list in storageCluster
// 2) save shard assignment
// 3) generate coordinator task
type StorageCluster interface {
	Start() error
	GetState() *models.StorageState
	// FlushDatabase submits the coordinator task for flushing memory database by name
	FlushDatabase(databaseName string) error
	// CreateShards creates shard by shard assignment.
	CreateShards(
		databaseName string,
		shardAssign *models.ShardAssignment,
		databaseOption option.DatabaseOption,
	) error
	// SubmitTask generates coordinator task
	SubmitTask(
		kind task.Kind,
		name string,
		params []task.ControllerTaskParam,
	) error
	// GetRepo returns current storage storageCluster's state repo
	GetRepo() state.Repository
	// Close closes storageCluster controller
	Close()
}

// storageCluster implements StorageCluster controller, master will maintain multi storage storageCluster
type storageCluster struct {
	ctx            context.Context
	cfg            config.StorageCluster
	taskController task.Controller
	storageRepo    state.Repository
	stateMgr       StateManager

	state *models.StorageState
	sm    discovery.StateMachine

	logger *logger.Logger
}

// newStorageCluster creates storageCluster controller, init active node list if exist node, must return storageCluster
func newStorageCluster(ctx context.Context,
	cfg config.StorageCluster,
	stateMgr StateManager,
	repoFactory state.RepositoryFactory,
	controllerFactory task.ControllerFactory) (cluster StorageCluster, err error) {
	//TODO need add config, and retry???
	cfg.Config.Timeout = ltoml.Duration(10 * time.Second)
	cfg.Config.DialTimeout = ltoml.Duration(5 * time.Second)
	var storageRepo state.Repository
	storageRepo, err = repoFactory.CreateRepo(cfg.Config)
	defer func() {
		if err != nil && storageRepo != nil {
			//TODO add log??
			_ = storageRepo.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	log := logger.GetLogger("coordinator", "storage")

	cluster = &storageCluster{
		ctx:            ctx,
		cfg:            cfg,
		taskController: controllerFactory.CreateController(ctx, storageRepo),
		storageRepo:    storageRepo,
		stateMgr:       stateMgr,
		state:          models.NewStorageState(cfg.Name),
		logger:         log,
	}

	log.Info("init storage cluster success", logger.String("storage", cfg.Name))
	return cluster, nil
}

func (c *storageCluster) Start() error {
	sm, err := c.stateMgr.GetStateMachineFactory().
		createStorageNodeStateMachine(c.cfg.Name, discovery.NewFactory(c.storageRepo))
	if err != nil {
		return err
	}
	c.sm = sm

	c.logger.Info("start storage cluster successfully", logger.String("storage", c.cfg.Name))
	return nil
}

func (c *storageCluster) GetState() *models.StorageState {
	return c.state
}

// GetRepo returns current storage storageCluster's state repo
func (c *storageCluster) GetRepo() state.Repository {
	return c.storageRepo
}

// FlushDatabase submits the coordinator task for flushing memory database by name
func (c *storageCluster) FlushDatabase(databaseName string) error {
	var params []task.ControllerTaskParam
	taskParam := &models.DatabaseFlushTask{DatabaseName: databaseName}
	for _, node := range c.state.LiveNodes {
		params = append(params, task.ControllerTaskParam{
			NodeID: node.Indicator(),
			Params: taskParam,
		})
	}
	// create create shard coordinator tasks
	if err := c.SubmitTask(constants.FlushDatabase, databaseName, params); err != nil {
		return err
	}
	c.logger.Info("submit flush database task", logger.String("storage", c.cfg.Name))
	return nil
}

// SaveShardAssign saves shard assignment, generates create shard task after saving successfully
func (c *storageCluster) CreateShards(
	databaseName string,
	shardAssign *models.ShardAssignment,
	databaseOption option.DatabaseOption,
) error {

	//TODO add retry? maybe active nodes not equals shard assign
	liveNodes := c.state.LiveNodes
	if len(liveNodes) == 0 {
		return fmt.Errorf("active node not found")
	}

	nodes := make(map[models.NodeID]*models.StatefulNode)
	for idx := range liveNodes {
		node := liveNodes[idx]
		nodes[node.ID] = &node
	}
	var tasks = make(map[models.NodeID]*models.CreateShardTask)

	for ID, shard := range shardAssign.Shards {
		for _, replicaID := range shard.Replicas {
			taskParam, ok := tasks[replicaID]
			if !ok {
				taskParam = &models.CreateShardTask{DatabaseName: databaseName}
				tasks[replicaID] = taskParam
			}
			taskParam.ShardIDs = append(taskParam.ShardIDs, ID)
			taskParam.DatabaseOption = databaseOption
		}
	}
	var params []task.ControllerTaskParam
	for nodeID, taskParam := range tasks {
		node := nodes[nodeID]
		params = append(params, task.ControllerTaskParam{
			NodeID: node.Indicator(), //TODO need use node id?
			Params: taskParam,
		})
	}
	// create create shard coordinator tasks
	if err := c.SubmitTask(constants.CreateShard, databaseName, params); err != nil {
		return err
	}
	c.logger.Info("submit create task", logger.String("storage", c.cfg.Name))
	return nil
}

// SubmitTask submits coordinator task based on kind and params into related storage storageCluster,
// storage node will execute task if it care this task kind
func (c *storageCluster) SubmitTask(kind task.Kind, name string, params []task.ControllerTaskParam) error {
	return c.taskController.Submit(kind, name, params)
}

// Close stops watch, and cleanups storageCluster's metadata
func (c *storageCluster) Close() {
	c.logger.Info("close storage storageCluster state machine", logger.String("storageCluster", c.cfg.Name))
	if c.taskController != nil {
		// need close task controller of current storage storageCluster
		if err := c.taskController.Close(); err != nil {
			c.logger.Error("close task controller", logger.String("storageCluster", c.cfg.Name), logger.Error(err))
		}
	}
	if c.sm != nil {
		if err := c.sm.Close(); err != nil {
			c.logger.Error("close storage node state machine of storage cluster",
				logger.String("storage", c.cfg.Name), logger.Error(err), logger.Stack())
		}
	}
	if err := c.storageRepo.Close(); err != nil {
		c.logger.Error("close state repo of storage cluster",
			logger.String("storageCluster", c.cfg.Name), logger.Error(err), logger.Stack())
	}
}
