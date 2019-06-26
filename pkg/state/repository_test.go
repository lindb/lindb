package state

import (
	"testing"

	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"
	etcd "go.etcd.io/etcd/clientv3"
)

func TestNewRepo(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	cfg := etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	}

	var err = New("zk", cfg)
	assert.NotNil(t, err)
	assert.Nil(t, GetRepo())

	err = New("etcd", "error_config")
	assert.NotNil(t, err)
	assert.Nil(t, GetRepo())

	err = New("etcd", cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, GetRepo())

	GetRepo().Close()
}
