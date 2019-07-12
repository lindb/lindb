package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/coordinator/task"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/service"

	"go.uber.org/zap"
)

// Cluster represents storage cluster controller,
// 1) discovery active node list in cluster
// 2) save shard assignment
// 3) generate coordinator task
type Cluster interface {
	discovery.Listener
	// GetActiveNodes returns all active nodes
	GetActiveNodes() []models.Node
	// GetShardAssign returns shard assignment by database name, return not exist err if it not exist
	GetShardAssign(databaseName string) (*models.ShardAssignment, error)
	// SaveShardAssign saves shard assignment
	SaveShardAssign(databaseName string, shardAssign *models.ShardAssignment) error
	// SubmitTask generates coordinator task
	SubmitTask(kind task.Kind, name string, params []task.ControllerTaskParam) error
	// GetRepo returns current storage cluster's state repo
	GetRepo() state.Repository
	// Close closes cluster controller
	Close()
}

// cluster implements cluster controller, master will maintain multi storage cluster
type cluster struct {
	cfg                models.StorageCluster
	repo               state.Repository
	discovery          discovery.Discovery
	shardAssignService service.ShardAssignService
	controller         *task.Controller
	nodes              map[string]models.Node
	databases          map[string]*models.DatabaseCluster

	mutex sync.RWMutex
	log   *zap.Logger
}

// newCluster creates cluster controller, init active node list if exist node
func newCluster(ctx context.Context, cfg models.StorageCluster) (Cluster, error) {
	repo, err := state.NewRepo(cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("new state repo error when create cluster,error:%s", err)
	}
	cluster := &cluster{
		cfg:                cfg,
		repo:               repo,
		shardAssignService: service.NewShardAssignService(repo),
		controller:         task.NewController(ctx, repo),
		nodes:              make(map[string]models.Node),
		databases:          make(map[string]*models.DatabaseCluster),
		log:                logger.GetLogger(),
	}
	// init active nodes if exist
	nodeList, err := repo.List(ctx, constants.ActiveNodesPath)
	if err != nil {
		return nil, fmt.Errorf("get active nodes error:%s", err)
	}
	for _, node := range nodeList {
		cluster.addNode(node)
	}

	// new storage active node discovery
	cluster.discovery = discovery.NewDiscovery(repo, constants.ActiveNodesPath, cluster)
	if err := cluster.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery active storage nodes error:%s", err)
	}
	return cluster, nil
}

// OnCreate adds node into active node list when node online
func (c *cluster) OnCreate(key string, resource []byte) {
	c.addNode(resource)
}

// OnDelete remove node from active node list when node offline
func (c *cluster) OnDelete(key string) {
	name := pathutil.GetName(key)
	c.mutex.Lock()
	delete(c.nodes, name)
	c.mutex.Unlock()
}

func (c *cluster) Cleanup() {
	// do nothing
}

// GetRepo returns current storage cluster's state repo
func (c *cluster) GetRepo() state.Repository {
	return c.repo
}

// GetActiveNodes returns all active nodes
func (c *cluster) GetActiveNodes() []models.Node {
	var activeNodes []models.Node
	c.mutex.RLock()
	for _, v := range c.nodes {
		activeNodes = append(activeNodes, v)
	}
	c.mutex.RUnlock()
	return activeNodes
}

// GetShardAssign returns shard assignment by database name, return not exist err if it not exist
func (c *cluster) GetShardAssign(databaseName string) (*models.ShardAssignment, error) {
	return c.shardAssignService.Get(databaseName)
}

// SaveShardAssign saves shard assignment, generates create shard task after saving successfully
func (c *cluster) SaveShardAssign(databaseName string, shardAssign *models.ShardAssignment) error {
	if err := c.shardAssignService.Save(databaseName, shardAssign); err != nil {
		return err
	}
	var tasks = make(map[int]*models.CreateShardTask)

	for ID, shard := range shardAssign.Shards {
		for _, replicaID := range shard.Replicas {
			taskParam, ok := tasks[replicaID]
			if !ok {
				taskParam = &models.CreateShardTask{Database: databaseName}
				tasks[replicaID] = taskParam
			}
			taskParam.ShardIDs = append(taskParam.ShardIDs, ID)
			taskParam.ShardOption = shardAssign.Config.ShardOption
		}
	}
	var params []task.ControllerTaskParam
	for nodeID, taskParam := range tasks {
		node := shardAssign.Nodes[nodeID]
		params = append(params, task.ControllerTaskParam{
			NodeID: (&node).String(),
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
	return c.controller.Submit(kind, name, params)
}

// Close stops watch, and cleanups cluster's metadata
func (c *cluster) Close() {
	c.mutex.Lock()
	c.nodes = make(map[string]models.Node)
	c.databases = make(map[string]*models.DatabaseCluster)
	c.mutex.Unlock()

	c.discovery.Close()
	if err := c.repo.Close(); err != nil {
		c.log.Error("close state repo of storage cluster",
			zap.String("cluster", c.cfg.Name), zap.Error(err), zap.Stack("stack"))
	}
}

// addNode adds node into active node list
func (c *cluster) addNode(resource []byte) {
	node := models.Node{}
	if err := json.Unmarshal(resource, &node); err != nil {
		c.log.Error("discovery new storage node but unmarshal error",
			zap.String("data", string(resource)), zap.Error(err))
		return
	}

	c.mutex.Lock()
	c.nodes[node.String()] = node
	c.mutex.Unlock()
}
