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

package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/tsdb"
)

func TestNewDatabaseLifecycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	walMgr.EXPECT().Stop().MaxTimes(2)
	walMgr.EXPECT().Close().Return(nil)
	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().Close().MaxTimes(2)

	dbLifecycle := NewDatabaseLifecycle(context.TODO(), repo, walMgr, engine)

	var wait sync.WaitGroup
	wait.Add(1)
	ch := make(chan struct{})
	go func() {
		<-ch
		dbLifecycle.Shutdown()
		wait.Done()
	}()

	dbLifecycle.Startup()
	ch <- struct{}{}
	wait.Wait()

	dbLifecycle = NewDatabaseLifecycle(context.TODO(), repo, walMgr, engine)

	walMgr.EXPECT().Close().Return(fmt.Errorf("err"))
	dbLifecycle.Shutdown()
}

func TestDatabaseLifecycle_ttlTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		config.SetGlobalStorageConfig(config.NewDefaultStorageBase())
		ctrl.Finish()
	}()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().Compact().AnyTimes()
	family.EXPECT().Evict().AnyTimes()
	family.EXPECT().Indicator().Return("ttl_family").AnyTimes()
	tsdb.GetFamilyManager().AddFamily(family)
	defer func() {
		tsdb.GetFamilyManager().RemoveFamily(family)
	}()

	repo := state.NewMockRepository(ctrl)
	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	walMgr.EXPECT().Close()
	walMgr.EXPECT().Stop()
	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().Close()

	dbLifecycle := NewDatabaseLifecycle(context.TODO(), repo, walMgr, engine)
	ch := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		dbLifecycle.Shutdown()
		ch <- struct{}{}
	}()

	dbLifecycle1 := dbLifecycle.(*databaseLifecycle)
	cfg := config.NewDefaultStorageBase()
	cfg.TTLTaskInterval = ltoml.Duration(time.Millisecond * 10)
	config.SetGlobalStorageConfig(cfg)
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	engine.EXPECT().TTL().AnyTimes()
	engine.EXPECT().EvictSegment().AnyTimes()
	dbLifecycle1.ttlTask()
	<-ch
}

func TestDatabaseLifecycle_dropDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	repo := state.NewMockRepository(ctrl)
	walMgr := replica.NewMockWriteAheadLogManager(ctrl)
	engine := tsdb.NewMockEngine(ctrl)

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "list active database err",
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
		},
		{
			name: "active database empty",
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "drop database",
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func([]byte, []byte)) error {
						fn([]byte(constants.GetDatabaseAssignPath("test")), []byte{})
						return nil
					})
				activeDatabases := map[string]struct{}{"test": {}}
				gomock.InOrder(
					walMgr.EXPECT().StopDatabases(activeDatabases),
					engine.EXPECT().DropDatabases(activeDatabases),
					walMgr.EXPECT().DropDatabases(activeDatabases),
				)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			dbLifecycle := NewDatabaseLifecycle(context.TODO(), repo, walMgr, engine)
			dbLifecycle1 := dbLifecycle.(*databaseLifecycle)
			if tt.prepare != nil {
				tt.prepare()
			}
			dbLifecycle1.tryDropDatabases()
		})
	}
}
