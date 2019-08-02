package tsdb

import (
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/series"
)

//go:generate mockgen -source=./scan.go -destination=./scan_mock.go -package=tsdb -self_package=github.com/lindb/lindb/tsdb

type ScanContext struct {
	// required
	MetricID  uint32
	FieldIDs  []uint16
	TimeRange timeutil.TimeRange

	// optional, if SeriesIDSet is nil, just search metric level data
	SeriesIDSet *series.MultiVerSeriesIDSet
}

type Scanner interface {
	HasNext() bool
	Next() field.MultiTimeSeries
	Close()
}
