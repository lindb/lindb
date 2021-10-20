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

package handler

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
)

// ReplicaHandler implements replica.ReplicaServiceServer interface for handling replica rpc request.
type ReplicaHandler struct {
	walMgr replica.WriteAheadLogManager

	logger *logger.Logger
}

// NewReplicaHandler creates a replica handler.
func NewReplicaHandler(
	walMgr replica.WriteAheadLogManager,
) *ReplicaHandler {
	return &ReplicaHandler{
		walMgr: walMgr,
		logger: logger.GetLogger("storage", "ReplicaRPC"),
	}
}

// GetReplicaAckIndex returns current replica ack index.
func (r *ReplicaHandler) GetReplicaAckIndex(_ context.Context,
	request *protoReplicaV1.GetReplicaAckIndexRequest,
) (*protoReplicaV1.GetReplicaAckIndexResponse, error) {
	p, err := r.getOrCreatePartition(
		request.Database,
		models.ShardID(request.Shard),
		request.FamilyTime,
		models.NodeID(request.Leader))
	if err != nil {
		r.logger.Error("get or create wal partition err, when do get replica ack index", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &protoReplicaV1.GetReplicaAckIndexResponse{
		AckIndex: p.ReplicaAckIndex(),
	}, nil
}

// Reset resets replica index.
func (r *ReplicaHandler) Reset(_ context.Context,
	request *protoReplicaV1.ResetIndexRequest,
) (*protoReplicaV1.ResetIndexResponse, error) {
	p, err := r.getOrCreatePartition(
		request.Database,
		models.ShardID(request.Shard),
		request.FamilyTime,
		models.NodeID(request.Leader))
	if err != nil {
		r.logger.Error("get or create wal partition err, when do reset replica index", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	p.ResetReplicaIndex(request.AppendIndex)
	return &protoReplicaV1.ResetIndexResponse{}, nil
}

// Replica does replica request, and writes data.
func (r *ReplicaHandler) Replica(server protoReplicaV1.ReplicaService_ReplicaServer) error {
	replicaState, err := r.getReplicaStateFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get replica state err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := r.getOrCreatePartition(
		replicaState.Database,
		replicaState.ShardID,
		replicaState.FamilyTime,
		replicaState.Leader)
	if err != nil {
		r.logger.Error("get or create wal partition err, when do replica", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	err = p.BuildReplicaForFollower(replicaState.Leader, replicaState.Follower)
	if err != nil {
		r.logger.Error("build replica replica err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	r.logger.Info("build replica stream channel successful", logger.String("replica", replicaState.String()))
	// handle replica request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("receive replica request err", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp := &protoReplicaV1.ReplicaResponse{}
		r.logger.Debug("receive write ahead log replica log",
			logger.Any("from", replicaState.Leader), logger.Int64("index", req.ReplicaIndex))
		// write replica wal log
		appendedIdx, err := p.ReplicaLog(req.ReplicaIndex, req.Record)

		resp.ReplicaIndex = req.ReplicaIndex
		resp.AckIndex = appendedIdx

		if err != nil {
			resp.Err = err.Error()
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// getReplicaStateFromCtx gets replica relationship metadata from rpc context.
func (r *ReplicaHandler) getReplicaStateFromCtx(ctx context.Context) (replicatorState models.ReplicaState, err error) {
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

// getOrCreatePartition returns write ahead log's partition if exist, else creates a new partition.
func (r *ReplicaHandler) getOrCreatePartition(
	database string,
	shardID models.ShardID,
	familyTime int64,
	leader models.NodeID,
) (replica.Partition, error) {
	wal := r.walMgr.GetOrCreateLog(database)
	p, err := wal.GetOrCreatePartition(shardID, familyTime, leader)
	if err != nil {
		return nil, err
	}
	return p, nil
}
