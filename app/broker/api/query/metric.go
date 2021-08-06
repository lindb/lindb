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

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
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

	ctx, cancel := context.WithTimeout(context.Background(), m.deps.BrokerCfg.Query.Timeout.Duration())
	defer cancel()

	metricQuery := m.deps.QueryFactory.NewMetricQuery(ctx, param.Database, param.SQL)
	resultSet, err := metricQuery.WaitResponse()
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, resultSet)
}
