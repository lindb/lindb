package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/server"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
)

var storageCfgPath = "./storage.toml"
var test *testing.T

type testStorageRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestStorageRuntime(t *testing.T) {
	check.Suite(&testStorageRuntimeSuite{})
	test = t
	check.TestingT(t)
}

func (ts *testStorageRuntimeSuite) TestStorageRun(c *check.C) {
	defer util.RemoveDir(storageCfgPath)
	// test run fail
	storage := NewStorageRuntime(storageCfgPath)
	err := storage.Run()
	if err == nil {
		c.Fail()
	}
	c.Assert(server.Failed, check.Equals, storage.State())

	// test normal storage run
	cfg := config.Storage{
		Server: config.Server{
			Port: 9999,
			TTL:  1,
		},
		Coordinator: state.Config{
			Namespace: "/test/storage",
			Endpoints: ts.Cluster.Endpoints,
		},
	}
	util.EncodeToml(storageCfgPath, &cfg)
	storage = NewStorageRuntime(storageCfgPath)
	err = storage.Run()
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(server.Running, check.Equals, storage.State())
	// wait register success
	time.Sleep(500 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := fmt.Sprintf("%s/%s", discovery.ActiveNodesPath, runtime.node.String())

	nodeBytes, _ := runtime.repo.Get(context.TODO(), nodePath)
	nodeInfo := models.Node{}
	json.Unmarshal(nodeBytes, &nodeInfo)

	c.Assert(runtime.node, check.Equals, nodeInfo)

	time.Sleep(500 * time.Millisecond)

	storage.Stop()
	c.Assert(server.Terminated, check.Equals, storage.State())
}
