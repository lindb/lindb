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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

func TestClusterStateMachine_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controllerFactory := task.NewMockControllerFactory(ctrl)
	storageService := service.NewMockStorageStateService(ctrl)
	shardAssignService := service.NewMockShardAssignService(ctrl)

	repoFactory := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	discoverFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	cluster := NewMockCluster(ctrl)
	discoverFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	clusterFactory := NewMockClusterFactory(ctrl)

	// list exist storage cluster err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)

	assert.NotNil(t, err)

	// register discovery err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)
	assert.NotNil(t, err)

	// normal case
	repo.EXPECT().List(gomock.Any(), gomock.Any()).
		Return([]state.KeyValue{
			{Key: "test1", Value: encoding.JSONMarshal(&models.StorageState{Name: "test1"})},
			{Key: "unmarshal_err", Value: []byte{1, 2, 3}},
			{Key: "test2", Value: encoding.JSONMarshal(&models.StorageState{Name: "test2"})},
			{Key: "test3", Value: encoding.JSONMarshal(&models.StorageState{Name: "test3"})},
			{Key: "error", Value: encoding.JSONMarshal(&models.StorageState{})},
		}, nil).AnyTimes()
	repo1 := state.NewMockRepository(ctrl)
	repo1.EXPECT().Close().Return(nil)

	gomock.InOrder(
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(repo1, nil),
		clusterFactory.EXPECT().newCluster(gomock.Any()).Return(cluster, fmt.Errorf("err")),
		cluster.EXPECT().Close(),
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(state.NewMockRepository(ctrl), nil),
		clusterFactory.EXPECT().newCluster(gomock.Any()).Return(cluster, nil),
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err")),
	)
	discovery1.EXPECT().Discovery().Return(nil)

	stateMachine, err := NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)

	assert.Nil(t, err)
	assert.NotNil(t, stateMachine)
	assert.Equal(t, 1, len(stateMachine.GetAllCluster()))
	assert.Equal(t, cluster, stateMachine.GetCluster("test2"))
	assert.Nil(t, stateMachine.GetCluster("test1"))

	// OnDelete
	cluster.EXPECT().Close()
	stateMachine.OnDelete("/test/data/test2")
	assert.Equal(t, 0, len(stateMachine.GetAllCluster()))

	// OnCreate
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(state.NewMockRepository(ctrl), nil)
	clusterFactory.EXPECT().newCluster(gomock.Any()).Return(cluster, nil)
	stateMachine.OnCreate("/test/data/test1", encoding.JSONMarshal(&models.StorageState{Name: "test1"}))

	cluster.EXPECT().Close()
	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
}

func TestClusterStateMachine_collect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controllerFactory := task.NewMockControllerFactory(ctrl)
	storageService := service.NewMockStorageStateService(ctrl)
	shardAssignService := service.NewMockShardAssignService(ctrl)

	repoFactory := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	discoverFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Discovery().Return(nil)
	discoverFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	clusterFactory := NewMockClusterFactory(ctrl)

	// list exist storage cluster err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	sm, err := NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	sm1 := sm.(*clusterStateMachine)
	cluster := NewMockCluster(ctrl)
	sm1.clusters["test"] = cluster
	cluster.EXPECT().CollectStat().Return(nil, fmt.Errorf("err"))
	cluster.EXPECT().CollectStat().Return(&models.StorageClusterStat{}, nil).AnyTimes()
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	sm1.interval = 300 * time.Millisecond
	sm1.timer.Reset(100 * time.Millisecond)

	time.Sleep(time.Second)
}
