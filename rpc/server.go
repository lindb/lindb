package rpc

import (
	"io"
	"net"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc/pkg/batch"
	"github.com/eleme/lindb/rpc/pkg/common"
)

// Server defines methods should be implemented by RPC service, any error the methods will result in a fatal
type Server interface {
	Init()
	Listen()
	Register()
	Serve()
}

// Handler is a func to handle batch request and return batch response
type Handler func(request *batch.BatchRequest_Request) (*batch.BatchResponse_Response, error)

// BaseServer handlers StreamBatchRequest and dispatched request to registered handler according to requestType
type BaseServer struct {
	sync.Mutex
	// binding ": + port"
	port     string
	Gserver  *grpc.Server
	listener net.Listener
	// registered RequestType Handler pair
	handlerMap map[common.RequestType]Handler
}

// NewBaseServer returns a BaseServer
func NewBaseServer(port string) *BaseServer {
	return &BaseServer{
		port:       port,
		handlerMap: make(map[common.RequestType]Handler),
	}
}

// Listen declares a local address on port
func (s *BaseServer) Listen() {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		logger.GetLogger().Fatal("failed to listen",
			zap.String("port", s.port),
			zap.Stack("stack"),
			zap.Error(err))
	}
	s.listener = lis

	s.Gserver = grpc.NewServer()

	batch.RegisterBatchServiceServer(s.Gserver, s)

	logger.GetLogger().Info("register server")
}

// Serve serves on declared port
func (s *BaseServer) Serve() {
	if err := s.Gserver.Serve(s.listener); err != nil {
		logger.GetLogger().Fatal("failed to serve",
			zap.String("port", s.port),
			zap.Stack("stack"),
			zap.Error(err))
	}
}

// RegisterHandler registers a Handler for a requestType
func (s *BaseServer) RegisterHandler(requestType common.RequestType, handler Handler) {
	if handler == nil {
		logger.GetLogger().Fatal("handler must not be nil", zap.Stack("stack"))
	}
	s.Lock()
	defer s.Unlock()

	_, ok := s.handlerMap[requestType]
	if ok {
		logger.GetLogger().Fatal("Handler already registered",
			zap.String("requestType", requestType.String()),
			zap.Stack("stack"))
	}
	s.handlerMap[requestType] = handler
}

// StreamBatchRequest implements grpc method to handle stream batch request
func (s *BaseServer) StreamBatchRequest(stream batch.BatchService_StreamBatchRequestServer) error {
	for {
		in, err := stream.Recv()
		logger.GetLogger().Debug("receive stream request")
		if err == io.EOF {
			return nil
		}

		if err != nil {
			logger.GetLogger().Error("receive stream request error", zap.Stack("stack"), zap.Error(err))
			return err
		}

		resps := s.handleRequest(in)

		if err = stream.Send(resps); err != nil {
			return err
		}
	}
}

// handleRequest dispatches requests to handlers according to requestType
func (s *BaseServer) handleRequest(batchRequest *batch.BatchRequest) *batch.BatchResponse {
	length := len(batchRequest.RequestIDs)
	resps := make([]*batch.BatchResponse_Response, length, length)

	for i := range batchRequest.RequestIDs {
		req := batchRequest.Requests[i]
		handler, ok := s.handlerMap[req.RequestTyp]
		if !ok {

			resps[i] = &batch.BatchResponse_Response{
				RequestType: req.RequestTyp,
				Response:    nil,
			}
			continue
		}

		res, _ := handler(req)
		resps[i] = res
	}

	return &batch.BatchResponse{
		RequestIDs: batchRequest.RequestIDs,
		Responses:  resps,
	}
}
