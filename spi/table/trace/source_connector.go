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
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/lindb/arrow/pkg/arrow/builder"
	tracespkg "github.com/lindb/arrow/pkg/traces"

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
		engine:        s.engine,
		table:         table,
		partitionIDs:  partitions,
		predicate:     predicate,
		outputColumns: outputColumns,
	}
}

type sourceConnector struct {
	engine storage.Engine

	table spi.TableHandle

	partitionIDs  []int
	partitions    []*Partition
	outputColumns []arrow.Field

	predicate tree.Expression
}

// Run implements spi.SourceConnector.
// It performs a trace-ID point lookup and writes matching spans into a columnar
// RecordBatch whose schema is driven by the requested outputColumns.
func (sc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	tableScan := sc.buildTableScan()
	if tableScan == nil {
		return
	}
	sc.partitions = sc.findPartitions(tableScan, sc.partitionIDs)
	if len(sc.partitions) == 0 {
		return
	}

	// Extract the hex trace ID from the predicate (e.g. WHERE trace_id = '...')
	traceIDHex := sc.extractTraceID()
	if traceIDHex == "" {
		return
	}
	rawTraceID, err := hex.DecodeString(traceIDHex)
	if err != nil {
		return
	}

	// Build output schema and per-column appenders from the requested columns.
	rb := builder.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema(sc.outputColumns, nil))
	defer rb.Release()

	appenders := sc.buildColumnAppenders(rb)

	for _, partition := range sc.partitions {
		for _, segment := range partition.segments {
			logSegment := segment.(*tracestore.Segment)
			walBatches, err := logSegment.GetTrace(traceIDHex)
			if err != nil || len(walBatches) == 0 {
				continue
			}
			for _, msg := range walBatches {
				sc.appendMatchingSpans(msg, rawTraceID, appenders)
			}
		}
	}

	output <- rb.NewRecord()
}

// appendMatchingSpans reads one WAL batch (Arrow IPC), filters spans by trace ID,
// and appends matching rows to the output builders.
func (sc *sourceConnector) appendMatchingSpans(
	msg []byte, rawTraceID []byte, appenders []func(*tracespkg.TraceReader, int),
) {
	reader, err := tracespkg.NewTraceReader(msg)
	if err != nil {
		fmt.Printf("trace reader err: %v\n", err)
		return
	}
	defer reader.Release()

	for i := 0; i < reader.NumOfRows(); i++ {
		if !bytes.Equal(reader.TraceID(i), rawTraceID) {
			continue
		}
		for _, appender := range appenders {
			appender(reader, i)
		}
	}
}

// buildColumnAppenders creates one appender per output column using the trace schema registry.
func (sc *sourceConnector) buildColumnAppenders(rb *builder.RecordBuilder) []func(*tracespkg.TraceReader, int) {
	appenders := make([]func(*tracespkg.TraceReader, int), 0, len(sc.outputColumns))
	for _, col := range sc.outputColumns {
		if appender, ok := BuildTraceColumnAppender(col.Name, rb); ok {
			appenders = append(appenders, appender)
		}
	}
	return appenders
}

// extractTraceID reads the trace ID hex string from the predicate.
func (sc *sourceConnector) extractTraceID() string {
	expr, ok := sc.predicate.(*tree.ComparisonExpression)
	if !ok {
		return ""
	}
	evalCtx := expression.NewEvalContext(context.TODO())
	traceID, _ := expression.EvalString(evalCtx, expr.Right)
	return traceID
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
