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
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shirou/gopsutil/disk"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
)

const clusterName = "test"

func TestStorageAPI_GetStorageClusterState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	stateMachine := broker.NewMockStorageStateMachine(ctrl)
	api := NewStorageAPI(&deps.HTTPDeps{
		Repo: repo,
		StateMachines: &coordinator.BrokerStateMachines{
			StorageSM: stateMachine,
		},
		Ctx: context.Background(),
		BrokerCfg: &config.BrokerBase{
			HTTP: config.HTTP{
				ReadTimeout: ltoml.Duration(time.Second)},
			Coordinator: config.RepoState{
				Timeout: ltoml.Duration(time.Second * 5)},
		},
	})
	r := gin.New()
	api.Register(r)

	// cluster name not input
	resp := mock.DoRequest(t, r, http.MethodGet, StorageStatePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// database get err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodGet, StorageStatePath+"?name=test", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	api.ListDatabase = func() (databases []*models.Database, err error) {
		return []*models.Database{{
			Name:          "test-db",
			Cluster:       "test-db",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}, {
			Name:          "test-db-2",
			Cluster:       "test-2",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}, {
			Name:          "test-db-3",
			Cluster:       "test",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}}, nil
	}
	sa := models.ShardAssignment{
		Name: "test-db",
		Nodes: map[int]*models.Node{
			1: {IP: "1.1.1.1", Port: 2890},
			2: {IP: "1.1.1.1", Port: 9000},
			5: {IP: "1.1.1.2", Port: 9000},
			6: {IP: "1.1.1.3", Port: 9000},
		},
		Shards: map[int]*models.Replica{1: {Replicas: []int{1, 2}}, 2: {Replicas: []int{5, 6}}},
	}
	saData, _ := json.Marshal(sa)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "", Value: saData}}, nil)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodGet, StorageStatePath+"?name=test", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// get ok
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "", Value: saData}}, nil)
	activeNode := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}}
	clusterStat := models.StorageClusterStat{
		Name: "test",
		Nodes: []*models.NodeStat{{Node: activeNode, System: models.SystemStat{
			DiskUsageStat: &disk.UsageStat{
				Total:       10,
				Used:        10,
				UsedPercent: 10,
			},
		}}},
		NodeStatus:    models.NodeStatus{},
		ReplicaStatus: models.ReplicaStatus{},
		Capacity:      disk.UsageStat{},
	}
	data, _ := json.Marshal(clusterStat)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data, nil)
	storageState := models.NewStorageState()
	storageState.Name = clusterName
	storageState.AddActiveNode(&activeNode)
	stateMachine.EXPECT().List().Return([]*models.StorageState{storageState})
	resp = mock.DoRequest(t, r, http.MethodGet, StorageStatePath+"?name=test-db", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestStorageAPI_ListStorageClusterState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	stateMachine := broker.NewMockStorageStateMachine(ctrl)
	api := NewStorageAPI(&deps.HTTPDeps{
		Repo: repo,
		StateMachines: &coordinator.BrokerStateMachines{
			StorageSM: stateMachine,
		},
		Ctx: context.Background(),
		BrokerCfg: &config.BrokerBase{
			HTTP: config.HTTP{
				ReadTimeout: ltoml.Duration(time.Second)},
			Coordinator: config.RepoState{
				Timeout: ltoml.Duration(time.Second * 5)},
		},
	})
	r := gin.New()
	api.Register(r)

	// get error
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "", Value: nil}}, nil)
	api.ListDatabase = func() (databases []*models.Database, err error) {
		return nil, io.ErrClosedPipe
	}
	resp := mock.DoRequest(t, r, http.MethodGet, ListStorageStatePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	//	// get shard assign err
	//	databaseService.EXPECT().List().Return([]*models.Database{{
	//		Name:          "test-db",
	//		Cluster:       "test",
	//		NumOfShard:    10,
	//		ReplicaFactor: 1,
	//	}, {
	//		Name:          "test-db-2",
	//		Cluster:       "test-2",
	//		NumOfShard:    10,
	//		ReplicaFactor: 1,
	//	}, {
	//		Name:          "test-db-3",
	//		Cluster:       "test",
	//		NumOfShard:    10,
	//		ReplicaFactor: 1,
	//	}, {
	//		Name:          "test-db-4",
	//		Cluster:       "test-3",
	//		NumOfShard:    10,
	//		ReplicaFactor: 1,
	//	},
	//	}, nil).AnyTimes()
	//	shardAssignService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	//	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageStatePath, "")
	//	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	//
	sas := []models.ShardAssignment{
		{
			Name: "test-db",
			Nodes: map[int]*models.Node{
				1: {IP: "1.1.1.1", Port: 2890},
				2: {IP: "1.1.1.1", Port: 9000},
				5: {IP: "1.1.1.2", Port: 9000},
				6: {IP: "1.1.1.3", Port: 9000},
			},
			Shards: map[int]*models.Replica{1: {Replicas: []int{1, 2}}, 2: {Replicas: []int{5, 6}}},
		},
		{
			Name: "test-db-2",
			Nodes: map[int]*models.Node{
				1: {IP: "1.1.1.1", Port: 2890},
				2: {IP: "1.1.1.1", Port: 9000},
				5: {IP: "1.1.1.2", Port: 9000},
				6: {IP: "1.1.1.3", Port: 9000},
			},
			Shards: map[int]*models.Replica{1: {Replicas: []int{1, 2}}, 2: {Replicas: []int{5, 6}}},
		}}
	var keys []state.KeyValue
	for _, sa := range sas {
		data, _ := json.Marshal(sa)
		keys = append(keys, state.KeyValue{Key: "", Value: data})
	}
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(keys, nil)

	api.ListDatabase = func() (databases []*models.Database, err error) {
		return []*models.Database{{
			Name:          "test-db",
			Cluster:       "test-db",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}, {
			Name:          "test-db-2",
			Cluster:       "test-2",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}, {
			Name:          "test-db-3",
			Cluster:       "test",
			NumOfShard:    10,
			ReplicaFactor: 1,
		}}, nil
	}
	activeNode := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}}
	clusterStat := models.StorageClusterStat{
		Name: "test",
		Nodes: []*models.NodeStat{{Node: activeNode, System: models.SystemStat{
			DiskUsageStat: &disk.UsageStat{
				Total:       10,
				Used:        10,
				UsedPercent: 10,
			},
		}}},
		NodeStatus:    models.NodeStatus{},
		ReplicaStatus: models.ReplicaStatus{},
		Capacity:      disk.UsageStat{},
	}
	data1 := encoding.JSONMarshal(&clusterStat)
	clusterStat.Name = "test-2"
	data2 := encoding.JSONMarshal(&clusterStat)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "/kv/test", Value: data1},
		{Key: "/kv/test-2", Value: data2},
	}, nil)
	storageState := models.NewStorageState()
	storageState.Name = clusterName
	storageState.AddActiveNode(&activeNode)
	stateMachine.EXPECT().List().Return([]*models.StorageState{storageState})
	resp = mock.DoRequest(t, r, http.MethodGet, ListStorageStatePath, "")
	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestNewStorageAPI_nodeIsAlive(t *testing.T) {
	assert.False(t, nodeIsAlive(nil, "test"))
	assert.False(t, nodeIsAlive(&models.StorageState{
		ActiveNodes: map[string]*models.ActiveNode{}}, "test"))
}
