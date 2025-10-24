package store

import (
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

type Segment interface {
	io.Closer
	Partition() Partition

	SegmentTimeRange() timeutil.TimeRange
	GetOrCreateWAL(leader models.NodeID) (WriteAheadLog, error)

	Retain()
	Release()

	ValidateSequence(leader models.NodeID, seq int64) bool
	CommitSequence(leader models.NodeID, seq int64)
	AckSequence(leader models.NodeID, fn func(seq int64))

	Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error)

	Flush() error
}
