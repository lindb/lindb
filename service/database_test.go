package service

import (
	"testing"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/state"
)

func TestCreateDatabase(t *testing.T) {
	//TODO mock it???
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})

	db := New()
	err := db.Create(option.Database{
		Name:          "test",
		NumOfShard:    12,
		ReplicaFactor: 3,
	})
	assert.Nil(t, err)
}

func TestDatabaseService_Get(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})

	db := New()
	err := db.Create(option.Database{
		Name:          "test",
		NumOfShard:    12,
		ReplicaFactor: 3,
	})
	assert.Nil(t, err)
	database, err := db.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", database.Name)
}
