package log

import (
	"fmt"
	"path"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
)

type shard struct {
	base.Shard

	database *Database
	dir      string
}

func NewShard(id models.ShardID, database *Database) (store.Shard, error) {
	shardPath := store.ShardPath(database.Name(), id)
	if err := fileutil.MkDirIfNotExist(shardPath); err != nil {
		return nil, err
	}
	s := &shard{
		Shard: base.Shard{
			ID:                  id,
			CalcPartitionTimeFn: intervalCalc.CalcSegmentTime,
			Partitions:          store.NewPartitions(),
		},
		database: database,
		dir:      shardPath,
	}
	s.CreatePartitionFn = s.createPartition

	partitions, err := fileutil.ListDir(path.Join(shardPath, store.PartitionDir, minuteInterval.Type().String()))
	if err != nil {
		return nil, err
	}
	for _, partition := range partitions {
		fmt.Println(partition)
		partitionTime, err := intervalCalc.ParseSegmentTime(partition)
		if err != nil {
			return nil, err
		}
		_, err = s.GetOrCreatePartition(partitionTime)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *shard) Database() store.Database {
	return s.database
}

func (s *shard) Path() string {
	return s.dir
}

func (s *shard) createPartition(timestamp int64) (store.Partition, error) {
	return NewPartition(timestamp, s)
}
