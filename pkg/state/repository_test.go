package state

import (
	"testing"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"
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
