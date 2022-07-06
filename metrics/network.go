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
	Accept         *linmetric.BoundCounter // accept total count
	AcceptFailures *linmetric.BoundCounter // accept failure
	ActiveConn     *linmetric.BoundGauge   // current active connections
	Read           *linmetric.BoundCounter // read total count
	ReadBytes      *linmetric.BoundCounter // read byte size
	ReadFailures   *linmetric.BoundCounter // read failure
	Write          *linmetric.BoundCounter // write total count
	WriteBytes     *linmetric.BoundCounter // write byte size
	WriteFailures  *linmetric.BoundCounter // write failure
	Close          *linmetric.BoundCounter // close total count
	CloseFailures  *linmetric.BoundCounter // close failure
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
	MsgReceivedDuration *linmetric.BoundHistogram // receive msg duration, include receive total count/handle duration
	MsgSentDuration     *linmetric.BoundHistogram // send msg duration, include send total count
}

// GRPCServerStatistics represents grpc server statistics.
type GRPCServerStatistics struct {
	Panics *linmetric.BoundCounter // panic when grpc server handle request
}

// NewConnStatistics creates tcp connection statistics.
func NewConnStatistics(r *linmetric.Registry, addr string) *ConnStatistics {
	tcpScope := r.NewScope("lindb.traffic.tcp", "addr", addr)
	return &ConnStatistics{
		Accept:         tcpScope.NewCounter("accept_conns"),
		AcceptFailures: tcpScope.NewCounter("accept_failures"),
		ActiveConn:     tcpScope.NewGauge("active_conns"),
		Read:           tcpScope.NewCounter("reads"),
		ReadBytes:      tcpScope.NewCounter("read_bytes"),
		ReadFailures:   tcpScope.NewCounter("read_failures"),
		Write:          tcpScope.NewCounter("writes"),
		WriteBytes:     tcpScope.NewCounter("write_bytes"),
		WriteFailures:  tcpScope.NewCounter("write_failures"),
		Close:          tcpScope.NewCounter("close_conns"),
		CloseFailures:  tcpScope.NewCounter("close_failures"),
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

// NewGRPCServerStatistics creates grpc server statistics.
func NewGRPCServerStatistics(registry *linmetric.Registry) *GRPCServerStatistics {
	scope := registry.NewScope("lindb.traffic.grpc_server")
	return &GRPCServerStatistics{
		Panics: scope.NewCounter("panics"),
	}
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
