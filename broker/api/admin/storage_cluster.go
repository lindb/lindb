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
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/service"
)

// StorageClusterAPI represents storage cluster admin rest api
type StorageClusterAPI struct {
	storageClusterService service.StorageClusterService
}

// NewStorageClusterAPI create storage cluster api
func NewStorageClusterAPI(storageClusterService service.StorageClusterService) *StorageClusterAPI {
	return &StorageClusterAPI{
		storageClusterService: storageClusterService,
	}
}

// Create creates config of storage cluster
func (s *StorageClusterAPI) Create(w http.ResponseWriter, r *http.Request) {
	storage := &config.StorageCluster{}
	err := api.GetJSONBodyFromRequest(r, storage)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = s.storageClusterService.Save(storage)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

// GetByName gets storage cluster by name
func (s *StorageClusterAPI) GetByName(w http.ResponseWriter, r *http.Request) {
	name, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	cluster, err := s.storageClusterService.Get(name)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, cluster)
}

// DeleteByName deletes storage cluster by name
func (s *StorageClusterAPI) DeleteByName(w http.ResponseWriter, r *http.Request) {
	name, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = s.storageClusterService.Delete(name)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

// List lists all storage clusters
func (s *StorageClusterAPI) List(w http.ResponseWriter, r *http.Request) {
	clusters, err := s.storageClusterService.List()
	if err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, clusters)
}
