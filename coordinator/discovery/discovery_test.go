package discovery

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	"github.com/coreos/pkg/capnslog"
	"github.com/stretchr/testify/assert"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

func TestNodeList(t *testing.T) {
	clus := mock.StartEtcdCluster(t)
	defer clus.Terminate(t)
	repo, _ := state.NewRepo(state.Config{
		Endpoints: clus.Endpoints,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newDiscovery, err := NewDiscovery(ctx, repo, "/test/")

	if err != nil {
		t.Errorf("Discovery error :%s", err.Error())
	}
	if newDiscovery == nil {
		t.Error("Discovery is empty")
		return
	}
	nodeList := newDiscovery.NodeList()
	assert.Equal(t, 0, len(nodeList))

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	bytes, _ := json.Marshal(node)
	_ = repo.Put(context.TODO(), "/test/key1", bytes)
	_ = repo.Put(context.TODO(), "/test/key2", bytes)
	_ = repo.Put(context.TODO(), "/test/key3", bytes)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 3, len(newDiscovery.NodeList()))
}
