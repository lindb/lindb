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
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	StorageClusterPath = "/storage/cluster"
)

type storageClusterParam struct {
	ClusterName string `form:"name" binding:"required"`
}

// StorageClusterAPI represents storage cluster admin rest api
type StorageClusterAPI struct {
	deps   *depspkg.HTTPDeps
	logger *logger.Logger
}

// NewStorageClusterAPI create storage cluster api
func NewStorageClusterAPI(deps *depspkg.HTTPDeps) *StorageClusterAPI {
	return &StorageClusterAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "StorageClusterAPI"),
	}
}

// Register adds storage admin url route.
func (s *StorageClusterAPI) Register(route gin.IRoutes) {
	route.GET(StorageClusterPath, s.GetByName)
	route.DELETE(StorageClusterPath, s.DeleteByName)
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
