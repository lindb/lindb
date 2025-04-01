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
	"time"

	"github.com/samber/lo"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
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
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
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
	tableName     string
	predicate     tree.Expression
	outputColumns []types.ColumnMetadata
	colIdxs       []int
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
	for i, col := range p.outputColumns {
		if _, idx, exist := lo.FindIndexOf(schema.Columns, func(item types.ColumnMetadata) bool {
			return item.Name == col.Name
		}); exist {
			p.colIdxs[i] = idx
		}
	}
	if len(p.colIdxs) != len(p.outputColumns) {
		// FIXME: add panic?
		return
	}

	p.tableName = infoTable.Table
}

func (p *sourceConnector) Run(output chan<- *types.Page) {
	p.open()
	rows, err := p.reader.ReadData(p.ctx, p.tableName, p.predicate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("info schema: rows=%v\n", rows)
	page := types.NewPage()
	var columns []*types.Column
	outputs := make(map[string]int)
	for idx, output := range p.outputColumns {
		column := types.NewColumn()
		page.AppendColumn(output, column)
		columns = append(columns, column)
		outputs[output.Name] = idx
	}
	colIdxs := p.colIdxs
	for _, row := range rows {
		for idx, col := range columns {
			switch p.outputColumns[idx].DataType {
			case types.DTString:
				col.AppendString(row[colIdxs[idx]].String())
			case types.DTFloat:
				col.AppendFloat(row[colIdxs[idx]].Float())
			case types.DTInt:
				col.AppendInt(row[colIdxs[idx]].Int())
			case types.DTTimestamp:
				col.AppendTimestamp(time.UnixMilli(row[colIdxs[idx]].Int()))
			case types.DTDuration:
				col.AppendDuration(row[colIdxs[idx]].Duration())
			}
		}
	}

	// send result set
	output <- page

	// close(output)
}
