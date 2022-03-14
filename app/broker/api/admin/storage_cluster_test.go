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

package admin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

func TestStorageClusterAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare
	repoFct := state.NewMockRepositoryFactory(ctrl)
	mockRepo := state.NewMockRepository(ctrl)
	stateMgr := broker.NewMockStateManager(ctrl)
	api := NewStorageClusterAPI(&deps.HTTPDeps{
		Ctx:         context.Background(),
		Repo:        mockRepo,
		RepoFactory: repoFct,
		StateMgr:    stateMgr,
		BrokerCfg: &config.Broker{
			BrokerBase: config.BrokerBase{
				HTTP: config.HTTP{
					ReadTimeout: ltoml.Duration(time.Second)}},
			Coordinator: config.RepoState{
				Timeout: ltoml.Duration(time.Second * 5)},
		},
	})
	r := gin.New()
	api.Register(r)

	// build test cases
	tests := []struct {
		name    string
		method  string
		url     string
		reqBody string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			"get storage param invalid",
			http.MethodGet,
			StorageClusterPath + "?name=",
			``,
			nil,
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			"get storage failure",
			http.MethodGet,
			StorageClusterPath + "?name=xxx",
			``,
			func() {
				mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			"get storage successfully, but unmarshal failure",
			http.MethodGet,
			StorageClusterPath + "?name=xxx",
			``,
			func() {
				mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			"get storage successfully",
			http.MethodGet,
			StorageClusterPath + "?name=xxx",
			``,
			func() {
				mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("{}"), nil)
			},
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			"delete storage param invalid",
			http.MethodDelete,
			StorageClusterPath,
			``,
			nil,
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			"delete storage failure",
			http.MethodDelete,
			StorageClusterPath + "?name=test1",
			``,
			func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), constants.GetStorageClusterConfigPath("test1")).
					Return(io.ErrClosedPipe)
			},
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			"delete storage successfully",
			http.MethodDelete,
			StorageClusterPath + "?name=test1",
			``,
			func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), constants.GetStorageClusterConfigPath("test1")).
					Return(nil)
			},
			func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, resp.Code)
			},
		},
	}

	// run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			resp := mock.DoRequest(t, r, tt.method, tt.url, tt.reqBody)
			if tt.assert != nil {
				tt.assert(resp)
			}
		})
	}
}
