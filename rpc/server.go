package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/lindb/lindb/pkg/logger"
)

// TCPServer represents a tcp server using grpc
type TCPServer interface {
	// Start starts tcp server
	Start() error
	// GetServer returns the grpc server
	GetServer() *grpc.Server
	// Stops stops tpc server
	Stop()
}

// server represents grpc server
type server struct {
	bindAddress string
	gs          *grpc.Server

	logger *logger.Logger
}

// NewTCPServer creates the tcp server
func NewTCPServer(bindAddress string) TCPServer {
	return &server{
		bindAddress: bindAddress,
		gs:          grpc.NewServer(),
		logger:      logger.GetLogger("rpc", "Server"),
	}
}

// Start listens the bind address and serves grpc server
func (s *server) Start() error {
	lis, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}

	s.logger.Info("rpc server start serving", logger.String("address", s.bindAddress))
	return s.gs.Serve(lis)
}

// GetServer returns the grpc server
func (s *server) GetServer() *grpc.Server {
	return s.gs
}

// Stop stops the grpc server
func (s *server) Stop() {
	s.gs.Stop()
}
