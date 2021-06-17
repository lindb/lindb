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
	"fmt"
	"net/http"
	"testing"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStorageClusterAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageClusterService := service.NewMockStorageClusterService(ctrl)

	api := NewStorageClusterAPI(&deps.HTTPDeps{
		StorageClusterSrv: storageClusterService,
	})
	r := gin.New()
	api.Register(r)

	// get request error
	resp := mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, "{}")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	cfg := config.StorageCluster{
		Name: "test1",
	}
	storageClusterService.EXPECT().Save(gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	storageClusterService.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodPost, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	storageClusterService.EXPECT().Get(gomock.Any()).Return(&cfg, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=test1", `{"name":"test1"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	storageClusterService.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath+"?name=test1", `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	resp = mock.DoRequest(t, r, http.MethodGet, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	storageClusterService.EXPECT().List().Return([]*config.StorageCluster{&cfg}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusOK, resp.Code)
	storageClusterService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	storageClusterService.EXPECT().Delete(gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath+"?name=test1", ``)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath, `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	storageClusterService.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodDelete, StorageClusterPath+"?name=test1", `{"name":"test1"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
