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

func TestReplicaAPI_GetReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		doRequest = defaultClient.Do
		ctrl.Finish()
	}()

	stateMgr := broker.NewMockStateManager(ctrl)
	api := NewReplicaAPI(&deps.HTTPDeps{StateMgr: stateMgr})
	r := gin.New()
	api.Register(r)

	// case 1: params invalid
	resp := mock.DoRequest(t, r, http.MethodGet, ReplicaPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 2: storage not found
	stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
	resp = mock.DoRequest(t, r, http.MethodGet, ReplicaPath+"?storageName=test&db=db", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 3: live node is empty
	stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
		LiveNodes: nil}, true)
	resp = mock.DoRequest(t, r, http.MethodGet, ReplicaPath+"?storageName=test&db=db", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 4: fetch err
	stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
		LiveNodes: map[models.NodeID]models.StatefulNode{1: {
			StatelessNode: models.StatelessNode{
				HostIP:   "1.1.1.1",
				HTTPPort: 8080,
			},
			ID: 1,
		}}}, true).AnyTimes()
	doRequest = func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("err")
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ReplicaPath+"?storageName=test&db=db", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	// case 5: fetch ok, resp data invalid
	doRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: io.NopCloser(strings.NewReader("a"))}, nil
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ReplicaPath+"?storageName=test&db=db", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 5: data is ok
	doRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: io.NopCloser(strings.NewReader("[]"))}, nil
	}
	resp = mock.DoRequest(t, r, http.MethodGet, ReplicaPath+"?storageName=test&db=db", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}
