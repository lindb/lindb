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

package root

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/pkg/state"
)

func TestNewBrokerCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	stateMgr := NewMockStateManager(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	cfg := &config.BrokerCluster{
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
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			}, wantErr: true,
		},
		{
			name: "create broker cluster successfully",
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
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

			sc, err := newBrokerCluster(cfg, stateMgr, repoFct)
			if ((err != nil) != tt.wantErr && sc == nil) || (!tt.wantErr && sc == nil) {
				t.Errorf("newBrokerCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				assert.NotNil(t, sc.GetState())
			}
		})
	}
}

func TestBrokerCluster_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	cfg := &config.BrokerCluster{
		Config: &config.RepoState{Namespace: "test"},
	}
	repo := state.NewMockRepository(ctrl)
	stateMgr := NewMockStateManager(ctrl)
	fct := &stateMachineFactory{
		ctx: context.TODO(),
	}
	stateMgr.EXPECT().GetStateMachineFactory().Return(fct).AnyTimes()
	sc := &brokerCluster{
		stateMgr:   stateMgr,
		cfg:        cfg,
		brokerRepo: repo,
		logger:     logger.GetLogger("Master", "Test"),
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

func TestBrokerCluster_close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	sm := discovery.NewMockStateMachine(ctrl)
	repo := state.NewMockRepository(ctrl)
	sc := &brokerCluster{
		cfg:        &config.BrokerCluster{Config: &config.RepoState{Namespace: "test"}},
		sm:         sm,
		brokerRepo: repo,
		logger:     logger.GetLogger("Root", "Test"),
	}
	sm.EXPECT().Close().Return(fmt.Errorf("err"))
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	sc.Close()
}
