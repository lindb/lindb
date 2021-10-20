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
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	DatabasePath     = "/database"
	ListDatabasePath = "/database/list"
)

// DatabaseAPI represents database admin rest api
type DatabaseAPI struct {
	deps   *deps.HTTPDeps
	logger *logger.Logger
}

// NewDatabaseAPI creates database api instance
func NewDatabaseAPI(deps *deps.HTTPDeps) *DatabaseAPI {
	return &DatabaseAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "DatabaseAPI"),
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

// Save creates the database config if there is no database
// config with the name database.Name, otherwise update the config
func (d *DatabaseAPI) Save(c *gin.Context) {
	database := &models.Database{}
	if err := c.ShouldBind(&database); err != nil {
		http.Error(c, err)
		return
	}
	if err := d.saveDataBase(database); err != nil {
		http.Error(c, err)
		return
	}
	http.NoContent(c)
}

func (d *DatabaseAPI) saveDataBase(database *models.Database) error {
	if len(database.Storage) == 0 {
		//TODO add golang tag?
		return fmt.Errorf("storage name cannot be empty")
	}
	if database.NumOfShard <= 0 {
		return fmt.Errorf("num of shard must be > 0")
	}
	if database.ReplicaFactor <= 0 {
		return fmt.Errorf("replica factor must be > 0")
	}
	opt := database.Option
	// validate time series engine option
	if err := opt.Validate(); err != nil {
		return err
	}
	// set default value
	(&opt).Default()

	data := encoding.JSONMarshal(database)

	ctx, cancel := d.deps.WithTimeout()
	defer cancel()
	d.logger.Info("Saving Database", logger.String("config", string(data)))
	return d.deps.Repo.Put(ctx, constants.GetDatabaseConfigPath(database.Name), data)
}

// List returns all database configs
func (d *DatabaseAPI) List(c *gin.Context) {
	dbs, err := d.ListDataBase()
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, dbs)
}

func (d *DatabaseAPI) ListDataBase() ([]*models.Database, error) {
	ctx, cancel := d.deps.WithTimeout()
	defer cancel()

	var result []*models.Database
	data, err := d.deps.Repo.List(ctx, constants.DatabaseConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		db := &models.Database{}
		err = encoding.JSONUnmarshal(val.Value, db)
		if err != nil {
			d.logger.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			db.Desc = db.String()
			result = append(result, db)
		}
	}
	return result, nil
}
