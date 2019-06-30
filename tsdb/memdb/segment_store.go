package memdb

import (
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/bufioutil"
)

// segmentStore stores field from a baseTime.
type segmentStore struct {
	writer   bufioutil.BufioWriter
	baseTime int64
}

// newSegmentStore returns a new segmentStore
func newSegmentStore(baseTime int64) *segmentStore {
	return &segmentStore{
		baseTime: baseTime,
	}
}

func (t *segmentStore) Write(slotTime int32, field models.Field) error {
	// todo: implement this
	_ = t.writer
	return nil
}
