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

package metadata

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

func TestExploreAPI_Explore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewExploreAPI(&deps.HTTPDeps{})
	r := gin.New()
	api.Register(r)
	resp := mock.DoRequest(t, r, http.MethodGet, ExplorePath, "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestExploreAPI_ExploreRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	master := coordinator.NewMockMasterController(ctrl)
	stateMgr := masterpkg.NewMockStateManager(ctrl)
	master.EXPECT().GetStateManager().Return(stateMgr).AnyTimes()
	repo := state.NewMockRepository(ctrl)
	api := NewExploreAPI(&deps.HTTPDeps{
		Repo:   repo,
		Master: master,
		Ctx:    context.Background(),
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
	// case 1: param err
	resp := mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// not found
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=broker&type=LiveNode1", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)

	// case 2: walk entry err
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=broker&type=LiveNode", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 3: walk entry value format err
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, prefix string, fn func(key, value []byte)) error {
			fn([]byte("key"), []byte("value"))
			return nil
		})
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=broker&type=LiveNode", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 4: walk entry ok
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, prefix string, fn func(key, value []byte)) error {
			fn([]byte("key"), encoding.JSONMarshal(&models.StatelessNode{HostIP: "1.1.1.1"}))
			return nil
		})
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=broker&type=LiveNode", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 6: explore master
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, prefix string, fn func(key, value []byte)) error {
			fn([]byte("key"), encoding.JSONMarshal(&models.Master{ElectTime: 11}))
			return nil
		})
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=master&type=Master", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 7: explore storage, storage name is nil
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=storage&type=LiveNode", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 8: explore storage, current master, state not found
	master.EXPECT().IsMaster().Return(true).MaxTimes(2)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=storage&type=LiveNode1&storageName=test", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	stateMgr.EXPECT().GetStorageCluster("test").Return(nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=storage&type=LiveNode&storageName=test", "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// case 9: explore storage, current master, ok
	master.EXPECT().IsMaster().Return(true)
	storage := masterpkg.NewMockStorageCluster(ctrl)
	stateMgr.EXPECT().GetStorageCluster(gomock.Any()).Return(storage).AnyTimes()
	storage.EXPECT().GetRepo().Return(repo)
	repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=storage&type=LiveNode&storageName=test", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	// case 10: explore storage, current is not master, need forward to master
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	}))
	//hack
	_ = backend.Listener.Close()
	l, err := net.Listen("tcp", "127.0.0.1:8089")
	assert.NoError(t, err)
	backend.Listener = l
	// Start the server.
	backend.Start()

	master.EXPECT().IsMaster().Return(false)
	master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: 8089}})
	resp = mock.DoRequest(t, r, http.MethodGet, ExploreRepoPath+"?role=storage&type=LiveNode&storageName=test", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}
