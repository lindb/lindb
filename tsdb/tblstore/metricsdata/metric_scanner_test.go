package metricsdata

import (
	"testing"

	"github.com/lindb/lindb/pkg/encoding"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
)

func TestMetricScanner_Next(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case 1: series id not exist
	s := newMetricScanner(nil, roaring.BitmapOf(10).GetContainer(0), nil)
	s.Scan(1)
	// case 2: read series data
	r := NewMockReader(ctrl)
	r.EXPECT().readSeriesData(gomock.Any())
	encoder := encoding.NewFixedOffsetEncoder()
	encoder.Add(100)
	data := encoder.MarshalBinary()
	seriesOffsets := encoding.NewFixedOffsetDecoder(data)
	s = newMetricScanner(r, roaring.BitmapOf(10).GetContainer(0), seriesOffsets)
	s.Scan(10)
}
