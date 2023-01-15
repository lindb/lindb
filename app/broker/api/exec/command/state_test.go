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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

func TestState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := broker.NewMockStateManager(ctrl)
	master := coordinator.NewMockMasterController(ctrl)
	deps := &depspkg.HTTPDeps{
		StateMgr: stateMgr,
		Master:   master,
	}

	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "unknown metadata statement type",
			statement: &stmt.State{},
		},
		{
			name:      "master not found",
			statement: &stmt.State{Type: stmt.Master},
			prepare: func() {
				master.EXPECT().GetMaster().Return(nil)
			},
		},
		{
			name:      "found master",
			statement: &stmt.State{Type: stmt.Master},
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{})
			},
		},
		{
			name:      "show broker alive node",
			statement: &stmt.State{Type: stmt.BrokerAlive},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "1.1.1.1",
					HTTPPort: 8080,
				}})
			},
		},
		{
			name:      "show storage alive node",
			statement: &stmt.State{Type: stmt.StorageAlive},
			prepare: func() {
				stateMgr.EXPECT().GetStorageList().Return([]*models.StorageState{})
			},
		},
		{
			name:      "show memory database state, but storage not found",
			statement: &stmt.State{Type: stmt.MemoryDatabase, StorageName: "a", Database: "b"},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
		},
		{
			name:      "show memory database state",
			statement: &stmt.State{Type: stmt.MemoryDatabase, StorageName: "a", Database: "b"},
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("[]"))
				}))
				u, err := url.Parse(svr.URL)
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
			},
		},
		{
			name:      "show replication state, but storage not found",
			statement: &stmt.State{Type: stmt.Replication, StorageName: "a", Database: "b"},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
		},
		{
			name:      "show replication state, alive node empty",
			statement: &stmt.State{Type: stmt.Replication, StorageName: "a", Database: "b"},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
					LiveNodes: nil}, true)
			},
		},
		{
			name:      "show replication state, but fetch state failure",
			statement: &stmt.State{Type: stmt.Replication, StorageName: "a", Database: "b"},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
					LiveNodes: map[models.NodeID]models.StatefulNode{1: {
						StatelessNode: models.StatelessNode{
							HostIP:   "127.0.01", // mock host err
							HTTPPort: 8080,
						},
						ID: 1,
					}}}, true)
			},
		},
		{
			name:      "show replication state, but fetch state failure",
			statement: &stmt.State{Type: stmt.Replication, StorageName: "a", Database: "b"},
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("[]"))
				}))
				u, err := url.Parse(svr.URL)
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
			},
		},
		{
			name:      "show broker metric, no alive node",
			statement: &stmt.State{Type: stmt.BrokerMetric, MetricNames: []string{"a", "b"}},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return(nil)
			},
		},
		{
			name:      "show broker metric successfully",
			statement: &stmt.State{Type: stmt.BrokerMetric, MetricNames: []string{"a", "b"}},
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`{"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}, {
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}})
			},
		},
		{
			name:      "show storage metric, storage name empty",
			statement: &stmt.State{Type: stmt.StorageMetric, StorageName: "", MetricNames: []string{"a", "b"}},
			wantErr:   true,
		},
		{
			name:      "show storage metric, storage not found",
			statement: &stmt.State{Type: stmt.StorageMetric, StorageName: "a", MetricNames: []string{"a", "b"}},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
		},
		{
			name:      "show storage metric, storage no alive node",
			statement: &stmt.State{Type: stmt.StorageMetric, StorageName: "a", MetricNames: []string{"a", "b"}},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{}, true)
			},
		},
		{
			name:      "show storage metric successfully",
			statement: &stmt.State{Type: stmt.StorageMetric, StorageName: "a", MetricNames: []string{"a", "b"}},
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).
					Return(&models.StorageState{LiveNodes: map[models.NodeID]models.StatefulNode{1: {}, 2: {}}}, true)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := StateCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("StateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
