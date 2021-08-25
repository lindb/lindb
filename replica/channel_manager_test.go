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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
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

	cm := NewChannelManager(context.TODO(), nil)

	_, err := cm.CreateChannel("database", 2, 2)
	assert.Error(t, err)

	ch1, err := cm.CreateChannel("database", 3, 0)
	assert.NoError(t, err)

	ch111, err := cm.CreateChannel("database", 3, 0)
	assert.NoError(t, err)
	assert.Equal(t, ch111, ch1)

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

	cm := NewChannelManager(context.TODO(), nil)
	err := cm.Write("database", nil)
	assert.Error(t, err)

	dbChannel := NewMockDatabaseChannel(ctrl)
	dbChannel.EXPECT().Stop()
	cm1 := cm.(*channelManager)
	cm1.databaseChannelMap.Store("database", dbChannel)
	dbChannel.EXPECT().Write(gomock.Any()).Return(nil)
	err = cm.Write("database", &protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
		{Namespace: "xx"},
	}})
	assert.NoError(t, err)
	cm.Close()
}
