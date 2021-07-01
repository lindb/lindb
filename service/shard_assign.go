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
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./shard_assign.go -destination=./shard_assign_mock.go -package service

var (
	shardAssignServiceLogger = logger.GetLogger("service", "ShardAssignService")
)

// ShardAssignService represents database shard assignment maintain
// Master generates assignment, then storing into related storage cluster's state repo.
// Storage node will create the time series engine based on related shard assignment.
type ShardAssignService interface {
	// List returns all database's shard assignments
	List() ([]*models.ShardAssignment, error)
	// Get gets shard assignment by given database name, if not exist return ErrNotExist
	Get(databaseName string) (*models.ShardAssignment, error)
	// Save saves shard assignment for given database name, if fail return error
	Save(databaseName string, shardAssign *models.ShardAssignment) error
}

// shardAssignService implements shard assign service interface
type shardAssignService struct {
	ctx  context.Context
	repo state.Repository
}

// NewShardAssignService creates shard assign service
func NewShardAssignService(ctx context.Context, repo state.Repository) ShardAssignService {
	return &shardAssignService{
		ctx:  ctx,
		repo: repo,
	}
}

// List returns all database's shard assignments
func (s *shardAssignService) List() ([]*models.ShardAssignment, error) {
	data, err := s.repo.List(s.ctx, constants.DatabaseAssignPath)
	if err != nil {
		return nil, err
	}

	var result []*models.ShardAssignment
	for _, val := range data {
		shardAssign := &models.ShardAssignment{}
		err = encoding.JSONUnmarshal(val.Value, shardAssign)
		if err != nil {
			shardAssignServiceLogger.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			result = append(result, shardAssign)
		}
	}
	return result, nil
}

// Get gets shard assignment by given database name, if not exist return ErrNotExist
func (s *shardAssignService) Get(databaseName string) (*models.ShardAssignment, error) {
	data, err := s.repo.Get(s.ctx, constants.GetDatabaseAssignPath(databaseName))
	if err != nil {
		return nil, err
	}
	shardAssign := &models.ShardAssignment{}
	if err := json.Unmarshal(data, shardAssign); err != nil {
		return nil, err
	}
	return shardAssign, nil
}

// Save saves shard assignment for given database name, if fail return error
func (s *shardAssignService) Save(databaseName string, shardAssign *models.ShardAssignment) error {
	data, _ := json.Marshal(shardAssign)
	return s.repo.Put(s.ctx, constants.GetDatabaseAssignPath(databaseName), data)
}
