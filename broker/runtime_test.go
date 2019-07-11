package broker

import (
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/pkg/server"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
)

var brokerCfgPath = "./broker.toml"
var test *testing.T

type testBrokerRuntimeSuite struct {
	mock.RepoTestSuite
}

func TestBrokerRuntime(t *testing.T) {
	check.Suite(&testBrokerRuntimeSuite{})
	test = t
	check.TestingT(t)
}

func (ts *testBrokerRuntimeSuite) TestBrokerRun(c *check.C) {
	defer util.RemoveDir(brokerCfgPath)
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
	}
	util.EncodeToml(brokerCfgPath, &cfg)
	broker = NewBrokerRuntime(brokerCfgPath)
	err = broker.Run()
	if err != nil {
		c.Fatal(err)
	}
	// wait run finish
	time.Sleep(500 * time.Millisecond)

	c.Assert(server.Running, check.Equals, broker.State())

	broker.Stop()
	c.Assert(server.Terminated, check.Equals, broker.State())
}
