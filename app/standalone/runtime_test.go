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

package standalone

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/server/v3/embed"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/state"
)

func newDefaultStandaloneConfig(_ *testing.T) config.Standalone {
	saCfg := config.Standalone{
		Query:       *config.NewDefaultQuery(),
		Coordinator: *config.NewDefaultCoordinator(),
		StorageBase: *config.NewDefaultStorageBase(),
		BrokerBase:  *config.NewDefaultBrokerBase(),
		Logging:     *logger.NewDefaultSetting(),
		ETCD:        *config.NewDefaultETCD(),
		Monitor:     *config.NewDefaultMonitor(),
	}
	dir := "."
	saCfg.StorageBase.TSDB.Dir = filepath.Join(dir, "data")
	saCfg.StorageBase.WAL.Dir = filepath.Join(dir, "wal")
	saCfg.StorageBase.GRPC.Port = 3901
	saCfg.StorageBase.HTTP.Port = 3902
	saCfg.StorageBase.WAL.RemoveTaskInterval = ltoml.Duration(10 * time.Minute)
	config.SetGlobalStorageConfig(&saCfg.StorageBase)
	return saCfg
}

func TestRuntime_New(t *testing.T) {
	defer func() {
		assert.NoError(t, fileutil.RemoveDir("data"))
	}()
	cfg := newDefaultStandaloneConfig(t)
	standalone := NewStandaloneRuntime("test-version", &cfg, true)
	assert.NotNil(t, standalone)
	assert.NotNil(t, standalone.Config())
	assert.Equal(t, "standalone", standalone.Name())
}

func TestRuntime_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		assert.NoError(t, fileutil.RemoveDir("data"))
		ctrl.Finish()
	}()

	repoFct := state.NewMockRepositoryFactory(ctrl)
	s := server.NewMockService(ctrl)
	init := bootstrap.NewMockClusterInitializer(ctrl)
	pusher := monitoring.NewMockNativePusher(ctrl)
	pusher.EXPECT().Stop().AnyTimes()
	cases := []struct {
		name    string
		prepare func(cfg *config.Standalone)
		wantErr bool
	}{
		{
			name: "start etcd server failure",
			prepare: func(_ *config.Standalone) {
				startEtcdFn = func(inCfg *embed.Config) (e *embed.Etcd, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create broker state failure",
			prepare: func(_ *config.Standalone) {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{

			name: "clean up master state failure",
			prepare: func(_ *config.Standalone) {
				repo := state.NewMockRepository(ctrl)
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{

			name: "create storage state failure",
			prepare: func(_ *config.Standalone) {
				repo := state.NewMockRepository(ctrl)
				repo.EXPECT().Close().Return(nil)
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "list storage alive nodes failure",
			prepare: func(_ *config.Standalone) {
				repo := state.NewMockRepository(ctrl)
				repo.EXPECT().Close().Return(nil).MaxTimes(2)
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "delete storage alive node failure",
			prepare: func(_ *config.Standalone) {
				repo := state.NewMockRepository(ctrl)
				gomock.InOrder(
					repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil),
					repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
					repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil),
					repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "/a/b"}}, nil),
					repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")),
					repo.EXPECT().Close().Return(fmt.Errorf("err")),
					repo.EXPECT().Close().Return(nil),
				)
			},
			wantErr: true,
		},
		{
			name: "run broker failure",
			prepare: func(_ *config.Standalone) {
				mockCleanState(ctrl, repoFct)
				s.EXPECT().Run().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "run storage failure",
			prepare: func(_ *config.Standalone) {
				mockCleanState(ctrl, repoFct)
				s.EXPECT().Run().Return(nil)
				s.EXPECT().Run().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "run storage failure",
			prepare: func(_ *config.Standalone) {
				mockCleanState(ctrl, repoFct)
				s.EXPECT().Run().Return(nil)
				s.EXPECT().Run().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "init internal database failure",
			prepare: func(_ *config.Standalone) {
				mockCleanState(ctrl, repoFct)
				s.EXPECT().Run().Return(nil).MaxTimes(2)
				init.EXPECT().InitInternalDatabase(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: false,
		},
		{
			name: "init internal database successfully",
			prepare: func(_ *config.Standalone) {
				mockCleanState(ctrl, repoFct)
				s.EXPECT().Run().Return(nil).MaxTimes(2)
				init.EXPECT().InitInternalDatabase(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				startEtcdFn = embed.StartEtcd
			}()
			cfg := newDefaultStandaloneConfig(t)
			ctx, cancel := context.WithCancel(context.TODO())
			s.EXPECT().Stop().MaxTimes(2)
			r := &runtime{
				embedEtcd:   true,
				ctx:         ctx,
				cancel:      cancel,
				cfg:         &cfg,
				repoFactory: repoFct,
				broker:      s,
				storage:     s,
				initializer: init,
				pusher:      pusher,
				delayInit:   time.Millisecond,
			}
			if tt.prepare != nil {
				tt.prepare(r.cfg)
			}
			err := r.Run()
			time.Sleep(100 * time.Millisecond)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				assert.Equal(t, server.Failed, r.State())
			} else {
				assert.Equal(t, server.Running, r.State())
			}
			r.Stop()
			assert.Equal(t, server.Terminated, r.State())
		})
	}
}

func mockCleanState(ctrl *gomock.Controller, repoFct *state.MockRepositoryFactory) {
	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Close().Return(nil).MaxTimes(2)
	repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "/a/b"}}, nil)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
}
