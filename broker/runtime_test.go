package broker

import (
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

type testBrokerRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestBrokerRuntime(t *testing.T) {
	check.Suite(&testBrokerRuntimeSuite{})
	check.TestingT(t)
}

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
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
	broker := NewBrokerRuntime(cfg)
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
