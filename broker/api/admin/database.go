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
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/service"
)

// DatabaseAPI represents database admin rest api
type DatabaseAPI struct {
	databaseService service.DatabaseService
}

// NewDatabaseAPI creates database api instance
func NewDatabaseAPI(databaseService service.DatabaseService) *DatabaseAPI {
	return &DatabaseAPI{
		databaseService: databaseService,
	}
}

// GetByName gets a database config by the name.
func (d *DatabaseAPI) GetByName(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	database, err := d.databaseService.Get(databaseName)
	if err != nil {
		api.NotFound(w)
		return
	}
	api.OK(w, database)
}

// Save creates the database config if there is no database
// config with the name database.Name, otherwise update the config
func (d *DatabaseAPI) Save(w http.ResponseWriter, r *http.Request) {
	database := &models.Database{}
	err := api.GetJSONBodyFromRequest(r, database)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = d.databaseService.Save(database)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

// List returns all database configs
func (d *DatabaseAPI) List(w http.ResponseWriter, r *http.Request) {
	dbs, err := d.databaseService.List()
	if err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, dbs)
}
