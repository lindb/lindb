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
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/models"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	protoStorageV1 "github.com/lindb/lindb/proto/gen/v1/storage"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

/**
case replica seq match:

case replica seq not match:
*/

var (
	node = &models.StatelessNode{
		HostIP:   "127.0.0.1",
		GRPCPort: 123,
	}
	database = "database"
	shardID  = models.ShardID(0)
)

func TestWriter_Next(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	engine := tsdb.NewMockEngine(ctl)
	s := replication.NewMockSequence(ctl)
	db := tsdb.NewMockDatabase(ctl)
	shard := tsdb.NewMockShard(ctl)
	shard.EXPECT().GetOrCreateSequence(gomock.Any()).Return(s, nil)
	db.EXPECT().GetShard(gomock.Any()).Return(shard, true)
	engine.EXPECT().GetDatabase(gomock.Any()).Return(db, true)

	seq := int64(5)
	s.EXPECT().GetHeadSeq().Return(seq)

	writer := NewWriter(engine)

	ctx := mockContext(database, shardID, node)
	resp, err := writer.Next(ctx, &protoStorageV1.NextSeqRequest{
		ShardID:  int32(shardID),
		Database: database})
	assert.NoError(t, err)
	assert.Equal(t, seq, resp.Seq)

	// not metadata
	ctx = context.TODO()
	_, err = writer.Next(ctx, &protoStorageV1.NextSeqRequest{
		Database: database,
		ShardID:  int32(shardID),
	})
	assert.Error(t, err)

	ctx = mockContext(database, shardID, node)
	engine.EXPECT().GetDatabase(gomock.Any()).Return(nil, false)
	_, err = writer.Next(ctx, &protoStorageV1.NextSeqRequest{
		ShardID:  int32(shardID),
		Database: database})
	assert.Error(t, err)
}

func TestWriter_Reset(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	engine := tsdb.NewMockEngine(ctl)
	s := replication.NewMockSequence(ctl)
	db := tsdb.NewMockDatabase(ctl)
	shard := tsdb.NewMockShard(ctl)
	shard.EXPECT().GetOrCreateSequence(gomock.Any()).Return(s, nil)
	db.EXPECT().GetShard(gomock.Any()).Return(shard, true)
	engine.EXPECT().GetDatabase(gomock.Any()).Return(db, true)

	seq := int64(5)
	s.EXPECT().SetHeadSeq(seq).Return()

	writer := NewWriter(engine)

	ctx := mockContext(database, shardID, node)
	_, err := writer.Reset(ctx, &protoStorageV1.ResetSeqRequest{
		Database: database,
		ShardID:  int32(shardID),
		Seq:      seq,
	})
	assert.NoError(t, err)

	// not metadata
	ctx = context.TODO()
	_, err = writer.Reset(ctx, &protoStorageV1.ResetSeqRequest{
		Database: database,
		ShardID:  int32(shardID),
		Seq:      seq,
	})
	assert.Error(t, err)

	ctx = mockContext(database, shardID, node)
	engine.EXPECT().GetDatabase(gomock.Any()).Return(db, true)
	db.EXPECT().GetShard(gomock.Any()).Return(nil, false)
	_, err = writer.Reset(ctx, &protoStorageV1.ResetSeqRequest{
		Database: database,
		ShardID:  int32(shardID),
		Seq:      seq,
	})
	assert.Error(t, err)
}

func TestWriter_Write_Fail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	engine := tsdb.NewMockEngine(ctl)

	writer := NewWriter(engine)

	// metadata err
	writeServer := protoStorageV1.NewMockWriteService_WriteServer(ctl)
	writeServer.EXPECT().Context().Return(context.TODO())
	err := writer.Write(writeServer)
	assert.Error(t, err)

	// no shard
	ctx := mockContext(database, shardID, node)
	writeServer.EXPECT().Context().Return(ctx)
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(nil, false)
	err = writer.Write(writeServer)
	assert.Error(t, err)

	// get sequence err
	writeServer.EXPECT().Context().Return(ctx).AnyTimes()
	shard := tsdb.NewMockShard(ctl)
	engine.EXPECT().GetShard(gomock.Any(), gomock.Any()).Return(shard, true).AnyTimes()
	shard.EXPECT().GetOrCreateSequence(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = writer.Write(writeServer)
	assert.Error(t, err)

	// stream eof
	s := replication.NewMockSequence(ctl)
	shard.EXPECT().GetOrCreateSequence(gomock.Any()).Return(s, nil).AnyTimes()
	writeServer.EXPECT().Recv().Return(nil, io.EOF)
	err = writer.Write(writeServer)
	assert.Nil(t, err)

	// internal error
	writeServer.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	err = writer.Write(writeServer)
	assert.Error(t, err)

	// no replica
	writeServer.EXPECT().Recv().Return(&protoStorageV1.WriteRequest{}, nil)
	writeServer.EXPECT().Recv().Return(nil, io.EOF)
	err = writer.Write(writeServer)
	assert.Nil(t, err)

	// replica index not match
	writeServer.EXPECT().Recv().Return(&protoStorageV1.WriteRequest{Replicas: []*protoStorageV1.Replica{{Seq: int64(10)}}}, nil)
	s.EXPECT().GetHeadSeq().Return(int64(8))
	err = writer.Write(writeServer)
	assert.Error(t, err)

	writeServer.EXPECT().Recv().Return(&protoStorageV1.WriteRequest{Replicas: []*protoStorageV1.Replica{{Seq: int64(10)}}}, nil)
	s.EXPECT().GetHeadSeq().Return(int64(9)).MaxTimes(2)
	s.EXPECT().SetHeadSeq(gomock.Any())
	s.EXPECT().GetAckSeq().Return(int64(8))
	writeServer.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	err = writer.Write(writeServer)
	assert.Error(t, err)
}

func TestWriter_handle_replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	engine := tsdb.NewMockEngine(ctrl)

	writer := NewWriter(engine)
	shard := tsdb.NewMockShard(ctrl)
	writer.handleReplica(shard, &protoStorageV1.Replica{Seq: int64(10), Data: []byte{1, 2, 3}})

	buf := &bytes.Buffer{}
	compressBuf := snappy.NewBufferedWriter(buf)
	_, _ = compressBuf.Write([]byte{1, 2, 4})
	_ = compressBuf.Flush()
	writer.handleReplica(shard, &protoStorageV1.Replica{Seq: int64(10), Data: buf.Bytes()})

	metricList := &protoMetricsV1.MetricList{
		Metrics: []*protoMetricsV1.Metric{{Name: "test"}},
	}
	data, _ := metricList.Marshal()

	buf = &bytes.Buffer{}
	compressBuf = snappy.NewBufferedWriter(buf)
	_, _ = compressBuf.Write(data)
	_ = compressBuf.Flush()
	shard.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
	writer.handleReplica(shard, &protoStorageV1.Replica{Seq: int64(10), Data: buf.Bytes()})
}

func TestWrite_parse_ctx(t *testing.T) {
	_, _, _, err := parseCtx(context.TODO())
	assert.NotNil(t, err)

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs("metaKeyLogicNode", "1.1.1.1:9000"))
	_, _, _, err = parseCtx(ctx)
	assert.NotNil(t, err)
	ctx = metadata.NewIncomingContext(context.TODO(), metadata.Pairs("metaKeyLogicNode", "1.1.1.1:9000", "metaKeyDatabase", "db"))
	_, _, _, err = parseCtx(ctx)
	assert.NotNil(t, err)
	db, shard, node, err := parseCtx(mockContext("db", models.ShardID(10), &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, db, "db")
	assert.Equal(t, shard, models.ShardID(10))
	assert.Equal(t, &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}, node)
}

func mockContext(db string, shardID models.ShardID, node models.Node) context.Context {
	return rpc.CreateIncomingContext(context.TODO(), db, shardID, node)
}
