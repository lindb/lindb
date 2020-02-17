package tsdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestDatabase_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
		newMetadataFunc = metadb.NewMetadata
		newKVStoreFunc = kv.NewStore
		newShardFunc = newShard
	}()
	// case 1: create kv store err
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err := newDatabase("db", testPath, &databaseConfig{
		Option: option.DatabaseOption{},
	})
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 2: create family err
	kvStore := kv.NewMockStore(ctrl)
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return kvStore, nil
	}
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	db, err = newDatabase("db", testPath, &databaseConfig{
		Option: option.DatabaseOption{},
	})
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 3: new metadata err
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	newMetadataFunc = func(ctx context.Context, name, parent string, tagFamily kv.Family) (metadata metadb.Metadata, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = newDatabase("db", testPath, &databaseConfig{
		Option: option.DatabaseOption{},
	})
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 4: create shard err
	newMetadataFunc = metadb.NewMetadata
	newShardFunc = func(db Database, shardID int32, shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = newDatabase("db", testPath, &databaseConfig{
		ShardIDs: []int32{1, 2, 3},
		Option:   option.DatabaseOption{},
	})
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 5: create db success
	newShardFunc = newShard
	db, err = newDatabase("db", testPath, &databaseConfig{
		ShardIDs: []int32{1, 2, 3},
		Option:   option.DatabaseOption{Interval: "10s"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, 3, db.NumOfShards())
	kvStore.EXPECT().Close().Return(nil).AnyTimes() // include shard close
	err = db.Close()
	assert.NoError(t, err)
	// case 6: close metadata err when create db
	// case 4: create shard err
	metadata := metadb.NewMockMetadata(ctrl)
	newMetadataFunc = func(ctx context.Context, name, parent string, tagFamily kv.Family) (metadb.Metadata, error) {
		return metadata, nil
	}
	newShardFunc = func(db Database, shardID int32, shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, fmt.Errorf("err")
	}
	metadata.EXPECT().Close().Return(fmt.Errorf("err"))
	db, err = newDatabase("db", testPath, &databaseConfig{
		ShardIDs: []int32{1, 2, 3},
		Option:   option.DatabaseOption{},
	})
	assert.Error(t, err)
	assert.Nil(t, db)
}

func Test_Database_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := kv.NewMockStore(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().Close().Return(nil)
	db := &database{
		metadata:  metadata,
		metaStore: mockStore}
	mockStore.EXPECT().Close().Return(nil).AnyTimes()

	// mock shard close error
	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(fmt.Errorf("error"))
	db.shards.Store(int32(1), mockShard)
	assert.Nil(t, db.Close())

	assert.Nil(t, db.IDGetter())
	assert.Nil(t, db.MetricMetaSuggester())
}

func Test_Database_FlushMeta(t *testing.T) {
	db := &database{
		isFlushing: *atomic.NewBool(false)}
	db.isFlushing.Store(true)
	assert.Nil(t, db.FlushMeta())

	db.Range(func(key, value interface{}) bool {
		assert.Fail(t, "")
		return true
	})
}
