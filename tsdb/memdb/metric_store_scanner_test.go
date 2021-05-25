package memdb

import (
	"testing"

	"github.com/lindb/roaring"
)

func TestMetricStoreScanner_Next(t *testing.T) {
	// case 1: series not exist
	s := newMetricStoreScanner(roaring.BitmapOf(10, 100).GetContainer(0),
		nil, nil)
	s.Scan(200)
}
