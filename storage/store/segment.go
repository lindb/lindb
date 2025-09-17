package store

import (
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/wal"
)

type Segment interface {
	io.Closer

	GetOrCreateWAL(leader models.NodeID) (wal.WriteAheadLog, error)
}
