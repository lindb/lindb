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

package trace

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/lindb/arrow/pkg/arrow/builder"
	"github.com/lindb/common/pkg/encoding"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/utils"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/store"
	tracestore "github.com/lindb/lindb/storage/trace"
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
		engine:       s.engine,
		table:        table,
		partitionIDs: partitions,

		predicate: predicate,
	}
}

type sourceConnector struct {
	engine storage.Engine

	table spi.TableHandle

	partitionIDs []int
	partitions   []*Partition

	predicate tree.Expression
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

	schema := arrow.NewSchema([]arrow.Field{
		{Name: "callstack", Type: arrow.BinaryTypes.Binary},
	}, nil)
	rb := builder.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	msgColumn := rb.Fields()[0].(*array.BinaryBuilder)

	expr, ok := sc.predicate.(*tree.ComparisonExpression)
	var traceID string
	if ok {
		evalCtx := expression.NewEvalContext(context.TODO())
		traceID, _ = expression.EvalString(evalCtx, expr.Right)
	}
	// TODO: check err

	for _, partition := range sc.partitions {
		for _, segment := range partition.segments {
			logSegment := segment.(*tracestore.Segment)
			logData, err := logSegment.GetTrace(traceID)
			if err != nil {
			} else if len(logData) > 0 {
				for _, msg := range logData {
					FilterTracesByTraceID(traceID, msg, msgColumn)
				}
			}
		}
	}

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
		db:        db,
		timeRange: logTable.GetTimeRange(),
	}
}

func (sc *sourceConnector) findPartitions(tableScan *TableScan, partitionIDs []int) (partitions []*Partition) {
	utils.FindSegments(tableScan.db, partitionIDs, sc.table.GetInterval(), tableScan.timeRange,
		func(shard store.Shard, partition store.Partition, segments []store.Segment) {
			partitions = append(partitions, &Partition{
				tableScan: tableScan,
				shard:     shard,
				segments:  segments,
			})
		})
	return
}

func FilterTracesByTraceID(traceID string, msg []byte, column *array.BinaryBuilder) {
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(msg); err != nil {
		return
	}
	traces := req.Traces()
	resourceSpans := traces.ResourceSpans()

	if resourceSpans.Len() == 0 {
		return
	}

	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		callStack := TranslateResourceSpans(rs, traceID)
		if callStack != nil {
			column.Append(encoding.JSONMarshal(callStack))
		}
	}
}
