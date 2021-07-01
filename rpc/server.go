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
	"net"

	"google.golang.org/grpc"

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

func NewGRPCServer(bindAddress string) GRPCServer {
	return &grpcServer{
		bindAddress: bindAddress,
		logger:      logger.GetLogger("rpc", "GRPCServer"),
		gs:          grpc.NewServer(),
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
