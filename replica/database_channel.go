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

	"github.com/lithammer/go-jump-consistent-hash"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source=./database_channel.go -destination=./database_channel_mock.go -package=replica

// for testing
var (
	createChannel = newChannel
)

// DatabaseChannel represents the database level replication channel
type DatabaseChannel interface {
	// Write writes the metric data into channel's buffer
	Write(metricList *protoMetricsV1.MetricList) error
	// CreateChannel creates the shard level replication channel by given shard id
	CreateChannel(numOfShard int32, shardID models.ShardID) (Channel, error)
	Stop()
}

type databaseChannel struct {
	databaseCfg   models.Database
	ahead         *atomic.Int64
	behind        *atomic.Int64
	ctx           context.Context
	fct           rpc.ClientStreamFactory
	numOfShard    atomic.Int32
	shardChannels sync.Map
	mutex         sync.Mutex
}

// newDatabaseChannel creates a new database replication channel
func newDatabaseChannel(ctx context.Context,
	databaseCfg models.Database, numOfShard int32,
	fct rpc.ClientStreamFactory,
) (DatabaseChannel, error) {
	ch := &databaseChannel{
		databaseCfg: databaseCfg,
		ctx:         ctx,
		fct:         fct,
	}

	var ahead timeutil.Interval
	var behind timeutil.Interval
	_ = ahead.ValueOf(databaseCfg.Option.Ahead)
	_ = behind.ValueOf(databaseCfg.Option.Behind)
	ch.ahead = atomic.NewInt64(ahead.Int64())
	ch.behind = atomic.NewInt64(behind.Int64())
	if ch.ahead.Load() <= 0 {
		ch.ahead.Store(constants.MetricMaxBehindDuration)
	}
	if ch.behind.Load() <= 0 {
		ch.behind.Store(constants.MetricMaxAheadDuration)
	}

	ch.numOfShard.Store(numOfShard)
	return ch, nil
}

// Write writes the metric data into channel's buffer
func (dc *databaseChannel) Write(metricList *protoMetricsV1.MetricList) (err error) {
	now := fasttime.UnixMilliseconds()
	behind := dc.behind.Load()
	ahead := dc.ahead.Load()

	// sharding metrics to shards
	numOfShard := dc.numOfShard.Load()
	for _, metric := range metricList.Metrics {
		timestamp := metric.Timestamp

		// check metric timestamp if in acceptable time range
		if (behind > 0 && timestamp < now-behind) ||
			(ahead > 0 && timestamp > now+ahead) {
			//TODO need add metric
			continue
		}
		hash := tag.XXHashOfKeyValues(metric.Tags)

		idx := int(jump.Hash(hash, numOfShard))
		// set tags hash code for storage side reuse
		// !!!IMPORTANT: storage side will use this hash for writeTask
		metric.TagsHash = hash
		shardID := models.ShardID(idx)
		channel, ok := dc.getChannelByShardID(shardID)
		if !ok {
			err = errChannelNotFound
			// broker error, do not return to client
			log.Error("channel not found",
				logger.String("database", dc.databaseCfg.Name),
				logger.Any("shardID", shardID))
			continue
		}
		if err = channel.Write(metric); err != nil {
			log.Error("channel writeTask data error",
				logger.String("database", dc.databaseCfg.Name),
				logger.Any("shardID", shardID))
		}
	}
	//TODO if need return nil?
	return err
}

// CreateChannel creates the shard level replication channel by given shard id
func (dc *databaseChannel) CreateChannel(numOfShard int32, shardID models.ShardID) (Channel, error) {
	channel, ok := dc.getChannelByShardID(shardID)
	if !ok {
		dc.mutex.Lock()
		defer dc.mutex.Unlock()

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
			dc.shardChannels.Store(shardID, ch)
			return ch, nil
		}
	}
	return channel, nil
}

func (dc *databaseChannel) Stop() {
	dc.shardChannels.Range(func(key, channel interface{}) bool {
		ch, ok := channel.(Channel)
		if ok {
			ch.Stop()
		}
		return true
	})
}

// getChannelByShardID gets the replica channel by shard id
func (dc *databaseChannel) getChannelByShardID(shardID models.ShardID) (Channel, bool) {
	channel, ok := dc.shardChannels.Load(shardID)
	if !ok {
		return nil, ok
	}
	ch, ok := channel.(Channel)
	if !ok {
		return nil, ok
	}
	return ch, true
}
