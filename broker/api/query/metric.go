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

package query

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/http"
)

var (
	MetricQueryPath = "/query/metric"
)

// MetricAPI represents the metric query api
type MetricAPI struct {
	deps *deps.HTTPDeps
}

// NewMetricAPI creates the metric query api
func NewMetricAPI(deps *deps.HTTPDeps) *MetricAPI {
	return &MetricAPI{
		deps: deps,
	}
}

// Register adds metric query url route.
func (m *MetricAPI) Register(route gin.IRoutes) {
	route.GET(MetricQueryPath, m.Search)
}

// Search searches the metric data based on database and sql.
func (m *MetricAPI) Search(c *gin.Context) {
	var param struct {
		Database string `form:"db" binding:"required"`
		SQL      string `form:"sql" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	//FIXME add timeout cfg
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	exec := m.deps.ExecutorFct.NewBrokerExecutor(ctx, param.Database, param.SQL,
		m.deps.StateMachines.ReplicaStatusSM, m.deps.StateMachines.NodeSM, m.deps.StateMachines.DatabaseSM,
		m.deps.JobManager)
	exec.Execute()

	brokerExecutor := exec.(parallel.BrokerExecutor)
	exeCtx := brokerExecutor.ExecuteContext()

	//FIXME timeout logic use select
	resultCh := exeCtx.ResultCh()
	for result := range resultCh {
		exeCtx.Emit(result)
	}

	resultSet, err := exeCtx.ResultSet()
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, resultSet)
}
