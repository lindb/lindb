package tsdb

import (
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/timeutil"
	"github.com/eleme/lindb/tsdb/series"
)

//go:generate mockgen -source=./scan.go -destination=./scan_mock.go -package=tsdb -self_package=github.com/eleme/lindb/tsdb

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
