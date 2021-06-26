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

package storage

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

func TestStorageCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	databaseOption := option.DatabaseOption{
		Interval: "10s",
	}
	factory := NewClusterFactory()
	storage := config.StorageCluster{
		Config: config.RepoState{Namespace: "storage"},
	}
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	storageService := service.NewMockStorageStateService(ctrl)
	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	controller := task.NewMockController(ctrl)
	controller.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	controllerFactory := task.NewMockControllerFactory(ctrl)
	controllerFactory.EXPECT().CreateController(gomock.Any(), gomock.Any()).Return(controller).AnyTimes()
	shardAssignService := service.NewMockShardAssignService(ctrl)
	cfg := clusterCfg{
		storageStateService: storageService,
		cfg:                 storage,
		repo:                repo,
		factory:             discoveryFactory,
		controllerFactory:   controllerFactory,
		shardAssignService:  shardAssignService,
		logger:              logger.GetLogger("coordinator", "storage-test"),
	}
	_, err := factory.newCluster(cfg)
	assert.NotNil(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "/node1", Value: encoding.JSONMarshal(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 4000}})},
		{Key: "/node_err", Value: []byte{1, 1, 1, 1}},
		{Key: "/node2", Value: encoding.JSONMarshal(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 4000}})},
	}, nil).AnyTimes()

	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = factory.newCluster(cfg)
	assert.NotNil(t, err)

	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	discovery1.EXPECT().Discovery().Return(nil)
	cluster, err := factory.newCluster(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// get active nodes
	assert.Equal(t, 2, len(cluster.GetActiveNodes()))
	// OnCreate
	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	cluster.OnCreate("/active/node/1",
		encoding.JSONMarshal(&models.ActiveNode{Node: models.Node{IP: "1.1.1.4", Port: 4000}}))
	cluster.OnCreate("/active/node/2", []byte{1, 2, 3})
	assert.Equal(t, 3, len(cluster.GetActiveNodes()))

	// OnDelete
	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	cluster.OnDelete("/active/nodes/1.1.1.2:4000")
	assert.Equal(t, 2, len(cluster.GetActiveNodes()))

	// get shard assign
	shardAssignService.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
	shardAssign, err := cluster.GetShardAssign("test")
	assert.Nil(t, shardAssign)
	assert.NotNil(t, err)
	shardAssignService.EXPECT().Get(gomock.Any()).Return(models.NewShardAssignment("test"), nil)
	shardAssign, err = cluster.GetShardAssign("test")
	assert.NotNil(t, shardAssign)
	assert.Nil(t, err)

	// save shard assignment
	shardAssign = models.NewShardAssignment("test")
	shardAssign.AddReplica(1, 1)
	shardAssign.AddReplica(2, 1)
	shardAssign.AddReplica(3, 2)
	shardAssign.AddReplica(4, 2)
	shardAssign.Nodes[1] = &models.Node{IP: "1.1.1.1", Port: 8000}
	shardAssign.Nodes[2] = &models.Node{IP: "1.1.1.2", Port: 8000}
	// save shard assign err
	shardAssignService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = cluster.SaveShardAssign("test", shardAssign, databaseOption)
	assert.NotNil(t, err)
	// submit task err
	shardAssignService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = cluster.SaveShardAssign("test", shardAssign, databaseOption)
	assert.NotNil(t, err)
	// success
	shardAssignService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = cluster.SaveShardAssign("test", shardAssign, databaseOption)
	assert.Nil(t, err)

	// test submit task
	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err1 := cluster.SubmitTask("test", "test", nil)
	assert.NotNil(t, err1)
	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err1 = cluster.SubmitTask("test", "test", nil)
	assert.Nil(t, err1)

	assert.Equal(t, repo, cluster.GetRepo())

	discovery1.EXPECT().Close()
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	cluster.Close()

}

func TestCluster_CollectStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := NewClusterFactory()
	storage := config.StorageCluster{
		Config: config.RepoState{Namespace: "storage"},
	}
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	storageService := service.NewMockStorageStateService(ctrl)
	repo := state.NewMockRepository(ctrl)
	discovery1.EXPECT().Discovery().Return(nil)

	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	controller := task.NewMockController(ctrl)
	controller.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	controllerFactory := task.NewMockControllerFactory(ctrl)
	controllerFactory.EXPECT().CreateController(gomock.Any(), gomock.Any()).Return(controller).AnyTimes()
	cfg := clusterCfg{
		storageStateService: storageService,
		cfg:                 storage,
		repo:                repo,
		factory:             discoveryFactory,
		controllerFactory:   controllerFactory,
		logger:              logger.GetLogger("coordinator", "storage-test"),
	}
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	cluster1, err := factory.newCluster(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cluster1)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	stat, err := cluster1.CollectStat()
	assert.Error(t, err)
	assert.Nil(t, stat)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{Key: "/test/test", Value: []byte{1, 1}}}, nil)
	stat, err = cluster1.CollectStat()
	assert.Error(t, err)
	assert.Nil(t, stat)

	activeNode := models.ActiveNode{
		Node: models.Node{IP: "1.1.1.1", Port: 9000},
	}
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "/test/1.1.1.1:9000", Value: encoding.JSONMarshal(&models.NodeStat{Node: activeNode})},
		{Key: "/test/test-2", Value: encoding.JSONMarshal(&models.NodeStat{Node: models.ActiveNode{
			Node: models.Node{IP: "1.1.1.2", Port: 9000},
		}})},
	}, nil)
	cluster2 := cluster1.(*cluster)
	cluster2.clusterState = &models.StorageState{
		Name:        "/test",
		ActiveNodes: map[string]*models.ActiveNode{"1.1.1.1:9000": &activeNode},
	}
	stat, err = cluster1.CollectStat()
	assert.NoError(t, err)
	assert.NotNil(t, stat)
}

func TestCluster_FlushDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := NewClusterFactory()
	storage := config.StorageCluster{
		Config: config.RepoState{Namespace: "storage"},
	}
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	storageService := service.NewMockStorageStateService(ctrl)
	repo := state.NewMockRepository(ctrl)
	discovery1.EXPECT().Discovery().Return(nil)

	storageService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	controller := task.NewMockController(ctrl)
	controller.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	controllerFactory := task.NewMockControllerFactory(ctrl)
	controllerFactory.EXPECT().CreateController(gomock.Any(), gomock.Any()).Return(controller).AnyTimes()
	cfg := clusterCfg{
		storageStateService: storageService,
		cfg:                 storage,
		repo:                repo,
		factory:             discoveryFactory,
		controllerFactory:   controllerFactory,
		logger:              logger.GetLogger("coordinator", "storage-test"),
	}
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	cluster1, err := factory.newCluster(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, cluster1)

	cluster2 := cluster1.(*cluster)
	cluster2.mutex.Lock()
	cluster2.clusterState.AddActiveNode(&models.ActiveNode{
		Node: models.Node{IP: "1.1.1.1", Port: 9000},
	})
	cluster2.mutex.Unlock()
	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = cluster1.FlushDatabase("test")
	assert.NoError(t, err)

	controller.EXPECT().Submit(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = cluster1.FlushDatabase("test")
	assert.Error(t, err)
}
