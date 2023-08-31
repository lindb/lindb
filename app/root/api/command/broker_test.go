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

package command

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/encoding"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestBrokerCommand_NotFound(t *testing.T) {
	rs, err := BrokerCommand(context.TODO(), nil, nil, &stmt.Broker{Type: stmt.BrokerOpUnknown})
	assert.Nil(t, rs)
	assert.Nil(t, err)
}

func TestBrokerCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	stateMgr := root.NewMockStateManager(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo:        repo,
		RepoFactory: repoFct,
		StateMgr:    stateMgr,
	}
	cfg := `{"config":{"namespace":"test","timeout":10,"dialTimeout":10,`
	cfg += `"leaseTTL":10,"endpoints":["http://localhost:2379"]}}`

	cases := []struct {
		name    string
		body    *stmt.Broker
		prepare func()
		wantErr bool
	}{
		{
			name: "create broker json err",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: "xx",
			},
			wantErr: true,
		},
		{
			name: "create broker, config validate failure",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: `{"config":{}}`,
			},
			wantErr: true,
		},
		{
			name: "create broker successfully, broker not exist",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						if err := check([]byte{1, 2, 3}); err != nil {
							return false, err
						}
						return true, nil
					})
			},
		},
		{
			name: "create broker successfully, broker exist",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						data := []byte(cfg)
						broker := &config.BrokerCluster{}
						err := encoding.JSONUnmarshal(data, broker)
						assert.NoError(t, err)
						data = encoding.JSONMarshal(broker)
						if err := check(data); err != nil {
							return false, err
						}
						return true, nil
					})
			},
		},
		{
			name: "create broker failure with err",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create broker failure",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "create broker repo failure",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create broker, close repo failure",
			body: &stmt.Broker{
				Type:  stmt.BrokerOpCreate,
				Value: cfg,
			},
			prepare: func() {
				repoFct.EXPECT().CreateBrokerRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "show brokers, get brokers failure",
			body: &stmt.Broker{
				Type: stmt.BrokerOpShow,
			},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "show brokers, list broker successfully, but unmarshal failure",
			body: &stmt.Broker{
				Type: stmt.BrokerOpShow,
			},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte("[]")}}, nil)
			},
			wantErr: true,
		},
		{
			name: "show brokers successfully",
			body: &stmt.Broker{
				Type: stmt.BrokerOpShow,
			},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetBrokerState("xxx").Return(models.BrokerState{}, true)
			},
		},
		{
			name: "show brokers successfully,but state not found",
			body: &stmt.Broker{
				Type: stmt.BrokerOpShow,
			},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetBrokerState("xxx").Return(models.BrokerState{}, false)
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := BrokerCommand(context.TODO(), deps, nil, tt.body)
			if ((err != nil) != tt.wantErr && rs == nil) || (!tt.wantErr && rs == nil) {
				t.Errorf("BrokerCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
