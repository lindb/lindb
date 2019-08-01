package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/replication"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/storage"
	"github.com/eleme/lindb/service"
	"github.com/eleme/lindb/tsdb"
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

/**
stream replica -> 01234
stream replica -> 012 (duplicate reset to 5)
stream replica -> 34
stream replica -> 56789
*/
func TestWriter_Write(t *testing.T) {
	tmp := path.Join(os.TempDir(), "test_write")
	if err := fileutil.RemoveDir(tmp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := fileutil.RemoveDir(tmp); err != nil {
			t.Error(err)
		}
	}()

	mockCtl := gomock.NewController(t)

	node := models.Node{
		IP:   "1.1.1.1",
		Port: 12345,
	}

	ws := mockStream(mockCtl, "db", 0, node)

	sm, err := replication.NewSequenceManager(tmp)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriter(mockStorage(mockCtl, "db", 0, mockShard(mockCtl)), sm)

	err = writer.Write(ws)
	if err != nil {
		t.Fatal("should be nil")
	}

}

func mockStream(ctl *gomock.Controller, db string, shardID uint32, node models.Node) storage.WriteService_WriteServer {
	ctx := mockContext(db, shardID, node)
	mockWriteServer := storage.NewMockWriteService_WriteServer(ctl)

	gomock.InOrder(
		mockWriteServer.EXPECT().Context().Return(ctx).Times(1),
		mockWriteServer.EXPECT().Recv().DoAndReturn(func() (*storage.WriteRequest, error) {
			rq, str := buildWriteRequest(0, 5)
			fmt.Println("recv:" + str)
			return rq, nil
		}),
		mockWriteServer.EXPECT().Recv().DoAndReturn(func() (*storage.WriteRequest, error) {
			rq, str := buildWriteRequest(0, 3)
			fmt.Println("recv:" + str)
			return rq, nil
		}),
		mockWriteServer.EXPECT().Send(&storage.WriteResponse{
			Seq: &storage.WriteResponse_ResetSeq{
				ResetSeq: 5,
			},
		}).Return(nil),
		mockWriteServer.EXPECT().Recv().DoAndReturn(func() (*storage.WriteRequest, error) {
			rq, str := buildWriteRequest(3, 5)
			fmt.Println("recv:" + str)
			return rq, nil
		}),
		mockWriteServer.EXPECT().Recv().DoAndReturn(func() (*storage.WriteRequest, error) {
			rq, str := buildWriteRequest(5, 10)
			fmt.Println("recv:" + str)
			return rq, nil
		}),
		mockWriteServer.EXPECT().Recv().DoAndReturn(func() (*storage.WriteRequest, error) {
			fmt.Println("recv: eof")
			return nil, io.EOF
		}),
	)

	return mockWriteServer
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

func mockContext(db string, shardID uint32, node models.Node) context.Context {
	return rpc.CreateIncomingContext(context.TODO(), db, shardID, node)
}
