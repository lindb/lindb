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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := broker.NewMockStateManager(ctrl)
	repo := state.NewMockRepository(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	deps := &depspkg.HTTPDeps{
		StateMgr:    stateMgr,
		Repo:        repo,
		RepoFactory: repoFct,
	}
	cfg := `{"config":{"namespace":"test","timeout":10,"dialTimeout":10,`
	cfg += `"leaseTTL":10,"endpoints":["http://localhost:2379"]}}`

	mockSrv := func(data []byte) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Add("content-type", "application/json")
			_, _ = rw.Write(data)
		}))
		u, err := url.Parse(server.URL)
		assert.NoError(t, err)
		p, err := strconv.Atoi(u.Port())
		assert.NoError(t, err)
		stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				},
				ID: 1,
			}}}, true)
	}
	databaseCfgData := encoding.JSONMarshal(map[string]models.DatabaseConfig{
		"test": {
			ShardIDs: []models.ShardID{1, 2, 3},
			Option:   &option.DatabaseOption{},
		},
	})
	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "unknown storage op type",
			statement: &stmt.Storage{Type: stmt.StorageOpUnknown},
		},
		{
			name:      "show storages, get storages failure",
			statement: &stmt.Storage{Type: stmt.StorageOpShow},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "show storages, list storage successfully, but unmarshal failure",
			statement: &stmt.Storage{Type: stmt.StorageOpShow},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte("[]")}}, nil)
			},
			wantErr: true,
		},
		{
			name:      "show storages successfully",
			statement: &stmt.Storage{Type: stmt.StorageOpShow},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetStorage("xxx").Return(nil, true)
			},
		},
		{
			name:      "show storages successfully,but state not found",
			statement: &stmt.Storage{Type: stmt.StorageOpShow},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetStorage("xxx").Return(nil, false)
			},
		},
		{
			name:      "create storage json err",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: `xx`},
			wantErr:   true,
		},
		{
			name:      "create storage, config validate failure",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: `{"config":{}}`},
			wantErr:   true,
		},
		{
			name:      "create storage successfully, storage not exist",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
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
			name:      "create storage successfully, storage exist",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						cfg1 := strings.ReplaceAll(cfg, `\"`, `"`)
						data := []byte(cfg1)
						storage := &config.StorageCluster{}
						err := encoding.JSONUnmarshal(data, storage)
						assert.NoError(t, err)
						data = encoding.JSONMarshal(storage)
						if err := check(data); err != nil {
							return false, err
						}
						return true, nil
					})
			},
		},
		{
			name:      "create storage failure with err",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "create storage failure",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name:      "create storage repo failure",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "create storage, close repo failure",
			statement: &stmt.Storage{Type: stmt.StorageOpCreate, Value: cfg},
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "recover storage, but storage not found",
			statement: &stmt.Storage{Type: stmt.StorageOpRecover, Value: "test"},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
			wantErr: true,
		},
		{
			name:      "recover storage, but get database config failure",
			statement: &stmt.Storage{Type: stmt.StorageOpRecover, Value: "test"},
			prepare: func() {
				mockSrv([]byte("abc"))
			},
		},
		{
			name:      "recover storage, but recover shard assignment failure",
			statement: &stmt.Storage{Type: stmt.StorageOpRecover, Value: "test"},
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "recover storage, but recover database schema failure",
			statement: &stmt.Storage{Type: stmt.StorageOpRecover, Value: "test"},
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "recover storage successfully",
			statement: &stmt.Storage{Type: stmt.StorageOpRecover, Value: "test"},
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := StorageCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("StorageCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
