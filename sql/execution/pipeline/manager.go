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

package pipeline

import (
	"sync"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/operator"
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
