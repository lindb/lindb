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
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb"
)

var (
	MemoryDatabase = "/state/tsdb/memory"
)

// TSDBAPI represents tsdb internal state rest api.
type TSDBAPI struct {
	logger logger.Logger
}

// NewTSDBAPI creates a tsdb state api instance.
func NewTSDBAPI() *TSDBAPI {
	return &TSDBAPI{
		logger: logger.GetLogger("Storage", "TSDBAPI"),
	}
}

// Register adds the route for tsdb state api.
func (db *TSDBAPI) Register(route gin.IRoutes) {
	route.GET(MemoryDatabase, db.GetMemoryDatabaseState)
}

// GetMemoryDatabaseState returns memory database
func (db *TSDBAPI) GetMemoryDatabaseState(c *gin.Context) {
	var param struct {
		DB string `form:"db" binding:"required"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	var rs []models.DataFamilyState
	tsdb.GetFamilyManager().WalkEntry(func(family tsdb.DataFamily) {
		if param.DB == family.Shard().Database().Name() {
			rs = append(rs, family.GetState())
		}
	})
	httppkg.OK(c, rs)
}
