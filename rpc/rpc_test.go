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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
)

var (
	node = models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	}
)

func TestClientConnFactory(t *testing.T) {
	fct := GetBrokerClientConnFactory()

	conn1, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)

	conn11, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)

	conn2, err := fct.GetClientConn(&models.StatelessNode{
		HostIP:   "1.1.1.1",
		GRPCPort: 456,
	})
	assert.NoError(t, err)

	assert.Same(t, conn1, conn11)
	assert.NotSame(t, conn1, conn2)

	// test close
	err = fct.CloseClientConn(&models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	})
	assert.NoError(t, err)
}

func TestClientStreamFactory_CreateTaskClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	go ctrl.Finish()

	handler := protoCommonV1.NewMockTaskServiceServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	factory := NewClientStreamFactory(ctx, &models.StatelessNode{HostIP: "127.0.0.2", GRPCPort: 9000}, GetStorageClientConnFactory())
	target := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 9000}

	client, err := factory.CreateTaskClient(&target)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	grpcServer := NewGRPCServer(config.GRPC{Port: 9000, ConnectTimeout: ltoml.Duration(time.Second)}, linmetric.StorageRegistry)
	protoCommonV1.RegisterTaskServiceServer(grpcServer.GetServer(), handler)
	go func() {
		_ = grpcServer.Start()
	}()

	// wait server start finish
	time.Sleep(10 * time.Millisecond)

	_, _ = factory.CreateTaskClient(&target)

	time.Sleep(10 * time.Millisecond)
	grpcServer.Stop()
}
