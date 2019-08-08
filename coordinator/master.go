package coordinator

import (
	"context"
	"sync"

	coCtx "github.com/lindb/lindb/coordinator/context"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/elect"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

// Master represents all metadata/state controller, only has one active master in broker cluster.
// Master will control all storage cluster metadata, update state, then notify each broker node.
type Master interface {
	// Start starts master do election master, if success build master context,
	// starts state machine do cluster coordinate such metadata, cluster state etc.
	Start() error
	// IsMaster returns current node if is master
	IsMaster() bool
	// Stop stops master if current node is master, cleanup master context and stops state machine
	Stop()
}

// master implements master interface
type master struct {
	node             models.Node
	repo             state.Repository
	discoveryFactory discovery.Factory

	ctx    context.Context
	cancel context.CancelFunc

	elect          elect.Election
	masterCtx      *coCtx.MasterContext
	taskController task.Controller
	clusterFactory storage.ClusterFactory
	repoFactory    state.RepositoryFactory

	storageStateService service.StorageStateService
	shardAssignService  service.ShardAssignService

	mutex sync.Mutex

	log *logger.Logger
}

// NewMaster create master for current node
func NewMaster(
	repo state.Repository,
	node models.Node,
	ttl int64,
	taskController task.Controller,
	discoveryFactory discovery.Factory,
	repoFactory state.RepositoryFactory,
	clusterFactory storage.ClusterFactory,
	storageStateService service.StorageStateService,
	shardAssignService service.ShardAssignService,
) Master {
	ctx, cancel := context.WithCancel(context.Background())
	m := &master{
		repo:                repo,
		taskController:      taskController,
		discoveryFactory:    discoveryFactory,
		clusterFactory:      clusterFactory,
		repoFactory:         repoFactory,
		storageStateService: storageStateService,
		shardAssignService:  shardAssignService,
		node:                node,
		ctx:                 ctx,
		cancel:              cancel,
		log:                 logger.GetLogger("coordinator/master"),
	}
	// create master election
	m.elect = elect.NewElection(ctx, repo, node, ttl, m)
	return m
}

// OnFailOver invoked after master electing, current node become a new master
func (m *master) OnFailOver() {
	m.mutex.Lock()

	stateMachine := &coCtx.StateMachine{}
	storageCluster, err := storage.NewClusterStateMachine(m.ctx, m.repo,
		m.taskController, m.discoveryFactory, m.clusterFactory, m.repoFactory,
		m.storageStateService, m.shardAssignService)
	if err != nil {
		//TODO modify
		m.log.Error("start storage cluster state machine error", logger.Error(err))
		return
	}
	stateMachine.StorageCluster = storageCluster

	databaseAdmin, err := database.NewAdminStateMachine(m.ctx, m.discoveryFactory, storageCluster)
	if err != nil {
		m.log.Error("start database admin state machine error", logger.Error(err))
		return
	}
	stateMachine.DatabaseAdmin = databaseAdmin

	m.masterCtx = coCtx.NewMasterContext(stateMachine)

	//FIXME resign if init master context
	defer m.mutex.Unlock()
}

// OnResignation invoked current node is master, before re-electing
func (m *master) OnResignation() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.masterCtx.Close()
	m.masterCtx = nil
}

// IsMaster returns current node if is master
func (m *master) IsMaster() bool {
	return m.elect.IsMaster()
}

// Start starts master do election master, if success build master context,
// starts state machine do cluster coordinate such metadata, cluster state etc.
func (m *master) Start() error {
	m.elect.Initialize()
	m.elect.Elect()

	return nil
}

// Stop stops master if current node is master, cleanup master context and stops state machine
func (m *master) Stop() {
	// close master elect
	m.elect.Close()

	m.cancel()
}
