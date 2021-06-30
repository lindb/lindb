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

package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestAdminStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)

	discovery1.EXPECT().Discovery(false).Return(fmt.Errorf("err"))
	_, err := NewShardAssignmentStateMachine(context.TODO(), factory, nil)
	assert.NotNil(t, err)

	storageCluster := storage.NewMockClusterStateMachine(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery(false).Return(nil)
	stateMachine, err := NewShardAssignmentStateMachine(context.TODO(), factory, storageCluster)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, stateMachine)

	stateMachine.OnCreate("/data/db1", []byte{1, 1, 1})

	data, _ := json.Marshal(&models.Database{})
	stateMachine.OnCreate("/data/db1", data)

	data, _ = json.Marshal(&models.Database{Name: "db1"})
	storageCluster.EXPECT().GetCluster("").Return(nil)
	stateMachine.OnCreate("/data/db1", data)

	data, _ = json.Marshal(&models.Database{
		Name:    "db1",
		Cluster: "db1_cluster1",
	})
	storageCluster.EXPECT().GetCluster("db1_cluster1").Return(nil)
	stateMachine.OnCreate("/data/db1", data)

	cluster := storage.NewMockCluster(ctrl)
	storageCluster.EXPECT().GetCluster("db1_cluster1").Return(cluster).AnyTimes()
	cluster.EXPECT().GetShardAssign("db1").Return(nil, fmt.Errorf("err"))
	stateMachine.OnCreate("/data/db1", data)

	cluster.EXPECT().GetShardAssign("db1").Return(nil, state.ErrNotExist).AnyTimes()
	cluster.EXPECT().GetActiveNodes().Return(nil)
	stateMachine.OnCreate("/data/db1", data)

	cluster.EXPECT().GetActiveNodes().Return(prepareStorageCluster())
	stateMachine.OnCreate("/data/db1", data)

	data, _ = json.Marshal(&models.Database{
		Name:          "db1",
		Cluster:       "db1_cluster1",
		NumOfShard:    10,
		ReplicaFactor: 3,
	})

	cluster.EXPECT().SaveShardAssign("db1", gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	cluster.EXPECT().GetActiveNodes().Return(prepareStorageCluster())
	stateMachine.OnCreate("/data/db1", data)

	cluster.EXPECT().SaveShardAssign("db1", gomock.Any(), gomock.Any()).Return(nil)
	cluster.EXPECT().GetActiveNodes().Return(prepareStorageCluster())
	stateMachine.OnCreate("/data/db1", data)

	stateMachine.OnDelete("mock")
	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
	_ = stateMachine.Close()
}

func prepareStorageCluster() []*models.ActiveNode {
	return []*models.ActiveNode{
		{Node: models.Node{IP: "127.0.0.1", Port: 2080}},
		{Node: models.Node{IP: "127.0.0.2", Port: 2080}},
		{Node: models.Node{IP: "127.0.0.3", Port: 2080}},
		{Node: models.Node{IP: "127.0.0.4", Port: 2080}},
		{Node: models.Node{IP: "127.0.0.5", Port: 2080}},
	}
}
