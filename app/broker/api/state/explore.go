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
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

// for testing
var (
	defaultClient = http.Client{Timeout: time.Second * 10}
	doRequest     = defaultClient.Do
)

var (
	ExplorePath         = "/state/explore"
	ExploreLiveNodePath = "/state/explore/alive"
)

// ExploreAPI represents internal state explore rest api.
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
	route.GET(ExploreLiveNodePath, d.ExploreLiveNode)
}

// ExploreLiveNode explores live nodes for given role.
func (d *ExploreAPI) ExploreLiveNode(c *gin.Context) {
	var param struct {
		Role string `form:"role" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	switch param.Role {
	case constants.BrokerRole:
		httppkg.OK(c, d.deps.StateMgr.GetLiveNodes())
		return
	case constants.StorageRole:
		httppkg.OK(c, d.deps.StateMgr.GetStorageList())
		return
	}
	httppkg.NotFound(c)
}

// Explore explores the state of cluster by given params.
// returns internal state metric.
func (d *ExploreAPI) Explore(c *gin.Context) {
	var param struct {
		Role        string   `form:"role" binding:"required"`
		Names       []string `form:"names" binding:"required"`
		StorageName string   `form:"storageName"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	switch param.Role {
	case constants.BrokerRole:
		liveNodes := d.deps.StateMgr.GetLiveNodes()
		var nodes []models.Node
		for idx := range liveNodes {
			nodes = append(nodes, &liveNodes[idx])
		}
		d.fetchStateData(c, nodes)
		return
	case constants.StorageRole:
		if param.StorageName == "" {
			httppkg.Error(c, fmt.Errorf("storage name cannot be empty"))
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
		return
	}
	httppkg.NotFound(c)
}

// fetchStateData fetches the state metric from each live nodes.
func (d *ExploreAPI) fetchStateData(c *gin.Context, nodes []models.Node) {
	size := len(nodes)
	if size == 0 {
		httppkg.NotFound(c)
		return
	}
	q := c.Request.URL.Query()
	delete(q, "role")
	params := q.Encode()
	result := make([]map[string][]*models.StateMetric, size)
	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			req, _ := http.NewRequest(http.MethodGet, node.HTTPAddress(), nil)
			req.URL.Path = c.Request.URL.Path + "/current"
			req.URL.RawQuery = params
			var metric map[string][]*models.StateMetric
			if err := get(req, func(body io.Reader) error {
				return json.NewDecoder(body).Decode(&metric)
			}); err == nil {
				result[i] = metric
			}
		}()
	}
	wait.Wait()
	rs := make(map[string][]*models.StateMetric)
	for _, metricList := range result {
		if metricList == nil {
			continue
		}
		for name, list := range metricList {
			l, ok := rs[name]
			if ok {
				l = append(l, list...)
				rs[name] = l
			} else {
				rs[name] = list
			}
		}
	}
	httppkg.OK(c, rs)
}

// get does http get request, then returns the internal metric for given target node.
func get(req *http.Request, decoder func(body io.Reader) error) error {
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := doRequest(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return decoder(resp.Body)
}
