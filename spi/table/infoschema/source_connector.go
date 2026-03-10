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

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

// sourceConnectorProvider creates sourceConnector instances for information-schema tables.
type sourceConnectorProvider struct {
	metadataMgr meta.MetadataManager
}

// NewSourceConnectorProvider returns a SourceConnectorProvider backed by the given MetadataManager.
func NewSourceConnectorProvider(
	metadataMgr meta.MetadataManager,
) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		metadataMgr: metadataMgr,
	}
}

// CreateSourceConnector builds a sourceConnector for the given table handle and query parameters.
// outputColumns defines the subset of columns the query projects; the connector will return only those fields.
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

// sourceConnector reads data from an information-schema table and projects it to the requested output columns.
type sourceConnector struct {
	ctx    context.Context
	reader Reader

	table         spi.TableHandle
	tableHandle   *TableHandle
	predicate     tree.Expression
	outputColumns []arrow.Field // columns projected by the query (subset of the full table schema)
}

// Run reads the full information-schema record, projects it to outputColumns, and sends the result to output.
func (sc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	if len(sc.outputColumns) == 0 {
		return
	}

	infoTable, ok := sc.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("information schema provider not support table handle<%T>", sc.table))
	}
	sc.tableHandle = infoTable

	// Read the full record for this information-schema table.
	record, err := sc.reader.ReadData(sc.ctx, sc.tableHandle, sc.predicate)
	if err != nil {
		panic(err)
	}
	if record == nil || record.NumRows() == 0 {
		return
	}
	defer record.Release()

	// Project the full record down to only the columns requested by the query.
	projected := sc.projectRecord(record)
	if projected != nil && projected.NumRows() > 0 {
		output <- projected
	}
}

// projectRecord returns a new RecordBatch that contains only the fields listed
// in outputColumns, in the order they appear in outputColumns.
// Columns are matched by field name. Missing columns are filled with nulls.
func (sc *sourceConnector) projectRecord(record arrow.RecordBatch) arrow.RecordBatch {
	// Build an index from field name → column index in the source record.
	srcSchema := record.Schema()
	nameToIdx := make(map[string]int, srcSchema.NumFields())
	for i, f := range srcSchema.Fields() {
		nameToIdx[f.Name] = i
	}

	// Collect the projected columns, retaining each source column we reference.
	cols := make([]arrow.Array, len(sc.outputColumns))
	for i, f := range sc.outputColumns {
		if srcIdx, ok := nameToIdx[f.Name]; ok {
			col := record.Column(srcIdx)
			cols[i] = col
		} else {
			// Output column not present in source: fill with nulls.
			builder := array.NewBuilder(memory.DefaultAllocator, f.Type)
			defer builder.Release()
			for range record.NumRows() {
				builder.AppendNull()
			}
			cols[i] = builder.NewArray()
		}
	}

	outSchema := arrow.NewSchema(sc.outputColumns, nil)
	return array.NewRecordBatch(outSchema, cols, record.NumRows())
}
