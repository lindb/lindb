package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	"gopkg.in/check.v1"
)

type testRegistrySuite struct {
	mock.RepoTestSuite
}

var testRegistryPath = "/test/registry"

func TestRegistry(t *testing.T) {
	check.Suite(&testRegistrySuite{})
	check.TestingT(t)
}

func (ts *testRegistrySuite) TestRegister(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	registry := NewRegistry(repo, testRegistryPath, 100)

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	err := registry.Register(node)
	if err != nil {
		c.Fatal(err)
	}
	// wait register success
	time.Sleep(500 * time.Millisecond)

	nodePath := fmt.Sprintf("%s/%s", testRegistryPath, node.Indicator())
	nodeBytes, _ := repo.Get(context.TODO(), nodePath)
	nodeInfo := models.Node{}
	_ = json.Unmarshal(nodeBytes, &nodeInfo)
	c.Assert(node, check.Equals, nodeInfo)

	// test re-register
	_ = repo.Delete(context.TODO(), nodePath)
	_, err = repo.Get(context.TODO(), nodePath)
	c.Assert(err, check.NotNil)
	// wait register success
	time.Sleep(500 * time.Millisecond)
	nodeBytes, _ = repo.Get(context.TODO(), nodePath)
	_ = json.Unmarshal(nodeBytes, &nodeInfo)
	c.Assert(node, check.Equals, nodeInfo)

	_ = registry.Close()
	time.Sleep(time.Second)
	_, err = repo.Get(context.TODO(), nodePath)
	c.Assert(err, check.NotNil)
}

func (ts *testRegistrySuite) TestDeregister(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	registry := NewRegistry(repo, testRegistryPath, 100)
	defer func() {
		_ = registry.Close()
	}()

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	err := registry.Register(node)
	if err != nil {
		c.Fatal(err)
	}
	// wait register success
	time.Sleep(500 * time.Millisecond)

	nodePath := fmt.Sprintf("%s/%s", testRegistryPath, node.Indicator())
	nodeBytes, _ := repo.Get(context.TODO(), nodePath)
	nodeInfo := models.Node{}
	_ = json.Unmarshal(nodeBytes, &nodeInfo)
	c.Assert(node, check.Equals, nodeInfo)

	_ = registry.Deregister(node)
	time.Sleep(500 * time.Millisecond)
	_, err = repo.Get(context.TODO(), nodePath)
	c.Assert(err, check.NotNil)
}
