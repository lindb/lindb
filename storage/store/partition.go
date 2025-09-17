package store

import (
	"io"

	"github.com/lindb/lindb/pkg/timeutil"
)

type Partition interface {
	io.Closer

	PartitionTime() int64

	Path() string
	Shard() Shard

	GetOrCreateSegment(timestamp int64) (Segment, error)

	GetSegments(timeRange timeutil.TimeRange) []Segment
}
