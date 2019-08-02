package handler

import (
	"io"

	"go.uber.org/zap"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/storage"
)

type Query struct {
	fct    rpc.ServerStreamFactory
	logger *logger.Logger
}

func NewQuery(fct rpc.ServerStreamFactory) *Query {
	return &Query{
		fct:    fct,
		logger: logger.GetLogger("handler/query"),
	}
}

func (q *Query) Query(stream storage.QueryService_QueryServer) error {
	clientLogicNode, err := rpc.GetLogicNodeFromContext(stream.Context())
	if err != nil {
		return err
	}

	q.fct.Register(*clientLogicNode, stream)
	q.logger.Info("register query stream for node:" + clientLogicNode.Indicator())

	// when return, the stream is closed, Deregister the stream
	defer func() {
		q.fct.Deregister(*clientLogicNode)
		q.logger.Info("unregister query stream for node:" + clientLogicNode.Indicator())
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			q.logger.Info("query server stream close")
			return nil
		}
		if err != nil {
			q.logger.Error("query server stream error", zap.Error(err))
			return err
		}

		q.logger.Info("recv " + req.Msg)
	}
}
