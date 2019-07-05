package discovery

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/state"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/coreos/pkg/capnslog"
	"github.com/stretchr/testify/assert"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

func TestNodeList(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newDiscovery, err := NewDiscovery(ctx, "/test/")

	if err != nil {
		t.Errorf("Discovery error :%s", err.Error())
	}
	if newDiscovery == nil {
		t.Error("Discovery is empty")
		return
	}
	nodeList := newDiscovery.NodeList()
	assert.Equal(t, 0, len(nodeList))

	repo := state.GetRepo()
	node := Node{"127.0.0.1", 2080}
	bytes, _ := json.Marshal(node)
	_ = repo.Put(context.TODO(), "/test/key1", bytes)
	_ = repo.Put(context.TODO(), "/test/key2", bytes)
	_ = repo.Put(context.TODO(), "/test/key3", bytes)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 3, len(newDiscovery.NodeList()))
}
