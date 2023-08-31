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
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/lindb/common/pkg/http"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
)

// for testing
var (
	urlParseFn      = url.Parse
	urlParseQueryFn = url.ParseQuery
)

var (
	EnvPath = "/env"
)

// EnvAPI represents LinDB's env api.
type EnvAPI struct {
	monitor config.Monitor
	role    string
}

// NewEnvAPI creates a EnvAPI instance.
func NewEnvAPI(monitor config.Monitor, role string) *EnvAPI {
	return &EnvAPI{
		monitor: monitor,
		role:    role,
	}
}

// Register adds request env url route.
func (api *EnvAPI) Register(route gin.IRoutes) {
	route.GET(EnvPath, api.GetEnv)
}

// GetEnv returns LinDB's env vars.
func (api *EnvAPI) GetEnv(c *gin.Context) {
	monitor, err := api.getSelfMonitor()
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, &models.Env{Monitor: *monitor, Role: api.role})
}

// getSelfMonitor retruns LinDB's self-monitor vars.
func (api *EnvAPI) getSelfMonitor() (*models.Monitor, error) {
	monitorURL := api.monitor.URL

	u, err := urlParseFn(monitorURL)
	if err != nil {
		return nil, err
	}
	q, err := urlParseQueryFn(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return &models.Monitor{
		Database: q.Get("db"),
	}, nil
}
