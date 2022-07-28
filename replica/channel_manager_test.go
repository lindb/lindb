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
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
)

func TestChannelManager_GetChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Error(err)
		}
		ctrl.Finish()
	}()

	stateMgr := broker.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchShardStateChangeEvent(gomock.Any())
	cm := NewChannelManager(context.TODO(), nil, stateMgr)
	cm1 := cm.(*channelManager)

	_, err := cm1.CreateChannel(models.Database{Name: "database"}, 2, 2)
	assert.Error(t, err)

	opt := &option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}
	ch1, err := cm1.CreateChannel(models.Database{Name: "database",
		Option: opt,
	}, 3, 0)
	assert.NoError(t, err)

	ch111, err := cm1.CreateChannel(models.Database{Name: "database",
		Option: opt,
	}, 3, 0)
	assert.NoError(t, err)
	assert.Equal(t, ch111, ch1)

	cm.Close()
}

func TestChannelManager_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Error(err)
		}
		ctrl.Finish()
	}()

	stateMgr := broker.NewMockStateManager(ctrl)
	stateMgr.EXPECT().WatchShardStateChangeEvent(gomock.Any())
	cm := NewChannelManager(context.TODO(), nil, stateMgr)
	err := cm.Write(context.TODO(), "database", nil)
	assert.NoError(t, err)

	dbChannel := NewMockDatabaseChannel(ctrl)
	dbChannel.EXPECT().Stop()
	cm1 := cm.(*channelManager)
	cm1.insertDatabaseChannel("database", dbChannel)
	dbChannel.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	dbChannel.EXPECT().Stop().AnyTimes()
	err = cm.Write(context.TODO(), "database", nil)
	assert.NoError(t, err)

	rows := mockBrokerRows(t)

	err = cm.Write(context.TODO(), "database", rows)
	assert.NoError(t, err)
	err = cm.Write(context.TODO(), "database_not_exist", rows)
	assert.Error(t, err)

	cm1.insertDatabaseChannel("database2", dbChannel)
	cm1.insertDatabaseChannel("database3", dbChannel)
	cm.Close()
}

func TestChannelManager_handleShardStateChangeEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cm := &channelManager{
		logger: logger.GetLogger("Replica", "Test"),
	}
	cm.databaseChannels.value.Store(make(database2Channel))
	dbChannel := NewMockDatabaseChannel(ctrl)
	cm.insertDatabaseChannel("database", dbChannel)

	cases := []struct {
		name      string
		db        models.Database
		shards    map[models.ShardID]models.ShardState
		liveNodes map[models.NodeID]models.StatefulNode
		prepare   func()
	}{
		{
			name: "shard empty",
		},
		{
			name: "shard num wrong",
			db:   models.Database{Name: "database"},
			shards: map[models.ShardID]models.ShardState{
				3: {ID: 3},
			},
		},
		{
			name: "sync shard state successfully",
			db:   models.Database{Name: "database"},
			shards: map[models.ShardID]models.ShardState{
				0: {ID: 0},
			},
			prepare: func() {
				shardCh := NewMockShardChannel(ctrl)
				dbChannel.EXPECT().CreateChannel(gomock.Any(), gomock.Any()).Return(shardCh, nil)
				shardCh.EXPECT().SyncShardState(gomock.Any(), gomock.Any())
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			cm.handleShardStateChangeEvent(tt.db, tt.shards, tt.liveNodes)
		})
	}
}

func TestChannelManager_gcWriteFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	config.SetGlobalBrokerConfig(&config.BrokerBase{Write: config.Write{
		GCTaskInterval: ltoml.Duration(time.Millisecond * 10),
	}})

	ctx, cancel := context.WithCancel(context.TODO())
	cm := &channelManager{
		ctx:    ctx,
		logger: logger.GetLogger("Replica", "Test"),
	}
	cm.databaseChannels.value.Store(make(database2Channel))
	dbChannel := NewMockDatabaseChannel(ctrl)
	cm.insertDatabaseChannel("database", dbChannel)
	dbChannel.EXPECT().garbageCollect().AnyTimes()

	cm.gcWriteFamilies()

	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)
}

func mockBrokerRows(t *testing.T) *metric.BrokerBatchRows {
	converter := metric.NewProtoConverter()
	var brokerRow metric.BrokerRow
	assert.NoError(t, converter.ConvertTo(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}, &brokerRow))
	rows := metric.NewBrokerBatchRows()
	assert.NoError(t, rows.TryAppend(func(row *metric.BrokerRow) error {
		return nil
	}))
	return rows
}
