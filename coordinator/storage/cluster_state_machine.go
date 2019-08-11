package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./cluster_state_machine.go -destination=./cluster_state_machine_mock.go -package=storage

// ClusterStateMachine represents storage cluster control when node is master,
// watches cluster config change event, then create/delete related storage cluster controller.
type ClusterStateMachine interface {
	discovery.Listener
	// GetCluster returns cluster controller for maintain the metadata of storage cluster
	GetCluster(name string) Cluster
	// GetAllCluster returns all cluster controller
	GetAllCluster() []Cluster
	// Close closes state machine, cleanup and close all cluster controller
	Close() error
}

// clusterStateMachine implements storage cluster state machine,
// maintain cluster controller for controlling cluster's metadata
type clusterStateMachine struct {
	repo      state.Repository
	discovery discovery.Discovery
	ctx       context.Context
	cancel    context.CancelFunc

	storageStateService service.StorageStateService
	shardAssignService  service.ShardAssignService
	clusterFactory      ClusterFactory
	discoveryFactory    discovery.Factory
	repoFactory         state.RepositoryFactory
	controllerFactory   task.ControllerFactory

	clusters map[string]Cluster

	mutex sync.RWMutex
	log   *logger.Logger
}

// NewClusterStateMachine create state machine, init cluster controller if exist, watch change event
func NewClusterStateMachine(
	ctx context.Context,
	repo state.Repository,
	controllerFactory task.ControllerFactory,
	discoveryFactory discovery.Factory,
	clusterFactory ClusterFactory,
	repoFactory state.RepositoryFactory,
	storageStateService service.StorageStateService,
	shardAssignService service.ShardAssignService) (ClusterStateMachine, error) {
	log := logger.GetLogger("cluster/state/machine")
	c, cancel := context.WithCancel(ctx)
	stateMachine := &clusterStateMachine{
		repo:                repo,
		ctx:                 c,
		cancel:              cancel,
		clusterFactory:      clusterFactory,
		discoveryFactory:    discoveryFactory,
		repoFactory:         repoFactory,
		storageStateService: storageStateService,
		controllerFactory:   controllerFactory,
		shardAssignService:  shardAssignService,
		clusters:            make(map[string]Cluster),
		log:                 log,
	}
	clusterList, err := repo.List(c, constants.StorageClusterConfigPath)
	if err != nil {
		return nil, fmt.Errorf("get storage cluster list error:%s", err)
	}

	// init exist cluster list
	for _, cluster := range clusterList {
		stateMachine.addCluster(cluster)
	}
	// new storage config discovery
	stateMachine.discovery = discoveryFactory.CreateDiscovery(constants.StorageClusterConfigPath, stateMachine)
	if err := stateMachine.discovery.Discovery(); err != nil {
		return nil, fmt.Errorf("discovery storage cluster config error:%s", err)
	}
	log.Info("storage cluster state machine started")
	return stateMachine, nil
}

// OnCreate creates and starts cluster controller when receive create event
func (c *clusterStateMachine) OnCreate(key string, resource []byte) {
	c.log.Info("storage cluster be created", logger.String("key", key))
	c.addCluster(resource)
}

// OnDelete deletes cluster controller from cache, closes it
func (c *clusterStateMachine) OnDelete(key string) {
	name := pathutil.GetName(key)
	c.mutex.Lock()
	c.deleteCluster(name)
	c.mutex.Unlock()
}

// GetCluster returns cluster controller for maintain the metadata of storage cluster
func (c *clusterStateMachine) GetCluster(name string) Cluster {
	var cluster Cluster
	c.mutex.RLock()
	cluster = c.clusters[name]
	c.mutex.RUnlock()
	return cluster
}

// GetAllCluster returns all cluster controller
func (c *clusterStateMachine) GetAllCluster() []Cluster {
	var clusters []Cluster
	c.mutex.RLock()
	for _, v := range c.clusters {
		clusters = append(clusters, v)
	}
	c.mutex.RUnlock()
	return clusters
}

// Close closes state machine, cleanup and close all cluster controller
func (c *clusterStateMachine) Close() error {
	// 1) close listen for storage cluster config change
	c.discovery.Close()
	// 2) cleanup clusters and release resource
	c.mutex.Lock()
	c.cleanupCluster()
	c.mutex.Unlock()

	c.cancel()
	return nil
}

// cleanupCluster cleanups cluster controller
func (c *clusterStateMachine) cleanupCluster() {
	for _, v := range c.clusters {
		v.Close()
	}
}

// addCluster creates and starts cluster controller, if success cache it
func (c *clusterStateMachine) addCluster(resource []byte) {
	cfg := models.StorageCluster{}
	if err := json.Unmarshal(resource, &cfg); err != nil {
		c.log.Error("discovery new storage config but unmarshal error",
			logger.String("data", string(resource)), logger.Error(err))
		return
	}
	if len(cfg.Name) == 0 {
		c.log.Error("cluster name is empty", logger.Any("cfg", cfg))
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// shutdown old cluster state machine if exist
	c.deleteCluster(cfg.Name)

	repo, err := c.repoFactory.CreateRepo(cfg.Config)
	if err != nil {
		c.log.Error("new state repo error when create cluster",
			logger.Any("cfg", cfg), logger.Error(err))
		return
	}
	clusterCfg := clusterCfg{
		ctx:                 c.ctx,
		cfg:                 cfg,
		storageStateService: c.storageStateService,
		repo:                repo,
		controllerFactory:   c.controllerFactory,
		factory:             discovery.NewFactory(repo),
		shardAssignService:  c.shardAssignService,
	}
	cluster, err := c.clusterFactory.newCluster(clusterCfg)
	if err != nil {
		// IMPORTANT!!!!!!!: need clean cluster cfg resource when new cluster fail
		if cluster != nil {
			cluster.Close()
		}
		(&clusterCfg).clean()
		c.log.Error("create storage cluster error",
			logger.Any("cfg", cfg), logger.Error(err))
		return
	}
	c.clusters[cfg.Name] = cluster
}

// deleteCluster deletes the cluster if exist
func (c *clusterStateMachine) deleteCluster(name string) {
	cluster, ok := c.clusters[name]
	if ok {
		// need cleanup cluster resource
		cluster.Close()
		delete(c.clusters, name)
	}
}
