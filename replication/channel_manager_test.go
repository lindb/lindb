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
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

var replicationConfig = config.ReplicationChannel{
	Dir:                "/tmp/broker/replication",
	DataSizeLimit:      int64(128),
	RemoveTaskInterval: ltoml.Duration(time.Minute),
	ReportInterval:     ltoml.Duration(time.Second),
	FlushInterval:      ltoml.Duration(0),
	CheckFlushInterval: ltoml.Duration(100 * time.Millisecond),
	BufferSize:         2,
}

func TestChannelManager_GetChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Error(err)
		}
		ctrl.Finish()
	}()

	replicatorStateReport := NewMockReplicatorStateReport(ctrl)
	replicatorStateReport.EXPECT().Report(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()

	replicationConfig.Dir = dirPath
	cm := NewChannelManager(replicationConfig, nil, replicatorStateReport)

	_, err := cm.CreateChannel("database", 2, 2)
	assert.Error(t, err)

	ch1, err := cm.CreateChannel("database", 3, 0)
	assert.NoError(t, err)

	ch111, err := cm.CreateChannel("database", 3, 0)
	assert.NoError(t, err)
	assert.Equal(t, ch111, ch1)

	defer func() {
		mkdir = fileutil.MkDirIfNotExist
	}()
	mkdir = func(path string) error {
		return fmt.Errorf("err")
	}
	_, err = cm.CreateChannel("database-err", 3, 1)
	assert.Error(t, err)

	cm1 := cm.(*channelManager)
	cm1.databaseChannelMap.Store("database-value-err", "test")
	c, ok := cm1.getDatabaseChannel("database-value-err")
	assert.False(t, ok)
	assert.Nil(t, c)

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

	replicatorStateReport := NewMockReplicatorStateReport(ctrl)
	replicatorStateReport.EXPECT().Report(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()

	replicationConfig.Dir = dirPath
	cm := NewChannelManager(replicationConfig, nil, replicatorStateReport)
	err := cm.Write("database", nil)
	assert.Error(t, err)

	dbChannel := NewMockDatabaseChannel(ctrl)
	cm1 := cm.(*channelManager)
	cm1.databaseChannelMap.Store("database", dbChannel)
	dbChannel.EXPECT().Write(gomock.Any()).Return(nil)
	err = cm.Write("database", &pb.MetricList{Metrics: []*pb.Metric{
		{Namespace: "xx"},
	}})
	assert.NoError(t, err)
	cm.Close()
}

func TestChannelManager_ReportState(t *testing.T) {
	ctrl := gomock.NewController(t)
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Error(err)
		}
		ctrl.Finish()
	}()

	replicatorStateReport := NewMockReplicatorStateReport(ctrl)
	replicatorStateReport.EXPECT().Report(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()

	replicationConfig.Dir = dirPath
	cm := NewChannelManager(replicationConfig, nil, replicatorStateReport)
	time.Sleep(2 * time.Second)
	cm.Close()
	// waiting close complete
	time.Sleep(400 * time.Millisecond)

	dbChannel := NewMockDatabaseChannel(ctrl)
	cm1 := cm.(*channelManager)
	cm1.databaseChannelMap.Store("database", dbChannel)
	dbChannel.EXPECT().ReplicaState().Return([]models.ReplicaState{{}}).AnyTimes()
	cm1.reportState()
}

func TestChannelManager_SyncReplicatorState(t *testing.T) {
	ctrl := gomock.NewController(t)
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Error(err)
		}
		ctrl.Finish()
	}()

	replicatorStateReport := NewMockReplicatorStateReport(ctrl)
	replicatorStateReport.EXPECT().Report(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()

	replicationConfig.Dir = dirPath
	cm := NewChannelManager(replicationConfig, nil, replicatorStateReport)
	cm.SyncReplicatorState()

	dbChannel := NewMockDatabaseChannel(ctrl)
	cm1 := cm.(*channelManager)
	cm1.databaseChannelMap.Store("database", dbChannel)
	dbChannel.EXPECT().ReplicaState().Return([]models.ReplicaState{{}}).AnyTimes()
	cm1.reportState()
	cm.Close()
}
