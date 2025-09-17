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

	Partitions map[int64]store.Partition
	mutex      sync.Mutex
}

func (s *Shard) GetOrCreatePartition(timestamp int64) (store.Partition, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	partition, ok := s.Partitions[s.CalcPartitionTimeFn(timestamp)]
	if ok {
		return partition, nil
	}

	partition, err := s.CreatePartitionFn(timestamp)
	if err != nil {
		return nil, err
	}

	s.Partitions[partition.PartitionTime()] = partition

	return partition, nil
}

// ShardID implements store.Shard.
func (s *Shard) ShardID() models.ShardID {
	return s.ID
}

func (s *Shard) GetPartitions(timeRange timeutil.TimeRange) (partitions []store.Partition) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, partition := range s.Partitions {
		partitions = append(partitions, partition)
	}
	return
}

func (s *Shard) Close() error {
	for _, partition := range s.Partitions {
		partition.Close()
	}
	return nil
}
