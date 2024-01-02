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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
)

func TestStorageStateMachineAPI_Explore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	statMgr := storage.NewMockStateManager(ctrl)
	api := NewStorageStateMachineAPI(statMgr)
	r := gin.New()
	api.Register(r)

	cases := []struct {
		name    string
		reqBody string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			name:    "param invalid",
			reqBody: ``,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "type not match",
			reqBody: `type=unknown`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "shard assignment",
			reqBody: `type=` + constants.ShardAssignment,
			prepare: func() {
				statMgr.EXPECT().GetShardAssignments().Return([]*models.ShardAssignment{
					{Name: "test2"},
					{Name: "test1"},
				})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "live node",
			reqBody: `type=` + constants.LiveNode,
			prepare: func() {
				statMgr.EXPECT().GetLiveNodes().Return([]models.StatefulNode{
					{StatelessNode: models.StatelessNode{HostIP: "1.1.1.2", HTTPPort: 8080}},
					{StatelessNode: models.StatelessNode{HostIP: "1.1.1.1", HTTPPort: 8080}},
				})
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
			resp := mock.DoRequest(t, r, http.MethodGet, "/state/machine/explore?"+tt.reqBody, "")
			if tt.assert != nil {
				tt.assert(resp)
			}
		})
	}
}
