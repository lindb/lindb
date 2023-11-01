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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
)

func TestNewWriteStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fct := NewMockClientStreamFactory(ctrl)

	// case 1: create write service cli err
	fct.EXPECT().CreateWriteServiceClient(gomock.Any()).Return(nil, fmt.Errorf("err"))
	stream, err := NewWriteStream(context.TODO(), nil, "test", &models.ShardState{}, 1, fct)
	assert.Error(t, err)
	assert.Nil(t, stream)

	// case 2: create write cli err
	writeSrv := protoWriteV1.NewMockWriteServiceClient(ctrl)
	fct.EXPECT().CreateWriteServiceClient(gomock.Any()).Return(writeSrv, nil).AnyTimes()
	writeSrv.EXPECT().Write(gomock.Any()).Return(nil, fmt.Errorf("err"))
	stream, err = NewWriteStream(context.TODO(), nil, "test", &models.ShardState{}, 1, fct)
	assert.Error(t, err)
	assert.Nil(t, stream)

	// case 3: create instance success
	cli := protoWriteV1.NewMockWriteService_WriteClient(ctrl)
	writeSrv.EXPECT().Write(gomock.Any()).Return(cli, nil)
	cli.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()
	cli.EXPECT().Context().Return(context.TODO()).AnyTimes()
	stream, err = NewWriteStream(context.TODO(), &models.StatefulNode{}, "test", &models.ShardState{}, 1, fct)
	assert.NoError(t, err)
	assert.NotNil(t, stream)

	cli.EXPECT().CloseSend().Return(nil)
	err = stream.Close()
	assert.NoError(t, err)
}

func TestWriteStream_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := protoWriteV1.NewMockWriteService_WriteClient(ctrl)
	stream := &writeStream{
		cli:    cli,
		closed: atomic.NewBool(true),
	}
	assert.Equal(t, io.EOF, stream.Send(nil))
	stream.closed.Store(false)
	cli.EXPECT().Send(gomock.Any()).Return(nil)
	assert.NoError(t, stream.Send(nil))
}

func TestWriteStream_Recv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := &writeStream{
		target: &models.StatefulNode{},
		closed: atomic.NewBool(false),
		logger: logger.GetLogger("RPC", "WriteStream"),
	}
	// case 1: panic
	stream.recvLoop()
	assert.True(t, stream.closed.Load())
	// case 2: context is done
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	cli := protoWriteV1.NewMockWriteService_WriteClient(ctrl)
	stream = &writeStream{
		cli:    cli,
		target: &models.StatefulNode{},
		closed: atomic.NewBool(false),
		logger: logger.GetLogger("RPC", "WriteStream"),
	}
	cli.EXPECT().Context().Return(ctx).MaxTimes(2)
	stream.recvLoop()
	assert.True(t, stream.closed.Load())
	// case 3: recv err
	stream = &writeStream{
		cli:    cli,
		closed: atomic.NewBool(false),
		target: &models.StatefulNode{},
		logger: logger.GetLogger("RPC", "WriteStream"),
	}
	cli.EXPECT().Context().Return(context.TODO()).AnyTimes()
	cli.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	cli.EXPECT().Recv().Return(&protoWriteV1.WriteResponse{Err: "err"}, nil)
	cli.EXPECT().Recv().Return(nil, io.EOF)
	stream.recvLoop()
}
