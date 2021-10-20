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

	"go.uber.org/atomic"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source=./channel_database.go -destination=./channel_database_mock.go -package=replica

// for testing
var (
	createChannel        = newChannel
	databaseChannelScope = linmetric.NewScope("lindb.replica.database")
	evictedCounterVec    = databaseChannelScope.NewCounterVec("metrics_out_of_time_range", "db")
)

// DatabaseChannel represents the database level replication channel
type DatabaseChannel interface {
	// Write writes the metric data into channel's buffer
	Write(ctx context.Context, brokerBatchRows *metric.BrokerBatchRows) error
	// CreateChannel creates the shard level replication channel by given shard id
	CreateChannel(numOfShard int32, shardID models.ShardID) (Channel, error)
	Stop()
}

type (
	shard2Channel map[models.ShardID]Channel
	shardChannels struct {
		value atomic.Value // readonly shard2Channel
		mu    sync.Mutex   // lock for modifying shard2Channel
	}
	databaseChannel struct {
		databaseCfg   models.Database
		ahead         *atomic.Int64
		behind        *atomic.Int64
		ctx           context.Context
		fct           rpc.ClientStreamFactory
		numOfShard    atomic.Int32
		shardChannels shardChannels
		interval      timeutil.Interval
		logger        *logger.Logger

		statistics struct {
			evictedCounter *linmetric.BoundCounter
		}
	}
)

// newDatabaseChannel creates a new database replication channel
func newDatabaseChannel(
	ctx context.Context,
	databaseCfg models.Database,
	numOfShard int32,
	fct rpc.ClientStreamFactory,
) (DatabaseChannel, error) {
	ch := &databaseChannel{
		databaseCfg: databaseCfg,
		ctx:         ctx,
		fct:         fct,
		logger:      logger.GetLogger("replica", "DatabaseChannel"),
	}
	ch.shardChannels.value.Store(make(shard2Channel))

	opt := databaseCfg.Option
	ahead, behind := (&opt).GetAcceptWritableRange()
	ch.ahead = atomic.NewInt64(ahead)
	ch.behind = atomic.NewInt64(behind)
	_ = ch.interval.ValueOf(databaseCfg.Option.Interval)

	ch.numOfShard.Store(numOfShard)
	ch.statistics.evictedCounter = evictedCounterVec.WithTagValues(databaseCfg.Name)
	return ch, nil
}

// Write writes the metric data into channel's buffer
func (dc *databaseChannel) Write(ctx context.Context, brokerBatchRows *metric.BrokerBatchRows) error {
	var err error

	behind := dc.behind.Load()
	ahead := dc.ahead.Load()

	evicted := brokerBatchRows.EvictOutOfTimeRange(behind, ahead)
	dc.statistics.evictedCounter.Add(float64(evicted))

	// sharding metrics to shards
	shardingIterator := brokerBatchRows.NewShardGroupIterator(dc.numOfShard.Load())
	for shardingIterator.HasRowsForNextShard() {
		shardIdx, familyIterator := shardingIterator.FamilyRowsForNextShard(dc.interval)
		shardID := models.ShardID(shardIdx)
		channel, ok := dc.getChannelByShardID(shardID)
		if !ok {
			err = errChannelNotFound
			// broker error, do not return to client
			dc.logger.Error("shardChannel not found",
				logger.String("database", dc.databaseCfg.Name),
				logger.Int("shardID", shardID.Int()))
			continue
		}
		for familyIterator.HasNextFamily() {
			familyTime, rows := familyIterator.NextFamily()
			familyChannel := channel.GetOrCreateFamilyChannel(familyTime)
			if err = familyChannel.Write(ctx, rows); err != nil {
				dc.logger.Error("failed writing rows to family channel",
					logger.String("database", dc.databaseCfg.Name),
					logger.Int("shardID", shardID.Int()),
					logger.Int("rows", len(rows)),
					logger.Int64("familyTime", familyTime),
					logger.Error(err))
			}
		}
	}
	//TODO if need return nil?
	return err
}

// CreateChannel creates the shard level replication channel by given shard id
func (dc *databaseChannel) CreateChannel(numOfShard int32, shardID models.ShardID) (Channel, error) {
	channel, ok := dc.getChannelByShardID(shardID)
	if !ok {
		dc.shardChannels.mu.Lock()
		defer dc.shardChannels.mu.Unlock()

		// double check
		channel, ok = dc.getChannelByShardID(shardID)
		if !ok {
			if numOfShard <= 0 || int32(shardID) >= numOfShard {
				return nil, errInvalidShardID
			}
			if numOfShard < dc.numOfShard.Load() {
				return nil, errInvalidShardNum
			}
			ch := createChannel(dc.ctx, dc.databaseCfg.Name, shardID, dc.fct)

			// cache shard level channel
			dc.insertShardChannel(shardID, ch)
			return ch, nil
		}
	}
	return channel, nil
}

func (dc *databaseChannel) Stop() {
	dc.shardChannels.mu.Lock()
	defer dc.shardChannels.mu.Unlock()

	channels := dc.shardChannels.value.Load().(shard2Channel)
	for _, channel := range channels {
		channel.Stop()
	}
}

// getChannelByShardID gets the replica channel by shard id
func (dc *databaseChannel) getChannelByShardID(shardID models.ShardID) (Channel, bool) {
	ch, ok := dc.shardChannels.value.Load().(shard2Channel)[shardID]
	return ch, ok
}

func (dc *databaseChannel) insertShardChannel(newShardID models.ShardID, newChannel Channel) {
	oldMap := dc.shardChannels.value.Load().(shard2Channel)
	newMap := make(shard2Channel)
	for shardID, channel := range oldMap {
		newMap[shardID] = channel
	}
	newMap[newShardID] = newChannel
	dc.shardChannels.value.Store(newMap)
}
