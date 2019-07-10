package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	clus := mock.StartEtcdCluster(t)
	defer clus.Terminate(t)
	repo, _ := state.NewRepo(state.Config{
		Endpoints: clus.Endpoints,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := models.Node{IP: "127.0.0.1", Port: 2080}
	register := NewRegister(repo, "/test/node1", node, 1)
	err := register.Register(ctx)
	assert.Nil(t, err)
	nodeBytes, _ := json.Marshal(node)
	getBytes, _ := repo.Get(context.TODO(), "/test/node1")
	assert.True(t, bytes.Equal(nodeBytes, getBytes))
}

func TestUnRegister(t *testing.T) {
	clus := mock.StartEtcdCluster(t)
	defer clus.Terminate(t)
	repo, _ := state.NewRepo(state.Config{
		Endpoints: clus.Endpoints,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node := models.Node{IP: "127.0.0.1", Port: 2080}
	register := NewRegister(repo, "/test/node1", node, 1)
	err := register.Register(ctx)
	assert.Nil(t, err)
	nodeBytes, _ := json.Marshal(node)
	time.Sleep(2 * time.Second)
	getBytes, _ := repo.Get(context.TODO(), "/test/node1")
	assert.True(t, bytes.Equal(nodeBytes, getBytes))
	// unregister the node will be deleted
	_ = register.UnRegister(context.TODO())
	_, err = repo.Get(context.TODO(), "/test/node1")
	assert.NotNil(t, err)
}
