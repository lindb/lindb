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

package infoschema

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/lindb/arrow/pkg/arrow/builder"
	"github.com/samber/lo"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

type sourceConnectorProvider struct {
	metadataMgr meta.MetadataManager
}

func NewSourceConnectorProvider(
	metadataMgr meta.MetadataManager,
) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		metadataMgr: metadataMgr,
	}
}

func (p *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int,
	columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []arrow.Field, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	return &sourceConnector{
		ctx:           ctx,
		table:         table,
		predicate:     predicate,
		outputColumns: outputColumns,
		reader:        NewReader(p.metadataMgr),
	}
}

type sourceConnector struct {
	ctx    context.Context
	reader Reader

	table         spi.TableHandle
	tableHandle   *TableHandle
	predicate     tree.Expression
	outputColumns []arrow.Field
	colIdxs       []int

	rb *builder.RecordBuilder
}

func (p *sourceConnector) open() {
	infoTable, ok := p.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("information schema provider not support table handle<%T>", p.table))
	}
	schema, ok := GetTableSchema(infoTable.Table)
	if !ok {
		panic(fmt.Errorf("information table schema not found: %s", infoTable.Table))
	}
	p.colIdxs = make([]int, len(p.outputColumns))
	fields := schema.Fields()
	for i, col := range p.outputColumns {
		if _, idx, exist := lo.FindIndexOf(fields, func(item arrow.Field) bool {
			return item.Name == col.Name
		}); exist {
			p.colIdxs[i] = idx
		}
	}
	if len(p.colIdxs) != len(p.outputColumns) {
		panic("output columns not found in table schema")
	}
	p.rb = builder.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema(p.outputColumns, nil))
	p.tableHandle = infoTable
}

func (p *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	defer func() {
		if p.rb != nil {
			p.rb.Release()
		}
	}()

	p.open()

	rows, err := p.reader.ReadData(p.ctx, p.tableHandle, p.predicate)
	if err != nil {
		panic(err)
	}
	colIdxs := p.colIdxs
	fields := p.rb.Fields()
	for _, row := range rows {
		for idx, field := range fields {
			switch col := field.(type) {
			case *array.StringBuilder:
				col.Append(row[colIdxs[idx]].String())
			case *array.Float64Builder:
				col.Append(row[colIdxs[idx]].Float())
			case *array.Int64Builder:
				col.Append(row[colIdxs[idx]].Int())
			case *array.TimestampBuilder:
				col.Append(arrow.Timestamp(row[colIdxs[idx]].Int()))
			case *array.DurationBuilder:
				col.Append(arrow.Duration(row[colIdxs[idx]].Duration().Milliseconds()))
			}
		}
	}

	// send result set
	output <- p.rb.NewRecord()
}
