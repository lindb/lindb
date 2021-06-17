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

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	// for testing
	httpGet = http.Get
	// FlushDatabasePath represents database flush api path.
	FlushDatabasePath = "/database/flush"
)

// DatabaseFlusherAPI represents the memory database flush by manual.
type DatabaseFlusherAPI struct {
	deps *deps.HTTPDeps

	logger *logger.Logger
}

// NewDatabaseFlusherAPI create database flusher api.
func NewDatabaseFlusherAPI(deps *deps.HTTPDeps) *DatabaseFlusherAPI {
	return &DatabaseFlusherAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "databaseFlusherAPI"),
	}
}

// Register adds database flush admin url route.
func (df *DatabaseFlusherAPI) Register(route gin.IRoutes) {
	route.PUT(FlushDatabasePath, df.SubmitFlushTask)
}

// SubmitFlushTask submits the task which does flush job over memory database
func (df *DatabaseFlusherAPI) SubmitFlushTask(c *gin.Context) {
	var param struct {
		Cluster  string `json:"cluster" binding:"required"`
		Database string `json:"database" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	if df.deps.Master.IsMaster() {
		// if current node is master, submits the flush task
		if err := df.deps.Master.FlushDatabase(param.Cluster, param.Database); err != nil {
			httppkg.Error(c, err)
			return
		}
	} else {
		// if current node is not master, need forward to master node
		masterNode := df.deps.Master.GetMaster().Node
		resp, err := httpGet(fmt.Sprintf("http://%s:%d"+c.Request.RequestURI, masterNode.IP, masterNode.Port))
		if resp != nil {
			if resp.Body != nil {
				if err := resp.Body.Close(); err != nil {
					df.logger.Error("close http response body", logger.Error(err))
				}
			}

			if resp.StatusCode != http.StatusOK {
				httppkg.Error(c, fmt.Errorf("master handle error after forward"))
				return
			}
		}
		if err != nil {
			httppkg.Error(c, err)
			return
		}
	}
	httppkg.OK(c, "success")
}
