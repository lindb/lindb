package state

import (
	"context"
	"testing"
	"time"

	etcdcliv3 "go.etcd.io/etcd/clientv3"
	"gopkg.in/check.v1"

	"github.com/lindb/lindb/mock"
)

type testHeartbeatSuite struct {
	mock.RepoTestSuite
}

func TestHeartbeatSuite(t *testing.T) {
	check.Suite(&testHeartbeatSuite{})
	check.TestingT(t)
}

func (ts *testHeartbeatSuite) TestHeartBeat_keepalive_stop(c *check.C) {
	cfg := etcdcliv3.Config{
		Endpoints: ts.Cluster.Endpoints,
	}
	cli, err := etcdcliv3.New(cfg)
	if err != nil {
		c.Fatal(err)
	}
	key := "/test/heartbeat"
	heartbeat := newHeartbeat(cli, key, []byte("value"), 0, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ok, err := heartbeat.grantKeepAliveLease(ctx)
	c.Assert(ok, check.Equals, true)
	if err != nil {
		c.Fatal(err)
	}
	go func() {
		heartbeat.keepAlive(ctx)
	}()

	_, err = cli.Get(ctx, key)
	if err != nil {
		c.Fatal(err)
	}
	_ = cli.Close()
	time.Sleep(time.Second)
}
