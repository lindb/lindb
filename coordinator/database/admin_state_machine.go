package database

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/coordinator/storage"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

// AdminStateMachine is database config controller,
// creates shard assignment based on config and active nodes related storage cluster.
// runtime watches database change event, maintain shard assignment and create related coordinator task.
type AdminStateMachine interface {
	discovery.Listener

	// Close closes admin state machine, stops watch change event
	Close() error
}

// adminStateMachine implement admin state machine interface.
// all metadata change will store related storage cluster.
type adminStateMachine struct {
	repo           state.Repository
	storageCluster storage.ClusterStateMachine
	discovery      discovery.Discovery

	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewAdminStateMachine creates admin state machine instance
func NewAdminStateMachine(ctx context.Context, repo state.Repository,
	storageCluster storage.ClusterStateMachine) (AdminStateMachine, error) {
	c, cancel := context.WithCancel(ctx)
	// new admin state machine instance
	stateMachine := &adminStateMachine{
		repo:           repo,
		storageCluster: storageCluster,
		ctx:            c,
		cancel:         cancel,
		log:            logger.GetLogger("database/admin/state/machine"),
	}
	// new database config discovery
	stateMachine.discovery = discovery.NewDiscovery(repo, constants.DatabaseConfigPath, stateMachine)
	if err := stateMachine.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery database config error:%s", err)
	}
	return stateMachine, nil
}

// OnCreate creates shard assignment when receive database create event
func (sm *adminStateMachine) OnCreate(key string, resource []byte) {
	cfg := models.Database{}
	if err := json.Unmarshal(resource, &cfg); err != nil {
		sm.log.Error("discovery database create but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}

	if len(cfg.Name) == 0 {
		sm.log.Error("database name cannot be empty", logger.String("data", string(resource)))
		return
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, clusterCfg := range cfg.Clusters {
		cluster := sm.storageCluster.GetCluster(clusterCfg.Name)
		if cluster == nil {
			sm.log.Error("storage cluster not exist",
				logger.String("cluster", clusterCfg.Name))
			continue
		}
		shardAssign, err := cluster.GetShardAssign(cfg.Name)
		if err != nil && err != state.ErrNotExist {
			sm.log.Error("get shard assign error", logger.Error(err), logger.Stack())
			return
		}
		// build shard assignment for creation database, generate related coordinator task
		if shardAssign == nil {
			if err := sm.createShardAssignment(cfg.Name, cluster, clusterCfg); err != nil {
				sm.log.Error("create shard assignment error",
					logger.String("data", string(resource)), logger.Error(err))
			}
		}

	}
	//} else if len(shardAssign.Shards) != cfg.NumOfShard {
	//TODO need implement modify database shard num.
}

func (sm *adminStateMachine) OnDelete(key string) {
	//TODO impl delete database???
	//panic("implement me")
}

// Cleanup does cleanup operation when receive event
func (sm *adminStateMachine) Cleanup() {
	//TODO
	//panic("implement me")
}

// Close closes admin state machine, stops watch change event
func (sm *adminStateMachine) Close() error {
	sm.discovery.Close()
	sm.cancel()
	return nil
}

// createShardAssignment creates shard assignment for spec cluster
// 1) generate shard assignment
// 2) save shard assignment into related storage cluster
// 3) submit create shard coordinator task(storage node will execute it when receive task event)
func (sm *adminStateMachine) createShardAssignment(databaseName string,
	cluster storage.Cluster, clusterCfg models.DatabaseCluster) error {
	activeNodes := cluster.GetActiveNodes()
	if len(activeNodes) == 0 {
		return fmt.Errorf("active node not found")
	}
	//TODO need calc resource and pick related node for store data
	var nodes = make(map[int]models.Node)
	for idx, node := range activeNodes {
		nodes[idx] = node
	}

	var nodeIDs []int
	for idx := range nodes {
		nodeIDs = append(nodeIDs, idx)
	}

	// generate shard assignment based on node ids and config
	shardAssign, err := ShardAssignment(nodeIDs, clusterCfg)
	if err != nil {
		return err
	}
	// set nodes and config, storage node will use it when execute create shard task
	shardAssign.Nodes = nodes
	shardAssign.Config = clusterCfg

	// save shard assignment into related storage cluster
	if err := cluster.SaveShardAssign(databaseName, shardAssign); err != nil {
		return err
	}
	return nil
}

// getNodes returns all active nodes by cluster name
func (sm *adminStateMachine) getNodes(clusterName string) (map[int]models.Node, error) {
	cluster := sm.storageCluster.GetCluster(clusterName)
	if cluster == nil {
		return nil, fmt.Errorf("stroage cluster not exist")
	}
	activeNodes := cluster.GetActiveNodes()
	if len(activeNodes) == 0 {
		return nil, fmt.Errorf("active node not found")
	}
	var nodes = make(map[int]models.Node)
	for idx, node := range activeNodes {
		nodes[idx] = node
	}
	return nodes, nil
}
