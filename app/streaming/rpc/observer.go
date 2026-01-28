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
	"context"
	"fmt"
	"io"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/streaming"
)

// ObserverHandler implements replica.ReplicaServiceServer interface for handling observer rpc request.
type ObserverHandler struct {
	protoReplicaV1.UnimplementedReplicaServiceServer

	logger logger.Logger
}

// NewObserverHandler creates a observer handler.
func NewObserverHandler() protoReplicaV1.ReplicaServiceServer {
	return &ObserverHandler{
		logger: logger.GetLogger("Streaming", "observerRPC"),
	}
}

// GetReplicaAckIndex returns current replica ack index.
func (r *ObserverHandler) GetReplicaAckIndex(_ context.Context,
	request *protoReplicaV1.GetReplicaAckIndexRequest,
) (*protoReplicaV1.GetReplicaAckIndexResponse, error) {
	// TODO: persist observer ack index??
	return &protoReplicaV1.GetReplicaAckIndexResponse{
		AckIndex: request.CurrentIndex, // NOTE: just return current index.
	}, nil
}

// ResetIndex resets replica index.
func (r *ObserverHandler) ResetIndex(_ context.Context,
	request *protoReplicaV1.ResetIndexRequest,
) (*protoReplicaV1.ResetIndexResponse, error) {
	// just return successful response
	return &protoReplicaV1.ResetIndexResponse{}, nil
}

// Replica does replica request, and writes data.
func (r *ObserverHandler) Replica(server protoReplicaV1.ReplicaService_ReplicaServer) error {
	replicaState, err := r.getReplicaStateFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get observer state err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ds, ok := streaming.GetManager().GetDataSource(replicaState.Database)
	if !ok {
		r.logger.Error("data source for database not found", logger.String("database", replicaState.Database))
		return status.Error(codes.NotFound, fmt.Sprintf("data source for database %s not found", replicaState.Database))
	}

	r.logger.Info("build observer stream channel successful", logger.String("observer", replicaState.String()))
	// handle replica request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("receive observer request err", logger.String("observer", replicaState.String()), logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp := &protoReplicaV1.ReplicaResponse{}

		if err := ds.Produce(req.Record); err != nil {
			r.logger.Error("publish event err", logger.String("observer", replicaState.String()), logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp.ReplicaIndex = req.ReplicaIndex
		resp.AckIndex = req.ReplicaIndex

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// getReplicaStateFromCtx gets replica relationship metadata from rpc context.
func (r *ObserverHandler) getReplicaStateFromCtx(ctx context.Context) (replicatorState models.ReplicaState, err error) {
	replicaStateData, err := rpc.GetStringFromContext(ctx, constants.RPCMetaReplicaState)
	if err != nil {
		return
	}
	err = encoding.JSONUnmarshal([]byte(replicaStateData), &replicatorState)
	if err != nil {
		return
	}
	return
}
