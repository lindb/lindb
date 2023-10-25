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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
)

func TestRootStateMachineAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := root.NewMockStateManager(ctrl)

	deps := &depspkg.HTTPDeps{
		StateMgr: stateMgr,
	}
	api := NewRootStateMachineAPI(deps)
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
			name: "live nodes",
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
			name: "database config",
			req:  `role=3&type=` + constants.DatabaseConfig,
			prepare: func() {
				stateMgr.EXPECT().GetDatabases().Return([]models.LogicDatabase{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "broker state list",
			req:  `role=3&type=` + constants.BrokerState,
			prepare: func() {
				stateMgr.EXPECT().GetBrokerStates().Return([]models.BrokerState{
					{Name: "test1"},
					{Name: "test2"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "broker state found",
			req:  `role=3&brokerName=test&type=` + constants.BrokerState,
			prepare: func() {
				stateMgr.EXPECT().GetBrokerState(gomock.Any()).Return(models.BrokerState{}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "broker state not found",
			req:  `role=3&brokerName=test&type=` + constants.BrokerState,
			prepare: func() {
				stateMgr.EXPECT().GetBrokerState(gomock.Any()).Return(models.BrokerState{}, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
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
