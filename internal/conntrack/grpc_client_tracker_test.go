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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

func TestGRPCClientTracker_UnaryClientInterceptor(t *testing.T) {
	tracker := NewGRPCClientTracker(linmetric.BrokerRegistry)
	fn := tracker.UnaryClientInterceptor()
	err := fn(context.TODO(), "service/method", nil, nil, nil,
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return fmt.Errorf("err")
		}, nil)
	assert.Error(t, err)

	err = fn(context.TODO(), "service/method", nil, nil, nil,
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}, nil)
	assert.NoError(t, err)
}

func TestGRPCClientTracker_StreamClientInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tracker := NewGRPCClientTracker(linmetric.BrokerRegistry)
	cliStream := NewMockClientStream(ctrl)
	fn := tracker.StreamClientInterceptor()
	cs, err := fn(context.TODO(), &grpc.StreamDesc{StreamName: "test"},
		nil, "service/method", func(ctx context.Context, desc *grpc.StreamDesc,
			cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return nil, fmt.Errorf("err")
		})
	assert.Error(t, err)
	assert.Nil(t, cs)

	cs, err = fn(context.TODO(), &grpc.StreamDesc{StreamName: "test"},
		nil, "service/method", func(ctx context.Context, desc *grpc.StreamDesc,
			cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return cliStream, nil
		})
	assert.NoError(t, err)
	assert.NotNil(t, cs)
}

func TestWrappedClientStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cliStream := NewMockClientStream(ctrl)
	wrap := &wrappedClientStream{
		ClientStream: cliStream,
		statistics: metrics.NewGRPCStreamClientStatistics(linmetric.BrokerRegistry,
			"type", "service", "method"),
	}

	msg := "msg"

	// send ok
	cliStream.EXPECT().SendMsg(gomock.Any()).Return(nil)
	err := wrap.SendMsg(&msg)
	assert.NoError(t, err)
	// send failure
	cliStream.EXPECT().SendMsg(gomock.Any()).Return(fmt.Errorf("err"))
	err = wrap.SendMsg(&msg)
	assert.Error(t, err)

	// receive ok
	cliStream.EXPECT().RecvMsg(gomock.Any()).Return(nil)
	err = wrap.RecvMsg(&msg)
	assert.NoError(t, err)
	// receive failure
	cliStream.EXPECT().RecvMsg(gomock.Any()).Return(fmt.Errorf("err"))
	err = wrap.RecvMsg(&msg)
	assert.Error(t, err)
}
