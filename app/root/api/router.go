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

package api

import (
	"github.com/gin-gonic/gin"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/monitoring"
)

// API represents root http api.
type API struct {
	execute *ExecuteAPI
	env     *monitoring.EnvAPI
}

// NewAPI creates root http api.
func NewAPI(deps *depspkg.HTTPDeps) *API {
	return &API{
		execute: NewExecuteAPI(deps),
		env:     monitoring.NewEnvAPI(deps.Cfg.Monitor, []string{constants.RootRole}),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	v1 := router.Group(constants.APIVersion1)
	// execute lin query language statement
	api.execute.Register(v1)

	api.env.Register(v1)
}
