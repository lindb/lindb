package metric

import (
	"fmt"

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
	Database        string             `json:"database"`
	Namespace       string             `json:"namespace"`
	Metric          string             `json:"metric"`
	TimeRange       timeutil.TimeRange `json:"timeRange"`
	Interval        timeutil.Interval  `json:"interval"`
	StorageInterval timeutil.Interval  `json:"storageInterval"`
	IntervalRatio   int                `json:"intervalRatio"`
}

func (t *MetricTableHandle) String() string {
	return fmt.Sprintf("%s:%s:%s", t.Database, t.Namespace, t.Metric)
}

type MetricScanSplit struct {
	LowSeriesIDsContainer roaring.Container
	Fields                field.Metas
	ResultSet             []flow.FilterResultSet
	MinSeriesID           uint16
	MaxSeriesID           uint16
	HighSeriesID          uint16
}
