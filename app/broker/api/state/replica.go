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
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/models"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	ReplicaPath = "/state/replica"
)

// ReplicaAPI represents internal wal replica state rest api.
type ReplicaAPI struct {
	deps   *deps.HTTPDeps
	logger *logger.Logger
}

// NewReplicaAPI creates replica api instance.
func NewReplicaAPI(deps *deps.HTTPDeps) *ReplicaAPI {
	return &ReplicaAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "NewReplicaAPI"),
	}
}

// Register adds replica state url route.
func (d *ReplicaAPI) Register(route gin.IRoutes) {
	route.GET(ReplicaPath, d.GetReplicaState)
}

// GetReplicaState returns wal replica state.
func (d *ReplicaAPI) GetReplicaState(c *gin.Context) {
	var param struct {
		StorageName string `form:"storageName" binding:"required"`
		DB          string `form:"db" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	storage, ok := d.deps.StateMgr.GetStorage(param.StorageName)
	if !ok {
		httppkg.NotFound(c)
		return
	}
	liveNodes := storage.LiveNodes
	var nodes []models.Node
	for id := range liveNodes {
		n := liveNodes[id]
		nodes = append(nodes, &n)
	}
	d.fetchStateData(c, nodes)
}

// fetchStateData fetches the state metric from each live nodes.
func (d *ReplicaAPI) fetchStateData(c *gin.Context, nodes []models.Node) {
	size := len(nodes)
	if size == 0 {
		httppkg.NotFound(c)
		return
	}
	q := c.Request.URL.Query()
	params := q.Encode()
	result := make([][]models.FamilyLogReplicaState, size)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			req, _ := http.NewRequest(http.MethodGet, node.HTTPAddress(), nil)
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = params
			var state []models.FamilyLogReplicaState
			if err := get(req, func(body io.Reader) error {
				return json.NewDecoder(body).Decode(&state)
			}); err == nil {
				result[i] = state
			}
		}()
	}
	wait.Wait()
	rs := make(map[string][]models.FamilyLogReplicaState)
	for idx := range nodes {
		rs[nodes[idx].Indicator()] = result[idx]
	}
	httppkg.OK(c, rs)
}
