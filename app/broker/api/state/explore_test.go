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
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
)

func TestExploreAPI_Explore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		doRequest = defaultClient.Do
		ctrl.Finish()
	}()

	stateMgr := broker.NewMockStateManager(ctrl)
	api := NewExploreAPI(&deps.HTTPDeps{StateMgr: stateMgr})
	r := gin.New()
	api.Register(r)

	// case 1: params invalid
	resp := mock.DoRequest(t, r, http.MethodGet, ExplorePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Broker1&names=cpu", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 2: fetch err
	stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
		HostIP:   "1.1.1.1",
		HTTPPort: 8080,
	}}).AnyTimes()
	doRequest = func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("err")
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Broker&names=cpu", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	// case 3: fetch ok, resp data invalid
	doRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: io.NopCloser(strings.NewReader("a"))}, nil
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Broker&names=cpu", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 4: broker success
	buf := io.NopCloser(strings.NewReader(`{
"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))
	doRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: buf}, nil
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Broker&names=cpu", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 5: storage name is nil, err
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Storage&names=cpu", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 6: storage name not exist
	stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Storage&names=cpu&storageName=xx", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 6: storage no live node
	stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{}, true)
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Storage&names=cpu&storageName=xx", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)

	// case 7: storage success
	doRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: io.NopCloser(strings.NewReader(`{
"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))}, nil
	}
	stateMgr.EXPECT().GetStorage(gomock.Any()).
		Return(&models.StorageState{LiveNodes: map[models.NodeID]models.StatefulNode{1: {}, 2: {}}}, true)
	resp = mock.DoRequest(t, r, http.MethodGet, ExplorePath+"?role=Storage&names=cpu&storageName=xx", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestExploreAPI_ExploreLiveNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		doRequest = defaultClient.Do
		ctrl.Finish()
	}()

	stateMgr := broker.NewMockStateManager(ctrl)
	api := NewExploreAPI(&deps.HTTPDeps{StateMgr: stateMgr})
	r := gin.New()
	api.Register(r)

	// case 1: params invalid
	resp := mock.DoRequest(t, r, http.MethodGet, ExploreLiveNodePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreLiveNodePath+"?role=Broker1", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 2: broker ok
	stateMgr.EXPECT().GetLiveNodes().Return(nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreLiveNodePath+"?role=Broker", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 3: storage ok
	stateMgr.EXPECT().GetStorageList().Return(nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreLiveNodePath+"?role=Storage", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}
