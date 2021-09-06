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
	"fmt"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/conntrack"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./server.go -destination=./server_mock.go -package=rpc

type GRPCServer interface {
	// Start starts grpc server
	Start() error
	// Stop stops grpc server
	Stop()
	// GetServer returns the grpc server
	GetServer() *grpc.Server
}

type grpcServer struct {
	bindAddress string
	logger      *logger.Logger
	gs          *grpc.Server
}

func NewGRPCServer(cfg config.GRPC) GRPCServer {
	log := logger.GetLogger("rpc", "GRPCServer")
	grpcServerTracker := conntrack.NewGRPCServerTracker()
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			//TODO add metric
			log.Error("panic trigger when handle rpc request", logger.Any("err", p), logger.Stack())
			return status.Errorf(codes.Internal, "panic triggered: %v", p)
		}),
	}
	return &grpcServer{
		logger:      log,
		bindAddress: fmt.Sprintf(":%d", cfg.Port),
		gs: grpc.NewServer(
			grpc.ConnectionTimeout(cfg.ConnectTimeout.Duration()),
			grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
				grpcServerTracker.StreamServerInterceptor(),
				grpcrecovery.StreamServerInterceptor(opts...),
			)),
			grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
				grpcServerTracker.UnaryServerInterceptor(),
				grpcrecovery.UnaryServerInterceptor(opts...),
			)),
			grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams),
		),
	}
}

// Start listens the bind address and serves grpc tcpServer,
// block the caller, return fatal error or non-nil error if server is not stop gracefully.
func (s *grpcServer) Start() error {
	lis, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}

	s.logger.Info("GRPCServer start serving", logger.String("address", s.bindAddress))

	return s.gs.Serve(lis)
}

// GetServer returns the grpc tcpServer
func (s *grpcServer) GetServer() *grpc.Server {
	return s.gs
}

// Stop stops the grpc tcpServer immediately, will cause Start() return non-nil error.
func (s *grpcServer) Stop() {
	// Gracefully stop will wait for all the connection close, not we want.
	s.gs.Stop()
}
