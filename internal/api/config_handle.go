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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/http"
)

var (
	ConfigPath = "/config"
)

// ConfigAPI represents current configuration explore rest api..
type ConfigAPI struct {
	node models.Node
	cfg  config.Configuration
}

// NewConfigAPI creates a ConfigAPI instance.
func NewConfigAPI(node models.Node, cfg config.Configuration) *ConfigAPI {
	return &ConfigAPI{
		node: node,
		cfg:  cfg,
	}
}

// Register adds config explore url route.
func (h *ConfigAPI) Register(route gin.IRoutes) {
	route.GET(ConfigPath, h.Configuration)
}

// Configuration returns current node's configuration.

// @Summary current node's configuration
// @Description return current node's configuration.
// @Tags State
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Router /config [get]
func (h *ConfigAPI) Configuration(c *gin.Context) {
	http.OK(c, map[string]interface{}{"node": h.node, "config": h.cfg.TOML()})
}
