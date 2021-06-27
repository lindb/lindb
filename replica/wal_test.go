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
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWriteAheadLogManager_GetOrCreateLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newWriteAheadLog = NewWriteAheadLog
		ctrl.Finish()
	}()

	newWriteAheadLog = func(cfg config.Replica,
		currentNodeID models.NodeID, database string,
		storageSrv service.StorageService,
		cliFct rpc.ClientStreamFactory,
	) WriteAheadLog {
		return NewMockWriteAheadLog(ctrl)
	}
	m := NewWriteAheadLogManager(config.Replica{}, 1, nil, nil)
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
	srv := service.NewMockStorageService(ctrl)
	l := NewWriteAheadLog(config.Replica{}, 1, "test", srv, nil)

	// case 1: shard not exist
	srv.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, false)
	p, err := l.GetOrCreatePartition(1)
	assert.Error(t, err)
	assert.Nil(t, p)
	// case 2: new log err
	newFanOutQueue = func(dirPath string, dataSizeLimit int64,
		removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, fmt.Errorf("err")
	}
	srv.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, true)
	p, err = l.GetOrCreatePartition(1)
	assert.Error(t, err)
	assert.Nil(t, p)
	// case 3: create log ok
	newFanOutQueue = func(dirPath string, dataSizeLimit int64,
		removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, nil
	}
	srv.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, true)
	p, err = l.GetOrCreatePartition(1)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	// case 4: return exist one
	p, err = l.GetOrCreatePartition(1)
	assert.NoError(t, err)
	assert.NotNil(t, p)
}
