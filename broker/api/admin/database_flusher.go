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
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/pkg/logger"
)

var adminLogger = logger.GetLogger("broker", "adminAPI")
var httpGet = http.Get

// DatabaseFlusherAPI represents the memory database flush by manual
type DatabaseFlusherAPI struct {
	master coordinator.Master
}

// NewDatabaseFlusherAPI create database flusher api
func NewDatabaseFlusherAPI(master coordinator.Master) *DatabaseFlusherAPI {
	return &DatabaseFlusherAPI{master: master}
}

// SubmitFlushTask submits the task which does flush job over memory database
func (df *DatabaseFlusherAPI) SubmitFlushTask(w http.ResponseWriter, r *http.Request) {
	cluster, err := api.GetParamsFromRequest("cluster", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	databaseName, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	if df.master.IsMaster() {
		// if current node is master, submits the flush task
		if err := df.master.FlushDatabase(cluster, databaseName); err != nil {
			api.Error(w, err)
			return
		}
	} else {
		// if current node is not master, need forward to master node
		masterNode := df.master.GetMaster().Node
		resp, err := httpGet(fmt.Sprintf("http://%s:%d"+r.RequestURI, masterNode.IP, masterNode.Port))
		if resp != nil {
			if resp.Body != nil {
				if err := resp.Body.Close(); err != nil {
					adminLogger.Error("close http response body", logger.Error(err))
				}
			}

			if resp.StatusCode != http.StatusOK {
				api.Error(w, fmt.Errorf("master handle error after forward"))
				return
			}
		}
		if err != nil {
			api.Error(w, err)
			return
		}
	}
	api.OK(w, "success")
}
