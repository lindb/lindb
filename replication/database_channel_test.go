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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/rpc"
)

func TestDatabaseChannel_new(t *testing.T) {
	defer func() {
		mkdir = fileutil.MkDirIfNotExist
	}()
	mkdir = func(path string) error {
		return fmt.Errorf("err")
	}
	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 10, nil)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func TestDatabaseChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 1, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	err = ch.Write(&protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
		{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		},
	}})
	assert.Equal(t, errChannelNotFound, err)

	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(models.ShardID(0), shardCh)

	shardCh.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
	err = ch.Write(&protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
		{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		},
	}})
	assert.Error(t, err)
}

func TestDatabaseChannel_CreateChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 4, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(models.ShardID(0), shardCh)
	shardCh2, err := ch.CreateChannel(int32(1), models.ShardID(0))
	assert.NoError(t, err)
	assert.Equal(t, shardCh, shardCh2)

	_, err = ch.CreateChannel(0, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(1, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(2, 1)
	assert.Equal(t, errInvalidShardNum, err)

	_, err = ch.CreateChannel(4, 1)
	assert.NoError(t, err)

	defer func() {
		createChannel = newChannel
	}()
	createChannel = func(cxt context.Context,
		cfg config.ReplicationChannel, database string, shardID models.ShardID,
		fct rpc.ClientStreamFactory,
	) (i Channel, e error) {
		return nil, fmt.Errorf("err")
	}

	_, err = ch.CreateChannel(4, 2)
	assert.Error(t, err)

	ch1.shardChannels.Store(int32(3), "test")
	c, ok := ch1.getChannelByShardID(models.ShardID(3))
	assert.False(t, ok)
	assert.Nil(t, c)
}

func TestDatabaseChannel_ReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 3, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)

	shardCh0 := NewMockChannel(ctrl)
	shardCh1 := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(models.ShardID(0), shardCh0)
	ch1.shardChannels.Store(models.ShardID(1), shardCh1)

	shardCh0.EXPECT().Targets().Return(nil)
	shardCh1.EXPECT().Targets().Return([]models.Node{
		&models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 12345},
		&models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 12345},
	})
	shardCh1.EXPECT().GetOrCreateReplicator(&models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 12345}).Return(nil, fmt.Errorf("err"))
	replicator := NewMockReplicator(ctrl)
	shardCh1.EXPECT().GetOrCreateReplicator(&models.StatelessNode{HostIP: "2.2.2.2", GRPCPort: 12345}).Return(replicator, nil)
	replicator.EXPECT().Database().Return("db")
	replicator.EXPECT().ShardID().Return(models.ShardID(1))
	replicator.EXPECT().Pending().Return(int64(0))
	replicator.EXPECT().ReplicaIndex().Return(int64(0))
	replicator.EXPECT().AckIndex().Return(int64(0))

	replicaState := ch.ReplicaState()
	assert.Len(t, replicaState, 1)
}
