package base

import (
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

type CreatePartitionFn func(timestamp int64) (store.Partition, error)

type CalcPartitionTimeFn func(timestamp int64) (partitionTime int64)

type Shard struct {
	ID                  models.ShardID
	CreatePartitionFn   CreatePartitionFn
	CalcPartitionTimeFn CalcPartitionTimeFn

	Partitions *store.Partitions // partition timestamp -> partition
	mutex      sync.Mutex
}

func (s *Shard) GetOrCreatePartition(timestamp int64) (store.Partition, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	partition, ok := s.Partitions.GetPartition(s.CalcPartitionTimeFn(timestamp))
	if ok {
		return partition, nil
	}

	partition, err := s.CreatePartitionFn(timestamp)
	if err != nil {
		return nil, err
	}

	s.Partitions.PutPartition(partition)

	return partition, nil
}

// ShardID implements store.Shard.
func (s *Shard) ShardID() models.ShardID {
	return s.ID
}

func (s *Shard) GetPartitions(interval timeutil.Interval, timeRange timeutil.TimeRange) (partitions []store.Partition) {
	s.mutex.Lock()
	partitions = s.Partitions.GetPartitions()
	s.mutex.Unlock()
	return
}

func (s *Shard) Close() error {
	var partitions []store.Partition

	s.mutex.Lock()
	partitions = s.Partitions.GetPartitions()
	s.mutex.Unlock()

	for _, partition := range partitions {
		partition.Close()
	}
	return nil
}
