package storage

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/server"
)

type testStorageRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestStorageRuntime(t *testing.T) {
	check.Suite(&testStorageRuntimeSuite{})
	check.TestingT(t)
}

func (ts *testStorageRuntimeSuite) TestStorageRun(c *check.C) {
	// test normal storage run
	cfg := config.Storage{StorageKernel: config.StorageKernel{
		GRPC: config.GRPC{
			Port: 9999,
			TTL:  1,
		},
		Engine: config.Engine{Dir: "/tmp/storage/data"},
		Coordinator: config.RepoState{
			Namespace: "/test/storage",
			Endpoints: ts.Cluster.Endpoints,
		},
		Replication: config.Replication{
			Dir: "/tmp/storage/replication",
		},
	}}
	storage := NewStorageRuntime(cfg)
	err := storage.Run()
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(server.Running, check.Equals, storage.State())
	// wait register success
	time.Sleep(200 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := constants.GetNodePath(constants.ActiveNodesPath, runtime.node.Indicator())
	nodeBytes, err := runtime.repo.Get(context.TODO(), nodePath)
	if err != nil {
		c.Fatal(err)
	}
	activeNodeMap := models.ActiveNodeMap{}
	_ = json.Unmarshal(nodeBytes, &activeNodeMap)

	c.Assert(runtime.node, check.Equals, *activeNodeMap.NodeMap[models.NodeTypeRPC])
	c.Assert("storage", check.Equals, storage.Name())

	_ = storage.Stop()
	c.Assert(server.Terminated, check.Equals, storage.State())
}
