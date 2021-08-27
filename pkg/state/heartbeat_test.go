// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package state

import (
	"context"
	"testing"
	"time"

	etcdcliv3 "go.etcd.io/etcd/clientv3"
	"gopkg.in/check.v1"

	"github.com/lindb/lindb/internal/mock"
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
	heartbeat := newHeartbeat(cli, key, []byte("value"), 1, true)
	ctx, cancel := context.WithCancel(context.Background())
	ok, err := heartbeat.grantKeepAliveLease(ctx)
	c.Assert(ok, check.Equals, true)
	if err != nil {
		c.Fatal(err)
	}

	_, err = cli.Get(ctx, key)
	if err != nil {
		c.Fatal(err)
	}
	// close heartbeat
	cancel()

	// next term, exist
	heartbeat = newHeartbeat(cli, key, []byte("value"), 1, true)
	ctx, cancel = context.WithCancel(context.Background())
	ok, _ = heartbeat.grantKeepAliveLease(ctx)
	c.Assert(ok, check.Equals, false)

	// assert lease expired
	time.Sleep(time.Second * 2)
	ok, _ = heartbeat.grantKeepAliveLease(ctx)
	c.Assert(ok, check.Equals, true)

	cancel()

	_ = cli.Close()
}
