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
	"github.com/lindb/lindb/tsdb"
)

// ReplicaHandler implements replica.ReplicaServiceServer interface for handling replica rpc request.
type ReplicaHandler struct {
	walMgr replica.WriteAheadLogManager
	engine tsdb.Engine

	logger *logger.Logger
}

// NewReplicaHandler creates a replica handler.
func NewReplicaHandler(
	walMgr replica.WriteAheadLogManager,
	engine tsdb.Engine,
) *ReplicaHandler {
	return &ReplicaHandler{
		walMgr: walMgr,
		engine: engine,
		logger: logger.GetLogger("storage", "ReplicaRPC"),
	}
}

// GetReplicaAckIndex returns current replica ack index.
func (r *ReplicaHandler) GetReplicaAckIndex(ctx context.Context,
	req *protoReplicaV1.GetReplicaAckIndexRequest,
) (*protoReplicaV1.GetReplicaAckIndexResponse, error) {
	panic("implement me")
}

// Reset resets replica index.
func (r *ReplicaHandler) Reset(ctx context.Context,
	request *protoReplicaV1.ResetIndexRequest,
) (*protoReplicaV1.ResetIndexResponse, error) {
	panic("implement me")
}

// Replica does replica request, and writes data.
func (r *ReplicaHandler) Replica(server protoReplicaV1.ReplicaService_ReplicaServer) error {
	database, shardID, leader, follower, err := r.getFollowerInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := r.getOrCreatePartition(database, shardID)
	if err != nil {
		r.logger.Error("create wal partition err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	err = p.BuildReplicaForFollower(leader, follower)
	if err != nil {
		r.logger.Error("build replica replica err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	// handle write request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("get write request err", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp := &protoReplicaV1.ReplicaResponse{}
		// write replica wal log
		appendedIdx, err := p.ReplicaLog(req.ReplicaIndex, req.Record)

		resp.ReplicaIndex = appendedIdx

		if err != nil {
			resp.Err = err.Error()
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// Write does metric write request.
func (r *ReplicaHandler) Write(server protoReplicaV1.ReplicaService_WriteServer) error {
	database, shardState, err := r.getReplicasInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	//TODO need check leader?
	if len(shardState.Replica.Replicas) == 0 {
		return status.Error(codes.InvalidArgument, "replicas cannot be empty")
	}

	p, err := r.getOrCreatePartition(database, shardState.ID)
	if err != nil {
		r.logger.Error("create wal partition err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	err = p.BuildReplicaForLeader(shardState.Leader, shardState.Replica.Replicas)
	if err != nil {
		r.logger.Error("build replica replica err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}

	// handle write request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("get write request err", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp := &protoReplicaV1.WriteResponse{}
		// write wal log
		err = p.WriteLog(req.Record)

		if err != nil {
			resp.Err = err.Error()
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// getReplicasInfoFromCtx gets shard replica metadata from rpc context.
func (r *ReplicaHandler) getReplicasInfoFromCtx(ctx context.Context) (database string, shardState models.ShardState, err error) {
	database, err = rpc.GetStringFromContext(ctx, constants.RPCMetaKeyDatabase)
	if err != nil {
		return
	}
	shardStateData, err := rpc.GetStringFromContext(ctx, constants.RPCMetaKeyShardState)
	if err != nil {
		return
	}
	err = encoding.JSONUnmarshal([]byte(shardStateData), &shardState)
	if err != nil {
		return
	}
	return
}

// getFollowerInfoFromCtx gets follower metadata from rpc context.
func (r *ReplicaHandler) getFollowerInfoFromCtx(ctx context.Context) (database string, shardID models.ShardID,
	leader models.NodeID, replica models.NodeID, err error) {
	return
}

// getOrCreatePartition returns write ahead log's partition if exist, else creates a new partition.
func (r *ReplicaHandler) getOrCreatePartition(database string, shardID models.ShardID) (replica.Partition, error) {
	wal := r.walMgr.GetOrCreateLog(database)
	p, err := wal.GetOrCreatePartition(shardID)
	if err != nil {
		r.logger.Error("create wal partition err", logger.Error(err))
		return nil, err
	}
	return p, nil
}
