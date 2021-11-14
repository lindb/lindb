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

package monitoring

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/internal/linmetric"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/tag"
)

var (
	ExploreCurrentPath = "/state/explore/current"
)

// ExploreAPI represents monitoring metric explore rest api.
type ExploreAPI struct {
	globalKeyValues tag.Tags
	logger          *logger.Logger
}

// NewExploreAPI creates explore api instance.
func NewExploreAPI(globalKeyValues tag.Tags) *ExploreAPI {
	return &ExploreAPI{
		globalKeyValues: globalKeyValues,
		logger:          logger.GetLogger("monitoring", "ExploreAPI"),
	}
}

// Register adds explore url route.
func (d *ExploreAPI) Register(route gin.IRoutes) {
	route.GET(ExploreCurrentPath, d.ExploreCurrent)
}

// ExploreCurrent explores current node monitoring metric.
func (d *ExploreAPI) ExploreCurrent(c *gin.Context) {
	var param struct {
		Names []string `form:"names" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}

	// find metric by name from default metric registry
	rs := linmetric.FindMetricList(param.Names)
	globalKeyValues := d.globalKeyValues
	for _, metricList := range rs {
		for _, metric := range metricList {
			for _, kv := range globalKeyValues {
				// append global tags
				metric.Tags[string(kv.Key)] = string(kv.Value)
			}
		}
	}
	httppkg.OK(c, rs)
}
