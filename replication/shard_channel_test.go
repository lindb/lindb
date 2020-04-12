package replication

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

func TestChannel_New(t *testing.T) {
	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, "database", ch.Database())
	assert.Equal(t, int32(1), ch.ShardID())

	defer func() {
		newFanOutQueue = queue.NewFanOutQueue
	}()
	newFanOutQueue = func(dirPath string, dataFileSize int, removeTaskInterval time.Duration) (queue.FanOutQueue, error) {
		return nil, fmt.Errorf("err")
	}
	ch, err = newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func TestChannel_GetOrCreateReplicator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()
	target := models.Node{IP: "1.1.1.1", Port: 12345}
	r, err := ch.GetOrCreateReplicator(target)
	assert.NoError(t, err)
	assert.Equal(t, target, r.Target())

	r2, err := ch.GetOrCreateReplicator(target)
	assert.NoError(t, err)
	assert.Equal(t, r, r2)

	assert.Len(t, ch.Targets(), 1)
	assert.Equal(t, target, ch.Targets()[0])

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().GetOrCreateFanOut(gomock.Any()).Return(nil, fmt.Errorf("err"))
	ch1.q = fanout
	r2, err = ch.GetOrCreateReplicator(models.Node{IP: "err", Port: 12345})
	assert.Error(t, err)
	assert.Nil(t, r2)
	cancel()
	time.Sleep(300 * time.Millisecond)
}

func TestChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().Append(gomock.Any()).Return(int64(0), fmt.Errorf("err")).AnyTimes()
	ch1.q = fanout

	metric := &pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)
	err = ch.Write(metric)
	assert.NoError(t, err)

	cancel()
	time.Sleep(time.Millisecond * 600)

	ch, err = newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch1 = ch.(*channel)
	// ignore data, after closed
	chunk := NewMockChunk(ctrl)
	ch1.chunk = chunk
	// make sure chan is full
	ch1.ch <- []byte{1, 2}
	ch1.ch <- []byte{1, 2}
	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return([]byte{1, 2, 3}, nil)
	err = ch.Write(metric)
	assert.Error(t, err)
	time.Sleep(time.Millisecond * 500)
}

func TestChannel_checkFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := newChannel(ctx, replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	ch.Startup()

	ch1 := ch.(*channel)
	fanout := queue.NewMockFanOutQueue(ctrl)
	fanout.EXPECT().Append(gomock.Any()).Return(int64(0), fmt.Errorf("err")).AnyTimes()
	ch1.q = fanout

	metric := &pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)

	time.Sleep(time.Second)
	cancel()
	time.Sleep(300 * time.Millisecond)
}

func TestChannel_write_pending_before_close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	metric := &pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}
	err = ch.Write(metric)
	assert.NoError(t, err)

	ch1 := ch.(*channel)
	ch1.ch <- []byte{1, 2, 3}
	fanOut := queue.NewMockFanOutQueue(ctrl)
	fanOut.EXPECT().Append(gomock.Any()).Return(int64(0), fmt.Errorf("err")).AnyTimes()
	ch1.q = fanOut
	ch1.writePendingBeforeClose()
}

func TestChannel_chunk_marshal_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newChannel(context.TODO(), replicationConfig, "database", 1, nil)
	assert.NoError(t, err)
	chunk := NewMockChunk(ctrl)
	ch1 := ch.(*channel)
	ch1.chunk = chunk

	metric := &pb.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}
	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	err = ch.Write(metric)
	assert.Error(t, err)

	chunk.EXPECT().Append(gomock.Any())
	chunk.EXPECT().IsFull().Return(true)
	chunk.EXPECT().MarshalBinary().Return(nil, nil)
	err = ch.Write(metric)
	assert.NoError(t, err)

	chunk.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	ch1.flushChunk()
	chunk.EXPECT().MarshalBinary().Return(nil, nil)
	ch1.flushChunk()
}
