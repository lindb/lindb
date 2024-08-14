package metric

import (
	"fmt"
	"time"

	common_tileutil "github.com/lindb/common/pkg/timeutil"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
)

func init() {
	encoding.RegisterNodeType(MetricTableHandle{})
	spi.RegisterCreateTableHandleFn("metric", func(db, ns, name string) spi.TableHandle {
		return &MetricTableHandle{
			Database:  db,
			Namespace: ns,
			Metric:    name,
			TimeRange: timeutil.TimeRange{
				Start: time.Now().UnixMilli() - time.Hour.Milliseconds(),
				End:   time.Now().UnixMilli(),
			},
			Interval:        timeutil.Interval(10 * common_tileutil.OneSecond),
			StorageInterval: timeutil.Interval(10 * common_tileutil.OneSecond),
			IntervalRatio:   1,
		}
	})
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
	GroupingTags          tag.Metas
	ResultSet             []flow.FilterResultSet
	MinSeriesID           uint16
	MaxSeriesID           uint16
	HighSeriesID          uint16

	ShardExecuteContext *flow.ShardExecuteContext
}
