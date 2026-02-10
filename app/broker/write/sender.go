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

package write

import (
	"context"
	"errors"
	"io"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
	"github.com/lindb/lindb/rpc"
)

type Sender interface {
	io.Closer
	// Send sends arrow payloads to storage.
	Send(payloads []byte) error
}

type sender struct {
	ctx    context.Context
	cancel context.CancelFunc

	target models.Node
	state  string

	cli protoWriteV1.WriteService_WriteClient

	running *atomic.Bool

	logger logger.Logger
}

func NewSender(parent context.Context,
	target models.Node,
	state *models.SegmentState,
) (Sender, error) {
	ctx, cancel := context.WithCancel(parent)
	s := &sender{
		ctx:     ctx,
		cancel:  cancel,
		target:  target,
		state:   string(encoding.JSONMarshal(state)),
		running: atomic.NewBool(true),
		logger:  logger.GetLogger("Write", "Sender"),
	}
	// initialize streaming of sender
	if err := s.initialize(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *sender) initialize() error {
	conn, err := rpc.GetBrokerClientConnFactory().GetClientConn(s.target)
	if err != nil {
		return err
	}
	client := protoWriteV1.NewWriteServiceClient(conn)
	ctx := rpc.CreateOutgoingContextWithPairs(s.ctx, constants.RPCMetaKeySegmentState, s.state)
	writeCli, err := client.Write(ctx)
	if err != nil {
		return err
	}

	// set write client
	s.cli = writeCli

	go s.recvLoop()

	s.logger.Info("initialized sender",
		logger.String("target", s.target.Indicator()), logger.String("state", s.state))
	return nil
}

func (s *sender) Send(payloads []byte) error {
	if !s.running.Load() {
		return errors.New("sender is closed")
	}
	return s.cli.Send(&protoWriteV1.WriteRequest{
		Record: payloads,
	})
}

// Close closes the sender.
func (s *sender) Close() error {
	if s.running.CompareAndSwap(true, false) {
		defer s.cancel()
		s.logger.Info("closing sender", logger.String("target", s.target.Indicator()), logger.String("state", s.state))
		return s.cli.CloseSend()
	}
	return nil
}

func (s *sender) recvLoop() {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("panic in sender receive loop",
				logger.String("target", s.target.Indicator()),
				logger.String("state", s.state), logger.Any("err", err), logger.Stack())
		}
		if err := s.Close(); err != nil {
			s.logger.Error("failed to close sender", logger.String("target", s.target.Indicator()),
				logger.String("state", s.state), logger.Error(err))
		}
		s.logger.Info("stopping sender receive loop",
			logger.String("target", s.target.Indicator()), logger.String("state", s.state))
	}()
	for {
		select {
		case <-s.ctx.Done():
			if err := s.ctx.Err(); err != nil {
				s.logger.Error("sender context is canceled",
					logger.String("target", s.target.Indicator()),
					logger.String("state", s.state),
					logger.Error(err))
			}
			return
		default:
			resp, err := s.cli.Recv()
			if err != nil {
				s.logger.Error("failed to receive write response",
					logger.String("target", s.target.Indicator()),
					logger.String("state", s.state),
					logger.Error(err))
				if errors.Is(err, io.EOF) {
					return
				}
				continue
			}
			if resp.Err != "" {
				// just log storage returned write error
				s.logger.Error("received write error from storage",
					logger.String("target", s.target.Indicator()),
					logger.String("state", s.state),
					logger.String("err", resp.Err))
			}
		}
	}
}
