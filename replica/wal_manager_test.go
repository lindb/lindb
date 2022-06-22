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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
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
	// other db
	l = m.GetOrCreateLog("test-2")
	assert.NotNil(t, l)
}

func TestWAL_garbageCollectTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.TODO())
	log := NewMockWriteAheadLog(ctrl)
	wal := &writeAheadLogManager{
		ctx: ctx,
		databaseLogs: map[string]WriteAheadLog{
			"test": log,
		},
		cfg: config.WAL{RemoveTaskInterval: ltoml.Duration(time.Millisecond * 10)},
	}
	log.EXPECT().destroy().AnyTimes()
	wal.garbageCollectTask()

	time.Sleep(time.Millisecond * 50)
	cancel()
	time.Sleep(time.Millisecond * 50)
}

func TestWriteAheadLogManager_GetReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	log := NewMockWriteAheadLog(ctrl)
	mgr := &writeAheadLogManager{
		databaseLogs: map[string]WriteAheadLog{
			"test": log,
		},
	}
	log.EXPECT().getReplicaState().Return([]models.FamilyLogReplicaState{{}})
	s := mgr.GetReplicaState("test")
	assert.Len(t, s, 1)

	// db not exist
	s = mgr.GetReplicaState("test-not-exist")
	assert.Nil(t, s)
}

func TestWriteAheadLogManager_Recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		fileExistFn = fileutil.Exist
		listDirFn = fileutil.ListDir
		ctrl.Finish()
	}()
	log := NewMockWriteAheadLog(ctrl)
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
				log.EXPECT().recovery().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "recovery write ahead log successfully",
			prepare: func() {
				log.EXPECT().recovery().Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				fileExistFn = func(file string) bool {
					return true
				}
				listDirFn = func(path string) ([]string, error) {
					return []string{"test"}, nil
				}
			}()
			mgr := &writeAheadLogManager{
				databaseLogs: map[string]WriteAheadLog{
					"test": log,
				},
			}
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

func TestWriteAheadLogManager_Stop_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log1 := NewMockWriteAheadLog(ctrl)
	log1.EXPECT().Name().Return("test1").AnyTimes()
	log2 := NewMockWriteAheadLog(ctrl)
	log2.EXPECT().Name().Return("test2").AnyTimes()
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "stop channel then close log",
			prepare: func() {
				log1.EXPECT().Stop()
				log2.EXPECT().Stop()
				log1.EXPECT().Close().Return(fmt.Errorf("err"))
				log2.EXPECT().Close().Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mgr := &writeAheadLogManager{
				databaseLogs: map[string]WriteAheadLog{
					"test1": log1,
					"test2": log2,
				},
				logger: logger.GetLogger("Test", "WAL"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			mgr.Stop()
			_ = mgr.Close()
		})
	}
}

func TestMockWriteAheadLogMockRecorder_Drop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log1 := NewMockWriteAheadLog(ctrl)
	log1.EXPECT().Name().Return("test1").AnyTimes()
	log2 := NewMockWriteAheadLog(ctrl)
	log2.EXPECT().Name().Return("test2").AnyTimes()
	cases := []struct {
		name    string
		prepare func()
		assert  func(d *writeAheadLogManager)
	}{
		{
			name: "close log failure",
			prepare: func() {
				log2.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			assert: func(d *writeAheadLogManager) {
				assert.Len(t, d.databaseLogs, 2)
			},
		},
		{
			name: "drop log failure",
			prepare: func() {
				log2.EXPECT().Close().Return(nil)
				log2.EXPECT().Drop().Return(fmt.Errorf("err"))
			},
			assert: func(d *writeAheadLogManager) {
				assert.Len(t, d.databaseLogs, 2)
			},
		},
		{
			name: "drop database successfully",
			prepare: func() {
				log2.EXPECT().Close().Return(nil)
				log2.EXPECT().Drop().Return(nil)
			},
			assert: func(d *writeAheadLogManager) {
				assert.Len(t, d.databaseLogs, 1)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mgr := &writeAheadLogManager{
				databaseLogs: map[string]WriteAheadLog{
					"test1": log1,
					"test2": log2,
				},
				logger: logger.GetLogger("Test", "WAL"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			mgr.DropDatabases(map[string]struct{}{"test1": {}})
			tt.assert(mgr)
		})
	}
}

func TestMockWriteAheadLogMockRecorder_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log1 := NewMockWriteAheadLog(ctrl)
	log1.EXPECT().Name().Return("test1").AnyTimes()
	log2 := NewMockWriteAheadLog(ctrl)
	log2.EXPECT().Name().Return("test2").AnyTimes()
	cases := []struct {
		name    string
		prepare func()
		assert  func(d *writeAheadLogManager)
	}{
		{
			name: "stop database successfully",
			prepare: func() {
				log2.EXPECT().Stop()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mgr := &writeAheadLogManager{
				databaseLogs: map[string]WriteAheadLog{
					"test1": log1,
					"test2": log2,
				},
				logger: logger.GetLogger("Test", "WAL"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			mgr.StopDatabases(map[string]struct{}{"test1": {}})
		})
	}
}
