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
	"io"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
)

//go:generate mockgen -source=./write_stream.go -destination=./write_stream_mock.go -package=rpc

// WriteStream represents the channel which writes metric to storage based on grpc stream,
// and receives write response in background.
type WriteStream interface {
	io.Closer
	// Send sends metric data to storage.
	Send(data []byte) error
}

// writeStream implements WriteStream interface.
type writeStream struct {
	ctx    context.Context
	cancel context.CancelFunc

	target     models.Node
	database   string
	shardState *models.ShardState
	familyTime int64

	fct    ClientStreamFactory
	cli    protoWriteV1.WriteService_WriteClient
	closed *atomic.Bool

	logger *logger.Logger
}

// NewWriteStream creates a WriteStream instance, initialize grpc connection(stream) and receive response task.
func NewWriteStream(
	ctx context.Context,
	target models.Node,
	database string, shardState *models.ShardState, familyTime int64,
	fct ClientStreamFactory,
) (WriteStream, error) {
	c, cancel := context.WithCancel(ctx)
	s := &writeStream{
		ctx:        c,
		cancel:     cancel,
		target:     target,
		database:   database,
		shardState: shardState,
		familyTime: familyTime,
		fct:        fct,
		closed:     atomic.NewBool(false),
		logger:     logger.GetLogger("rpc", "WriteStream"),
	}

	// initialize write stream
	if err := s.initialize(); err != nil {
		return nil, err
	}
	return s, nil
}

// initialize grpc connection, then starts receive response task.
func (s *writeStream) initialize() error {
	writeService, err := s.fct.CreateWriteServiceClient(s.target)
	if err != nil {
		return err
	}

	// pass metadata(database/shard/family state) when create rpc connection.
	familyState := encoding.JSONMarshal(&models.FamilyState{
		Database:   s.database,
		Shard:      *s.shardState,
		FamilyTime: s.familyTime,
	})
	ctx := CreateOutgoingContextWithPairs(s.ctx, constants.RPCMetaKeyFamilyState, string(familyState))
	writeCli, err := writeService.Write(ctx)

	if err != nil {
		return err
	}

	// set write client
	s.cli = writeCli

	// start receive response task
	go s.recvLoop()

	s.logger.Info("initialize write client stream successfully",
		logger.String("database", s.database),
		logger.Any("shard", s.shardState.ID),
		logger.String("target", s.target.Indicator()))
	return nil
}

// Send sends metric data to storage.
func (s *writeStream) Send(data []byte) error {
	if s.closed.Load() {
		// if write stream is closed, return EOF err
		return io.EOF
	}
	return s.cli.Send(&protoWriteV1.WriteRequest{Record: data})
}

// Close closes send stream, and cancel stream context, server will stop receive write request under this stream.
func (s *writeStream) Close() error {
	defer s.cancel() // close stream context
	s.logger.Info("close write stream",
		logger.String("target", s.target.Indicator()))
	return s.cli.CloseSend()
}

// recvLoop is a loop to receive message from write stream.
// if stream context is done or receive io.EOF err, need mark stream is closed.
func (s *writeStream) recvLoop() {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("panic when receive response from write stream",
				logger.String("target", s.target.Indicator()),
				logger.Any("err", err),
				logger.Stack())
			s.closed.Store(true)
		}
	}()

	for {
		select {
		case <-s.cli.Context().Done():
			// stream is closed, return it.
			if err := s.cli.Context().Err(); err != nil {
				s.logger.Error("write stream context is canceled",
					logger.String("target", s.target.Indicator()),
					logger.Error(err))
			}
			s.closed.Store(true)
			return
		default:
			resp, err := s.cli.Recv()
			if err != nil {
				s.logger.Error("receive error from write stream",
					logger.String("target", s.target.Indicator()),
					logger.Error(err))
				if err == io.EOF {
					s.closed.Store(true)
					// stream is closed, return it.
					return
				}
				continue
			}
			if resp.Err != "" {
				// get err from response
				s.logger.Error("get err write response",
					logger.String("target", s.target.Indicator()),
					logger.String("err", resp.Err))
			}
		}
	}
}
