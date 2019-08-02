package replication

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/storage"
)

const (
	defaultBufferSize                 = 32
	defaultSegmentDataFileSizeLimit   = 128 * 1024 * 1024
	defaultRemoveTaskIntervalInSecond = 60
)

var replicationConfig = config.ReplicationChannel{
	Path:                       "/tmp/broker/replication",
	BufferSize:                 defaultBufferSize,
	SegmentFileSize:            defaultSegmentDataFileSizeLimit,
	RemoveTaskIntervalInSecond: defaultRemoveTaskIntervalInSecond,
}

// mock write client
type mockWriteClientFactory struct {
	client storage.WriteService_WriteClient
}

func (m *mockWriteClientFactory) LogicNode() models.Node {
	panic("implement me")
}

func (m *mockWriteClientFactory) CreateWriteClient(db string, shardID uint32,
	remote models.Node) (storage.WriteService_WriteClient, error) {
	return m.client, nil
}

func (m *mockWriteClientFactory) CreateQueryClient(remote models.Node) (storage.QueryService_QueryClient, error) {
	panic("implement me")
}

func newMockWriteClientFactory(client storage.WriteService_WriteClient) rpc.ClientStreamFactory {
	return &mockWriteClientFactory{
		client: client,
	}
}

var node = models.Node{
	IP:   "123",
	Port: 123,
}

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

/***
send 0~10
send 10~20
recv reset 0
send 0~10
send 10~20
*/
func TestReplicator_ResetOnce(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "TestReplica_RestOnce")
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}

	}()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWriteCli := storage.NewMockWriteService_WriteClient(ctl)

	resetSeq := int64(0)
	endSeq := int64(20)

	req1, _ := buildWriteRequest(resetSeq, resetSeq+10)
	req2, _ := buildWriteRequest(resetSeq+10, endSeq)
	req3, _ := buildWriteRequest(resetSeq, resetSeq+10)
	req4, _ := buildWriteRequest(resetSeq+10, endSeq)

	done := make(chan struct{})

	respReset := &storage.WriteResponse{
		Seq: &storage.WriteResponse_ResetSeq{
			ResetSeq: resetSeq,
		},
	}

	respNil := &storage.WriteResponse{
		Seq: nil,
	}

	recv1 := mockWriteCli.EXPECT().Recv().Do(func() {
		<-done
	}).Return(respReset, nil)

	mockWriteCli.EXPECT().Recv().Do(func() {
		<-done
	}).Return(respNil, nil).After(recv1).AnyTimes()

	send1 := mockWriteCli.EXPECT().Send(req1).Return(nil)

	send2 := mockWriteCli.EXPECT().Send(req2).Do(func(interface{}) {
		close(done)
	}).Return(nil).After(send1)

	send3 := mockWriteCli.EXPECT().Send(req3).Return(nil).After(send2)

	mockWriteCli.EXPECT().Send(req4).Return(nil).After(send3)

	q, err := queue.NewFanOutQueue(tmpDir, defaultSegmentDataFileSizeLimit, defaultRemoveTaskIntervalInSecond)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < int(endSeq); i++ {
		if _, err := q.Append([]byte(strconv.Itoa(i))); err != nil {
			t.Fatal(err)
		}
	}

	fo, err := q.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	rep, err := newReplicator(node, "cluster", "database", 0, fo, newMockWriteClientFactory(mockWriteCli))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rep.Cluster(), "cluster")
	assert.Equal(t, rep.Database(), "database")
	assert.Equal(t, rep.ShardID(), uint32(0))
	assert.Equal(t, rep.Target(), node)

	// too short sleep time will crash the test
	time.Sleep(500 * time.Millisecond)

	q.Close()
}

/***
send 0~10
recv reset 5
send 5~10
*/
func TestReplicator_Batch(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "TestReplica_Batch")
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}

	}()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWriteCli := storage.NewMockWriteService_WriteClient(ctl)

	begSeq := int64(0)
	resetSeq := int64(5)
	endSeq := int64(10)

	req1, _ := buildWriteRequest(begSeq, endSeq)
	req2, _ := buildWriteRequest(resetSeq, endSeq)

	done := make(chan struct{})

	respReset := &storage.WriteResponse{
		Seq: &storage.WriteResponse_ResetSeq{
			ResetSeq: resetSeq,
		},
	}

	respNil := &storage.WriteResponse{
		Seq: nil,
	}

	recv1 := mockWriteCli.EXPECT().Recv().Do(func() {
		<-done
	}).Return(respReset, nil)

	mockWriteCli.EXPECT().Recv().Do(func() {
		<-done
	}).Return(respNil, nil).After(recv1).AnyTimes()

	send1 := mockWriteCli.EXPECT().Send(req1).Do(func(interface{}) {
		close(done)
	}).Return(nil)

	mockWriteCli.EXPECT().Send(req2).Return(nil).After(send1)

	q, err := queue.NewFanOutQueue(tmpDir, defaultSegmentDataFileSizeLimit, defaultRemoveTaskIntervalInSecond)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < int(endSeq); i++ {
		if _, err := q.Append([]byte(strconv.Itoa(i))); err != nil {
			t.Fatal(err)
		}
	}

	fo, err := q.GetOrCreateFanOut("f1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = newReplicator(node, "cluster", "database", 0, fo, newMockWriteClientFactory(mockWriteCli))
	if err != nil {
		t.Fatal(err)
	}

	// too short sleep time will crash the test
	time.Sleep(500 * time.Millisecond)

	q.Close()
}

/***
getClient cli1
cli1 recv err
getClient err
getClient cli2
close
cli2 recv err and return
*/
func TestReplicator_RecvError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWriteCli1 := storage.NewMockWriteService_WriteClient(ctl)
	mockWriteCli1.EXPECT().Recv().Return(nil, errors.New("recv error 1"))

	done := make(chan struct{})
	mockWriteCli2 := storage.NewMockWriteService_WriteClient(ctl)
	mockWriteCli2.EXPECT().Recv().DoAndReturn(func() (*storage.WriteResponse, error) {
		<-done
		return nil, errors.New("recv error 2")
	})

	mockFanOut := queue.NewMockFanOut(ctl)
	mockFanOut.EXPECT().Consume().Return(queue.SeqNoNewMessageAvailable).AnyTimes()

	mockFct := rpc.NewMockClientStreamFactory(ctl)
	gomock.InOrder(
		mockFct.EXPECT().CreateWriteClient("database", uint32(0), node).Return(mockWriteCli1, nil),
		mockFct.EXPECT().CreateWriteClient("database", uint32(0), node).Return(nil, errors.New("get client error")),
		mockFct.EXPECT().CreateWriteClient("database", uint32(0), node).Return(mockWriteCli2, nil),
	)

	rep, err := newReplicator(node, "cluster", "database", 0, mockFanOut, mockFct)
	if err != nil {
		t.Fatal(err)
	}

	// re-conn will wait 1 second.
	time.Sleep(100*time.Millisecond + time.Second)

	rep.Stop()
	close(done)

	// wait for sendLoop recvLoop to exit.
	time.Sleep(time.Second)
}
