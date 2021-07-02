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

package replication

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

//go:generate mockgen -source=./channel_manager.go -destination=./channel_manager_mock.go -package=replication

// ErrCanceled is the error returned when writing data ctx canceled.
var ErrCanceled = errors.New("write data ctx done")

const (
	defaultReportInterval = 30 * time.Second
	defaultBufferSize     = 1024
)

var log = logger.GetLogger("replication", "ChannelManager")

// ChannelManager manages the construction, retrieving, closing for all channels.
type ChannelManager interface {
	// Write writes a MetricList, the manager handler the database, sharding things.
	Write(database string, list *pb.MetricList) error
	// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID,
	// numOfShard should be greater or equal than the origin setting, otherwise error is returned.
	// numOfShard is used eot calculate the shardID for a given hash.
	CreateChannel(database string, numOfShard, shardID int32) (Channel, error)
	// SyncReplicatorState syncs replicator state
	SyncReplicatorState()

	// Close closes all the channel.
	Close()
}

// channelManager implements ChannelManager.
type channelManager struct {
	// context passed to all Channel
	ctx context.Context
	// cancelFun to cancel context
	cancel context.CancelFunc
	// config
	cfg config.ReplicationChannel
	// factory to get rpc  write client
	fct rpc.ClientStreamFactory
	// for report replica state
	replicatorStateReport ReplicatorStateReport
	// channelID(database name)  -> Channel
	databaseChannelMap sync.Map
	// lock for channelMap
	lock4map  sync.Mutex
	syncState chan struct{}
	logger    *logger.Logger
}

// NewChannelManager returns a ChannelManager with dirPath and WriteClientFactory.
// WriteClientFactory makes it easy to mock rpc streamClient for test.
func NewChannelManager(cfg config.ReplicationChannel, fct rpc.ClientStreamFactory,
	replicatorStateReport ReplicatorStateReport) ChannelManager {
	ctx, cancel := context.WithCancel(context.Background())
	cm := &channelManager{
		ctx:                   ctx,
		cancel:                cancel,
		cfg:                   cfg,
		fct:                   fct,
		replicatorStateReport: replicatorStateReport,
		syncState:             make(chan struct{}),
		logger:                logger.GetLogger("replication", "channelManager"),
	}
	cm.scheduleStateReport()
	return cm
}

// Write writes a MetricList, the manager handler the database, sharding things.
func (cm *channelManager) Write(database string, metricList *pb.MetricList) error {
	databaseChannel, ok := cm.getDatabaseChannel(database)
	if !ok {
		return fmt.Errorf("database [%s] not found", database)
	}
	if metricList == nil || len(metricList.Metrics) == 0 {
		return fmt.Errorf("metrics is empty")
	}
	return databaseChannel.Write(metricList)
}

// CreateChannel creates a new channel or returns a existed channel for storage with specific database and shardID.
// NumOfShard should be greater or equal than the origin setting, otherwise error is returned.
func (cm *channelManager) CreateChannel(database string, numOfShard, shardID int32) (Channel, error) {
	if numOfShard <= 0 || shardID >= numOfShard {
		return nil, errors.New("numOfShard should be greater than 0 and shardID should less then numOfShard")
	}
	ch, ok := cm.getDatabaseChannel(database)
	if !ok {
		// double check, need lock
		cm.lock4map.Lock()
		defer cm.lock4map.Unlock()

		ch, ok = cm.getDatabaseChannel(database)
		if !ok {
			// if not exist, create database channel
			ch, err := newDatabaseChannel(cm.ctx, database, cm.cfg, numOfShard, cm.fct)
			if err != nil {
				return nil, err
			}
			// add to cache
			cm.databaseChannelMap.Store(database, ch)
			// create shard level channel
			return ch.CreateChannel(numOfShard, shardID)
		}
	}
	return ch.CreateChannel(numOfShard, shardID)
}

// SyncReplicatorState syncs replicator state
func (cm *channelManager) SyncReplicatorState() {
	cm.syncState <- struct{}{}
}

// Close closes all the channel.
func (cm *channelManager) Close() {
	cm.cancel()
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

// scheduleStateReport schedules a state report background job
func (cm *channelManager) scheduleStateReport() {
	interval := defaultReportInterval
	if cm.cfg.ReportInterval > 0 {
		interval = time.Duration(cm.cfg.ReportInterval)
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cm.reportState()
			case <-cm.syncState:
				cm.reportState()
			case <-cm.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// reportState reports the state of all replicators under current broker
func (cm *channelManager) reportState() {
	brokerState := models.BrokerReplicaState{
		ReportTime: timeutil.Now(),
	}
	cm.databaseChannelMap.Range(func(key, value interface{}) bool {
		channel, ok := value.(DatabaseChannel)
		if ok {
			replicas := channel.ReplicaState()
			if len(replicas) > 0 {
				brokerState.Replicas = append(brokerState.Replicas, replicas...)
			}
		}
		return true
	})
	if err := cm.replicatorStateReport.Report(&brokerState); err != nil {
		log.Error("report broker replicator state fail", logger.Error(err))
	}
}
