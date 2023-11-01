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

package state

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
)

func TestBrokerStateMachineAPI_Explore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := broker.NewMockStateManager(ctrl)
	masterStateMgr := masterpkg.NewMockStateManager(ctrl)
	master := coordinator.NewMockMasterController(ctrl)
	master.EXPECT().GetStateManager().Return(masterStateMgr).AnyTimes()

	deps := &depspkg.HTTPDeps{
		StateMgr: stateMgr,
		Master:   master,
	}
	cli := client.NewMockStateMachineCli(ctrl)
	api := NewBrokerStateMachineAPI(deps)
	api.cli = cli
	r := gin.New()
	api.Register(r)

	cases := []struct {
		name    string
		req     string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			name: "param invalid",
			req:  ``,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "role not match",
			req:  `role=9999&type=` + constants.StorageState,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "broker state, type not match",
			req:  `role=2&type=` + constants.ShardAssignment,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "broker state, live nodes",
			req:  `role=2&type=` + constants.LiveNode,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{
					{HostIP: "1.1.1.2", HTTPPort: 8080},
					{HostIP: "1.1.1.1", HTTPPort: 8080},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "broker state, storage state",
			req:  `role=2&type=` + constants.StorageState,
			prepare: func() {
				stateMgr.EXPECT().GetStorageList().Return([]*models.StorageState{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "broker state, database config",
			req:  `role=2&type=` + constants.DatabaseConfig,
			prepare: func() {
				stateMgr.EXPECT().GetDatabases().Return([]models.Database{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "master state, type not match",
			req:  `role=3&type=` + constants.LiveNode,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "master state, database config",
			req:  `role=3&type=` + constants.DatabaseConfig,
			prepare: func() {
				masterStateMgr.EXPECT().GetDatabases().Return([]models.Database{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "master state, storage state",
			req:  `role=3&type=` + constants.StorageState,
			prepare: func() {
				masterStateMgr.EXPECT().GetStorageStates().Return([]*models.StorageState{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "master state, storage config",
			req:  `role=3&type=` + constants.StorageConfig,
			prepare: func() {
				masterStateMgr.EXPECT().GetStorages().Return([]config.StorageCluster{
					{Config: &config.RepoState{Namespace: "test2"}},
					{Config: &config.RepoState{Namespace: "test1"}},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "master state, master",
			req:  `role=3&type=` + constants.Master,
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "master state, shard assignment",
			req:  `role=3&type=` + constants.ShardAssignment,
			prepare: func() {
				masterStateMgr.EXPECT().GetShardAssignments().Return([]models.ShardAssignment{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "storage state, storage not found",
			req:  `role=4&type=` + constants.ShardAssignment,
			prepare: func() {
				masterStateMgr.EXPECT().GetStorageCluster(gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "storage state, get storage live node failure",
			req:  `role=4&type=` + constants.ShardAssignment,
			prepare: func() {
				cluster := masterpkg.NewMockStorageCluster(ctrl)
				masterStateMgr.EXPECT().GetStorageCluster(gomock.Any()).Return(cluster)
				cluster.EXPECT().GetLiveNodes().Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "storage state, get storage state successfully",
			req:  `role=4&type=` + constants.ShardAssignment,
			prepare: func() {
				cluster := masterpkg.NewMockStorageCluster(ctrl)
				masterStateMgr.EXPECT().GetStorageCluster(gomock.Any()).Return(cluster)
				cluster.EXPECT().GetLiveNodes().Return([]models.StatefulNode{
					{StatelessNode: models.StatelessNode{HostIP: "1.1.1.2", HTTPPort: 8080}},
					{StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", HTTPPort: 8080}},
				}, nil)
				cli.EXPECT().FetchStateByNodes(gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			resp := mock.DoRequest(t, r, http.MethodGet, "/state/machine/explore?"+tt.req, "")
			if tt.assert != nil {
				tt.assert(resp)
			}
		})
	}
}
