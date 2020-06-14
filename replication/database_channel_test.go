package replication

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

func TestDatabaseChannel_new(t *testing.T) {
	defer func() {
		mkdir = fileutil.MkDirIfNotExist
	}()
	mkdir = func(path string) error {
		return fmt.Errorf("err")
	}
	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 10, nil)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func TestDatabaseChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 1, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	err = ch.Write(&pb.MetricList{Metrics: []*pb.Metric{
		{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			Fields: []*pb.Field{{
				Name:  "f1",
				Type:  pb.FieldType_Sum,
				Value: 1.0,
			}},
			Tags: map[string]string{"host": "1.1.1.1"},
		},
	}})
	assert.Equal(t, errChannelNotFound, err)

	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(int32(0), shardCh)

	shardCh.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
	err = ch.Write(&pb.MetricList{Metrics: []*pb.Metric{
		{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			Fields: []*pb.Field{{
				Name:  "f1",
				Type:  pb.FieldType_Sum,
				Value: 1.0,
			}},
			Tags: map[string]string{"host": "1.1.1.1"},
		},
	}})
	assert.Error(t, err)
}

func TestDatabaseChannel_CreateChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 4, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(int32(0), shardCh)
	shardCh2, err := ch.CreateChannel(int32(1), int32(0))
	assert.NoError(t, err)
	assert.Equal(t, shardCh, shardCh2)

	_, err = ch.CreateChannel(0, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(1, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(2, 1)
	assert.Equal(t, errInvalidShardNum, err)

	_, err = ch.CreateChannel(4, 1)
	assert.NoError(t, err)

	defer func() {
		createChannel = newChannel
	}()
	createChannel = func(cxt context.Context,
		cfg config.ReplicationChannel, database string, shardID int32,
		fct rpc.ClientStreamFactory,
	) (i Channel, e error) {
		return nil, fmt.Errorf("err")
	}

	_, err = ch.CreateChannel(4, 2)
	assert.Error(t, err)

	ch1.shardChannels.Store(int32(3), "test")
	c, ok := ch1.getChannelByShardID(int32(3))
	assert.False(t, ok)
	assert.Nil(t, c)
}

func TestDatabaseChannel_ReplicaState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(), "test-db", replicationConfig, 3, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)

	shardCh0 := NewMockChannel(ctrl)
	shardCh1 := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.shardChannels.Store(int32(0), shardCh0)
	ch1.shardChannels.Store(int32(1), shardCh1)

	shardCh0.EXPECT().Targets().Return(nil)
	shardCh1.EXPECT().Targets().Return([]models.Node{{IP: "1.1.1.1", Port: 12345}, {IP: "2.2.2.2", Port: 12345}})
	shardCh1.EXPECT().GetOrCreateReplicator(models.Node{IP: "1.1.1.1", Port: 12345}).Return(nil, fmt.Errorf("err"))
	replicator := NewMockReplicator(ctrl)
	shardCh1.EXPECT().GetOrCreateReplicator(models.Node{IP: "2.2.2.2", Port: 12345}).Return(replicator, nil)
	replicator.EXPECT().Database().Return("db")
	replicator.EXPECT().ShardID().Return(int32(1))
	replicator.EXPECT().Pending().Return(int64(0))
	replicator.EXPECT().ReplicaIndex().Return(int64(0))
	replicator.EXPECT().AckIndex().Return(int64(0))

	replicaState := ch.ReplicaState()
	assert.Len(t, replicaState, 1)
}
