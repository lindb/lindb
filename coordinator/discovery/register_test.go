package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := models.Node{IP: "127.0.0.1", Port: 2080}
	register := NewRegister("/test/node1", node, 1)
	err := register.Register(ctx)
	assert.Nil(t, err)
	nodeBytes, _ := json.Marshal(node)
	getBytes, _ := state.GetRepo().Get(context.TODO(), "/test/node1")
	assert.True(t, bytes.Equal(nodeBytes, getBytes))

}

func TestUnRegister(t *testing.T) {

	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := models.Node{IP: "127.0.0.1", Port: 2080}
	register := NewRegister("/test/node1", node, 1)
	err := register.Register(ctx)
	assert.Nil(t, err)
	nodeBytes, _ := json.Marshal(node)
	time.Sleep(2 * time.Second)
	getBytes, _ := state.GetRepo().Get(context.TODO(), "/test/node1")
	assert.True(t, bytes.Equal(nodeBytes, getBytes))
	// unregister the node will be deleted
	_ = register.UnRegister(context.TODO())
	_, err = state.GetRepo().Get(context.TODO(), "/test/node1")
	assert.NotNil(t, err)
}
