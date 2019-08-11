package database

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

	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err := NewAdminStateMachine(context.TODO(), factory, nil)
	assert.NotNil(t, err)

	storageCluster := storage.NewMockClusterStateMachine(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewAdminStateMachine(context.TODO(), factory, storageCluster)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, stateMachine)

	stateMachine.OnCreate("/data/db1", []byte{1, 1, 1})

	data, _ := json.Marshal(&models.Database{})
	stateMachine.OnCreate("/data/db1", data)

	data, _ = json.Marshal(&models.Database{Name: "db1"})
	stateMachine.OnCreate("/data/db1", data)

	data, _ = json.Marshal(&models.Database{
		Name:     "db1",
		Clusters: []models.DatabaseCluster{{Name: "db1_cluster1"}},
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
		Name: "db1",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "db1_cluster1",
				NumOfShard:    10,
				ReplicaFactor: 3,
			},
		},
	})

	cluster.EXPECT().SaveShardAssign("db1", gomock.Any()).Return(fmt.Errorf("err"))
	cluster.EXPECT().GetActiveNodes().Return(prepareStorageCluster())
	stateMachine.OnCreate("/data/db1", data)

	cluster.EXPECT().SaveShardAssign("db1", gomock.Any()).Return(nil)
	cluster.EXPECT().GetActiveNodes().Return(prepareStorageCluster())
	stateMachine.OnCreate("/data/db1", data)

	stateMachine.OnDelete("mock")
	discovery1.EXPECT().Close()
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
