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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

func TestGRPCServerTracker_UnaryServerInterceptor(t *testing.T) {
	tracker := NewGRPCServerTracker(linmetric.BrokerRegistry)
	fn := tracker.UnaryServerInterceptor()
	// handle failure
	_, err := fn(context.TODO(), nil, &grpc.UnaryServerInfo{FullMethod: "servier/method"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, fmt.Errorf("err")
		})
	assert.Error(t, err)

	// handle ok
	_, err = fn(context.TODO(), nil, &grpc.UnaryServerInfo{FullMethod: "servier/method"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
	assert.NoError(t, err)
}

func TestGRPCServerTracker_StreamServerInterceptor(t *testing.T) {
	tracker := NewGRPCServerTracker(linmetric.BrokerRegistry)
	fn := tracker.StreamServerInterceptor()
	err := fn(nil, nil, &grpc.StreamServerInfo{FullMethod: "servier/methdo"},
		func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		})
	assert.NoError(t, err)
}

func TestWrappedServerStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ss := NewMockServerStream(ctrl)
	wrap := &wrappedServerStream{
		ServerStream: ss,
		statistics: metrics.NewGRPCStreamServerStatistics(linmetric.BrokerRegistry,
			"type", "service", "method"),
	}
	msg := "msg"
	// receive msg ok
	ss.EXPECT().RecvMsg(gomock.Any()).Return(nil)
	err := wrap.RecvMsg(&msg)
	assert.NoError(t, err)
	// receive msg failure
	ss.EXPECT().RecvMsg(gomock.Any()).Return(fmt.Errorf("err"))
	err = wrap.RecvMsg(&msg)
	assert.Error(t, err)

	// sent msg ok
	ss.EXPECT().SendMsg(gomock.Any()).Return(nil)
	err = wrap.SendMsg(&msg)
	assert.NoError(t, err)
	// sent msg failure
	ss.EXPECT().SendMsg(gomock.Any()).Return(fmt.Errorf("err"))
	err = wrap.SendMsg(&msg)
	assert.Error(t, err)
}
