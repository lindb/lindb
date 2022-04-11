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
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWriteAheadLogManager_GetOrCreateLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newWriteAheadLog = NewWriteAheadLog
		getLogFn = getLog
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
	// other db
	l = m.GetOrCreateLog("test-2")
	assert.NotNil(t, l)
	// other db
	getLogFn = func(w *writeAheadLogManager, database string) (WriteAheadLog, bool) {
		return NewMockWriteAheadLog(ctrl), true
	}
	l = m.GetOrCreateLog("test-3")
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
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard.EXPECT().Database().Return(db).AnyTimes()
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, nil)
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.Error(t, err)
	assert.Nil(t, p)
	// case 3: create log ok
	log := queue.NewMockFanOutQueue(ctrl)
	newFanOutQueue = func(dirPath string, dataSizeLimit int64,
		removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return log, nil
	}
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().FamilyTime().Return(timeutil.Now()).AnyTimes()
	shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(family, nil)
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	// case 4: return exist one
	p, err = l.GetOrCreatePartition(1, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	// case 5: get replica state
	log.EXPECT().FanOutNames().Return([]string{"1"})
	fan := queue.NewMockFanOut(ctrl)
	log.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(fan, nil)
	fan.EXPECT().HeadSeq().Return(int64(1))
	fan.EXPECT().TailSeq().Return(int64(1))
	fan.EXPECT().Pending().Return(int64(1))
	log.EXPECT().HeadSeq().Return(int64(1))
	rs := l.getReplicaState()
	assert.Len(t, rs, 1)
}

func TestWAL_garbageCollectTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.TODO())
	wal := &writeAheadLogManager{
		ctx: ctx,
		cfg: config.WAL{RemoveTaskInterval: ltoml.Duration(time.Millisecond * 10)},
	}
	dbs := make(databaseLogs)
	log := NewMockWriteAheadLog(ctrl)
	log.EXPECT().destroy().AnyTimes()
	dbs["test"] = log
	wal.databaseLogs.Store(dbs)
	wal.garbageCollectTask()

	time.Sleep(time.Millisecond * 50)
	cancel()
	time.Sleep(time.Millisecond * 50)
}

func TestWriteAheadLogManager_GetReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		getLogFn = getLog
		ctrl.Finish()
	}()
	mgr := &writeAheadLogManager{}
	log := NewMockWriteAheadLog(ctrl)
	getLogFn = func(w *writeAheadLogManager, database string) (WriteAheadLog, bool) {
		return log, true
	}
	log.EXPECT().getReplicaState().Return([]models.FamilyLogReplicaState{{}})
	s := mgr.GetReplicaState("test")
	assert.Len(t, s, 1)

	getLogFn = func(w *writeAheadLogManager, database string) (WriteAheadLog, bool) {
		return nil, false
	}
	s = mgr.GetReplicaState("test")
	assert.Nil(t, s)
}

func TestWriteAheadLogManager_Recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		fileExistFn = fileutil.Exist
		listDirFn = fileutil.ListDir
		ctrl.Finish()
	}()
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "wal path not exist",
			prepare: func() {
				fileExistFn = func(file string) bool {
					return false
				}
			},
			wantErr: false,
		},
		{
			name: "list wal path failure",
			prepare: func() {
				listDirFn = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "recovery write ahead log failure",
			prepare: func() {
				log := NewMockWriteAheadLog(ctrl)
				getLogFn = func(w *writeAheadLogManager, database string) (WriteAheadLog, bool) {
					return log, true
				}
				log.EXPECT().recovery().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "recovery write ahead log successfully",
			prepare: func() {
				log := NewMockWriteAheadLog(ctrl)
				getLogFn = func(w *writeAheadLogManager, database string) (WriteAheadLog, bool) {
					return log, true
				}
				log.EXPECT().recovery().Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				getLogFn = getLog
				fileExistFn = func(file string) bool {
					return true
				}
				listDirFn = func(path string) ([]string, error) {
					return []string{"test"}, nil
				}
			}()
			mgr := &writeAheadLogManager{}
			mgr.databaseLogs.Store(make(databaseLogs))
			if tt.prepare != nil {
				tt.prepare()
			}
			err := mgr.Recovery()
			if (err != nil) != tt.wantErr {
				t.Errorf("Recovery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
