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
	"time"

	"google.golang.org/grpc"

	"github.com/lindb/lindb/internal/linmetric"
)

// GRPCClientTracker represents a collection of lin-metrics to be collected
// for a gRPC client conn.
type GRPCClientTracker struct {
	clientStreamMsgReceivedTimerVec *linmetric.DeltaHistogramVec
	clientStreamMsgSentTimerVec     *linmetric.DeltaHistogramVec
	clientStreamMsgReceivedVec      *linmetric.DeltaCounterVec
	clientStreamMsgSentVec          *linmetric.DeltaCounterVec
}

// NewGRPCClientTracker returns a metric tracker for grpc client.
func NewGRPCClientTracker(r *linmetric.Registry) *GRPCClientTracker {
	tracker := &GRPCClientTracker{}
	grpcClientScope := r.NewScope("lindb.traffic.grpc_client")
	tracker.clientStreamMsgReceivedVec = grpcClientScope.NewCounterVec(
		"msg_received", "grpc_type", "grpc_service", "grpc_method")
	tracker.clientStreamMsgSentVec = grpcClientScope.NewCounterVec(
		"msg_sent", "grpc_type", "grpc_service", "grpc_method")
	tracker.clientStreamMsgReceivedTimerVec = grpcClientScope.Scope("msg_received_duration").
		NewHistogramVec("grpc_type", "grpc_service", "grpc_method")
	tracker.clientStreamMsgSentTimerVec = grpcClientScope.Scope("msg_sent_duration").
		NewHistogramVec("grpc_type", "grpc_service", "grpc_method")
	return tracker
}

// UnaryClientInterceptor is a gRPC client-side interceptor for tracking Unary RPCs
func (tracker *GRPCClientTracker) UnaryClientInterceptor() func(
	ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		serviceName, methodName := splitMethodName(method)

		tracker.clientStreamMsgSentVec.WithTagValues(string(Unary), serviceName, methodName).Incr()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			tracker.clientStreamMsgReceivedVec.WithTagValues(string(Unary), serviceName, methodName).Incr()
		}
		return err
	}
}

// StreamClientInterceptor is a gRPC client-side interceptor  for tracking Streaming RPCs
func (tracker *GRPCClientTracker) StreamClientInterceptor() func(
	ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption,
) (grpc.ClientStream, error) {

	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		grpcType := string(clientStreamType(desc))
		serviceName, methodName := splitMethodName(method)
		return &wrappedClientStream{
			ClientStream:                 clientStream,
			clientStreamMsgReceived:      tracker.clientStreamMsgReceivedVec.WithTagValues(grpcType, serviceName, methodName),
			clientStreamMsgSent:          tracker.clientStreamMsgSentVec.WithTagValues(grpcType, serviceName, methodName),
			clientStreamMsgReceivedTimer: tracker.clientStreamMsgReceivedTimerVec.WithTagValues(grpcType, serviceName, methodName),
			clientStreamMsgSentTimer:     tracker.clientStreamMsgSentTimerVec.WithTagValues(grpcType, serviceName, methodName),
		}, nil
	}
}

// wrappedClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type wrappedClientStream struct {
	grpc.ClientStream
	clientStreamMsgReceived      *linmetric.BoundCounter
	clientStreamMsgSent          *linmetric.BoundCounter
	clientStreamMsgReceivedTimer *linmetric.BoundHistogram
	clientStreamMsgSentTimer     *linmetric.BoundHistogram
}

func (s *wrappedClientStream) SendMsg(m interface{}) error {
	start := time.Now()
	err := s.ClientStream.SendMsg(m)
	s.clientStreamMsgSentTimer.UpdateSince(start)
	if err == nil {
		s.clientStreamMsgSent.Incr()
	}
	return err
}

func (s *wrappedClientStream) RecvMsg(m interface{}) error {
	start := time.Now()
	err := s.ClientStream.RecvMsg(m)
	s.clientStreamMsgReceivedTimer.UpdateSince(start)
	if err == nil {
		s.clientStreamMsgReceived.Incr()
	}
	return err
}
