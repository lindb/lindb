package store

import (
	"io"

	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/timeutil"
)

type Partition interface {
	io.Closer

	PartitionTime() int64
	PartitionInterval() timeutil.Interval

	Path() string
	Shard() Shard

	GetOrCreateSegment(timestamp int64) (Segment, error)

	GetSegments(timeRange timeutil.TimeRange) []Segment
}

type Partitions struct {
	partitions map[int64]Partition // partition timestamp -> partition
}

func NewPartitions() *Partitions {
	return &Partitions{
		partitions: make(map[int64]Partition),
	}
}

func (ps *Partitions) GetPartition(timestamp int64) (Partition, bool) {
	p, ok := ps.partitions[timestamp]
	return p, ok
}

func (ps *Partitions) PutPartition(partition Partition) {
	ps.partitions[partition.PartitionTime()] = partition
}

func (ps *Partitions) GetPartitions() []Partition {
	return lo.Values(ps.partitions)
}
