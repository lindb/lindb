package context

import (
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/pkg/logger"
)

// StateMachine represents all state machine for master
type StateMachine struct {
	StorageCluster storage.ClusterStateMachine
	DatabaseAdmin  database.AdminStateMachine
}

// MasterContext represents master context, creates it after node elect master
type MasterContext struct {
	stateMachine *StateMachine
}

// NewMasterContext creates master context using state machine
func NewMasterContext(stateMachine *StateMachine) *MasterContext {
	return &MasterContext{
		stateMachine: stateMachine,
	}
}

// Close closes all state machines, releases resource that master used
func (m *MasterContext) Close() {
	log := logger.GetLogger("coordinator/context")
	if err := m.stateMachine.StorageCluster.Close(); err != nil {
		log.Error("close storage cluster state machine error", logger.Error(err), logger.Stack())
	}
	if err := m.stateMachine.DatabaseAdmin.Close(); err != nil {
		log.Error("close database admin state machine error", logger.Error(err), logger.Stack())
	}
}
