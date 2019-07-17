package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/eleme/lindb/pkg/logger"
)

type TCPServer interface {
	Start() error
	GetServer() *grpc.Server
	Stop()
}

type server struct {
	bindAddress string
	gs          *grpc.Server

	logger *logger.Logger
}

func NewTCPServer(bindAddress string) TCPServer {
	return &server{
		bindAddress: bindAddress,
		gs:          grpc.NewServer(),
		logger:      logger.GetLogger("rpc/server"),
	}
}

func (s *server) Start() error {
	lis, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}

	s.logger.Info("rpc server start serving")
	return s.gs.Serve(lis)
}

func (s *server) GetServer() *grpc.Server {
	return s.gs
}

func (s *server) Stop() {
	if s.gs != nil {
		s.gs.Stop()
	}
}
