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
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

func init() {
	spi.RegisterGetTableSchemaFn(spi.InfoSchema, func(db, ns, table string) (*types.TableSchema, error) {
		schema, ok := GetTableSchema(table)
		if !ok {
			return nil, fmt.Errorf("information table schema not found: %s", table)
		}
		return schema, nil
	})
}

func InitInfoSchema(metadataMgr meta.MetadataManager) {
	spi.RegisterSourceConnectorProvider(&TableHandle{}, NewSourceConnectorProvider(metadataMgr))
}

func GetTableSchema(name string) (schema *types.TableSchema, ok bool) {
	schema, ok = tables[strings.ToLower(name)]
	return
}

func GetShowSelectColumns(name string, start int) (columns []string) {
	schema, ok := GetTableSchema(name)
	if !ok {
		return
	}
	return lo.Map(schema.Columns[start:], func(item types.ColumnMetadata, index int) string {
		return item.Name
	})
}

var (
	envSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "instance", DataType: types.DTString, Hidden: true},
			{Name: "key", DataType: types.DTString},
			{Name: "value", DataType: types.DTString},
			{Name: "default", DataType: types.DTString},
		},
	}
	masterSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "http", DataType: types.DTInt},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTTimestamp},
			{Name: "elect_time", DataType: types.DTTimestamp},
		},
	}
	brokerSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTTimestamp},
			{Name: "uptime", DataType: types.DTDuration},
			{Name: "grpc", DataType: types.DTInt},
			{Name: "http", DataType: types.DTInt},
		},
	}
	storageSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "id", DataType: types.DTInt},
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTTimestamp},
			{Name: "uptime", DataType: types.DTDuration},
			{Name: "grpc", DataType: types.DTInt},
			{Name: "http", DataType: types.DTInt},
		},
	}
	replicationSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_schema", DataType: types.DTString},
			{Name: "node", DataType: types.DTString},
			{Name: "shard", DataType: types.DTInt},
			{Name: "family", DataType: types.DTString},
			{Name: "leader", DataType: types.DTInt},
			{Name: "replicator", DataType: types.DTString},
			{Name: "type", DataType: types.DTString},
			{Name: "append", DataType: types.DTInt},
			{Name: "consume", DataType: types.DTInt},
			{Name: "ack", DataType: types.DTInt},
			{Name: "pending", DataType: types.DTInt},
			{Name: "state", DataType: types.DTString},
			{Name: "error", DataType: types.DTString},
		},
	}
	memoryDatabaseSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_schema", DataType: types.DTString},
			{Name: "node", DataType: types.DTString},
			{Name: "shard", DataType: types.DTInt},
			{Name: "family", DataType: types.DTString},
			{Name: "state", DataType: types.DTString},
			{Name: "uptime", DataType: types.DTDuration},
			{Name: "mem_size", DataType: types.DTInt},
			{Name: "num_of_series", DataType: types.DTInt},
		},
	}
	enginesSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "engine", DataType: types.DTString},  // metric/log/trace
			{Name: "support", DataType: types.DTString}, // default/yes/no/disabled
		},
	}
	schemtatSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "schema_name", DataType: types.DTString},
			{Name: "engine", DataType: types.DTString},
		},
	}
	namespacesSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_schema", DataType: types.DTString},
			{Name: "namespace", DataType: types.DTString},
		},
	}
	tableNamesSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_schema", DataType: types.DTString},
			{Name: "namespace", DataType: types.DTString},
			{Name: "table_name", DataType: types.DTString},
		},
	}
	columnsSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_schema", DataType: types.DTString},
			{Name: "namespace", DataType: types.DTString},
			{Name: "table_name", DataType: types.DTString},
			{Name: "column_name", DataType: types.DTString},
			{Name: "data_type", DataType: types.DTString},
			{Name: "agg_type", DataType: types.DTString},
		},
	}
	metricsSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "role", DataType: types.DTString},
			{Name: "name", DataType: types.DTString},
			{Name: "tags", DataType: types.DTString},
			{Name: "field_name", DataType: types.DTString},
			{Name: "field_type", DataType: types.DTString},
			{Name: "field_value", DataType: types.DTFloat},
		},
	}
	metadataTypesSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "role", DataType: types.DTString},
			{Name: "type", DataType: types.DTString},
			{Name: "comment", DataType: types.DTString},
		},
	}
	metadatasSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "role", DataType: types.DTString},
			{Name: "type", DataType: types.DTString},
			{Name: "source", DataType: types.DTString},
			{Name: "data", DataType: types.DTString},
		},
	}
	functionsSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "name", DataType: types.DTString},
			{Name: "template", DataType: types.DTString},
		},
	}
	snippetsSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "name", DataType: types.DTString},
			{Name: "template", DataType: types.DTString},
		},
	}

	// tables represents the schema of tables.
	tables = map[string]*types.TableSchema{
		constants.TableEnv:             envSchema,
		constants.TableMaster:          masterSchema,
		constants.TableBrokers:         brokerSchema,
		constants.TableStorages:        storageSchema,
		constants.TableReplications:    replicationSchema,
		constants.TableMemoryDatabases: memoryDatabaseSchema,
		constants.TableEngines:         enginesSchema,
		constants.TableSchemata:        schemtatSchema,
		constants.TableMetrics:         metricsSchema,
		constants.TableNamespaces:      namespacesSchema,
		constants.TableTableNames:      tableNamesSchema,
		constants.TableColumns:         columnsSchema,
		constants.TableMetadataTypes:   metadataTypesSchema,
		constants.TableMetadatas:       metadatasSchema,
		constants.TableFunctions:       functionsSchema,
		constants.TableSnippets:        snippetsSchema,
	}
)
