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

package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc/proto/common"
)

var (
	node = models.Node{
		IP:   "127.0.0.1",
		Port: 123,
	}
	database = "database"
	shardID  = int32(0)
)

func TestClientConnFactory(t *testing.T) {
	node1 := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}

	node2 := models.Node{
		IP:   "1.1.1.1",
		Port: 456,
	}

	fct := GetClientConnFactory()

	conn1, err := fct.GetClientConn(node1)
	if err != nil {
		t.Fatal(err)
	}

	conn11, err := fct.GetClientConn(node1)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := fct.GetClientConn(node2)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, conn1 == conn11)
	assert.False(t, conn1 == conn2)
}

func TestContext(t *testing.T) {
	node := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}
	ctx := CreateIncomingContext(context.TODO(), database, shardID, node)

	n, err := GetLogicNodeFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, &node)

	db, err := GetDatabaseFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, db, database)

	sID, err := GetShardIDFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, shardID, sID)
}

func TestClientStreamFactory(t *testing.T) {
	target := models.Node{
		IP:   "127.0.0.1",
		Port: 1234,
	}
	fct := NewClientStreamFactory(node)
	_, err := fct.CreateWriteServiceClient(target)
	assert.Nil(t, err)

	assert.Equal(t, fct.LogicNode(), node)

	// stream client will dail the target address, it's no easy to test
}

func TestClientStreamFactory_CreateTaskClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	go ctrl.Finish()

	handler := common.NewMockTaskServiceServer(ctrl)

	factory := NewClientStreamFactory(models.Node{IP: "127.0.0.2", Port: 9000})
	target := models.Node{IP: "127.0.0.1", Port: 9000}

	client, err := factory.CreateTaskClient(target)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	grpcServer := NewGRPCServer(":9000")
	common.RegisterTaskServiceServer(grpcServer.GetServer(), handler)
	go func() {
		_ = grpcServer.Start()
	}()

	// wait server start finish
	time.Sleep(10 * time.Millisecond)

	_, _ = factory.CreateTaskClient(target)

	time.Sleep(10 * time.Millisecond)
	grpcServer.Stop()
}
