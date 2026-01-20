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
	context "context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/pipeline"
)

type ResultSetService struct {
	protoCommandV1.UnimplementedResultSetServiceServer

	logger logger.Logger
}

func NewResultSetService() protoCommandV1.ResultSetServiceServer {
	return &ResultSetService{
		logger: logger.GetLogger("RPC", "resultSet"),
	}
}

func (srv *ResultSetService) ResultSet(ctx context.Context,
	request *protoCommandV1.ResultSetRequest,
) (*protoCommandV1.ResultSetResponse, error) {
	resultSet := &model.TaskResultSet{}
	if err := encoding.JSONUnmarshal(request.Payload, resultSet); err != nil {
		// TODO: send handle error?
		return nil, err
	}

	srv.logger.Info("receive task result set", logger.Any("requestID", resultSet.TaskID.RequestID),
		logger.Int("TaskID", resultSet.TaskID.ID), logger.Int("nodeID", int(resultSet.Node)))

	sourceOperator := pipeline.DriverManager.GetSourceOperator(resultSet.TaskID, resultSet.Node)
	if sourceOperator != nil {
		if len(resultSet.Page) != 0 {
			page, err := types.UnmarshalPage(resultSet.Page)
			if err != nil {
				fmt.Println("unmarshal page error:", err)
				panic(err)
			}
			sourceOperator.Receive(page)
		}
		// FIXME: handle error
		if resultSet.NoMore {
			// current task no more splits
			sourceOperator.Complete()
		}
	} else {
		srv.logger.Warn("source operator not found", logger.Any("requestID", resultSet.TaskID.RequestID),
			logger.Int("TaskID", resultSet.TaskID.ID), logger.Int("nodeID", int(resultSet.Node)))
	}
	return &protoCommandV1.ResultSetResponse{}, nil
}
