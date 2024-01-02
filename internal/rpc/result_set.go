package rpc

import (
	context "context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/execution/pipeline"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/spi"
)

type ResultSetService struct {
	logger logger.Logger
}

func NewResultSetService() protoCommandV1.ResultSetServiceServer {
	return &ResultSetService{
		logger: logger.GetLogger("RPC", "resultSet"),
	}
}

func (srv *ResultSetService) ResultSet(ctx context.Context, request *protoCommandV1.ResultSetRequest) (*protoCommandV1.ResultSetResponse, error) {
	resultSet := &model.TaskResultSet{}
	if err := encoding.JSONUnmarshal(request.Payload, resultSet); err != nil {
		return nil, err
	}

	srv.logger.Debug("receive task result set", logger.String("requestID", resultSet.TaskID.RequestID),
		logger.Int("TaskID", resultSet.TaskID.ID), logger.Int("nodeID", int(resultSet.NodeID)))

	fmt.Println(string(request.Payload))

	sourceOperator := pipeline.DriverManager.GetSourceOperator(resultSet.TaskID, resultSet.NodeID)
	if sourceOperator != nil {
		sourceOperator.AddSplit(&spi.BinarySplit{
			Page: resultSet.Page,
		})

		if resultSet.NoMore {
			// current task no more splits
			sourceOperator.NoMoreSplits()
		}
	}
	return &protoCommandV1.ResultSetResponse{}, nil
}
