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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/service"

	"github.com/golang/mock/gomock"
)

func TestStorageClusterAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageClusterService := service.NewMockStorageClusterService(ctrl)

	api := NewStorageClusterAPI(storageClusterService)

	// get request error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/storage/cluster",
		RequestBody:    []byte{1, 2, 3},
		HandlerFunc:    api.Create,
		ExpectHTTPCode: http.StatusInternalServerError,
	})

	cfg := config.StorageCluster{
		Name: "test1",
	}
	storageClusterService.EXPECT().Save(gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/storage/cluster",
		RequestBody:    cfg,
		HandlerFunc:    api.Create,
		ExpectHTTPCode: 204,
	})
	storageClusterService.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/storage/cluster",
		RequestBody:    cfg,
		HandlerFunc:    api.Create,
		ExpectHTTPCode: 500,
	})
	storageClusterService.EXPECT().Get(gomock.Any()).Return(&cfg, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 200,
		ExpectResponse: cfg,
	})
	storageClusterService.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})

	storageClusterService.EXPECT().List().Return([]*config.StorageCluster{&cfg}, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/storage/cluster",
		HandlerFunc:    api.List,
		ExpectHTTPCode: 200,
		ExpectResponse: []config.StorageCluster{cfg},
	})
	storageClusterService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/storage/cluster",
		HandlerFunc:    api.List,
		ExpectHTTPCode: 500,
	})

	storageClusterService.EXPECT().Delete(gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodDelete,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.DeleteByName,
		ExpectHTTPCode: 204,
	})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodDelete,
		URL:            "/storage/cluster",
		HandlerFunc:    api.DeleteByName,
		ExpectHTTPCode: 500,
	})
	storageClusterService.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodDelete,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.DeleteByName,
		ExpectHTTPCode: 500,
	})
}
