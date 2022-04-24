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
	"github.com/lindb/lindb/metrics"
)

//go:generate mockgen  -destination=./server_stream_mock.go -package=conntrack google.golang.org/grpc ServerStream

// GRPCServerTracker represents a collection of lin-metrics to be collected
// for a gRPC server.
type GRPCServerTracker struct {
	r          *linmetric.Registry
	statistics *metrics.GRPCUnaryStatistics
}

// NewGRPCServerTracker returns a metric tracker for grpc server.
func NewGRPCServerTracker(r *linmetric.Registry) *GRPCServerTracker {
	return &GRPCServerTracker{
		r:          r,
		statistics: metrics.NewGRPCUnaryServerStatistics(r),
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor for tracking Unary RPCs.
func (tracker *GRPCServerTracker) UnaryServerInterceptor() func(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		serviceName, methodName := splitMethodName(info.FullMethod)
		resp, err := handler(ctx, req)
		tracker.statistics.Duration.WithTagValues(serviceName, methodName).UpdateSince(start)
		if err != nil {
			tracker.statistics.Failures.WithTagValues(serviceName, methodName).Incr()
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
			statistics:   metrics.NewGRPCStreamServerStatistics(tracker.r, string(streamRPCType(info)), serviceName, methodName),
		})
	}
}

// wrappedServerStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type wrappedServerStream struct {
	grpc.ServerStream
	statistics *metrics.GRPCStreamStatistics
}

func (s *wrappedServerStream) SendMsg(m interface{}) error {
	start := time.Now()
	err := s.ServerStream.SendMsg(m)
	s.statistics.MsgSentDuration.UpdateSince(start)
	if err != nil {
		s.statistics.MsgSentFailures.Incr()
	}
	return err
}

func (s *wrappedServerStream) RecvMsg(m interface{}) error {
	start := time.Now()
	err := s.ServerStream.RecvMsg(m)
	s.statistics.MsgReceivedDuration.UpdateSince(start)
	if err != nil {
		s.statistics.MsgReceivedFailures.Incr()
	}
	return err
}
