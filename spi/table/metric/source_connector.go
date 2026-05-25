// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package metric

import (
	"context"
	"errors"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/common/field"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	sfield "github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/utils"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
	metricstore "github.com/lindb/lindb/storage/metric"
	"github.com/lindb/lindb/storage/store"
)

var log = logger.GetLogger("Metric", "SourceConnector")

// sourceConnectorProvider is a factory that creates sourceConnector instances
// bound to a specific storage engine.
type sourceConnectorProvider struct {
	engine storage.Engine
}

// NewSourceConnectorProvider creates a SourceConnectorProvider backed by the given engine.
func NewSourceConnectorProvider(engine storage.Engine) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		engine: engine,
	}
}

// CreateSourceConnector builds a fully configured sourceConnector.
// The connector is lazy: no I/O happens until Run() is called.
func (p *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitionIDs []int,
	columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []arrow.Field, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	return &sourceConnector{
		ctx:           NewExecutionContext(ctx),
		engine:        p.engine,
		table:         table,
		partitionIDs:  partitionIDs,
		columnMapping: columnMapping,
		predicate:     predicate,
		outputColumns: outputColumns,
		assignments:   assignments,
	}
}

// sourceConnector orchestrates a single metric query:
//  1. Resolve schema & output columns → TableScan
//  2. Evaluate predicate tag filters (column-value lookup)
//  3. Fan out PartitionScan for each storage partition (parallelism via ExecutorPool)
//  4. Funnel results through reduceCh into a Reducer that emits Arrow RecordBatches
type sourceConnector struct {
	ctx *ExecutionContext

	engine storage.Engine

	table         spi.TableHandle
	partitionIDs  []int
	predicate     tree.Expression
	outputColumns []arrow.Field
	assignments   []*spi.ColumnAssignment
	columnMapping map[string]string

	reduceCh   chan any
	partitions []*Partition
}

// Run executes the metric query pipeline and writes results to output.
// It returns immediately after launching all tasks; the Reducer drains
// results synchronously on the caller's goroutine.
func (psc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	tableScan := psc.buildTableScan()
	if tableScan == nil {
		return
	}
	psc.partitions = psc.findPartitions(tableScan, psc.partitionIDs)
	if len(psc.partitions) == 0 {
		return
	}

	// Pre-evaluate the predicate: walk the expression tree to load tag value IDs
	// for each filter condition before partition scanning begins.
	if psc.predicate != nil {
		tableScan.predicate = psc.predicate
		lookup := NewColumnValuesLookVisitor(psc.ctx.ctx, tableScan)
		_ = psc.predicate.Accept(nil, lookup)
		if lookup.noResults {
			// A filter column or value was not found; return empty result early.
			return
		}
	}

	psc.reduceCh = make(chan any)

	// A separate goroutine waits for all inflight tasks to finish, then closes
	// reduceCh to signal the Reducer that no more data is coming.
	// This must run concurrently because the Reducer below blocks on the channel.
	// NOTE: goroutine is started after noResults check to avoid closing a nil channel.
	go func() {
		psc.ctx.Waiting()
		psc.close()
	}()

	// Launch one PartitionScan per partition on the shared MetaFetcher pool.
	for i := range psc.partitions {
		execute(psc.ctx, tableScan.db.ExecutorPool().MetaFetcher, func() {
			partitionScan := NewPartitionScan(psc.ctx, psc.partitions[i], psc.reduceCh)
			partitionScan.Run()
		})
	}

	// Reducer runs synchronously: it consumes from reduceCh until it is closed,
	// then builds the final Arrow RecordBatch and writes to output.
	reducer := NewReducer(psc.ctx, tableScan, psc.reduceCh, output)
	reducer.Run()
}

// close is called after all tasks complete. It signals the Reducer (by closing
// reduceCh) and releases field-data readers held by each partition.
func (psc *sourceConnector) close() {
	close(psc.reduceCh)

	for _, p := range psc.partitions {
		for _, rs := range p.fieldsData {
			rs.Close()
		}
	}
}

// getColumnName resolves a logical column name to its physical storage name
// using the provided mapping. Falls back to the original name when no mapping exists.
func getColumnName(name string, mapping map[string]string) string {
	if len(mapping) == 0 {
		return name
	}
	realName, ok := mapping[name]
	if !ok {
		return name
	}
	return realName
}

// buildTableScan constructs the TableScan that drives query execution.
// It resolves the metric schema, maps each output column to a field/tag/timestamp,
// and wires up the aggregation pipeline for every selected field.
// Returns nil when the metric does not exist (not-found is not an error).
func (psc *sourceConnector) buildTableScan() *TableScan {
	metricTable, ok := psc.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("metric provider not support table handle<%T>", psc.table))
	}
	db, ok := psc.engine.GetDatabase(metricTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, metricTable.Database))
	}

	metricID, schema, err := psc.getSchema(db.(*metricstore.Database), metricTable)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			// Metric doesn't exist yet; return empty result rather than an error.
			return nil
		}
		panic(err)
	}

	tableScan := &TableScan{
		db:            db.(*metricstore.Database),
		schema:        schema,
		metricID:      metricID,
		columnMapping: psc.columnMapping,
		outputs:       psc.outputColumns,
	}
	log.Info("buildTableScan: start",
		logger.Int("outputColumns", len(psc.outputColumns)),
		logger.Int("assignments", len(psc.assignments)))

	targetTimeRange, targetInterval := calcTimeRangeAndInterval(metricTable.TimeRange,
		metricTable.Interval, db.GetOption())

	// Pass 1: detect whether the timestamp column is selected so we can compute the
	// final interval before building field columns. Field columns read timeRange and
	// interval from tableScan during initialization (newColumn → initialize), so both
	// must be set before buildFieldColumn is called.
	for _, columnMeta := range psc.outputColumns {
		if arrow.TypeEqual(columnMeta.Type, arrow.FixedWidthTypes.Timestamp_ns) {
			tableScan.isTimestampSelected = true
			break
		}
	}

	// When the timestamp column is not selected (e.g. instant queries), collapse
	// the whole time range into a single interval so only one data point is produced.
	// Guard against Start == End after truncation (e.g. very narrow time window smaller
	// than the storage interval): fall back to targetInterval so the divisor is never zero.
	if !tableScan.isTimestampSelected {
		if collapsed := timeutil.Interval(targetTimeRange.End - targetTimeRange.Start); collapsed > 0 {
			targetInterval = collapsed
		}
		// else: keep targetInterval as computed by calcTimeRangeAndInterval
	}
	tableScan.timeRange = targetTimeRange
	tableScan.interval = targetInterval

	var (
		fields       sfield.Metas
		groupingTags tag.Metas
		fieldIndex   = uint8(0)
		numOfAggs    = 0
		// Track non-timestamp output columns; used to verify all columns were resolved.
		numOfDataColumns = len(psc.outputColumns)
	)

	// Pass 2: resolve each output column to a field, tag, or timestamp now that
	// tableScan.timeRange and tableScan.interval are fully initialized.
	// Output symbols are set to physical field names by push_aggregation so every
	// column should be resolvable via findFieldMeta (for fields) or tag lookup.
	for _, columnMeta := range psc.outputColumns {
		switch {
		case arrow.TypeEqual(columnMeta.Type, arrow.FixedWidthTypes.Timestamp_ns):
			// Timestamp column: mark it selected so the time axis is included in output.
			// Timestamp is not a "data" column, so exclude it from the resolution check.
			numOfDataColumns--

		case func() bool {
			// Inline field lookup to avoid calling lo.Find twice (matchField + findFieldMeta).
			fieldMeta, ok := psc.findFieldMeta(schema, columnMeta)
			if !ok {
				return false
			}
			fieldMeta.Index = fieldIndex

			ch, handleCount := psc.buildFieldColumn(tableScan, fieldMeta, numOfAggs)
			fields = append(fields, fieldMeta)
			tableScan.columns = append(tableScan.columns, ch)
			numOfAggs += handleCount
			fieldIndex++
			return true
		}():
			// Field column handled inside the closure above.

		default:
			// Tag column: collect for GROUP BY grouping.
			if tagKey, ok := lo.Find(schema.TagKeys, func(tagMeta tag.Meta) bool {
				return getColumnName(columnMeta.Name, psc.columnMapping) == tagMeta.Key &&
					arrow.TypeEqual(columnMeta.Type, arrow.BinaryTypes.String)
			}); ok {
				groupingTags = append(groupingTags, tagKey)
			} else {
				// Output column could not be resolved to a field or tag.
				// Log for debugging; decrement so the validation below doesn't count it.
				log.Info("buildTableScan: unresolved output column",
					logger.String("column", columnMeta.Name))
				numOfDataColumns--
			}
		}
	}

	// All non-timestamp output columns must resolve to either a field or a tag.
	if len(fields)+len(groupingTags) != numOfDataColumns {
		log.Info("buildTableScan: column mismatch, returning nil",
			logger.Int("fields", len(fields)),
			logger.Int("groupingTags", len(groupingTags)),
			logger.Int("expected", numOfDataColumns))
		return nil
	}

	tableScan.fields = fields
	if len(groupingTags) > 0 {
		tableScan.grouping = NewGrouping(db.(*metricstore.Database), groupingTags)
	}
	tableScan.numOfAggs = numOfAggs
	return tableScan
}

// findFieldMeta looks up the field.Meta for a given Arrow column.
func (psc *sourceConnector) findFieldMeta(schema *metric.Schema, columnMeta arrow.Field) (sfield.Meta, bool) {
	return lo.Find(schema.Fields, func(fieldMeta sfield.Meta) bool {
		return getColumnName(columnMeta.Name, psc.columnMapping) == fieldMeta.Name.String() &&
			!arrow.TypeEqual(columnMeta.Type, arrow.FixedWidthTypes.Timestamp_ns) &&
			!arrow.TypeEqual(columnMeta.Type, arrow.BinaryTypes.String)
	})
}

// buildFieldColumn creates a Column for a single metric field, wiring up downsampling
// and aggregation functions from the query assignments.
// Returns the column and the number of aggregation slots it occupies (one per handle).
// Histogram physical fields (bucket bounds, _sum, _count, etc.) are treated as regular
// fields with their native aggregation (sum); the broker's aggregation layer computes
// the final histogram function (quantile, avg, …) from the returned intermediate values.
func (psc *sourceConnector) buildFieldColumn(
	tableScan *TableScan,
	fieldMeta sfield.Meta,
	numOfAggs int,
) (Column, int) {
	// Find explicit aggregation handles for this field from the query plan.
	columnHandles := lo.Filter(psc.assignments, func(item *spi.ColumnAssignment, _ int) bool {
		return item.Column == fieldMeta.Name.String()
	})

	var handles []*ColumnHandle
	if len(columnHandles) == 0 {
		// No explicit handle: fall back to the field's native aggregation type.
		// Histogram aggregates as sum despite its type name being "histogram".
		funcName := tree.FuncName(fieldMeta.Type.String())
		if fieldMeta.Type == field.Histogram {
			funcName = tree.Sum
		}
		handles = []*ColumnHandle{{Downsampling: funcName, Aggregation: funcName}}
	}
	for _, columnHandle := range columnHandles {
		if handle, ok := columnHandle.Handler.(*ColumnHandle); ok {
			handles = append(handles, handle)
		}
	}

	var ch Column
	if fieldMeta.Type.IsExemplar() {
		// Exemplar fields use a fixed aggregation; the funcName is ignored.
		ch = newColumn(
			numOfAggs, tableScan, fieldMeta, handles, sfield.ExemplarAggregate,
			func(_ tree.FuncName) aggregateFunc[*models.Exemplar] {
				return sfield.ExemplarAggregate
			})
	} else {
		ch = newColumn(
			numOfAggs, tableScan, fieldMeta, handles, fieldMeta.Type.Aggregate,
			func(funcName tree.FuncName) aggregateFunc[float64] {
				return getAggFunc(funcName)
			})
	}
	return ch, len(handles)
}

// getSchema fetches the metric ID and full schema (fields + tag keys) from the metadata DB.
func (psc *sourceConnector) getSchema(db *metricstore.Database, table *TableHandle) (metric.ID, *metric.Schema, error) {
	metricID, err := db.MetaDB().GetMetricID(table.Namespace, table.Metric)
	if err != nil {
		return 0, nil, err
	}
	schema, err := db.MetaDB().GetSchema(metricID)
	if err != nil {
		return 0, nil, err
	}
	return metricID, schema, nil
}

// findPartitions maps the requested shard IDs to the concrete storage Partitions
// that overlap with the query's time range and interval.
func (psc *sourceConnector) findPartitions(tableScan *TableScan, shardIDs []int) (partitions []*Partition) {
	utils.FindSegments(tableScan.db, shardIDs, tableScan.interval, tableScan.timeRange,
		func(shard store.Shard, _ store.Partition, segments []store.Segment) {
			partitions = append(partitions, &Partition{
				tableScan: tableScan,
				shard:     shard.(*metricstore.Shard),
				segments: lo.Map(segments, func(item store.Segment, _ int) *metricstore.Segment {
					return item.(*metricstore.Segment)
				}),
			})
		})
	return
}
