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

package types

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/models"
)

type TableMetadata struct {
	Schema     *arrow.Schema
	Partitions map[models.InternalNode][]int

	SupportDynamicField bool
}

// tableArrowFields stores the pre-registered Arrow field list for each named table type.
// Populated via RegisterTableArrowFields during package init by each datasource package.
var tableArrowFields = map[string][]arrow.Field{}

// RegisterTableArrowFields registers the canonical Arrow fields for a named table type.
// Called from datasource packages (e.g. spi/table/log) during init().
func RegisterTableArrowFields(tableName string, fields []arrow.Field) {
	tableArrowFields[tableName] = fields
}

// GetTableArrowFields returns the registered Arrow fields for the named table type.
// Returns nil when no fields have been registered for that name.
func GetTableArrowFields(tableName string) []arrow.Field {
	return tableArrowFields[tableName]
}
