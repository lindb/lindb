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
	"github.com/gin-gonic/gin"

	httppkg "github.com/lindb/common/pkg/http"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb"
)

var (
	DatabaseCfgPath = "/state/metadata/local/database/config"
)

// MetadataAPI represents internal metadata state rest api.
type MetadataAPI struct {
	engine tsdb.Engine
}

// NewMetadataAPI creates a metadata api instance.
func NewMetadataAPI(engine tsdb.Engine) *MetadataAPI {
	return &MetadataAPI{
		engine: engine,
	}
}

// Register adds metadata api url route.
func (m *MetadataAPI) Register(route gin.IRoutes) {
	route.GET(DatabaseCfgPath, m.GetLocalAllDatabaseCfg)
}

// GetLocalAllDatabaseCfg returns the configuration map of all local databases.
func (m *MetadataAPI) GetLocalAllDatabaseCfg(c *gin.Context) {
	databases := m.engine.GetAllDatabases()
	cfgMap := make(map[string]models.DatabaseConfig)
	for name, db := range databases {
		cfgMap[name] = *db.GetConfig()
	}
	httppkg.OK(c, cfgMap)
}
