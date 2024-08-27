package metric

import (
	"fmt"
	"time"

	common_tileutil "github.com/lindb/common/pkg/timeutil"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/tsdb"
)

func init() {
	encoding.RegisterNodeType(TableHandle{})
	spi.RegisterCreateTableHandleFn("metric", func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
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

type TableHandle struct {
	Database        string             `json:"database"`
	Namespace       string             `json:"namespace"`
	Metric          string             `json:"metric"`
	TimeRange       timeutil.TimeRange `json:"timeRange"`
	Interval        timeutil.Interval  `json:"interval"`
	StorageInterval timeutil.Interval  `json:"storageInterval"`
	IntervalRatio   int                `json:"intervalRatio"`
}

func (t *TableHandle) String() string {
	return fmt.Sprintf("%s:%s:%s", t.Database, t.Namespace, t.Metric)
}

// scan by low series ids
type ScanSplit struct {
	LowSeriesIDsContainer roaring.Container
	tableScan             *TableScan
	groupingContext       flow.GroupingContext
	ResultSet             []flow.FilterResultSet

	MinSeriesID  uint16
	MaxSeriesID  uint16
	HighSeriesID uint16
}

type Partition struct {
	shard    tsdb.Shard
	families []tsdb.DataFamily
}
