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

	"github.com/stretchr/testify/assert"
	etcdcliv3 "go.etcd.io/etcd/client/v3"

	"github.com/lindb/lindb/internal/mock"
)

func TestHeartBeat_keepalive_stop(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:9797")
	defer cluster.Terminate(t)
	cfg := etcdcliv3.Config{
		Endpoints: cluster.Endpoints,
	}
	cli, err := etcdcliv3.New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cli)

	key := "/test/heartbeat"
	heartbeat := newHeartbeat(cli, key, []byte("value"), 1, true)
	ctx, cancel := context.WithCancel(context.Background())
	ok, err := heartbeat.grantKeepAliveLease(ctx)
	assert.NoError(t, err)
	assert.True(t, ok)

	_, err = cli.Get(ctx, key)
	assert.NoError(t, err)

	// close heartbeat
	cancel()

	// next term, exist
	heartbeat = newHeartbeat(cli, key, []byte("value"), 1, true)
	ctx, cancel = context.WithCancel(context.Background())
	ok, _ = heartbeat.grantKeepAliveLease(ctx)
	assert.False(t, ok)

	// assert lease expired
	time.Sleep(time.Second * 3)
	ok, _ = heartbeat.grantKeepAliveLease(ctx)
	assert.True(t, ok)

	cancel()

	_ = cli.Close()
}
