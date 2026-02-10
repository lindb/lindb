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

package write

import (
	"reflect"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
)

func (cm *channelManager) OnEvent(event discovery.MetaEvent) {
	switch e := event.(type) {
	case *models.ChangeShardStateEvent:
		cm.handleShardStateChangeEvent(e.DatabaseCfg, e.Shards, e.LiveNodes)
	// FIXME: add delete database logic
	default:
		cm.logger.Warn("unsupported shard channel manager event", logger.Any("eventType", reflect.TypeOf(event)))
	}
}

func (cm *channelManager) Subscribe(sub discovery.Subscriber) {
	panic("unimplemented")
}

func (cm *channelManager) Unsubscribe(sub discovery.Subscriber) {
	panic("unimplemented")
}

// handleShardStateChangeEvent handles shard state change event.
func (cm *channelManager) handleShardStateChangeEvent(
	databaseCfg models.Database,
	shards map[models.ShardID]models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) {
	numOfShard := len(shards)
	for _, shardState := range shards {
		shardID := shardState.ID
		ch, err := cm.CreateChannel(databaseCfg, int32(numOfShard), shardID)
		if err != nil {
			cm.logger.Error("create shard write shardChannel", logger.String("db", databaseCfg.Name),
				logger.Any("shard", shardID), logger.Error(err))
		} else {
			ch.SyncShardState(shardState, liveNodes)
		}
	}
}
