package broker

import (
	"testing"
	"time"

	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

var brokerCfgPath = "./broker.toml"

type testBrokerRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestBrokerRuntime(t *testing.T) {
	check.Suite(&testBrokerRuntimeSuite{})
	check.TestingT(t)
}

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
	defer func() {
		_ = fileutil.RemoveDir(brokerCfgPath)
	}()
	// test run fail
	broker := NewBrokerRuntime(brokerCfgPath)
	err := broker.Run()
	if err == nil {
		c.Fail()
	}
	c.Assert(server.Failed, check.Equals, broker.State())

	// test normal broker run
	cfg := config.Broker{
		HTTP: config.HTTP{
			Port: 9999,
		},
		Coordinator: state.Config{
			Namespace: "/test/broker",
			Endpoints: ts.Cluster.Endpoints,
		},
		Server: config.Server{
			Port: 2881,
			TTL:  1,
		},
		ReplicationChannel: config.ReplicationChannel{
			Path:                       "/tmp/broker/replication",
			BufferSize:                 32,
			SegmentFileSize:            128 * 1024 * 1024,
			RemoveTaskIntervalInSecond: 60,
		},
	}
	_ = fileutil.EncodeToml(brokerCfgPath, &cfg)
	broker = NewBrokerRuntime(brokerCfgPath)
	err = broker.Run()
	if err != nil {
		c.Fatal(err)
	}
	// wait run finish
	time.Sleep(500 * time.Millisecond)

	c.Assert(server.Running, check.Equals, broker.State())

	_ = broker.Stop()
	c.Assert(server.Terminated, check.Equals, broker.State())
}
