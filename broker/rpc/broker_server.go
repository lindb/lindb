package rpc

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/broker"
	"github.com/eleme/lindb/rpc/proto/common"
)

type BrokerServer interface {
	Start() error
	Close()
}

type brokerSever struct {
	bindAddress string
	gs          *grpc.Server
	logger      *zap.Logger
}

func NewBrokerServer(bindAddress string) BrokerServer {
	return &brokerSever{
		bindAddress: bindAddress,
		logger:      logger.GetLogger(),
	}
}

func (bs *brokerSever) Start() error {
	lis, err := net.Listen("tcp", bs.bindAddress)
	if err != nil {
		return err
	}

	bs.gs = grpc.NewServer()

	broker.RegisterBrokerServiceServer(bs.gs, bs)

	bs.logger.Info("brokerServer start serving")
	return bs.gs.Serve(lis)
}

func (bs *brokerSever) WritePoints(ctx context.Context, request *common.Request) (*common.Response, error) {
	// todo: @XiaTianliang
	//bs.logger.Info(string(request.Data))
	return rpc.ResponseOK(), nil
}

func (bs *brokerSever) Close() {
	if bs.gs != nil {
		bs.gs.Stop()
	}
}
