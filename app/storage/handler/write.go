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
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/rpc"
)

// WriteHandler implements protoWriteV1.WriteServiceServer interface for handling write rpc request.
type WriteHandler struct {
	walMgr replica.WriteAheadLogManager

	logger *logger.Logger
}

// NewWriteHandler creates a write handler.
func NewWriteHandler(
	walMgr replica.WriteAheadLogManager,
) *WriteHandler {
	return &WriteHandler{
		walMgr: walMgr,
		logger: logger.GetLogger("storage", "WriteRPC"),
	}
}

// Write does metric write request.
func (r *WriteHandler) Write(server protoWriteV1.WriteService_WriteServer) error {
	familyState, err := r.getFamilyInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if len(familyState.Shard.Replica.Replicas) == 0 {
		return status.Error(codes.InvalidArgument, "replicas cannot be empty")
	}

	p, err := r.getOrCreatePartition(familyState.Database, familyState.Shard.ID, familyState.FamilyTime)
	if err != nil {
		r.logger.Error("get or create wal partition err, when do write", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	err = p.BuildReplicaForLeader(familyState.Shard.Leader, familyState.Shard.Replica.Replicas)
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

		resp := &protoWriteV1.WriteResponse{}
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

// getFamilyInfoFromCtx returns family state metadata from rpc context.
func (r *WriteHandler) getFamilyInfoFromCtx(ctx context.Context) (familyState models.FamilyState, err error) {
	familyStateDate, err := rpc.GetStringFromContext(ctx, constants.RPCMetaKeyFamilyState)
	if err != nil {
		return
	}
	err = encoding.JSONUnmarshal([]byte(familyStateDate), &familyState)
	if err != nil {
		return
	}
	return
}

// getOrCreatePartition returns write ahead log's partition if exist, else creates a new partition.
func (r *WriteHandler) getOrCreatePartition(
	database string,
	shardID models.ShardID,
	familyTime int64,
) (replica.Partition, error) {
	wal := r.walMgr.GetOrCreateLog(database)
	p, err := wal.GetOrCreatePartition(shardID, familyTime)
	if err != nil {
		return nil, err
	}
	return p, nil
}
