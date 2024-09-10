package pipeline

import (
	"sync"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

var (
	DriverManager *driverManager

	log = logger.GetLogger("Pipeline", "DriverManager")
)

func init() {
	DriverManager = NewDriverManager()
}

type driver struct {
	sources   map[plan.PlanNodeID]operator.SourceOperator
	requestID model.RequestID
}

func newDriver(requestID model.RequestID) *driver {
	return &driver{
		requestID: requestID,
		sources:   make(map[plan.PlanNodeID]operator.SourceOperator),
	}
}

func (o *driver) registerSourceOperator(source operator.SourceOperator) {
	_, ok := o.sources[source.GetSourceID()]
	if !ok {
		o.sources[source.GetSourceID()] = source
	} else {
		log.Warn("source operator exist", logger.Any("requestID", o.requestID),
			logger.Any("sourceNode", source.GetSourceID()))
	}
}

func (o *driver) getSourceOperator(nodeID plan.PlanNodeID) operator.SourceOperator {
	return o.sources[nodeID]
}

type driverManager struct {
	drivers map[model.RequestID]*driver

	lock sync.RWMutex
}

func NewDriverManager() *driverManager {
	return &driverManager{
		drivers: make(map[model.RequestID]*driver),
	}
}

func (mgr *driverManager) RegisterSourceOperator(taskID model.TaskID, source operator.SourceOperator) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	driver, ok := mgr.drivers[taskID.RequestID]
	if !ok {
		driver = newDriver(taskID.RequestID)
		driver.registerSourceOperator(source)
		mgr.drivers[taskID.RequestID] = driver
	} else {
		driver.registerSourceOperator(source)
	}
}

func (mgr *driverManager) Cleanup(requestID model.RequestID) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	delete(mgr.drivers, requestID)
}

func (mgr *driverManager) GetSourceOperator(taskID model.TaskID, nodeID plan.PlanNodeID) operator.SourceOperator {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	driver, ok := mgr.drivers[taskID.RequestID]
	if !ok {
		return nil
	}
	return driver.getSourceOperator(nodeID)
}
