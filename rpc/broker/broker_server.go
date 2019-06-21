package broker

import (
	"context"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/pkg/batch"
	brokerpb "github.com/eleme/lindb/rpc/pkg/broker"
	"github.com/eleme/lindb/rpc/pkg/common"
)

type Server struct {
	rpc.BaseServer
}

func NewBrokerServer(address string) *Server {
	return &Server{
		BaseServer: *rpc.NewBaseServer(address),
	}
}

func (ws *Server) Register() {
	logger.GetLogger().Info("Register BrokerServiceServer")
	brokerpb.RegisterBrokerServiceServer(ws.Gserver, ws)
}

func (ws *Server) Init() {
	ws.RegisterHandler(common.RequestType_WritePoints,
		func(request *batch.BatchRequest_Request) (response *batch.BatchResponse_Response, e error) {
			res, err := ws.WritePoints(context.TODO(), request.GetWritePoints())
			return &batch.BatchResponse_Response{
				RequestType: common.RequestType_WritePoints,
				Response:    &batch.BatchResponse_Response_WritePoints{WritePoints: res},
			}, err
		})
}

// WritePoints implements grpc method to write points to memory or disk
func (ws *Server) WritePoints(
	ctx context.Context,
	request *brokerpb.WritePointsRequest) (*brokerpb.WritePointsResponse, error) {

	logger.GetLogger().Info("receive points",
		zap.Any("points", request.Points))

	return &brokerpb.WritePointsResponse{
		Context: rpc.BuildResponseContext(""),
	}, nil
}
