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
	"errors"
	"fmt"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./channel_manager.go -destination=./channel_manager_mock.go -package=replica

const (
	defaultBufferSize = 1024
)

var log = logger.GetLogger("replica", "ChannelManager")

// ChannelManager manages the construction, retrieving, closing for all channels.
type ChannelManager interface {
	// Write writes a MetricList, the manager handler the database, sharding things.
	Write(database string, list *protoMetricsV1.MetricList) error
	// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID,
	// numOfShard should be greater or equal than the origin setting, otherwise error is returned.
	// numOfShard is used eot calculate the shardID for a given hash.
	CreateChannel(databaseCfg models.Database, numOfShard int32, shardID models.ShardID) (Channel, error)

	// Close closes all the channel.
	Close()
}

// channelManager implements ChannelManager.
type channelManager struct {
	// context passed to all Channel
	ctx context.Context
	// cancelFun to cancel context
	cancel context.CancelFunc
	// factory to get rpc  writeTask client
	fct rpc.ClientStreamFactory
	// channelID(database name)  -> Channel
	databaseChannelMap sync.Map
	// lock for channelMap
	lock4map sync.Mutex

	logger *logger.Logger
}

// NewChannelManager returns a ChannelManager with dirPath and WriteClientFactory.
// WriteClientFactory makes it easy to mock rpc streamClient for test.
func NewChannelManager(ctx context.Context, fct rpc.ClientStreamFactory) ChannelManager {
	ctx, cancel := context.WithCancel(ctx)
	cm := &channelManager{
		ctx:    ctx,
		cancel: cancel,
		fct:    fct,
		logger: logger.GetLogger("replica", "ChannelManager"),
	}
	return cm
}

// Write writes a MetricList, the manager handler the database, sharding things.
func (cm *channelManager) Write(database string, metricList *protoMetricsV1.MetricList) error {
	if metricList == nil || len(metricList.Metrics) == 0 {
		return nil
	}
	databaseChannel, ok := cm.getDatabaseChannel(database)
	if !ok {
		return fmt.Errorf("database [%s] not found", database)
	}
	return databaseChannel.Write(metricList)
}

// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID.
// NumOfShard should be greater or equal than the origin setting, otherwise error is returned.
func (cm *channelManager) CreateChannel(databaseCfg models.Database, numOfShard int32, shardID models.ShardID) (Channel, error) {
	if numOfShard <= 0 || int32(shardID) >= numOfShard {
		return nil, errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	}
	database := databaseCfg.Name
	ch, ok := cm.getDatabaseChannel(database)
	if !ok {
		// double check, need lock
		cm.lock4map.Lock()
		defer cm.lock4map.Unlock()

		ch, ok = cm.getDatabaseChannel(database)
		if !ok {
			// if not exist, create database channel
			ch, err := newDatabaseChannel(cm.ctx, databaseCfg, numOfShard, cm.fct)
			if err != nil {
				return nil, err
			}
			// add to cache
			cm.databaseChannelMap.Store(database, ch)

			cm.logger.Info("create shard write channel successfully",
				logger.String("db", database), logger.Any("shardID", shardID))

			// create shard level channel
			return ch.CreateChannel(numOfShard, shardID)
		}
	}
	return ch.CreateChannel(numOfShard, shardID)
}

// Close closes all the channel.
func (cm *channelManager) Close() {
	cm.cancel()
	cm.databaseChannelMap.Range(func(key, ch interface{}) bool {
		channel, ok := ch.(DatabaseChannel)
		if ok {
			channel.Stop()
		}
		return true
	})
}

// getDatabaseChannel gets the database channel by given database name
func (cm *channelManager) getDatabaseChannel(databaseName string) (DatabaseChannel, bool) {
	ch, ok := cm.databaseChannelMap.Load(databaseName)
	if !ok {
		return nil, ok
	}
	channel, ok := ch.(DatabaseChannel)
	if !ok {
		return nil, ok
	}
	return channel, true
}
