package rpc

import (
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"

	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./server.go -destination=./server_mock.go -package=rpc

// TCPServer represents a tcp tcpServer using grpc
type TCPServer interface {
	// Start starts tcp tcpServer
	Start() error
	// Stops stops tpc tcpServer
	Stop()
}

type TCPHandler interface {
	Handle(conn net.Conn) error
}

// tcpServer represents grpc tcpServer
type tcpServer struct {
	bindAddress string
	handler     TCPHandler
	lis         net.Listener
	logger      *logger.Logger
	onceClose   sync.Once
	inWorking   int32
	inShutDown  int32
}

// NewTCPServer creates the tcp tcpServer
func NewTCPServer(bindAddress string, handler TCPHandler) TCPServer {
	return &tcpServer{
		bindAddress: bindAddress,
		handler:     handler,
		logger:      logger.GetLogger("rpc", "TCPServer"),
	}
}

// Start listens the bind address and serves grpc tcpServer, block the caller
func (s *tcpServer) Start() error {
	lis, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}

	s.lis = lis
	// working now
	atomic.StoreInt32(&s.inWorking, 1)
	s.logger.Info("TCPServer start serving", logger.String("address", s.bindAddress))

	for {
		// Listen for an incoming connection.
		conn, err := lis.Accept()
		if err != nil {
			// has been shutdown
			if atomic.LoadInt32(&s.inShutDown) != 0 {
				return nil
			}
			s.logger.Error("TPCServer error when accepting", logger.Error(err))
			return err
		}
		//s.logger.Info("accept")

		// Handle connections in a new goroutine.
		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					s.logger.Error("close tcp conn err", logger.Error(err))
				}
			}()

			if err := s.handler.Handle(conn); err != nil {
				s.logger.Error("handler tcp conn err", logger.Error(err))
			}
		}()
	}
}

// Stop stops the grpc tcpServer
func (s *tcpServer) Stop() {
	s.onceClose.Do(
		func() {
			// not working
			if atomic.LoadInt32(&s.inWorking) == 0 {
				return
			}
			// shutdown
			atomic.StoreInt32(&s.inShutDown, 1)
			err := s.lis.Close()
			if err != nil {
				s.logger.Error("close TCPServer error", logger.Error(err))
			}
		})
}

type GRPCServer interface {
	TCPServer
	// GetServer returns the grpc tcpServer
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
