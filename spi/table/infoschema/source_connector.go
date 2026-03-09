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
		ctx:       ctx,
		table:     table,
		predicate: predicate,
		reader:    NewReader(p.metadataMgr),
	}
}

type sourceConnector struct {
	ctx    context.Context
	reader Reader

	table       spi.TableHandle
	tableHandle *TableHandle
	predicate   tree.Expression
}

func (p *sourceConnector) open() {
	infoTable, ok := p.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("information schema provider not support table handle<%T>", p.table))
	}
	p.tableHandle = infoTable
}

func (p *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	p.open()

	record, err := p.reader.ReadData(p.ctx, p.tableHandle, p.predicate)
	if err != nil {
		panic(err)
	}
	if record != nil && record.NumRows() > 0 {
		output <- record
	}
}
