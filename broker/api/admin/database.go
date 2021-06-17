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

package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/http"
)

var (
	DatabasePath     = "/database"
	ListDatabasePath = "/database/list"
)

// DatabaseAPI represents database admin rest api
type DatabaseAPI struct {
	deps *deps.HTTPDeps
}

// NewDatabaseAPI creates database api instance
func NewDatabaseAPI(deps *deps.HTTPDeps) *DatabaseAPI {
	return &DatabaseAPI{
		deps: deps,
	}
}

// Register adds database admin url route.
func (d *DatabaseAPI) Register(route gin.IRoutes) {
	route.POST(DatabasePath, d.Save)
	route.GET(DatabasePath, d.GetByName)
	route.GET(ListDatabasePath, d.List)
}

// GetByName gets a database config by the name.
func (d *DatabaseAPI) GetByName(c *gin.Context) {
	var param struct {
		DatabaseName string `form:"name" binding:"required"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	database, err := d.deps.DatabaseSrv.Get(param.DatabaseName)
	if err != nil {
		http.NotFound(c)
		return
	}
	http.OK(c, database)
}

// Save creates the database config if there is no database
// config with the name database.Name, otherwise update the config
func (d *DatabaseAPI) Save(c *gin.Context) {
	database := &models.Database{}
	err := c.ShouldBind(&database)
	if err != nil {
		http.Error(c, err)
		return
	}
	err = d.deps.DatabaseSrv.Save(database)
	if err != nil {
		http.Error(c, err)
		return
	}
	http.NoContent(c)
}

// List returns all database configs
func (d *DatabaseAPI) List(c *gin.Context) {
	dbs, err := d.deps.DatabaseSrv.List()
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, dbs)
}
