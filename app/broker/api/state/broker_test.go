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
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

func TestBrokerAPI_ListBrokersStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	stateMgr := broker.NewMockStateManager(ctrl)
	api := NewBrokerAPI(&deps.HTTPDeps{
		Repo:     repo,
		StateMgr: stateMgr,
		Ctx:      context.Background(),
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

	// get stat list err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	resp := mock.DoRequest(t, r, http.MethodGet, BrokerStatePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// decoding stat err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{
			Key:   "/test/1.1.1.1:2080",
			Value: []byte{1, 2, 3},
		},
	}, nil)
	node := models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 2080}
	nodes := []models.StatelessNode{node}
	stateMgr.EXPECT().GetLiveNodes().Return(nodes)
	system := models.SystemStat{
		CPUs: 100,
	}
	resp = mock.DoRequest(t, r, http.MethodGet, BrokerStatePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// success
	stateMgr.EXPECT().GetLiveNodes().Return(nodes)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{
			Key: "/test/1.1.1.1:2080",
			Value: encoding.JSONMarshal(&models.NodeStat{
				Node:   &node,
				System: system,
			}),
		},
		{
			Key: "/test/1.1.1.2:2080",
			Value: encoding.JSONMarshal(&models.NodeStat{
				Node:   &node,
				System: system,
			}),
		},
	}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, BrokerStatePath, "")
	assert.Equal(t, http.StatusOK, resp.Code)
}
