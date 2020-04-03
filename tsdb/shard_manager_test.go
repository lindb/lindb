package tsdb

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestShardManager_AddShard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	shard := NewMockShard(ctrl)
	shard.EXPECT().ShardInfo().Return("shardInfo").AnyTimes()
	GetShardManager().AddShard(shard)

	GetShardManager().WalkEntry(func(theShard Shard) {
	})
	GetShardManager().RemoveShard(shard)
}
