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

package cluser

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator"
)

// MasterAPI represents query cluster master state
type MasterAPI struct {
	master coordinator.Master
}

// NewMasterAPI creates the master api
func NewMasterAPI(master coordinator.Master) *MasterAPI {
	return &MasterAPI{
		master: master,
	}
}

// GetMaster returns the current cluster's master
func (m *MasterAPI) GetMaster(w http.ResponseWriter, r *http.Request) {
	master := m.master.GetMaster()
	if master == nil {
		api.NotFound(w)
	} else {
		api.OK(w, master)
	}
}
