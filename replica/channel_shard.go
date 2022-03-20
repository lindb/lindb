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

package replica

import (
	"context"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./channel_shard.go -destination=./channel_shard_mock.go -package=replica

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
	SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)

	// GetOrCreateFamilyChannel musts picks the family channel by given family time.
	GetOrCreateFamilyChannel(familyTime int64) FamilyChannel

	Stop()

	garbageCollect(ahead, behind int64)
}

// channel implements Channel.
type channel struct {
	// context to close channel
	ctx context.Context
	cfg config.Write

	database string
	shardID  models.ShardID
	fct      rpc.ClientStreamFactory

	families   *familyChannelSet // send channel for each family time
	shardState models.ShardState
	liveNodes  map[models.NodeID]models.StatefulNode

	mutex sync.Mutex

	logger *logger.Logger
}

// newChannel returns a new channel with specific attribution.
func newChannel(
	ctx context.Context,
	database string,
	shardID models.ShardID,
	fct rpc.ClientStreamFactory,
) Channel {
	c := &channel{
		ctx:      ctx,
		cfg:      config.GlobalBrokerConfig().Write, // TODO
		database: database,
		shardID:  shardID,
		families: newFamilyChannelSet(),
		fct:      fct,
		logger:   logger.GetLogger("replica", "ShardChannel"),
	}

	// TODO need add family gc task
	return c
}

func (c *channel) SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.shardState.Leader != shardState.Leader {
		// leader change, need notify sender
		c.shardState = shardState
		c.liveNodes = liveNodes
		families := c.families.Entries()
		for _, family := range families {
			family.leaderChanged(c.shardState, c.liveNodes)
		}
		c.logger.Info("shard leader changed, need switch leader sender",
			logger.String("db", c.database),
			logger.Any("shardID", c.shardID))
	}
}

// GetOrCreateFamilyChannel returns family channel by given family time.
func (c *channel) GetOrCreateFamilyChannel(familyTime int64) FamilyChannel {
	familyChannel, exist := c.families.GetFamilyChannel(familyTime)
	if exist {
		return familyChannel
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// double check
	familyChannel, exist = c.families.GetFamilyChannel(familyTime)
	if exist {
		return familyChannel
	}
	familyChannel = newFamilyChannel(c.ctx, c.cfg, c.database, c.shardID, familyTime, c.fct, c.shardState, c.liveNodes)
	c.families.InsertFamily(familyTime, familyChannel)
	return familyChannel
}

func (c *channel) Stop() {
	families := c.families.Entries()
	for _, family := range families {
		family.Stop()
	}
}

func (c *channel) garbageCollect(ahead, behind int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	families := c.families.Entries()
	needRemovedFamilies := make(map[int64]struct{})
	for _, family := range families {
		if family.isExpire(ahead, behind) {
			c.logger.Info("family channel is expire, need stop it",
				logger.String("database", c.database),
				logger.Any("shard", c.shardID),
				logger.String("family", timeutil.FormatTimestamp(family.FamilyTime(), timeutil.DataTimeFormat4)))
			needRemovedFamilies[family.FamilyTime()] = struct{}{}
		}
	}
	removedFamilies := c.families.RemoveFamilies(needRemovedFamilies)
	// stop family after remove, just stop removed family.
	// maybe family will be used before remove.
	for _, family := range removedFamilies {
		family.Stop()
	}
}
