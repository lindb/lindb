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

// for testing
var (
	getFamilyFn = getFamily
)

// ShardChannel represents a place to buffer the data for a specific cluster, database, shardID.
type ShardChannel interface {
	// SyncShardState syncs shard state after state event changed.
	SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
	// GetOrCreateFamilyChannel musts picks the family shardChannel by given family time.
	GetOrCreateFamilyChannel(familyTime int64) FamilyChannel
	// Stop stops shard shardChannel.
	Stop()

	// garbageCollect recycles expired write family.
	garbageCollect(ahead, behind int64)
}

// shardChannel implements ShardChannel.
type shardChannel struct {
	// context to close shardChannel
	ctx context.Context
	cfg config.Write

	database string
	shardID  models.ShardID
	fct      rpc.ClientStreamFactory

	families   *familyChannelSet // send shardChannel for each family time
	shardState models.ShardState
	liveNodes  map[models.NodeID]models.StatefulNode

	mutex sync.Mutex

	logger *logger.Logger
}

// newShardChannel returns a new shardChannel with specific attribution.
func newShardChannel(
	ctx context.Context,
	database string,
	shardID models.ShardID,
	fct rpc.ClientStreamFactory,
) ShardChannel {
	return &shardChannel{
		ctx:      ctx,
		cfg:      config.GlobalBrokerConfig().Write,
		database: database,
		shardID:  shardID,
		families: newFamilyChannelSet(),
		fct:      fct,
		logger:   logger.GetLogger("Replica", "ShardChannel"),
	}
}

func (c *shardChannel) SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
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

// GetOrCreateFamilyChannel returns family shardChannel by given family time.
func (c *shardChannel) GetOrCreateFamilyChannel(familyTime int64) FamilyChannel {
	familyChannel, exist := c.families.GetFamilyChannel(familyTime)
	if exist {
		return familyChannel
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// double check
	familyChannel, exist = getFamilyFn(c.families, familyTime)
	if exist {
		return familyChannel
	}
	familyChannel = newFamilyChannel(c.ctx, c.cfg, c.database, c.shardID, familyTime, c.fct, c.shardState, c.liveNodes)
	c.families.InsertFamily(familyTime, familyChannel)

	return familyChannel
}

// Stop stops shard shardChannel.
func (c *shardChannel) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	families := c.families.Entries()
	for _, family := range families {
		family.Stop(10 * timeutil.OneSecond)
	}
}

// garbageCollect recycles expired write family.
func (c *shardChannel) garbageCollect(ahead, behind int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	families := c.families.Entries()
	needRemovedFamilies := make(map[int64]struct{})
	for _, family := range families {
		if family.isExpire(ahead, behind) {
			c.logger.Info("family shardChannel is expire, need stop it",
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
		family.Stop(10 * timeutil.OneSecond)
	}
}

// getFamily returns family channel by family time.
func getFamily(families *familyChannelSet, familyTime int64) (FamilyChannel, bool) {
	return families.GetFamilyChannel(familyTime)
}
