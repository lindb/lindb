package broker

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/hostutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

type testBrokerRuntimeSuite struct {
	mock.RepoTestSuite
	t *testing.T
}

func TestBrokerRuntime(t *testing.T) {
	check.Suite(&testBrokerRuntimeSuite{t: t})
	check.TestingT(t)
}

var cfg = config.Broker{
	Monitor: config.Monitor{
		RuntimeReportInterval: ltoml.Duration(10 * time.Second),
		SystemReportInterval:  ltoml.Duration(10 * time.Second),
	},
	BrokerBase: config.BrokerBase{
		HTTP: config.HTTP{
			Port: 9999,
		},
		Coordinator: config.RepoState{
			Namespace: "/test/broker",
		},
		GRPC: config.GRPC{
			Port: 2881,
			TTL:  1,
		},

		ReplicationChannel: config.ReplicationChannel{
			Dir:                "/tmp/broker/replication",
			DataSizeLimit:      128,
			RemoveTaskInterval: ltoml.Duration(time.Minute),
			CheckFlushInterval: ltoml.Duration(time.Second),
			FlushInterval:      ltoml.Duration(time.Second * 5),
			BufferSize:         128,
		},
	}}

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
	cfg.BrokerBase.Coordinator.Endpoints = ts.Cluster.Endpoints

	broker := NewBrokerRuntime("test-version", cfg)
	err := broker.Run()
	if err != nil {
		c.Fatal(err)
	}
	// wait run finish
	time.Sleep(500 * time.Millisecond)

	c.Assert(server.Running, check.Equals, broker.State())
	c.Assert("broker", check.Equals, broker.Name())

	_ = broker.Stop()
	c.Assert(server.Terminated, check.Equals, broker.State())
}

func (ts *testBrokerRuntimeSuite) TestBrokerRun_GetHost_Err(c *check.C) {
	defer func() {
		getHostIP = hostutil.GetHostIP
		hostName = os.Hostname
	}()
	broker := NewBrokerRuntime("test-version", cfg)
	getHostIP = func() (string, error) {
		return "ip1", fmt.Errorf("err")
	}
	err := broker.Run()
	c.Assert(err, check.NotNil)

	getHostIP = func() (string, error) {
		return "ip2", nil
	}
	hostName = func() (string, error) {
		return "host", fmt.Errorf("err")
	}
	err = broker.Run()
	if err != nil {
		c.Fatal(err)
	}
}

func (ts *testBrokerRuntimeSuite) TestBroker_Run_Err(c *check.C) {
	ctrl := gomock.NewController(ts.t)
	defer ctrl.Finish()

	defer func() {
		if err := recover(); err != nil {
			assert.NotNil(ts.t, err)
		}
	}()

	broker := NewBrokerRuntime("test-version", cfg)
	b := broker.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	b.repoFactory = repoFactory
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := broker.Run()
	assert.Error(ts.t, err)
	_ = broker.Stop()

	repo := state.NewMockRepository(ctrl)
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(repo, nil).AnyTimes()
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = broker.Run()
	assert.Error(ts.t, err)
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	registry := discovery.NewMockRegistry(ctrl)
	b.registry = registry
	registry.EXPECT().Close().Return(fmt.Errorf("err"))
	_ = broker.Stop()
}
