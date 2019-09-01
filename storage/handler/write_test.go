package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/storage"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/tsdb"
)

/**
case replica seq match:

case replica seq not match:
*/

var (
	node = models.Node{
		IP:   "127.0.0.1",
		Port: 123,
	}
	database = "database"
	shardID  = int32(0)
)

func buildWriteRequest(seqBegin, seqEnd int64) (*storage.WriteRequest, string) {
	replicas := make([]*storage.Replica, seqEnd-seqBegin)
	for i := seqBegin; i < seqEnd; i++ {
		replicas[i-seqBegin] = &storage.Replica{
			Seq:  i,
			Data: []byte(strconv.Itoa(int(i))),
		}
	}
	wr := &storage.WriteRequest{
		Replicas: replicas,
	}
	return wr, fmt.Sprintf("[%d,%d)", seqBegin, seqEnd)
}

func TestWriter_Next(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	sm := replication.NewMockSequenceManager(ctl)
	s := replication.NewMockSequence(ctl)

	seq := int64(5)
	s.EXPECT().GetHeadSeq().Return(seq)
	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true)

	writer := NewWriter(nil, sm)

	ctx := mockContext(database, shardID, node)
	resp, err := writer.Next(ctx, &storage.NextSeqRequest{
		ShardID:  shardID,
		Database: database})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, seq, resp.Seq)

	// not metadata
	ctx = context.TODO()
	_, err = writer.Next(ctx, &storage.NextSeqRequest{
		Database: database,
		ShardID:  shardID,
	})
	assert.NotNil(t, err)

	ctx = mockContext(database, shardID, node)
	sm.EXPECT().GetSequence(database, shardID, node).Return(nil, false)
	sm.EXPECT().CreateSequence(database, shardID, node).Return(nil, fmt.Errorf("err"))
	_, err = writer.Next(ctx, &storage.NextSeqRequest{
		Database: database,
		ShardID:  shardID,
	})
	assert.NotNil(t, err)
}

func TestWriter_Reset(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	sm := replication.NewMockSequenceManager(ctl)
	s := replication.NewMockSequence(ctl)

	seq := int64(5)
	s.EXPECT().SetHeadSeq(seq).Return()
	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true)

	writer := NewWriter(nil, sm)

	ctx := mockContext(database, shardID, node)
	_, err := writer.Reset(ctx, &storage.ResetSeqRequest{
		Database: database,
		ShardID:  shardID,
		Seq:      seq,
	})
	if err != nil {
		t.Fatal(err)
	}

	// not metadata
	ctx = context.TODO()
	_, err = writer.Reset(ctx, &storage.ResetSeqRequest{
		Database: database,
		ShardID:  shardID,
		Seq:      seq,
	})
	assert.NotNil(t, err)

	ctx = mockContext(database, shardID, node)
	sm.EXPECT().GetSequence(database, shardID, node).Return(nil, false)
	sm.EXPECT().CreateSequence(database, shardID, node).Return(nil, fmt.Errorf("err"))
	_, err = writer.Reset(ctx, &storage.ResetSeqRequest{
		Database: database,
		ShardID:  shardID,
		Seq:      seq,
	})
	assert.NotNil(t, err)
}

func TestWriter_Write_Fail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	sm := replication.NewMockSequenceManager(ctl)
	storageSRV := service.NewMockStorageService(ctl)

	writer := NewWriter(storageSRV, sm)
	stream := storage.NewMockWriteService_WriteServer(ctl)
	stream.EXPECT().Context().Return(context.TODO())
	err := writer.Write(stream)
	assert.NotNil(t, err)

	sm.EXPECT().GetSequence(database, shardID, node).Return(nil, false)
	sm.EXPECT().CreateSequence(database, shardID, node).Return(nil, fmt.Errorf("err"))
	stream.EXPECT().Context().Return(mockContext(database, shardID, node))
	err = writer.Write(stream)
	assert.NotNil(t, err)

	s := replication.NewMockSequence(ctl)
	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true)
	stream.EXPECT().Context().Return(mockContext(database, shardID, node))
	storageSRV.EXPECT().GetShard(database, shardID).Return(nil)
	err = writer.Write(stream)
	assert.NotNil(t, err)

	shard := tsdb.NewMockShard(ctl)
	stream.EXPECT().Context().Return(mockContext(database, shardID, node)).AnyTimes()
	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true).AnyTimes()
	storageSRV.EXPECT().GetShard(database, shardID).Return(shard).AnyTimes()
	stream.EXPECT().Recv().Return(nil, io.EOF)
	_ = writer.Write(stream)
	assert.Nil(t, nil)

	gomock.InOrder(
		stream.EXPECT().Recv().Return(&storage.WriteRequest{}, nil),
		stream.EXPECT().Recv().Return(nil, fmt.Errorf("err")),
	)
	err = writer.Write(stream)
	assert.NotNil(t, err)
}

func TestWriter_Write_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	sm := replication.NewMockSequenceManager(ctl)
	s := replication.NewMockSequence(ctl)

	var (
		seqBeg int64 = 5
		seqEnd int64 = 10
	)

	for i := seqBeg; i < seqEnd; i++ {
		s.EXPECT().GetHeadSeq().Return(i)
		s.EXPECT().SetHeadSeq(i + 1).Return()
	}

	s.EXPECT().GetHeadSeq().Return(seqEnd)
	s.EXPECT().Synced().Return(false)

	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true)

	stom := mockStorage(ctl, database, shardID, mockShard(ctl))

	writer := NewWriter(stom, sm)

	ctx := mockContext(database, shardID, node)

	stream := storage.NewMockWriteService_WriteServer(ctl)
	stream.EXPECT().Context().Return(ctx)

	wr1, _ := buildWriteRequest(seqBeg, seqEnd)
	stream.EXPECT().Recv().Return(wr1, nil)

	stream.EXPECT().Send(&storage.WriteResponse{
		CurSeq: seqEnd - 1,
	}).Return(nil)

	stream.EXPECT().Recv().Return(nil, errors.New("recv error"))

	err := writer.Write(stream)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestWriter_WriteSeqNotMatch(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	sm := replication.NewMockSequenceManager(ctl)
	s := replication.NewMockSequence(ctl)

	var (
		seqBeg int64 = 5
		seqEnd int64 = 10
	)

	// wrong seq
	s.EXPECT().GetHeadSeq().Return(seqEnd + 1)

	sm.EXPECT().GetSequence(database, shardID, node).Return(s, true)

	stom := mockStorage(ctl, database, shardID, mockShard(ctl))

	writer := NewWriter(stom, sm)

	ctx := mockContext(database, shardID, node)

	stream := storage.NewMockWriteService_WriteServer(ctl)
	stream.EXPECT().Context().Return(ctx)

	wr1, _ := buildWriteRequest(seqBeg, seqEnd)
	stream.EXPECT().Recv().Return(wr1, nil)

	err := writer.Write(stream)
	if err == nil {
		t.Fatal("should be error")
	}
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
	db, shard, node, err := parseCtx(mockContext("db", int32(10), models.Node{IP: "1.1.1.1", Port: 9000}))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, db, "db")
	assert.Equal(t, shard, int32(10))
	assert.Equal(t, models.Node{IP: "1.1.1.1", Port: 9000}, *node)
}

func mockStorage(ctl *gomock.Controller, db string, shardID int32, shard tsdb.Shard) service.StorageService {
	mockStorage := service.NewMockStorageService(ctl)
	mockStorage.EXPECT().GetShard(db, shardID).Return(shard)
	return mockStorage
}

func mockShard(ctl *gomock.Controller) tsdb.Shard {
	mockShard := tsdb.NewMockShard(ctl)
	mockShard.EXPECT().Write(gomock.Any()).Return(nil).AnyTimes()
	return mockShard
}

func mockContext(db string, shardID int32, node models.Node) context.Context {
	return rpc.CreateIncomingContext(context.TODO(), db, shardID, node)
}
