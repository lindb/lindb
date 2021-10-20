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
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWriteAheadLogManager_GetOrCreateLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newWriteAheadLog = NewWriteAheadLog
		ctrl.Finish()
	}()

	newWriteAheadLog = func(_ context.Context, cfg config.WAL,
		currentNodeID models.NodeID, database string,
		engine tsdb.Engine,
		cliFct rpc.ClientStreamFactory,
		_ storage.StateManager,
	) WriteAheadLog {
		return NewMockWriteAheadLog(ctrl)
	}
	m := NewWriteAheadLogManager(context.TODO(), config.WAL{RemoveTaskInterval: ltoml.Duration(time.Minute)},
		1, nil, nil, nil)
	// create new
	l := m.GetOrCreateLog("test")
	assert.NotNil(t, l)
	// return exist
	l = m.GetOrCreateLog("test")
	assert.NotNil(t, l)
}

func TestWriteAheadLog_GetOrCreatePartition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newFanOutQueue = queue.NewFanOutQueue
		ctrl.Finish()
	}()
	engine := tsdb.NewMockEngine(ctrl)
	l := NewWriteAheadLog(context.TODO(), config.WAL{RemoveTaskInterval: ltoml.Duration(time.Minute)},
		1, "test", engine, nil, nil)

	// case 1: shard not exist
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, false)
	p, err := l.GetOrCreatePartition(1, 1, 1)
	assert.Error(t, err)
	assert.Nil(t, p)
	// case 2: new log err
	newFanOutQueue = func(dirPath string, dataSizeLimit int64,
		removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, fmt.Errorf("err")
	}
	shard := tsdb.NewMockShard(ctrl)
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, nil)
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.Error(t, err)
	assert.Nil(t, p)
	// case 3: create log ok
	newFanOutQueue = func(dirPath string, dataSizeLimit int64,
		removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, nil
	}
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
	shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, nil)
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	// case 4: return exist one
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p)
}
