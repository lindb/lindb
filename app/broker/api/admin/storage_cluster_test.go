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
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

func TestStorageClusterAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := state.NewMockRepository(ctrl)
	api := NewStorageClusterAPI(&deps.HTTPDeps{
		Ctx:  context.Background(),
		Repo: mockRepo,
		BrokerCfg: &config.BrokerBase{
			HTTP: config.HTTP{
				ReadTimeout: ltoml.Duration(time.Second)},
			Coordinator: config.RepoState{
				Timeout: ltoml.Duration(time.Second * 5)},
		},
	})
	r := gin.New()
	api.Register(r)

	// get request error
	resp := mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, "{}")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	mockRepo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	mockRepo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// GetByName
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=", ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// repo get error
	mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=xxx", ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// unmarshal error
	mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=xxx", ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// ok
	mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("{}"), nil)
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=xxx", ``)
	assert.Equal(t, http.StatusOK, resp.Code)

	// List
	// unmarshal error
	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
		[]state.KeyValue{{Key: "", Value: []byte("[]")}}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// list error
	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
		[]state.KeyValue{{Key: "", Value: []byte("[]")}}, io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// list ok
	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
		[]state.KeyValue{{Key: "", Value: []byte(`{"name": "xxx", "config": {}}`)}}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	// DeleteByName
	// bind error
	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath, ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// delete error
	mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath+"?name=test1", ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// delete ok
	mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath+"?name=test1", "")
	assert.Equal(t, http.StatusNoContent, resp.Code)
}
