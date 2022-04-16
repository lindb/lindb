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

package master

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

func TestNewStorageCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	stateMgr := NewMockStateManager(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	cfg := &config.StorageCluster{
		Config: &config.RepoState{Namespace: "test"},
	}

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create repo failure",
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create storage cluster successfully",
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			sc, err := newStorageCluster(context.TODO(), cfg, stateMgr, repoFct)
			if ((err != nil) != tt.wantErr && sc == nil) || (!tt.wantErr && sc == nil) {
				t.Errorf("newStorageCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				assert.Equal(t, cfg, sc.GetConfig())
				assert.NotNil(t, sc.GetState())
				assert.NotNil(t, sc.GetRepo())
			}
		})
	}
}

func TestStorageCluster_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cfg := &config.StorageCluster{
		Config: &config.RepoState{Namespace: "test"},
	}
	repo := state.NewMockRepository(ctrl)
	stateMgr := NewMockStateManager(ctrl)
	fct := &StateMachineFactory{
		ctx: context.TODO(),
	}
	stateMgr.EXPECT().GetStateMachineFactory().Return(fct).AnyTimes()
	sc := &storageCluster{
		ctx:         context.TODO(),
		stateMgr:    stateMgr,
		cfg:         cfg,
		storageRepo: repo,
		logger:      logger.GetLogger("test", "test"),
	}
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := sc.Start()
	assert.Error(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	ch := make(<-chan *state.Event)
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch)
	err = sc.Start()
	assert.NoError(t, err)
}

func TestStorageCluster_listLiveNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	repo := state.NewMockRepository(ctrl)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
		rs      []models.StatefulNode
	}{
		{
			name: "list nodes failure",
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
			rs:      nil,
		},
		{
			name: "list nodes failure",
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Value: []byte{1, 2, 3}}}, nil)
			},
			wantErr: true,
			rs:      nil,
		},
		{
			name: "list nodes successfully",
			prepare: func() {
				n := &models.StatefulNode{ID: 12}
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Value: encoding.JSONMarshal(n)}}, nil)
			},
			wantErr: false,
			rs:      []models.StatefulNode{{ID: 12}},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sc := &storageCluster{
				storageRepo: repo,
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := sc.GetLiveNodes()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.rs, rs)
		})
	}
}

func TestStorageCluster_close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	sm := discovery.NewMockStateMachine(ctrl)
	repo := state.NewMockRepository(ctrl)
	sc := &storageCluster{
		cfg:         &config.StorageCluster{Config: &config.RepoState{Namespace: "test"}},
		sm:          sm,
		storageRepo: repo,
		logger:      logger.GetLogger("test", "test"),
	}
	sm.EXPECT().Close().Return(fmt.Errorf("err"))
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	sc.Close()
}

func TestStorageCluster_SaveDatabaseAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	repo := state.NewMockRepository(ctrl)
	sc := &storageCluster{
		cfg:         &config.StorageCluster{Config: &config.RepoState{Namespace: "test"}},
		storageRepo: repo,
		logger:      logger.GetLogger("test", "test"),
	}
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err := sc.SaveDatabaseAssignment(models.NewShardAssignment("test"), &option.DatabaseOption{})
	assert.Error(t, err)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = sc.SaveDatabaseAssignment(models.NewShardAssignment("test"), &option.DatabaseOption{})
	assert.NoError(t, err)
}

func TestStorageCluster_FlushDatabase(t *testing.T) {
	sc := &storageCluster{}
	assert.Panics(t, func() {
		_ = sc.FlushDatabase("test")
	})
}

func TestStorageCluster_DropDatabaseAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	repo := state.NewMockRepository(ctrl)
	sc := &storageCluster{
		cfg:         &config.StorageCluster{Config: &config.RepoState{Namespace: "test"}},
		storageRepo: repo,
		logger:      logger.GetLogger("test", "test"),
	}
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err := sc.DropDatabaseAssignment("test")
	assert.Error(t, err)

	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	err = sc.DropDatabaseAssignment("test")
	assert.NoError(t, err)
}
