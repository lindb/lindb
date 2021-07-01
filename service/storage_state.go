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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_state.go -destination=./storage_state_mock.go -package service

// StorageStateService represents storage cluster state maintain
type StorageStateService interface {
	// Save saves newest storage state for cluster name
	Save(clusterName string, storageState *models.StorageState) error
	// Get gets current storage state for given cluster name, if not exist return ErrNotExist
	Get(clusterName string) (*models.StorageState, error)
}

// storageStateService implements storage state service interface.
// broker need use storage state for write/read operation.
type storageStateService struct {
	ctx  context.Context
	repo state.Repository
}

// NewStorageStateService creates storage state service
func NewStorageStateService(ctx context.Context, repo state.Repository) StorageStateService {
	return &storageStateService{
		ctx:  ctx,
		repo: repo,
	}
}

// Save saves newest storage state for cluster name
func (s *storageStateService) Save(clusterName string, storageState *models.StorageState) error {
	data, _ := json.Marshal(storageState)
	if err := s.repo.Put(s.ctx, constants.GetStorageClusterNodeStatePath(clusterName), data); err != nil {
		return err
	}
	return nil
}

// Get gets current storage state for given cluster name, if not exist return ErrNotExist
func (s *storageStateService) Get(clusterName string) (*models.StorageState, error) {
	data, err := s.repo.Get(s.ctx, constants.GetStorageClusterNodeStatePath(clusterName))
	if err != nil {
		return nil, err
	}
	storageState := &models.StorageState{}
	err = json.Unmarshal(data, storageState)
	if err != nil {
		return nil, err
	}
	return storageState, err
}
