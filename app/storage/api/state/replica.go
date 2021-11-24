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

	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replica"
)

var (
	ReplicaPath = "/state/replica"
)

// ReplicaAPI represents internal replica state rest api.
type ReplicaAPI struct {
	walMgr replica.WriteAheadLogManager
	logger *logger.Logger
}

// NewReplicaAPI creates replica state api instance.
func NewReplicaAPI(walMgr replica.WriteAheadLogManager) *ReplicaAPI {
	return &ReplicaAPI{
		walMgr: walMgr,
		logger: logger.GetLogger("storage", "ReplicaAPI"),
	}
}

// Register adds explore url route.
func (d *ReplicaAPI) Register(route gin.IRoutes) {
	route.GET(ReplicaPath, d.GetReplicaState)
}

// GetReplicaState returns replica state by given database's name.
func (d *ReplicaAPI) GetReplicaState(c *gin.Context) {
	var param struct {
		DB string `form:"db" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	rs := d.walMgr.GetReplicaState(param.DB)
	httppkg.OK(c, rs)
}
