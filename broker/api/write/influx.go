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

package write //nolint:dupl

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/constants"
	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/ingestion/influx"
	"github.com/lindb/lindb/pkg/http"
)

var (
	InfluxWritePath = "/write/influx"
)

// InfluxWriter processes Influxdb line protocol.
type InfluxWriter struct {
	deps *deps.HTTPDeps
}

// NewInfluxWriter creates influx writer.
func NewInfluxWriter(deps *deps.HTTPDeps) *InfluxWriter {
	return &InfluxWriter{
		deps: deps,
	}
}

// Register adds influx write url route.
func (iw *InfluxWriter) Register(route gin.IRoutes) {
	route.PUT(InfluxWritePath, iw.Write)
}

func (iw *InfluxWriter) Write(c *gin.Context) {
	var param struct {
		Database  string `form:"db" binding:"required"`
		Namespace string `form:"ns"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	if param.Namespace == "" {
		param.Namespace = constants.DefaultNamespace
	}
	enrichedTags, err := ingestCommon.ExtractEnrichTags(c.Request)
	if err != nil {
		http.Error(c, err)
		return
	}
	metricList, err := influx.Parse(c.Request, enrichedTags, param.Namespace)
	if err != nil {
		http.Error(c, err)
		return
	}
	if err := iw.deps.CM.Write(param.Database, metricList); err != nil {
		http.Error(c, err)
		return
	}
	http.NoContent(c)
}
