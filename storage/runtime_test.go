package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

type testStorageRuntimeSuite struct {
	mock.RepoTestSuite
	t *testing.T
}

func TestStorageRuntime(t *testing.T) {
	check.Suite(&testStorageRuntimeSuite{t: t})
	check.TestingT(t)
}

var cfg = config.Storage{StorageKernel: config.StorageKernel{
	GRPC: config.GRPC{
		Port: 9999,
		TTL:  1,
	},
	Engine: config.Engine{Dir: "/tmp/storage/data"},
	Coordinator: config.RepoState{
		Namespace: "/test/storage",
	},
	Replication: config.Replication{
		Dir: "/tmp/storage/replication",
	},
}}

func (ts *testStorageRuntimeSuite) TestStorageRun(c *check.C) {
	// test normal storage run
	cfg.Coordinator.Endpoints = ts.Cluster.Endpoints
	cfg.GRPC.Port = 9999
	storage := NewStorageRuntime("test-version", cfg)
	err := storage.Run()
	assert.NoError(ts.t, err)
	c.Assert(server.Running, check.Equals, storage.State())
	// wait register success
	time.Sleep(500 * time.Millisecond)

	runtime, _ := storage.(*runtime)
	nodePath := constants.GetNodePath(constants.ActiveNodesPath, runtime.node.Indicator())
	nodeBytes, err := runtime.repo.Get(context.TODO(), nodePath)
	assert.NoError(ts.t, err)

	nodeInfo := models.ActiveNode{}
	_ = json.Unmarshal(nodeBytes, &nodeInfo)

	c.Assert(runtime.node, check.Equals, nodeInfo.Node)
	c.Assert("storage", check.Equals, storage.Name())

	_ = storage.Stop()
	c.Assert(server.Terminated, check.Equals, storage.State())
	time.Sleep(500 * time.Millisecond)
}

func (ts *testStorageRuntimeSuite) TestBrokerRun_GetHost_Err(c *check.C) {
	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
	}()
	cfg.GRPC.Port = 8889
	storage := NewStorageRuntime("test-version", cfg)
	getHostIP = func() (string, error) {
		return "test-ip", fmt.Errorf("err")
	}
	err := storage.Run()
	assert.NotNil(ts.t, err)

	getHostIP = func() (string, error) {
		return "ip", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	cfg.GRPC.Port = 8887
	cfg.Coordinator.Endpoints = ts.Cluster.Endpoints
	storage = NewStorageRuntime("test-version", cfg)
	err = storage.Run()
	assert.Nil(ts.t, err)
	err = storage.Stop()
	assert.NoError(ts.t, err)
}
func (ts *testStorageRuntimeSuite) TestBrokerRun_Err(c *check.C) {
	ctrl := gomock.NewController(ts.t)
	defer ctrl.Finish()

	cfg.GRPC.Port = 8886
	storage := NewStorageRuntime("test-version", cfg)
	s := storage.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := s.Run()
	assert.Error(ts.t, err)

	registry := discovery.NewMockRegistry(ctrl)
	s.registry = registry
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	repo := state.NewMockRepository(ctrl)
	s.repo = repo
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	err = s.Stop()
	assert.NoError(ts.t, err)
}
