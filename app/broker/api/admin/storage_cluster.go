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

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	StorageClusterPath     = "/storage/cluster"
	ListStorageClusterPath = "/storage/cluster/list"
)

type storageClusterParam struct {
	ClusterName string `form:"name" binding:"required"`
}

// StorageClusterAPI represents storage cluster admin rest api
type StorageClusterAPI struct {
	deps   *deps.HTTPDeps
	logger *logger.Logger
}

// NewStorageClusterAPI create storage cluster api
func NewStorageClusterAPI(deps *deps.HTTPDeps) *StorageClusterAPI {
	return &StorageClusterAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "StorageClusterAPI"),
	}
}

// Register adds storage admin url route.
func (s *StorageClusterAPI) Register(route gin.IRoutes) {
	route.POST(StorageClusterPath, s.Create)
	route.GET(StorageClusterPath, s.GetByName)
	route.DELETE(StorageClusterPath, s.DeleteByName)
	route.GET(ListStorageClusterPath, s.List)
}

// Create creates config of storage cluster
func (s *StorageClusterAPI) Create(c *gin.Context) {
	storage := &config.StorageCluster{}
	err := c.ShouldBind(&storage)
	if err != nil {
		http.Error(c, err)
		return
	}
	data := encoding.JSONMarshal(storage)
	ctx, cancel := s.deps.WithTimeout()
	defer cancel()
	s.logger.Info("Creating storage cluster", logger.String("config", string(data)))
	if err := s.deps.Repo.Put(ctx, constants.GetStorageClusterConfigPath(storage.Name), data); err != nil {
		http.Error(c, err)
		return
	}
	http.NoContent(c)
}

// GetByName gets storage cluster by name
func (s *StorageClusterAPI) GetByName(c *gin.Context) {
	param := storageClusterParam{}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	ctx, cancel := s.deps.WithTimeout()
	defer cancel()
	data, err := s.deps.Repo.Get(ctx, constants.GetStorageClusterConfigPath(param.ClusterName))
	if err != nil {
		http.Error(c, err)
		return
	}
	storageCluster := &config.StorageCluster{}
	err = encoding.JSONUnmarshal(data, storageCluster)
	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, storageCluster)
}

// DeleteByName deletes storage cluster by name
func (s *StorageClusterAPI) DeleteByName(c *gin.Context) {
	param := storageClusterParam{}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	ctx, cancel := s.deps.WithTimeout()
	defer cancel()
	if err = s.deps.Repo.Delete(ctx, constants.GetStorageClusterConfigPath(param.ClusterName)); err != nil {
		http.Error(c, err)
		return
	}
	http.NoContent(c)
}

// List lists all storage clusters
func (s *StorageClusterAPI) List(c *gin.Context) {
	ctx, cancel := s.deps.WithTimeout()
	defer cancel()
	data, err := s.deps.Repo.List(ctx, constants.StorageConfigPath)
	if err != nil {
		http.Error(c, err)
		return
	}
	stateMgr := s.deps.StateMgr
	var storages []models.Storage
	for _, val := range data {
		storage := models.Storage{}
		err = encoding.JSONUnmarshal(val.Value, &storage)
		if err != nil {
			s.logger.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			_, ok := stateMgr.GetStorage(storage.Name)
			if ok {
				storage.Status = models.StorageStatusReady
			} else {
				storage.Status = models.StorageStatusInitialize
				//TODO check storage un-health
			}
			storages = append(storages, storage)
		}
	}

	if err != nil {
		http.Error(c, err)
		return
	}
	http.OK(c, storages)
}
