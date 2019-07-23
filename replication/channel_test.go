package replication

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/rpc/proto/storage"
)

func mockWriteClient(ctl *gomock.Controller) *storage.MockWriteService_WriteClient {
	mockWriteCli := storage.NewMockWriteService_WriteClient(ctl)
	mockWriteCli.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	mockWriteCli.EXPECT().Recv().Return(
		&storage.WriteResponse{
			Seq: &storage.WriteResponse_AckSeq{
				AckSeq: 0,
			},
		}, nil).AnyTimes()

	return mockWriteCli
}

func TestNewChannel(t *testing.T) {

}

func TestChannelManager_GetChannel(t *testing.T) {
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Fatal(err)
		}
	}()

	replicationConfig.Path = dirPath

	ctl := gomock.NewController(t)
	cm := NewChannelManager(replicationConfig, &mockWriteClientFactory{client: mockWriteClient(ctl)})
	defer ctl.Finish()

	if _, err := cm.GetChannel("cluster", "database", 0); err == nil {
		t.Fatal("should be error")
	}

	_, err := cm.CreateChannel("cluster", "database", 2, 2)
	if err == nil {
		t.Fatal("should be error")
	}

	ch1, err := cm.CreateChannel("cluster", "database", 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	ch11, err := cm.GetChannel("cluster", "database", 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cm.CreateChannel("cluster", "database", 2, 1)
	if err == nil {
		t.Fatal(" should be error")
	}

	ch111, err := cm.CreateChannel("cluster", "database", 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ch11, ch1)
	assert.Equal(t, ch111, ch1)

	_, err = cm.GetChannel("cluster", "database", 1)
	if err == nil {
		t.Fatal("should be error")
	}

	cm.Close()
}

func TestChannel_GetOrCreateReplicator(t *testing.T) {
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Fatal(err)
		}
	}()

	replicationConfig.Path = dirPath

	ctl := gomock.NewController(t)
	cm := NewChannelManager(replicationConfig, &mockWriteClientFactory{client: mockWriteClient(ctl)})
	defer ctl.Finish()

	ch, err := cm.CreateChannel("cluster", "database", 2, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(ch.Targets()), 0)

	assert.Equal(t, ch.Cluster(), "cluster")
	assert.Equal(t, ch.Database(), "database")
	assert.Equal(t, ch.ShardID(), uint32(0))

	node := models.Node{
		IP:   "1.1.1.1",
		Port: 8000,
	}
	rep1, err := ch.GetOrCreateReplicator(node)
	if err != nil {
		t.Fatal(err)
	}

	rep11, err := ch.GetOrCreateReplicator(node)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rep1, rep11)
	assert.Equal(t, len(ch.Targets()), 1)

	cm.Close()
}

func TestChannel_Write(t *testing.T) {
	dirPath := path.Join(os.TempDir(), "test_channel_manager")
	if err := os.RemoveAll(dirPath); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(dirPath); err != nil {
			t.Fatal(err)
		}
	}()

	replicationConfig.Path = dirPath

	ctl := gomock.NewController(t)
	cm := NewChannelManager(replicationConfig, &mockWriteClientFactory{client: mockWriteClient(ctl)})
	defer ctl.Finish()

	ch, err := cm.CreateChannel("cluster", "database", 2, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(ch.Targets()), 0)

	node := models.Node{
		IP:   "1.1.1.1",
		Port: 8000,
	}
	rep1, err := ch.GetOrCreateReplicator(node)
	if err != nil {
		t.Fatal(err)
	}

	if err := ch.Write(context.TODO(), []byte("123")); err != nil {
		t.Fatal(err)
	}

	// wait for replication
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, rep1.Pending(), int64(0))

	cm.Close()

}
