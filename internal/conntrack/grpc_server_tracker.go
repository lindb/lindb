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

	"google.golang.org/grpc"

	"github.com/lindb/lindb/internal/linmetric"
)

// GRPCServerTracker represents a collection of lin-metrics to be collected
// for a gRPC server.
type GRPCServerTracker struct {
	streamMsgReceivedVec *linmetric.DeltaCounterVec
	streamMsgSentVec     *linmetric.DeltaCounterVec
}

// NewGRPCServerTracker returns a metric tracker for grpc server.
func NewGRPCServerTracker() *GRPCServerTracker {
	tracker := &GRPCServerTracker{}
	grpcServerScope := linmetric.NewScope("lindb.traffic.grpc_server")
	tracker.streamMsgReceivedVec = grpcServerScope.NewDeltaCounterVec(
		"msg_received", "grpc_type", "grpc_service", "grpc_method")
	tracker.streamMsgSentVec = grpcServerScope.NewDeltaCounterVec(
		"msg_sent", "grpc_type", "grpc_service", "grpc_method")
	return tracker
}

// UnaryServerInterceptor is a gRPC server-side interceptor for tracking Unary RPCs.
func (tracker *GRPCServerTracker) UnaryServerInterceptor() func(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		serviceName, methodName := splitMethodName(info.FullMethod)
		tracker.streamMsgReceivedVec.WithTagValues(string(Unary), serviceName, methodName).Incr()
		resp, err := handler(ctx, req)
		if err == nil {
			tracker.streamMsgSentVec.WithTagValues(string(Unary), serviceName, methodName).Incr()
		}
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor for tracking Streaming RPCs.
func (tracker *GRPCServerTracker) StreamServerInterceptor() func(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		serviceName, methodName := splitMethodName(info.FullMethod)

		return handler(srv, &wrappedServerStream{
			ServerStream: ss,
			serverStreamMsgReceived: tracker.streamMsgReceivedVec.WithTagValues(
				string(streamRPCType(info)), serviceName, methodName),
			serverStreamMsgSent: tracker.streamMsgSentVec.WithTagValues(
				string(streamRPCType(info)), serviceName, methodName),
		})
	}
}

// wrappedServerStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type wrappedServerStream struct {
	grpc.ServerStream
	serverStreamMsgReceived *linmetric.BoundDeltaCounter
	serverStreamMsgSent     *linmetric.BoundDeltaCounter
}

func (s *wrappedServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.serverStreamMsgSent.Incr()
	}
	return err
}

func (s *wrappedServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.serverStreamMsgReceived.Incr()
	}
	return err
}
