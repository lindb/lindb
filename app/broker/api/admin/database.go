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

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	DatabasePath = "/database"
)

// DatabaseAPI represents database admin rest api
type DatabaseAPI struct {
	deps   *depspkg.HTTPDeps
	logger *logger.Logger
}

// NewDatabaseAPI creates database api instance
func NewDatabaseAPI(deps *depspkg.HTTPDeps) *DatabaseAPI {
	return &DatabaseAPI{
		deps:   deps,
		logger: logger.GetLogger("Broker", "DatabaseAPI"),
	}
}

// Register adds database admin url route.
func (d *DatabaseAPI) Register(route gin.IRoutes) {
	route.GET(DatabasePath, d.GetByName)
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
	database, err := d.getByName(param.DatabaseName)
	if err != nil {
		http.NotFound(c)
		return
	}
	http.OK(c, database)
}

func (d *DatabaseAPI) getByName(name string) (*models.Database, error) {
	ctx, cancel := d.deps.WithTimeout()
	defer cancel()

	configBytes, err := d.deps.Repo.Get(ctx, constants.GetDatabaseConfigPath(name))
	if err != nil {
		return nil, err
	}
	database := &models.Database{}
	err = encoding.JSONUnmarshal(configBytes, database)
	if err != nil {
		return nil, err
	}
	return database, nil
}
