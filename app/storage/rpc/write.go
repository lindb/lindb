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

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/storage"
)

// WriteHandler implements protoWriteV1.WriteServiceServer interface for handling write rpc request.
type WriteHandler struct {
	protoWriteV1.UnimplementedWriteServiceServer

	engine storage.Engine

	logger logger.Logger
}

// NewWriteHandler creates a write handler.
func NewWriteHandler(
	engine storage.Engine,
) protoWriteV1.WriteServiceServer {
	return &WriteHandler{
		engine: engine,
		logger: logger.GetLogger("Storage", "WriteRPC"),
	}
}

func (r *WriteHandler) WriteArrow(server protoWriteV1.WriteService_WriteArrowServer) error {
	segmentState, err := r.getSegmentInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if len(segmentState.Shard.Replica.Replicas) == 0 {
		return status.Error(codes.InvalidArgument, "replicas cannot be empty")
	}

	log, err := getOrCreateSegment(
		r.engine,
		segmentState.Database,
		segmentState.Shard.ID,
		segmentState.SegmentTime,
		segmentState.Shard.Leader)
	if err != nil {
		r.logger.Error("get or create wal partition err, when do write", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}

	var deserializer []*larrow.Deserializer
	var records []arrow.RecordBatch

	// handle write request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("receive write request err", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}
		resp := &protoWriteV1.WriteArrowResponse{}
		payloads := req.Payloads
		// NOTE: initialize deserializer and records when receive first request, and reuse them for later requests, avoid unnecessary memory allocation.
		// paylaods length must be same for each request, otherwise need to reinitialize deserializer and records, which is not expected.
		if len(deserializer) == 0 {
			deserializer = make([]*larrow.Deserializer, len(payloads))
			records = make([]arrow.RecordBatch, len(payloads))
			for _, payload := range payloads {
				deserializer[payload.SchemaIndex] = larrow.NewDeserializer()
			}
		} else if len(deserializer) != len(payloads) {
			err = fmt.Errorf("payloads length is not same for each request, expected %d, but got %d", len(deserializer), len(payloads))
			r.logger.Error("invalid payloads length", logger.Error(err))
			return status.Error(codes.InvalidArgument, err.Error())
		}

		for _, payload := range payloads {
			record, err := deserializer[payload.SchemaIndex].Deserialize(payload.Record)
			if err != nil {
				resp.Err = err.Error()
				r.logger.Error("deserialize record err", logger.Error(err))
				break
			}
			records[payload.SchemaIndex] = record
		}

		if err := log.WriteArrow(records); err != nil {
			resp.Err = err.Error()
			r.logger.Error("write wal log err", logger.Error(err))
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// Write does metric write request.
func (r *WriteHandler) Write(server protoWriteV1.WriteService_WriteServer) error {
	segmentState, err := r.getSegmentInfoFromCtx(server.Context())
	if err != nil {
		r.logger.Error("get param err", logger.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if len(segmentState.Shard.Replica.Replicas) == 0 {
		return status.Error(codes.InvalidArgument, "replicas cannot be empty")
	}

	log, err := getOrCreateSegment(
		r.engine,
		segmentState.Database,
		segmentState.Shard.ID,
		segmentState.SegmentTime,
		segmentState.Shard.Leader)
	if err != nil {
		r.logger.Error("get or create wal partition err, when do write", logger.Error(err))
		return status.Error(codes.Internal, err.Error())
	}

	// handle write request from stream
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			r.logger.Error("receive write request err", logger.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		resp := &protoWriteV1.WriteResponse{}
		// write wal log
		err = log.Write(req.Record)
		if err != nil {
			resp.Err = err.Error()
		}

		if err := server.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

// getSegmentInfoFromCtx returns segment state metadata from rpc context.
func (r *WriteHandler) getSegmentInfoFromCtx(ctx context.Context) (segmentState models.SegmentState, err error) {
	segmentStateDate, err := rpc.GetStringFromContext(ctx, constants.RPCMetaKeySegmentState)
	if err != nil {
		return
	}
	err = encoding.JSONUnmarshal([]byte(segmentStateDate), &segmentState)
	if err != nil {
		return
	}
	return
}
