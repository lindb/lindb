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

package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./cluster.go -destination=./cluster_mock.go -package=storage

// clusterCfg represents the config which creates cluster instance need
// IMPORTANT: need clean config's resource
type clusterCfg struct {
	ctx                 context.Context
	cfg                 config.StorageCluster
	storageStateService service.StorageStateService
	repo                state.Repository
	controllerFactory   task.ControllerFactory
	factory             discovery.Factory
	shardAssignService  service.ShardAssignService
	logger              *logger.Logger
}

// clean cleans the resource for cfg
func (cfg *clusterCfg) clean() {
	if err := cfg.repo.Close(); err != nil {
		cfg.logger.Error("close state repo of storage cluster",
			logger.String("cluster", cfg.cfg.Name), logger.Error(err), logger.Stack())
	}
}

// ClusterFactory represents a cluster create factory
type ClusterFactory interface {
	// newCluster creates cluster controller
	newCluster(cfg clusterCfg) (Cluster, error)
}

// clusterFactory implements ClusterFactory interface
type clusterFactory struct {
}

// NewClusterFactory creates a cluster factory
func NewClusterFactory() ClusterFactory {
	return &clusterFactory{}
}

// Cluster represents storage cluster controller,
// 1) discovery active node list in cluster
// 2) save shard assignment
// 3) generate coordinator task
type Cluster interface {
	discovery.Listener

	// GetActiveNodes returns all active nodes
	GetActiveNodes() []*models.ActiveNode

	// CollectStat collects storage cluster's stat
	CollectStat() (*models.StorageClusterStat, error)

	// GetShardAssign returns shard assignment by database name, return not exist err if it not exist
	GetShardAssign(databaseName string) (*models.ShardAssignment, error)

	// FlushDatabase submits the coordinator task for flushing memory database by name
	FlushDatabase(databaseName string) error

	// SaveShardAssign saves shard assignment
	SaveShardAssign(
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

	// GetRepo returns current storage cluster's state repo
	GetRepo() state.Repository

	// Close closes cluster controller
	Close()
}

// cluster implements cluster controller, master will maintain multi storage cluster
type cluster struct {
	cfg            clusterCfg
	discovery      discovery.Discovery
	taskController task.Controller

	clusterState *models.StorageState

	mutex sync.RWMutex

	logger *logger.Logger
}

// newCluster creates cluster controller, init active node list if exist node, must return cluster
func (f *clusterFactory) newCluster(cfg clusterCfg) (Cluster, error) {
	cluster := &cluster{
		cfg:          cfg,
		clusterState: models.NewStorageState(),
		logger:       cfg.logger,
	}
	// init active nodes if exist
	nodeList, err := cfg.repo.List(cfg.ctx, constants.ActiveNodesPath+"/data")
	if err != nil {
		return cluster, fmt.Errorf("get active nodes error:%s", err)
	}
	for _, node := range nodeList {
		_ = cluster.addNode(node.Value)
	}
	// set cluster name
	cluster.clusterState.Name = cfg.cfg.Name
	// saving new cluster state
	cluster.saveClusterState()

	// new storage active node discovery
	cluster.discovery = cfg.factory.CreateDiscovery(constants.ActiveNodesPath+"/data", cluster)
	if err := cluster.discovery.Discovery(); err != nil {
		return cluster, fmt.Errorf("discovery active storage nodes error:%s", err)
	}
	cluster.taskController = cfg.controllerFactory.CreateController(cfg.ctx, cfg.repo)

	cluster.logger.Info("init storage cluster success", logger.String("cluster", cluster.clusterState.Name))
	return cluster, nil
}

// OnCreate adds node into active node list when node online
func (c *cluster) OnCreate(_ string, resource []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.addNode(resource) {
		c.saveClusterState()
	}
}

// OnDelete remove node from active node list when node offline
func (c *cluster) OnDelete(key string) {
	_, name := filepath.Split(key)
	c.mutex.Lock()
	c.clusterState.RemoveActiveNode(name)
	c.mutex.Unlock()

	c.saveClusterState()

	c.logger.Info("peer storage is offline",
		logger.String("node", name),
	)
}

// GetRepo returns current storage cluster's state repo
func (c *cluster) GetRepo() state.Repository {
	return c.cfg.repo
}

// GetActiveNodes returns all active nodes
func (c *cluster) GetActiveNodes() []*models.ActiveNode {
	c.mutex.RLock()
	activeNodes := c.clusterState.GetActiveNodes()
	c.mutex.RUnlock()
	return activeNodes
}

// CollectStat collects storage cluster's stat
func (c *cluster) CollectStat() (*models.StorageClusterStat, error) {
	kvs, err := c.GetRepo().List(c.cfg.ctx, constants.StateNodesPath)
	if err != nil {
		return nil, err
	}
	// build result
	var result []*models.NodeStat
	for _, kv := range kvs {
		_, nodeID := filepath.Split(kv.Key)
		stat := &models.NodeStat{}
		if err := encoding.JSONUnmarshal(kv.Value, stat); err != nil {
			return nil, err
		}
		_, ok := c.clusterState.ActiveNodes[nodeID]
		if ok {
			result = append(result, stat)
		}
	}
	return &models.StorageClusterStat{Nodes: result}, err
}

// FlushDatabase submits the coordinator task for flushing memory database by name
func (c *cluster) FlushDatabase(databaseName string) error {
	var params []task.ControllerTaskParam
	taskParam := &models.DatabaseFlushTask{DatabaseName: databaseName}
	for _, node := range c.clusterState.ActiveNodes {
		params = append(params, task.ControllerTaskParam{
			NodeID: node.Node.Indicator(),
			Params: taskParam,
		})
	}
	// create create shard coordinator tasks
	if err := c.SubmitTask(constants.FlushDatabase, databaseName, params); err != nil {
		return err
	}
	return nil
}

// GetShardAssign returns shard assignment by database name, return not exist err if it not exist
func (c *cluster) GetShardAssign(databaseName string) (*models.ShardAssignment, error) {
	return c.cfg.shardAssignService.Get(databaseName)
}

// SaveShardAssign saves shard assignment, generates create shard task after saving successfully
func (c *cluster) SaveShardAssign(
	databaseName string,
	shardAssign *models.ShardAssignment,
	databaseOption option.DatabaseOption,
) error {
	if err := c.cfg.shardAssignService.Save(databaseName, shardAssign); err != nil {
		return err
	}
	var tasks = make(map[int]*models.CreateShardTask)

	for ID, shard := range shardAssign.Shards {
		for _, replicaID := range shard.Replicas {
			taskParam, ok := tasks[replicaID]
			if !ok {
				taskParam = &models.CreateShardTask{DatabaseName: databaseName}
				tasks[replicaID] = taskParam
			}
			taskParam.ShardIDs = append(taskParam.ShardIDs, int32(ID))
			taskParam.DatabaseOption = databaseOption
		}
	}
	var params []task.ControllerTaskParam
	for nodeID, taskParam := range tasks {
		node := shardAssign.Nodes[nodeID]
		params = append(params, task.ControllerTaskParam{
			NodeID: node.Indicator(),
			Params: taskParam,
		})
	}
	// create create shard coordinator tasks
	if err := c.SubmitTask(constants.CreateShard, databaseName, params); err != nil {
		return err
	}
	return nil
}

// SubmitTask submits coordinator task based on kind and params into related storage cluster,
// storage node will execute task if it care this task kind
func (c *cluster) SubmitTask(kind task.Kind, name string, params []task.ControllerTaskParam) error {
	return c.taskController.Submit(kind, name, params)
}

// Close stops watch, and cleanups cluster's metadata
func (c *cluster) Close() {
	c.logger.Info("close storage cluster state machine", logger.String("cluster", c.cfg.cfg.Name))
	if c.taskController != nil {
		// need close task controller of current storage cluster
		if err := c.taskController.Close(); err != nil {
			c.logger.Error("close task controller", logger.String("cluster", c.cfg.cfg.Name), logger.Error(err))
		}
	}
	if c.discovery != nil {
		c.discovery.Close()
	}

	(&c.cfg).clean()
}

// addNode adds node into active node list
func (c *cluster) addNode(resource []byte) bool {
	node := &models.ActiveNode{}
	if err := encoding.JSONUnmarshal(resource, node); err != nil {
		c.logger.Error("discovery new storage node but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return false
	}

	c.clusterState.AddActiveNode(node)
	c.logger.Info("peer storage is online",
		logger.String("node", node.Node.Indicator()),
		logger.Int64("nodeOnlineTime", node.OnlineTime),
	)
	return true
}

// saveClusterState saves a new storage cluster snapshot into state repo.
// master do cluster state control, broker node discovery new state snapshot.
func (c *cluster) saveClusterState() {
	name := c.cfg.cfg.Name
	//TODO need to retry when save state error
	if err := c.cfg.storageStateService.Save(name, c.clusterState); err != nil {
		c.logger.Error("save storage state error", logger.String("cluster", name), logger.Error(err))
	}
}
