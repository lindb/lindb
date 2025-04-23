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

package spi

import (
	"context"
	"fmt"
	"reflect"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

var sourceConnectorProviders = make(map[reflect.Type]SourceConnectorProvider)

func RegisterSourceConnectorProvider(table TableHandle, provider SourceConnectorProvider) {
	sourceConnectorProviders[reflect.TypeOf(table)] = provider
}

func GetSourceConnectorProvider(table TableHandle) SourceConnectorProvider {
	if prodiver, ok := sourceConnectorProviders[reflect.TypeOf(table)]; ok {
		return prodiver
	}
	panic(fmt.Sprintf("source connector provider not found by table handle type for '%s'", reflect.TypeOf(table)))
}

type SourceConnector interface {
	Run(output chan<- *types.Page)
}

type SourceConnectorProvider interface {
	CreateSourceConnector(ctx context.Context,
		table TableHandle, partitions []int, // table info
		columnMapping map[string]string,
		predicate tree.Expression, // predicate
		outputColumns []types.ColumnMetadata, assignments []*ColumnAssignment, // output
	) SourceConnector
}
