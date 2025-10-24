package store

import (
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

type Shard interface {
	io.Closer
	// Database returns the database.
	Database() Database
	// ShardID returns the shard id.
	ShardID() models.ShardID

	GetOrCreatePartition(timestamp int64) (Partition, error)

	GetPartitions(interval timeutil.Interval, timeRange timeutil.TimeRange) []Partition
}
