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

package conntrack

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/internal/linmetric"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
)

type testGRPCClientTracker struct {
	serverListener net.Listener
	server         *grpc.Server
	clientConn     *grpc.ClientConn
	testClient     protoCommonV1.TaskServiceClient
	ctx            context.Context
	cancel         context.CancelFunc
}

func (tracker *testGRPCClientTracker) prepare(t *testing.T) {
	var err error

	tracker.serverListener, err = net.Listen("tcp", "127.0.0.1:23423")
	assert.NoErrorf(t, err, "failed to listen on 23423")
	serverTracker := NewGRPCServerTracker(linmetric.StorageRegistry)
	tracker.server = grpc.NewServer(
		grpc.StreamInterceptor(serverTracker.StreamServerInterceptor()),
		grpc.UnaryInterceptor(serverTracker.UnaryServerInterceptor()),
	)

	up := make(chan struct{})
	go func() {
		up <- struct{}{}
		_ = tracker.server.Serve(tracker.serverListener)
	}()
	<-up
	time.Sleep(time.Second)

	clientTracker := NewGRPCClientTracker(linmetric.BrokerRegistry)
	tracker.clientConn, err = grpc.Dial(
		tracker.serverListener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(clientTracker.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(clientTracker.StreamClientInterceptor()),
	)
	assert.NoErrorf(t, err, "failed dialing to testGRPCClientTracker")

	tracker.ctx, tracker.cancel = context.WithTimeout(context.Background(), 2*time.Second)
}

func (tracker *testGRPCClientTracker) shutdown() {
	if tracker.cancel != nil {
		defer tracker.cancel()
	}
	if tracker.serverListener != nil {
		tracker.server.Stop()
		_ = tracker.serverListener.Close()
	}
	if tracker.clientConn != nil {
		_ = tracker.clientConn.Close()
	}
}

func Test_GRPC(t *testing.T) {
	tracker := testGRPCClientTracker{}
	tracker.prepare(t)
	defer tracker.shutdown()

	tracker.testClient = protoCommonV1.NewTaskServiceClient(tracker.clientConn)
	client, err := tracker.testClient.Handle(tracker.ctx)

	assert.NoError(t, err)
	assert.Nil(t, client.Send(&protoCommonV1.TaskRequest{}))
	_, err = client.Recv()
	assert.Error(t, err)
}
