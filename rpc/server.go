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
	// Stops stops grpc server
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
