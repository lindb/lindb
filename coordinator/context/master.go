package context

import (
	"go.uber.org/zap"

	"github.com/eleme/lindb/coordinator/database"
	"github.com/eleme/lindb/coordinator/storage"
	"github.com/eleme/lindb/pkg/logger"
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
	log := logger.GetLogger()
	if err := m.stateMachine.StorageCluster.Close(); err != nil {
		log.Error("close storage cluster state machine error", zap.Error(err), zap.Stack("stack"))
	}
	if err := m.stateMachine.DatabaseAdmin.Close(); err != nil {
		log.Error("close database admin state machine error", zap.Error(err), zap.Stack("stack"))
	}
}
