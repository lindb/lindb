package storage

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/server"
	"github.com/eleme/lindb/pkg/state"
)

var storageCfgPath = "./storage.toml"

type testStorageRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestStorageRuntime(t *testing.T) {
	check.Suite(&testStorageRuntimeSuite{})
	check.TestingT(t)
}

func (ts *testStorageRuntimeSuite) TestStorageRun(c *check.C) {
	defer func() {
		_ = fileutil.RemoveDir(storageCfgPath)
	}()
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
		Replication: config.Replication{
			Path: "/tmp/storage/replication",
		},
	}
	_ = fileutil.EncodeToml(storageCfgPath, &cfg)
	storage = NewStorageRuntime(storageCfgPath)
	err = storage.Run()
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(server.Running, check.Equals, storage.State())
	// wait register success
	time.Sleep(200 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := pathutil.GetNodePath(constants.ActiveNodesPath, runtime.node.Indicator())
	nodeBytes, err := runtime.repo.Get(context.TODO(), nodePath)
	if err != nil {
		c.Fatal(err)
	}
	nodeInfo := models.ActiveNode{}
	_ = json.Unmarshal(nodeBytes, &nodeInfo)

	c.Assert(runtime.node, check.Equals, nodeInfo.Node)

	_ = storage.Stop()
	c.Assert(server.Terminated, check.Equals, storage.State())
}
