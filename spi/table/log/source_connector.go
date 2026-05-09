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

package log

import (
	"context"
	"fmt"
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/arrow/builder"
	logspkg "github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/utils"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/log"
	logstore "github.com/lindb/lindb/storage/log"
	"github.com/lindb/lindb/storage/store"
)

type sourceConnectorProvider struct {
	engine storage.Engine
}

func NewSourceConnectorProvider(engine storage.Engine) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		engine: engine,
	}
}

// CreateSourceConnector implements spi.SourceConnectorProvider.
func (s *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int, columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []arrow.Field, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	return &sourceConnector{
		ctx:          ctx,
		engine:       s.engine,
		table:        table,
		partitionIDs: partitions,
		predicate:    predicate,

		assignments:   assignments,
		outputColumns: outputColumns,
	}
}

type sourceConnector struct {
	ctx    context.Context
	engine storage.Engine

	table        spi.TableHandle
	partitionIDs []int

	partitions []*Partition

	predicate tree.Expression

	outputColumns []arrow.Field
	assignments   []*spi.ColumnAssignment

	outputsHasTimestamp bool
	hasAggregate        bool

	aggregator Aggregator

	fieldKeys []uint32
	fields    []string

	logIDsBucket []*roaring.Bitmap // need create when stats logs(ouput time series data)
}

// Run implements spi.SourceConnector.
func (sc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	tableScan := sc.buildTableScan()
	if tableScan == nil {
		return
	}
	sc.partitions = sc.findPartitions(tableScan, sc.partitionIDs)
	if len(sc.partitions) == 0 {
		return
	}
	indexDB := tableScan.db.IndexDatabase()
	ns, err := indexDB.GetNamespaceID([]byte("ns"))
	if err != nil {
		panic(ns)
	}
	tableScan.nsID = ns

	if sc.predicate != nil {
		fieldLookup := NewFieldValuesLookupVisitor(sc.ctx, tableScan)
		fieldLookup.Visit(sc.ctx, sc.predicate)
		if fieldLookup.noResults {
			// A filter column or value was not found; return empty result early.
			return
		}
	}

	sc.initializeSearchContext(tableScan)

	if sc.hasAggregate {
		sc.aggregator.Initialize()
		sc.aggregator.Aggregate(output)
		return
	}

	rb := builder.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema([]arrow.Field{
		{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_ns},
		{Name: "_msg", Type: arrow.BinaryTypes.String},
		{Name: "fields", Type: arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String)},
	}, nil))
	defer rb.Release()

	timeColumn := rb.TimestampBuilder("timestamp")
	msgColumn := rb.StringBuilder("_msg")
	fieldsColumn := rb.MapBuilder("fields")
	fKey := fieldsColumn.KeyBuilder().(*array.StringBuilder)
	fValue := fieldsColumn.ItemBuilder().(*array.StringBuilder)

	total := 0
	sc.findLogs(tableScan, func(segment *log.Segment, logIDs *roaring.Bitmap) bool {
		scanner := log.NewScanner(segment, logIDs)
		defer scanner.Close()

		for scanner.HasNext() {
			if err := scanner.Next(func(reader *logspkg.Reader, rowNum int) {
				// read log data from reader
				timeColumn.Append(arrow.Timestamp(reader.Timestamp(rowNum)))
				msgColumn.Append(reader.Message(rowNum))
				reader.AttributesToMap(rowNum, fieldsColumn, fKey, fValue)
			}); err != nil {
				fmt.Printf("scan logs err:%v\n", err)
			}

			total++
			if total >= 1000 {
				// limit return
				return false
			}
		}
		return true
	})

	output <- rb.NewRecord()
}

func (sc *sourceConnector) buildTableScan() *TableScan {
	logTable, ok := sc.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("metric provider not support table handle<%T>", sc.table))
	}
	db, ok := sc.engine.GetDatabase(logTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, logTable.Database))
	}

	return &TableScan{
		db:        db.(*log.Database),
		timeRange: logTable.GetTimeRange(),
		interval:  logTable.GetInterval(),

		predicate: sc.predicate,
	}
}

func (sc *sourceConnector) findPartitions(tableScan *TableScan, partitionIDs []int) (partitions []*Partition) {
	utils.FindSegments(tableScan.db, partitionIDs, tableScan.interval, tableScan.timeRange,
		func(shard store.Shard, partition store.Partition, segments []store.Segment) {
			partitions = append(partitions, &Partition{
				tableScan:  tableScan,
				shard:      shard,
				paritition: partition,
				segments:   segments,
			})
		})
	return
}

func (sc *sourceConnector) findLogs(tableScan *TableScan,
	callback func(segment *log.Segment, logIDs *roaring.Bitmap) bool,
) {
	if !sc.hasAggregate {
		// sort partitions desc
		sort.Slice(sc.partitions, func(i, j int) bool {
			return sc.partitions[i].paritition.PartitionTime() > sc.partitions[j].paritition.PartitionTime()
		})
	}
	for _, partition := range sc.partitions {
		fmt.Println(partition.segments)
		if !sc.hasAggregate {
			// sort segments desc
			sort.Slice(partition.segments, func(i, j int) bool {
				return partition.segments[i].SegmentTimeRange().Start > partition.segments[j].SegmentTimeRange().Start
			})
		}

		logIDs := roaring.New()
		for _, segment := range partition.segments {
			logSegment := segment.(*logstore.Segment)
			fmt.Println("search logs...")
			logSegment.FindLogIDsByTimeRange(tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
				fmt.Println(logIDsFromStore)
				logIDs.Or(logIDsFromStore)
			})
			if sc.predicate != nil {
				rowLookup := NewRowLookupVisitor(tableScan, logSegment, logIDs)
				logIDsObj := rowLookup.Visit(sc.ctx, sc.predicate)
				if logIDsByPredicate, ok := logIDsObj.(*roaring.Bitmap); ok {
					logIDs.And(logIDsByPredicate)
				}
			}
			if logIDs == nil || logIDs.IsEmpty() {
				continue
			}

			if !callback(logSegment, logIDs) {
				return
			}

			logIDs.Clear()
		}
	}
}

func (sc *sourceConnector) initializeSearchContext(tableScan *TableScan) {
	sc.hasAggregate = lo.ContainsBy(sc.assignments, func(item *spi.ColumnAssignment) bool {
		if handle, ok := item.Handler.(*ColumnHandle); ok && handle.Aggregation != "" {
			return true
		}
		return false
	})

	indexDB := tableScan.db.IndexDatabase()
	lo.ForEach(sc.outputColumns, func(item arrow.Field, index int) {
		if item.Name == constants.TimestampColumnName {
			sc.outputsHasTimestamp = true
		} else if arrow.TypeEqual(item.Type, larrow.ExtensionTypes.Dynamic) {
			if len(sc.fieldKeys) == 1 {
				panic("too many grouping fields, only support one field")
			}
			fieldKey, err := indexDB.GetFieldKeyID(tableScan.nsID, []byte(item.Name))
			if err != nil {
				panic(fmt.Errorf("field:%s,err:=%w", item.Name, err))
			}
			sc.fieldKeys = append(sc.fieldKeys, fieldKey)
			sc.fields = append(sc.fields, item.Name)
		}
	})

	if sc.hasAggregate {
		if sc.outputsHasTimestamp {
			sc.aggregator = newAggregatorByTime(sc, tableScan)
		} else {
			sc.aggregator = newAggregatorByField(sc, tableScan)
		}
	}
}
