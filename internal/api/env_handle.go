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
	"github.com/lindb/common/pkg/http"

	"github.com/lindb/lindb/config"
)

var EnvPath = "/env"

// EnvAPI represents LinDB's env api.
type EnvAPI struct {
	envs []config.Env
}

// NewEnvAPI creates a EnvAPI instance.
func NewEnvAPI(envs []config.Env) *EnvAPI {
	return &EnvAPI{
		envs: envs,
	}
}

// Register adds request env url route.
func (api *EnvAPI) Register(route gin.IRoutes) {
	route.GET(EnvPath, api.GetEnv)
}

// GetEnv returns LinDB's env vars.
func (api *EnvAPI) GetEnv(c *gin.Context) {
	http.OK(c, api.envs)
}
