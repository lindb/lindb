package memdb

import (
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestMetricStoreScanner_Close(t *testing.T) {
	s := newMetricStoreScanner(nil, nil, nil)
	err := s.Close()
	assert.NoError(t, err)
}

func TestMetricStoreScanner_Next(t *testing.T) {
	// case 1: series not exist
	s := newMetricStoreScanner(roaring.BitmapOf(10, 100).GetContainer(0), nil, nil)
	s.Scan(200)
}
