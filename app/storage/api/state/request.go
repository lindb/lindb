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

package state

import (
	"github.com/gin-gonic/gin"

	httppkg "github.com/lindb/common/pkg/http"

	"github.com/lindb/lindb/query"
)

const (
	RequestPath  = "/state/request"
	RequestsPath = "/state/requests"
)

// RequestAPI represents lin query request stats related api.
type RequestAPI struct {
}

// NewRequestAPI creates a RequestAPI instance.
func NewRequestAPI() *RequestAPI {
	return &RequestAPI{}
}

// Register adds request api route.
func (r *RequestAPI) Register(route gin.IRoutes) {
	route.GET(RequestPath, r.GetRequestState)
	route.GET(RequestsPath, r.GetAllAliveRequests)
}

// GetRequestState returns request stats by given request id.
func (r *RequestAPI) GetRequestState(c *gin.Context) {
	var param struct {
		RequestID string `form:"requestId" binding:"required"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	pipeline := query.GetPipelineManager().GetPipeline(param.RequestID)
	if pipeline == nil {
		httppkg.NotFound(c)
		return
	}
	httppkg.OK(c, pipeline.Stats())
}

// GetAllAliveRequests returns all alive requests.
func (r *RequestAPI) GetAllAliveRequests(c *gin.Context) {
	httppkg.OK(c, query.GetPipelineManager().GetAllAlivePipelines())
}
