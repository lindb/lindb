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

	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
)

type sourceConnectorProvider struct {
	engine storage.Engine
}

func NewSourceConnectorProvider(engine storage.Engine) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		engine: engine,
	}
}

func (p *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitionIDs []int,
	columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
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

type sourceConnector struct {
	ctx *ExecutionContext

	engine storage.Engine

	table         spi.TableHandle
	partitionIDs  []int
	predicate     tree.Expression
	outputColumns []types.ColumnMetadata
	assignments   []*spi.ColumnAssignment
	columnMapping map[string]string

	reduceCh   chan any
	partitions []*Partition
}

func (psc *sourceConnector) Run(output chan<- *types.Page) {
	tableScan := psc.buildTableScan()
	if tableScan == nil {
		fmt.Println("table scan is nil")
		return
	}
	// find partitions
	psc.partitions = psc.findPartitions(tableScan, psc.partitionIDs)
	if len(psc.partitions) == 0 {
		return
	}

	go func() {
		// waiting data read complete
		psc.ctx.Waiting()
		// close connector
		psc.close()
	}()

	// lookup columns if predicate not nil
	if psc.predicate != nil {
		tableScan.predicate = psc.predicate
		lookup := NewColumnValuesLookVisitor(psc.ctx.ctx, tableScan)
		_ = psc.predicate.Accept(nil, lookup)
	}

	psc.reduceCh = make(chan any)

	for i := range psc.partitions {
		// fetch data for each partitions
		execute(psc.ctx, tableScan.db.ExecutorPool().MetaFetcher, func() {
			partitionScan := NewPartitionScan(psc.ctx, psc.partitions[i], psc.reduceCh)
			partitionScan.Run()
		})
	}

	// reduce task
	reducer := NewReducer(psc.ctx, tableScan, psc.reduceCh, output)
	reducer.Run()
}

func (psc *sourceConnector) close() {
	close(psc.reduceCh)
	fmt.Println("close rducer")

	// release the resources of query
	for _, p := range psc.partitions {
		for _, rs := range p.fieldsData {
			rs.Close()
		}
	}
}

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

func (psc *sourceConnector) buildTableScan() *TableScan {
	metricTable, ok := psc.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("metric provider not support table handle<%T>", psc.table))
	}
	db, ok := psc.engine.GetDatabase(metricTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, metricTable.Database))
	}
	// find table(metric) schema
	metricID, schema, err := psc.getSchema(db, metricTable)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil
		}
		// if isn't not found error, throw it
		panic(err)
	}
	// mapping fields for searching
	var fields field.Metas
	var columns []*column
	index := uint8(0)
	columnIndex := 0
	// mapping tags for grouping
	var groupingTags tag.Metas
	numOfOutputColumns := len(psc.outputColumns)
	isTimestampSelected := false
	lo.ForEach(psc.outputColumns, func(columnMeta types.ColumnMetadata, _ int) {
		if columnMeta.DataType == types.DTTimestamp {
			// timestamp
			isTimestampSelected = true
			numOfOutputColumns--
		} else if fieldMeta, ok := lo.Find(schema.Fields, func(fieldMeta field.Meta) bool {
			// field
			return getColumnName(columnMeta.Name, psc.columnMapping) == fieldMeta.Name.String() && columnMeta.DataType == types.DTTimeSeries
		}); ok {
			fieldMeta.Index = index
			fields = append(fields, fieldMeta)
			index++

			column := &column{meta: fieldMeta, offset: columnIndex}
			columns = append(columns, column)

			// find column handles for field
			columnHandles := lo.Filter(psc.assignments, func(item *spi.ColumnAssignment, index int) bool {
				return item.Column == fieldMeta.Name.String()
			})
			// not input aggregation func for this field
			if len(columnHandles) == 0 {
				// if column handle not found, set default aggregation using field aggregation type
				funcName := tree.FuncName(fieldMeta.Type.String()) // TODO: using same type
				column.handles = []*ColumnHandle{{Downsampling: funcName, Aggregation: funcName}}
			}
			for _, columnHandle := range columnHandles {
				if handle, ok := columnHandle.Handler.(*ColumnHandle); ok {
					column.handles = append(column.handles, handle)
					columnIndex++
				}
			}
		} else if tagKey, ok := lo.Find(schema.TagKeys, func(tagMeta tag.Meta) bool {
			// tag
			return getColumnName(columnMeta.Name, psc.columnMapping) == tagMeta.Key && columnMeta.DataType == types.DTString
		}); ok {
			groupingTags = append(groupingTags, tagKey)
		}
	})
	fmt.Printf("all fields=%v, group key=%v, select field=%v,output=%v\n", schema.Fields, groupingTags, fields, psc.outputColumns)

	if len(fields)+len(groupingTags) != numOfOutputColumns {
		// TODO: only check grouping keys
		// output columns size not match
		fmt.Println("output columns not match......")
		return nil
	}

	var grouping *Grouping
	if len(groupingTags) > 0 {
		grouping = NewGrouping(db, groupingTags)
	}
	maxOfRollups := 0
	numOfAggs := 0
	for _, column := range columns {
		// init column rollup and aggregation context
		column.init()
		if maxOfRollups < len(column.rollups) {
			maxOfRollups = len(column.rollups)
		}
		numOfAggs += len(column.aggs)
	}

	targetTimeRange, targetInterval := calcTimeRangeAndInterval(metricTable.TimeRange,
		metricTable.Interval, db.GetConfig()) // TODO: move to plan?
	fmt.Printf("time range=%v,interval=%v\n", targetTimeRange, targetInterval)
	if !isTimestampSelected {
		targetInterval = timeutil.Interval(targetTimeRange.End - targetTimeRange.Start)
	}

	return &TableScan{
		db:                  db,
		schema:              schema,
		metricID:            metricID,
		isTimestampSelected: isTimestampSelected,
		timeRange:           targetTimeRange,
		interval:            targetInterval,
		fields:              fields,
		columns:             columns,
		columnMapping:       psc.columnMapping,
		maxOfRollups:        maxOfRollups,
		numOfAggs:           numOfAggs,
		grouping:            grouping,
		outputs:             psc.outputColumns,
	}
}

// getSchema returns table schema based on table handle.
func (psc *sourceConnector) getSchema(db storage.Database, table *TableHandle) (metric.ID, *metric.Schema, error) {
	// find metric id(table id)
	metricID, err := db.MetaDB().GetMetricID(table.Namespace, table.Metric)
	if err != nil {
		return 0, nil, err
	}
	// find table schema
	schema, err := db.MetaDB().GetSchema(metricID)
	if err != nil {
		return 0, nil, err
	}
	return metricID, schema, nil
}

func (psc *sourceConnector) findPartitions(tableScan *TableScan, partitionIDs []int) (partitions []*Partition) {
	storageInterval := tableScan.db.GetConfig().Option.FindMatchSmallestInterval(tableScan.interval)
	for _, id := range partitionIDs {
		shard, ok := tableScan.db.GetShard(models.ShardID(id))
		if ok {
			families := shard.GetDataFamilies(storageInterval.Type(), tableScan.timeRange)
			if len(families) > 0 {
				partitions = append(partitions, &Partition{
					tableScan: tableScan,
					shard:     shard,
					families:  families,
				})
			}
		}
	}
	return
}
