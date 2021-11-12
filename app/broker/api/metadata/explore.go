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

package metadata

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

var (
	ExplorePath     = "/metadata/explore"
	ExploreRepoPath = "/metadata/explore/repo"
)

// ExploreAPI represents metadata explore rest api.
type ExploreAPI struct {
	deps   *deps.HTTPDeps
	logger *logger.Logger
}

// NewExploreAPI creates explore api instance.
func NewExploreAPI(deps *deps.HTTPDeps) *ExploreAPI {
	return &ExploreAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "ExploreAPI"),
	}
}

// Register adds explore url route.
func (d *ExploreAPI) Register(route gin.IRoutes) {
	route.GET(ExplorePath, d.Explore)
	route.GET(ExploreRepoPath, d.ExploreRepo)
}

// Explore returns explore define info.
func (d *ExploreAPI) Explore(c *gin.Context) {
	httppkg.OK(c, map[string]interface{}{
		"broker":  broker.StateMachinePaths,
		"master":  master.StateMachinePaths,
		"storage": storage.StateMachinePaths,
	})
}

// ExploreRepo explores state repository by role/type.
func (d *ExploreAPI) ExploreRepo(c *gin.Context) {
	var param struct {
		Role        string `form:"role" binding:"required"`
		Type        string `form:"type" binding:"required"`
		StorageName string `form:"storageName"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	var stateMachineInfo models.StateMachineInfo
	var ok bool
	switch param.Role {
	case "broker":
		stateMachineInfo, ok = broker.StateMachinePaths[param.Type]
	case "master":
		stateMachineInfo, ok = master.StateMachinePaths[param.Type]
	case "storage":
		if param.StorageName == "" {
			httppkg.Error(c, fmt.Errorf("storage name cannot be empty"))
			return
		}
		if d.deps.Master.IsMaster() {
			// if current node is master, explore storage data.
			stateMachineInfo, ok = storage.StateMachinePaths[param.Type]
			if !ok {
				httppkg.NotFound(c)
				return
			}
			stateMgr := d.deps.Master.GetStateManager()
			storageCluster := stateMgr.GetStorageCluster(param.StorageName)
			if storageCluster == nil {
				httppkg.NotFound(c)
				return
			}
			d.exploreData(c, storageCluster.GetRepo(), stateMachineInfo)
			return
		}
		// if current node is not master, reverse proxy to master
		masterNode := d.deps.Master.GetMaster()
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", masterNode.Node.HostIP, masterNode.Node.HTTPPort)
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = c.Request.URL.RawQuery
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	}
	if !ok {
		httppkg.NotFound(c)
		return
	}

	d.exploreData(c, d.deps.Repo, stateMachineInfo)
}

// exploreData explores state repository data by given path.
func (d *ExploreAPI) exploreData(c *gin.Context, repo state.Repository, stateMachineInfo models.StateMachineInfo) {
	ctx, cancel := d.deps.WithTimeout()
	defer cancel()
	var rs []interface{}
	err := repo.WalkEntry(ctx, stateMachineInfo.Path, func(key, value []byte) {
		r := stateMachineInfo.CreateState()
		err0 := encoding.JSONUnmarshal(value, r)
		if err0 != nil {
			d.logger.Warn("unmarshal metadata info err, ignore it",
				logger.String("key", string(key)),
				logger.String("data", string(value)))
			return
		}
		rs = append(rs, r)
	})
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	httppkg.OK(c, rs)
}
