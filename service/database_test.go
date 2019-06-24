package service

import (
	"testing"

	"github.com/coreos/etcd/integration"
	etcd "go.etcd.io/etcd/clientv3"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/state"
)

func TestCreateDatabase(t *testing.T) {
	//TODO mock it???
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	state.New("etcd", etcd.Config{
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
