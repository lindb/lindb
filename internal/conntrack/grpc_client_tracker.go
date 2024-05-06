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

//go:generate mockgen -destination=./client_stream_mock.go -package=conntrack google.golang.org/grpc ClientStream

// GRPCClientTracker represents a collection of lin-metrics to be collected
// for a gRPC client conn.
type GRPCClientTracker struct {
	r          *linmetric.Registry
	statistics *metrics.GRPCUnaryStatistics
}

// NewGRPCClientTracker returns a metric tracker for grpc client.
func NewGRPCClientTracker(r *linmetric.Registry) *GRPCClientTracker {
	return &GRPCClientTracker{
		r:          r,
		statistics: metrics.NewGRPCUnaryClientStatistics(r),
	}
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
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		tracker.statistics.Duration.WithTagValues(serviceName, methodName).UpdateSince(start)
		if err != nil {
			tracker.statistics.Failures.WithTagValues(serviceName, methodName).Incr()
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
			ClientStream: clientStream,
			statistics:   metrics.NewGRPCStreamClientStatistics(tracker.r, grpcType, serviceName, methodName),
		}, nil
	}
}

// wrappedClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type wrappedClientStream struct {
	grpc.ClientStream

	statistics *metrics.GRPCStreamStatistics
}

func (s *wrappedClientStream) SendMsg(m interface{}) error {
	start := time.Now()
	err := s.ClientStream.SendMsg(m)
	s.statistics.MsgSentDuration.UpdateSince(start)
	if err != nil {
		s.statistics.MsgSentFailures.Incr()
	}
	return err
}

func (s *wrappedClientStream) RecvMsg(m interface{}) error {
	start := time.Now()
	err := s.ClientStream.RecvMsg(m)
	s.statistics.MsgReceivedDuration.UpdateSince(start)
	if err != nil {
		s.statistics.MsgReceivedFailures.Incr()
	}
	return err
}
