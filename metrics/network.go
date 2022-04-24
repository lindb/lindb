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

package metrics

import "github.com/lindb/lindb/internal/linmetric"

// ConnStatistics represents tpc connection statistics.
type ConnStatistics struct {
	Accept       *linmetric.BoundCounter // accept total count
	AcceptErrors *linmetric.BoundCounter // accept failure
	ActiveConn   *linmetric.BoundGauge   // current active connections
	Read         *linmetric.BoundCounter // read total count
	ReadBytes    *linmetric.BoundCounter // read byte size
	ReadErrors   *linmetric.BoundCounter // read failure
	Write        *linmetric.BoundCounter // write total count
	WriteBytes   *linmetric.BoundCounter // write byte size
	WriteErrors  *linmetric.BoundCounter // write failure
	Close        *linmetric.BoundCounter // close total count
	CloseErrors  *linmetric.BoundCounter // close failure
}

// GRPCUnaryStatistics represents unary grpc client/server statistics.
type GRPCUnaryStatistics struct {
	Failures *linmetric.DeltaCounterVec   // handle msg failure
	Duration *linmetric.DeltaHistogramVec // handle msg duration
}

// GRPCStreamStatistics represents stream grpc client/server statistics.
type GRPCStreamStatistics struct {
	MsgReceivedFailures *linmetric.BoundCounter   // receive msg failure
	MsgSentFailures     *linmetric.BoundCounter   // send msg failure
	MsgReceivedDuration *linmetric.BoundHistogram // receive msg duration, include receive total count/server handle duration
	MsgSentDuration     *linmetric.BoundHistogram // send msg duration, include send total count
}

// NewConnStatistics creates tcp connection statistics.
func NewConnStatistics(r *linmetric.Registry, addr string) *ConnStatistics {
	tcpScope := r.NewScope("lindb.traffic.tcp", "addr", addr)
	return &ConnStatistics{
		Accept:       tcpScope.NewCounter("accept_conns"),
		AcceptErrors: tcpScope.NewCounter("accept_errors"),
		ActiveConn:   tcpScope.NewGauge("active_conns"),
		Read:         tcpScope.NewCounter("read_count"),
		ReadBytes:    tcpScope.NewCounter("read_bytes"),
		ReadErrors:   tcpScope.NewCounter("read_errors"),
		Write:        tcpScope.NewCounter("write_count"),
		WriteBytes:   tcpScope.NewCounter("write_bytes"),
		WriteErrors:  tcpScope.NewCounter("write_errors"),
		Close:        tcpScope.NewCounter("close_conns"),
		CloseErrors:  tcpScope.NewCounter("close_errors"),
	}
}

// NewGRPCUnaryClientStatistics creates unary grpc client statistics.
func NewGRPCUnaryClientStatistics(registry *linmetric.Registry) *GRPCUnaryStatistics {
	return newGRPCUnaryStatistics(registry, "lindb.traffic.grpc_client.unary")
}

// NewGRPCStreamClientStatistics creates stream grpc client statistics.
func NewGRPCStreamClientStatistics(registry *linmetric.Registry, grpcType, grpcService, grpcMethod string) *GRPCStreamStatistics {
	return newGPRCStreamStatistics(registry, "lindb.traffic.grpc_client.stream", grpcType, grpcService, grpcMethod)
}

// NewGRPCUnaryServerStatistics creates unary grpc server statistics.
func NewGRPCUnaryServerStatistics(registry *linmetric.Registry) *GRPCUnaryStatistics {
	return newGRPCUnaryStatistics(registry, "lindb.traffic.grpc_server.unary")
}

// NewGRPCStreamServerStatistics creates stream grpc server statistics.
func NewGRPCStreamServerStatistics(registry *linmetric.Registry, grpcType, grpcService, grpcMethod string) *GRPCStreamStatistics {
	return newGPRCStreamStatistics(registry, "lindb.traffic.grpc_server.stream", grpcType, grpcService, grpcMethod)
}

// newGPRCStreamStatistics creates grpc client/server stream statistics.
func newGPRCStreamStatistics(registry *linmetric.Registry, name, grpcType, grpcService, grpcMethod string) *GRPCStreamStatistics {
	scope := registry.NewScope(name)
	return &GRPCStreamStatistics{
		MsgReceivedFailures: scope.NewCounterVec("msg_received_failures", "grpc_type", "grpc_service", "grpc_method").
			WithTagValues(grpcType, grpcService, grpcMethod),
		MsgSentFailures: scope.NewCounterVec("msg_sent_failures", "grpc_type", "grpc_service", "grpc_method").
			WithTagValues(grpcType, grpcService, grpcMethod),
		MsgReceivedDuration: scope.Scope("received_duration").NewHistogramVec("grpc_type", "grpc_service", "grpc_method").
			WithTagValues(grpcType, grpcService, grpcMethod),
		MsgSentDuration: scope.Scope("sent_duration").NewHistogramVec("grpc_type", "grpc_service", "grpc_method").
			WithTagValues(grpcType, grpcService, grpcMethod),
	}
}

// newGRPCUnaryStatistics creates unary grpc client/server statistics.
func newGRPCUnaryStatistics(registry *linmetric.Registry, name string) *GRPCUnaryStatistics {
	scope := registry.NewScope(name)
	return &GRPCUnaryStatistics{
		Failures: scope.NewCounterVec("failures", "grpc_service", "grpc_method"),
		Duration: scope.Scope("duration").NewHistogramVec("grpc_service", "grpc_method"),
	}
}
