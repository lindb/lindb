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

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	logproto "github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
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
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	fmt.Printf("create log source connector,table=%s,partitions=%v\n", outputColumns, assignments)
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

	outputColumns []types.ColumnMetadata
	assignments   []*spi.ColumnAssignment

	outputsHasTimestamp bool
	hasAggregate        bool

	aggregator Aggregator

	fieldKeys []uint32
	fields    []string

	logIDsBucket []*roaring.Bitmap // need create when stats logs(ouput time series data)
}

// Run implements spi.SourceConnector.
func (sc *sourceConnector) Run(output chan<- *types.Page) {
	tableScan := sc.buildTableScan()
	if tableScan == nil {
		fmt.Println("table scan is nil")
		return
	}
	sc.partitions = sc.findPartitions(tableScan, sc.partitionIDs)
	if len(sc.partitions) == 0 {
		fmt.Printf("table partition is nil,ids=%v\n", sc.partitionIDs)
		return
	}
	indexDB := tableScan.db.IndexDatabase()
	ns, err := indexDB.GetNamespaceID([]byte("ns"))
	if err != nil {
		fmt.Printf("rr1=%v\n", err)
		panic(ns)
	}
	tableScan.nsID = ns

	if sc.predicate != nil {
		fieldLookup := NewFieldValuesLookupVisitor(sc.ctx, tableScan)
		fieldLookup.Visit(sc.ctx, sc.predicate)
	}

	sc.initializeSearchContext(tableScan)

	page := types.NewPage()
	if sc.hasAggregate {
		sc.aggregator.Initialize()
		sc.aggregator.Aggregate(output)
		return
	}

	timeColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTInt, Name: "timestamp"}, timeColumn)
	msgColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "_msg"}, msgColumn)
	fieldsColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTJSON, Name: "fields"}, fieldsColumn)

	total := 0
	sc.findLogs(tableScan, func(segment *log.Segment, logIDs *roaring.Bitmap) bool {
		fmt.Printf("logSegment=%v,log ids:%v,%v\n", segment, logIDs)
		it := logIDs.ReverseIterator()
		for it.HasNext() {
			logID := it.Next()
			logData, err := segment.GetLog(logID)
			if err != nil {
				fmt.Printf("get log err:%v\n", err)
			} else {
				log := &flatLogV1.Log{}
				log.Init(logData, flatbuffers.GetUOffsetT(logData))
				timeColumn.AppendInt(log.Timestamp())
				// TODO:
				msgColumn.AppendString(string(log.Message()))
				fIt := logproto.NewFieldIterator(log)
				fMap := make(map[string]string)
				for fIt.HasNext() {
					fMap[string(fIt.NextName())] = string(fIt.NextValue())
				}
				fieldsColumn.AppendJSON(encoding.JSONMarshal(fMap))
				total++

				if total >= 1000 {
					// limit return
					return false
				}
			}
		}

		return true
	})

	output <- page
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
		if !sc.hasAggregate {
			// sort segments desc
			sort.Slice(partition.segments, func(i, j int) bool {
				return partition.segments[i].SegmentTimeRange().Start > partition.segments[j].SegmentTimeRange().Start
			})
		}

		logIDs := roaring.New()
		for _, segment := range partition.segments {
			logSegment := segment.(*logstore.Segment)
			logSegment.FindLogIDsByTimeRange(tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
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
	lo.ForEach(sc.outputColumns, func(item types.ColumnMetadata, index int) {
		if item.DataType == types.DTTimestamp && item.Name == constants.TimestampColumnName {
			sc.outputsHasTimestamp = true
		} else if item.DataType == types.DTDynamic {
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
	fmt.Printf("hahs ..fields=%v\n", sc.outputsHasTimestamp)

	if sc.hasAggregate {
		if sc.outputsHasTimestamp {
			sc.aggregator = newAggregatorByTime(sc, tableScan)
		} else {
			sc.aggregator = newAggregatorByField(sc, tableScan)
		}
	}
}
