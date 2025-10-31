package streaming

import (
	"fmt"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

func init() {
	// register table/column handle
	encoding.RegisterNodeType(TableHandle{})
	encoding.RegisterNodeType(ColumnHandle{})

	spi.RegisterCreateTableFn(spi.Log, func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
			App:    db,
			Stream: name,
		}
	})

	spi.RegisterApplyAggregationFn(spi.Streaming,
		func(table spi.TableHandle, tableMeta *types.TableMetadata,
			aggregations []spi.ColumnAggregation,
		) *spi.ApplyAggregationResult {
			result := &spi.ApplyAggregationResult{}
			// FIXME: find downSampling agg
			for _, agg := range aggregations {
				result.ColumnAssignments = append(result.ColumnAssignments,
					&spi.ColumnAssignment{Column: agg.Column, Handler: &ColumnHandle{Aggregation: agg.AggFuncName}},
				)
			}
			return result
		})
}

type TableHandle struct {
	App    string
	Stream string
}

func (t *TableHandle) SetTimeRange(timeRange timeutil.TimeRange) {
}

func (t *TableHandle) GetTimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{}
}

func (t *TableHandle) SetInterval(interval timeutil.Interval) {
}

func (t *TableHandle) GetInterval() timeutil.Interval {
	// fake interval
	return timeutil.Interval(10_000)
}

func (t *TableHandle) Kind() spi.DatasourceKind {
	return spi.Log
}

func (t *TableHandle) String() string {
	return fmt.Sprintf("%s:%s", t.App, t.Stream)
}

type ColumnHandle struct {
	Aggregation tree.FuncName `json:"aggregation"`
}

func (c *ColumnHandle) String() string {
	return fmt.Sprintf("(aggregation=%s)", c.Aggregation)
}
