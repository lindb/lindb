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
	"sort"

	"github.com/gin-gonic/gin"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var (
	ExplorePath = "/state/machine/explore"
)

type Param struct {
	Type        string               `form:"type" binding:"required"`
	Role        stmtpkg.MetadataType `form:"role" binding:"required"`
	StorageName string               `form:"storageName"`
}

// BrokerStateMachineAPI represents state machine explore api.
type BrokerStateMachineAPI struct {
	deps   *depspkg.HTTPDeps
	cli    client.StateMachineCli
	logger *logger.Logger
}

// NewBrokerStateMachineAPI creates broker state machine api instance.
func NewBrokerStateMachineAPI(deps *depspkg.HTTPDeps) *BrokerStateMachineAPI {
	return &BrokerStateMachineAPI{
		deps:   deps,
		cli:    client.NewStateMachineCli(),
		logger: logger.GetLogger("Broker", "StateMachineAPI"),
	}
}

// Register adds state machine url route.
func (api *BrokerStateMachineAPI) Register(route gin.IRoutes) {
	route.GET(ExplorePath, api.Explore)
}

// Explore explores the state from state machine of broker/master/storage.
// @BasePath /api/v1
// @Summary explore the state from state machine.
// @Schemes
// @Description explores the state from state machine of current node.
// @Description 1. Broker State Machine;
// @Description 2. Master State Machine;
// @Description 3. Storage State Machine;
// @Tags Internal
// @Accept json
// @Param param body Param ture "param data"
// @Produce json
// @Success 200 {array} models.Database
// @Success 200 {array} models.StorageState
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal error"
// @Router /state/machine/explore [get]
func (api *BrokerStateMachineAPI) Explore(c *gin.Context) {
	param := &Param{}
	err := c.ShouldBindQuery(param)
	if err != nil {
		http.Error(c, err)
		return
	}
	switch param.Role {
	case stmtpkg.BrokerMetadata:
		api.exploreBroker(c, param)
	case stmtpkg.MasterMetadata:
		api.exploreMaster(c, param)
	case stmtpkg.StorageMetadata:
		stateMgr := api.deps.Master.GetStateManager()
		storageCluster := stateMgr.GetStorageCluster(param.StorageName)
		if storageCluster == nil {
			http.NotFound(c)
			return
		}
		liveNodes, err := storageCluster.GetLiveNodes()
		if err != nil {
			http.Error(c, err)
			return
		}
		var nodes []models.Node
		for idx := range liveNodes {
			nodes = append(nodes, &liveNodes[idx])
		}
		http.OK(c, api.cli.FetchStateByNodes(map[string]string{"type": param.Type}, nodes))
	default:
		http.NotFound(c)
	}
}

// exploreMaster explores the state from state machine of master.
func (api *BrokerStateMachineAPI) exploreMaster(c *gin.Context, param *Param) {
	switch param.Type {
	case constants.StorageState:
		api.writeStorageState(c, api.deps.Master.GetStateManager().GetStorageStates())
	case constants.StorageConfig:
		nodes := api.deps.Master.GetStateManager().GetStorages()
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Config.Namespace < nodes[j].Config.Namespace
		})
		http.OK(c, nodes)
	case constants.DatabaseConfig:
		api.writeDatabaseState(c, api.deps.Master.GetStateManager().GetDatabases())
	case constants.ShardAssignment:
		shardAssignments := api.deps.Master.GetStateManager().GetShardAssignments()
		sort.Slice(shardAssignments, func(i, j int) bool {
			return shardAssignments[i].Name < shardAssignments[j].Name
		})
		http.OK(c, shardAssignments)
	case constants.Master:
		// return master slice, because common logic read state from repo.
		http.OK(c, []*models.Master{api.deps.Master.GetMaster()})
	default:
		http.NotFound(c)
	}
}

// exploreMaster explores the state from state machine of broker.
func (api *BrokerStateMachineAPI) exploreBroker(c *gin.Context, param *Param) {
	switch param.Type {
	case constants.StorageState:
		api.writeStorageState(c, api.deps.StateMgr.GetStorageList())
	case constants.LiveNode:
		nodes := api.deps.StateMgr.GetLiveNodes()
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Indicator() < nodes[j].Indicator()
		})
		http.OK(c, nodes)
	case constants.DatabaseConfig:
		api.writeDatabaseState(c, api.deps.StateMgr.GetDatabases())
	default:
		http.NotFound(c)
	}
}

// writeDatabaseState writes response with database.
func (api *BrokerStateMachineAPI) writeDatabaseState(c *gin.Context, dbs []models.Database) {
	sort.Slice(dbs, func(i, j int) bool {
		return dbs[i].Name < dbs[j].Name
	})
	http.OK(c, dbs)
}

// writeDatabaseState writes response with storage.
func (api *BrokerStateMachineAPI) writeStorageState(c *gin.Context, storages []*models.StorageState) {
	sort.Slice(storages, func(i, j int) bool {
		return storages[i].Name < storages[j].Name
	})
	http.OK(c, storages)
}
