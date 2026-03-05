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

package stream

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

type sourceConnectorProvider struct{}

func NewSourceConnectorProvider() spi.SourceConnectorProvider {
	return &sourceConnectorProvider{}
}

func (s *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int, columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []arrow.Field, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	tableHandle := table.(*TableHandle)
	schema, err := GetManager().GetStreamManager(tableHandle.Database).GetTableMetadata(tableHandle.Database, "", tableHandle.Stream)
	if err != nil {
		panic(err)
	}

	inputHandle := input.GetManager().GetInputHandler(tableHandle.Database, tableHandle.Stream)
	connector := &sourceConnector{
		ctx: ctx,

		table:   tableHandle,
		input:   inputHandle,
		running: atomic.NewBool(true),

		schema:        schema.Schema,
		predicate:     predicate,
		outputColumns: outputColumns,

		inbound: operator.NewQueue(make(chan arrow.RecordBatch, 256)),

		logger: logger.GetLogger("CEP", "SourceConnector"),
	}
	connector.initialize()

	inputHandle.Subscribe(connector)

	return connector
}

type sourceConnector struct {
	ctx context.Context

	table   *TableHandle
	input   input.InputHandler
	running *atomic.Bool

	schema        *arrow.Schema
	outputColumns []arrow.Field
	refs          []int // column index in input record for output columns, -1 means not exist

	recordSchema *arrow.Schema

	inbound *operator.Queue

	filter    *filter
	predicate tree.Expression

	logger logger.Logger
}

func (sc *sourceConnector) initialize() {
	index := 0
	// TODO: maybe income the columns of income record is diff schema.
	fields := sc.schema.Fields()
	columnMap := lo.Associate(fields, func(item arrow.Field) (string, int) {
		i := index
		index++
		return item.Name, i
	})
	sc.refs = make([]int, len(sc.outputColumns))
	for i := range sc.outputColumns {
		colMeta := &sc.outputColumns[i]
		colIndex, ok := columnMap[colMeta.Name]
		if !ok {
			sc.refs[i] = -1
			continue
		}
		sc.refs[i] = colIndex
	}

	sc.recordSchema = arrow.NewSchema(sc.outputColumns, nil)

	if sc.predicate != nil {
		sc.filter = newFilter(sc.schema, sc.predicate, sc)
	}
}

func (sc *sourceConnector) Receive(event models.Event) {
	if !sc.running.Load() {
		sc.logger.Warn("source connector has been stopped, drop received event",
			logger.String("database", sc.table.Database),
			logger.String("table", sc.table.Stream))
		return
	}
	if record, ok := event.(arrow.RecordBatch); ok {
		sc.inbound.Produce(record)
	}
}

func (sc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	defer func() {
		if err := recover(); err != nil {
			sc.logger.Error("source connector panicked", logger.Any("error", err), logger.Stack())
		}
		if sc.running.CompareAndSwap(true, false) {
			// unsubscribe input handler
			sc.input.Unsubscribe(sc)
			sc.inbound.Close()

			sc.logger.Info("source connector stopped",
				logger.String("database", sc.table.Database),
				logger.String("table", sc.table.Stream))
		}
	}()

	for {
		record, ok := sc.inbound.Consume(sc.ctx)
		if !ok {
			break
		}
		sc.process(record, output)
	}
}

func (sc *sourceConnector) process(record arrow.RecordBatch, output chan<- arrow.RecordBatch) {
	defer record.Release()

	var mark []uint32
	if sc.filter != nil {
		result, err := sc.filter.eval(record)
		if err != nil {
			sc.logger.Error("failed to evaluate predicate", logger.Error(err))
			return
		}
		if result.IsEmpty() {
			return
		}
		mark = result.ToArray()
	}

	columns := make([]arrow.Array, len(sc.outputColumns))
	for i, ref := range sc.refs {
		if ref >= 0 {
			columns[i] = record.Column(ref)
		}
	}
	rs := array.NewRecordBatch(sc.recordSchema, columns, record.NumRows())
	output <- larrow.NewFilterableRecord(rs, mark)
}
