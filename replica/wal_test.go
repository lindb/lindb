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
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

func TestWriteAheadLog_GetOrCreatePartition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := tsdb.NewMockEngine(ctrl)
	shard := tsdb.NewMockShard(ctrl)

	cases := []struct {
		name    string
		prepare func(l *writeAheadLog)
		wantErr bool
	}{
		{
			name: "shard not exist",
			prepare: func(l *writeAheadLog) {
				engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, false)
			},
			wantErr: true,
		},
		{
			name: "get exist partition",
			prepare: func(l *writeAheadLog) {
				l.mutex.Lock()
				l.familyLogs[partitionKey{shardID: 1, familyTime: 1, leader: 1}] = NewMockPartition(ctrl)
				defer l.mutex.Unlock()
			},
			wantErr: true,
		},
		{
			name: "get data family failure",
			prepare: func(l *writeAheadLog) {
				engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
				shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "new log queue failure",
			prepare: func(l *writeAheadLog) {
				engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
				shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, nil)
				newFanOutQueue = func(dirPath string, dataSizeLimit int64) (q queue.FanOutQueue, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create partition successfully",
			prepare: func(l *writeAheadLog) {
				engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true)
				shard.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(nil, nil)
				newFanOutQueue = func(dirPath string, dataSizeLimit int64) (q queue.FanOutQueue, err error) {
					return nil, nil
				}
				NewPartitionFn = func(ctx context.Context, shard tsdb.Shard, family tsdb.DataFamily,
					currentNodeID models.NodeID, log queue.FanOutQueue,
					cliFct rpc.ClientStreamFactory, stateMgr storage.StateManager) Partition {
					return NewMockPartition(ctrl)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newFanOutQueue = queue.NewFanOutQueue
				NewPartitionFn = NewPartition
			}()
			l := NewWriteAheadLog(context.TODO(), config.WAL{RemoveTaskInterval: ltoml.Duration(time.Minute)},
				1, "test", engine, nil, nil)
			l1 := l.(*writeAheadLog)

			if tt.prepare != nil {
				tt.prepare(l1)
			}

			p, err := l.GetOrCreatePartition(1, 1, 1)
			if ((err != nil) != tt.wantErr && p == nil) || (!tt.wantErr && p == nil) {
				t.Errorf("GetOrCreatePartition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockWriteAheadLogManager_GetReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newFanOutQueue = queue.NewFanOutQueue
		ctrl.Finish()
	}()
	p1 := NewMockPartition(ctrl)
	key := partitionKey{
		shardID: 1,
	}
	wal := &writeAheadLog{
		database: "test",
		familyLogs: map[partitionKey]Partition{
			key: p1,
		},
	}
	p1.EXPECT().getReplicaState().Return(models.FamilyLogReplicaState{})
	state := wal.getReplicaState()
	assert.Len(t, state, 1)
	assert.Equal(t, "test", wal.Name())
}

func TestWriteAheadLog_recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := timeutil.Now() - timeutil.Now()%timeutil.OneHour
	engine := tsdb.NewMockEngine(ctrl)
	p := NewMockPartition(ctrl)
	cases := []struct {
		name    string
		prepare func(wal *writeAheadLog)
		wantErr bool
	}{
		{
			name: "list partition path failure",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "list shard path failure",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(p string) ([]string, error) {
					if p == wal.dir {
						return []string{"1"}, nil
					}
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "list family path failure",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(p string) ([]string, error) {
					if p == wal.dir {
						return []string{"1"}, nil
					}
					if p == path.Join(wal.dir, "1") {
						return []string{"1"}, nil
					}
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create partition failure",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(p string) ([]string, error) {
					if p == wal.dir {
						return []string{"1"}, nil
					}
					if p == path.Join(wal.dir, "1") {
						return []string{"2"}, nil
					}
					return []string{timeutil.FormatTimestamp(now, timeutil.DataTimeFormat4)}, nil
				}
				engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, false)
			},
			wantErr: true,
		},
		{
			name: "partition recovery failure",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(p string) ([]string, error) {
					if p == wal.dir {
						return []string{"1"}, nil
					}
					if p == path.Join(wal.dir, "1") {
						return []string{timeutil.FormatTimestamp(now, timeutil.DataTimeFormat4)}, nil
					}
					return []string{"1"}, nil
				}
				p.EXPECT().recovery(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "partition recovery successfully",
			prepare: func(wal *writeAheadLog) {
				listDirFn = func(p string) ([]string, error) {
					if p == wal.dir {
						return []string{"1"}, nil
					}
					if p == path.Join(wal.dir, "1") {
						return []string{timeutil.FormatTimestamp(now, timeutil.DataTimeFormat4)}, nil
					}
					return []string{"1"}, nil
				}
				p.EXPECT().recovery(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				listDirFn = fileutil.ListDir
			}()
			key := partitionKey{
				shardID:    1,
				familyTime: now,
				leader:     1,
			}
			wal := &writeAheadLog{
				dir:    "db",
				engine: engine,
				familyLogs: map[partitionKey]Partition{
					key: p,
				},
			}
			if tt.prepare != nil {
				tt.prepare(wal)
			}
			err := wal.recovery()
			if (err != nil) != tt.wantErr {
				t.Errorf("recovery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteAheadLog_destroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		removeDirFn = fileutil.RemoveDir
		ctrl.Finish()
	}()

	p1 := NewMockPartition(ctrl)
	p1.EXPECT().Path().Return("p1").AnyTimes()
	p2 := NewMockPartition(ctrl)
	p2.EXPECT().Path().Return("p2").AnyTimes()
	key1 := partitionKey{shardID: 1}
	key2 := partitionKey{shardID: 2}
	wal := &writeAheadLog{
		familyLogs: map[partitionKey]Partition{
			key1: p1,
			key2: p2,
		},
		logger: logger.GetLogger("Test", "WAL"),
	}
	p1.EXPECT().IsExpire().Return(false)
	p2.EXPECT().IsExpire().Return(true)
	p2.EXPECT().Stop()
	p2.EXPECT().Close().Return(fmt.Errorf("err"))
	removeDirFn = func(path string) error {
		return fmt.Errorf("err")
	}
	wal.destroy()

	assert.Len(t, wal.familyLogs, 1)
}

func TestWriteAheadLog_Stop_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p1 := NewMockPartition(ctrl)
	p1.EXPECT().Path().Return("p1").AnyTimes()
	p2 := NewMockPartition(ctrl)
	p2.EXPECT().Path().Return("p2").AnyTimes()
	key1 := partitionKey{shardID: 1}
	key2 := partitionKey{shardID: 2}
	wal := &writeAheadLog{
		familyLogs: map[partitionKey]Partition{
			key1: p1,
			key2: p2,
		},
		logger: logger.GetLogger("Test", "WAL"),
	}
	p1.EXPECT().Stop()
	p2.EXPECT().Stop()
	p1.EXPECT().Close().Return(fmt.Errorf("err"))
	p2.EXPECT().Close().Return(nil)
	wal.Stop()
	_ = wal.Close()
}

func TestWriteAheadLog_Drop(t *testing.T) {
	defer func() {
		removeDirFn = fileutil.RemoveDir
	}()
	wal := &writeAheadLog{}
	removeDirFn = func(path string) error {
		return fmt.Errorf("err")
	}
	assert.Error(t, wal.Drop())
}
