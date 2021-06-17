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

package cluster

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/pkg/http"
)

var (
	// MasterStatePath represents cluster master's state.
	MasterStatePath = "/cluster/master"
)

// MasterAPI represents query cluster master state.
type MasterAPI struct {
	deps *deps.HTTPDeps
}

// NewMasterAPI creates the master api.
func NewMasterAPI(deps *deps.HTTPDeps) *MasterAPI {
	return &MasterAPI{
		deps: deps,
	}
}

// Register adds master url route.
func (m *MasterAPI) Register(route gin.IRoutes) {
	route.GET(MasterStatePath, m.GetMasterState)
}

// GetMasterState returns the current cluster's master state.
func (m *MasterAPI) GetMasterState(c *gin.Context) {
	master := m.deps.Master.GetMaster()
	if master == nil {
		http.NotFound(c)
	} else {
		http.OK(c, master)
	}
}
