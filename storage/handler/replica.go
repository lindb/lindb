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
	panic("implement me")
}

// Write does metric write request.
func (r *replicaHandler) Write(server replicaRpc.ReplicaService_WriteServer) error {
	database, shardID, leader, replicas, err := r.getReplicaInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	wal := r.walMgr.GetOrCreateLog(database)
	p, err := wal.GetOrCreatePartition(shardID)
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

// getReplicaInfoFromCtx gets metadata from rpc context.
func (r *replicaHandler) getReplicaInfoFromCtx(ctx context.Context) (database string, shardID models.ShardID,
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
