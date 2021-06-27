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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
	"github.com/lindb/lindb/service"
)

// replicaHandler implements replica.ReplicaServiceServer interface for handling replica rpc request.
type replicaHandler struct {
	walMgr         replica.WriteAheadLogManager
	storageService service.StorageService

	logger *logger.Logger
}

// NewReplicaHandler creates a replica handler.
func NewReplicaHandler(walMgr replica.WriteAheadLogManager,
	storageService service.StorageService) replicaRpc.ReplicaServiceServer {
	return &replicaHandler{
		walMgr:         walMgr,
		storageService: storageService,
		logger:         logger.GetLogger("storage", "replicaRpc"),
	}
}

// GetReplicaAckIndex returns current replica ack index.
func (r *replicaHandler) GetReplicaAckIndex(ctx context.Context,
	req *replicaRpc.GetReplicaAckIndexRequest,
) (*replicaRpc.GetReplicaAckIndexResponse, error) {
	panic("implement me")
}

// Reset resets replica index.
func (r *replicaHandler) Reset(ctx context.Context,
	request *replicaRpc.ResetIndexRequest,
) (*replicaRpc.ResetIndexResponse, error) {
	panic("implement me")
}

// Replica does replica request, and writes data.
func (r *replicaHandler) Replica(server replicaRpc.ReplicaService_ReplicaServer) error {
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

		resp := &replicaRpc.ReplicaResponse{}
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
func (r *replicaHandler) Write(server replicaRpc.ReplicaService_WriteServer) error {
	database, shardID, leader, replicas, err := r.getReplicasInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if len(replicas) == 0 {
		return status.Error(codes.InvalidArgument, "replicas cannot be empty")
	}

	p, err := r.getOrCreatePartition(database, shardID)
	if err != nil {
		r.logger.Error("create wal partition err", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	err = p.BuildReplicaForLeader(leader, replicas)
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

		resp := &replicaRpc.WriteResponse{}
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
func (r *replicaHandler) getReplicasInfoFromCtx(ctx context.Context) (database string, shardID models.ShardID,
	leader models.NodeID, replicas []models.NodeID, err error) {
	database, err = rpc.GetDatabaseFromContext(ctx)
	if err != nil {
		return
	}
	shard, err0 := rpc.GetShardIDFromContext(ctx)
	if err0 != nil {
		err = err0
		return
	}
	shardID = models.ShardID(shard)

	leader, err = rpc.GetLeaderFromContext(ctx)
	if err != nil {
		return
	}
	replicas, err = rpc.GetReplicasFromContext(ctx)
	if err != nil {
		return
	}
	return
}

// getFollowerInfoFromCtx gets follower metadata from rpc context.
func (r *replicaHandler) getFollowerInfoFromCtx(ctx context.Context) (database string, shardID models.ShardID,
	leader models.NodeID, replica models.NodeID, err error) {
	database, err = rpc.GetDatabaseFromContext(ctx)
	if err != nil {
		return
	}
	shard, err0 := rpc.GetShardIDFromContext(ctx)
	if err0 != nil {
		err = err0
		return
	}
	shardID = models.ShardID(shard)

	leader, err = rpc.GetLeaderFromContext(ctx)
	if err != nil {
		return
	}
	replica, err = rpc.GetFollowerFromContext(ctx)
	if err != nil {
		return
	}
	return
}

// getOrCreatePartition returns write ahead log's partition if exist, else creates a new partition.
func (r *replicaHandler) getOrCreatePartition(database string, shardID models.ShardID) (replica.Partition, error) {
	wal := r.walMgr.GetOrCreateLog(database)
	p, err := wal.GetOrCreatePartition(shardID)
	if err != nil {
		r.logger.Error("create wal partition err", logger.Error(err))
		return nil, err
	}
	return p, nil
}
