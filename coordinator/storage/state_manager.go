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

package storage

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/tsdb"
)

type StateManager interface {
	OnShardAssignmentChange(key string, data []byte)
	OnDatabaseDelete(key string)
}

type stateManager struct {
	engine  tsdb.Engine
	current *models.StatefulNode

	logger *logger.Logger
}

func NewStateManager(current *models.StatefulNode, engine tsdb.Engine) StateManager {
	return &stateManager{
		current: current,
		engine:  engine,
		logger:  logger.GetLogger("storage", "StateManager"),
	}
}

func (s *stateManager) OnShardAssignmentChange(key string, data []byte) {
	s.logger.Info("shard assignment is changed",
		logger.String("key", key),
		logger.String("data", string(data)))
	param := models.DatabaseAssignment{}
	if err := encoding.JSONUnmarshal(data, &param); err != nil {
		return
	}
	if param.ShardAssignment == nil {
		return
	}
	var shardIDs []models.ShardID
	for shardID, replica := range param.ShardAssignment.Shards {
		if replica.Contain(s.current.ID) {
			shardIDs = append(shardIDs, shardID)
		}
	}
	if len(shardIDs) == 0 {
		return
	}
	if err := s.engine.CreateShards(
		param.ShardAssignment.Name,
		param.Option,
		shardIDs...,
	); err != nil {
		return
	}
}

func (s *stateManager) OnDatabaseDelete(key string) {
	panic("implement me")
}
