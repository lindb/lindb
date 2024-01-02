package metric

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func init() {
	encoding.RegisterNodeType(MetricTableHandle{})
}

type MetricTableHandle struct {
	Namespace       string             `json:"namespace"`
	Metric          string             `json:"metric"`
	GroupBy         []string           `json:"groupBy"`
	TimeRange       timeutil.TimeRange `json:"timeRange"`
	Interval        timeutil.Interval
	StorageInterval timeutil.Interval
	IntervalRatio   int
}

type MetricScanSplit struct {
	LowSeriesIDsContainer roaring.Container
	Fields                field.Metas
	ResultSet             []flow.FilterResultSet
	MinSeriesID           uint16
	MaxSeriesID           uint16
	HighSeriesID          uint16
}
