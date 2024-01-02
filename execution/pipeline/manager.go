package pipeline

import (
	"sync"

	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

var DriverManager *driverManager

func init() {
	DriverManager = NewDriverManager()
}

type driverManager struct {
	operators map[model.OperatorKey]operator.SourceOperator

	lock sync.RWMutex
}

func NewDriverManager() *driverManager {
	return &driverManager{
		operators: make(map[model.OperatorKey]operator.SourceOperator),
	}
}

func (mgr *driverManager) RegisterSourceOperator(taskID model.TaskID, source operator.SourceOperator) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	mgr.operators[model.OperatorKey{
		TaskID: taskID,
		NodeID: source.GetSourceID(),
	}] = source
}

func (mgr *driverManager) GetSourceOperator(taskID model.TaskID, nodeID plan.PlanNodeID) operator.SourceOperator {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	return mgr.operators[model.OperatorKey{
		TaskID: taskID,
		NodeID: nodeID,
	}]
}
