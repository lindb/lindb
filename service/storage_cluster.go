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

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_cluster.go -destination=./storage_service_mock.go -package service

// StorageClusterService defines storage cluster service interface
type StorageClusterService interface {
	// Save saves storage cluster config
	Save(storageCluster *config.StorageCluster) error
	// Delete deletes storage cluster config
	Delete(name string) error
	// Get storage cluster by given name, if not exist return ErrNotExist
	Get(name string) (*config.StorageCluster, error)
	// List lists all storage cluster config
	List() ([]*config.StorageCluster, error)
}

// storageClusterService implements storage cluster service interface
type storageClusterService struct {
	ctx  context.Context
	repo state.Repository
}

// NewStorageClusterService creates storage cluster service
func NewStorageClusterService(ctx context.Context, repo state.Repository) StorageClusterService {
	return &storageClusterService{
		ctx:  ctx,
		repo: repo,
	}
}

// Save saves storage cluster config
func (s *storageClusterService) Save(storageCluster *config.StorageCluster) error {
	if storageCluster.Name == "" {
		return fmt.Errorf("storage cluster name cannot be empty")
	}
	data, _ := json.Marshal(storageCluster)

	if err := s.repo.Put(s.ctx, constants.GetStorageClusterConfigPath(storageCluster.Name), data); err != nil {
		return err
	}
	return nil
}

// Delete deletes storage cluster config
func (s *storageClusterService) Delete(name string) error {
	return s.repo.Delete(s.ctx, constants.GetStorageClusterConfigPath(name))
}

// Get storage cluster by given name
func (s *storageClusterService) Get(name string) (*config.StorageCluster, error) {
	data, err := s.repo.Get(s.ctx, constants.GetStorageClusterConfigPath(name))
	if err != nil {
		return nil, err
	}
	storageCluster := &config.StorageCluster{}
	err = json.Unmarshal(data, storageCluster)
	if err != nil {
		return nil, err
	}
	return storageCluster, err
}

// List lists config of all storage clusters
func (s *storageClusterService) List() ([]*config.StorageCluster, error) {
	var result []*config.StorageCluster
	data, err := s.repo.List(s.ctx, constants.StorageClusterConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		storageCluster := &config.StorageCluster{}
		err = json.Unmarshal(val.Value, storageCluster)
		if err != nil {
			logger.GetLogger("service", "StorageCluster").
				Warn("unmarshal data error",
					logger.String("data", string(val.Value)))
		} else {
			result = append(result, storageCluster)
		}
	}
	return result, nil
}
