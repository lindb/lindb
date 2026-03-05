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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
)

func init() {
	spi.RegisterGetTableSchemaFn(spi.InfoSchema, func(db, ns, table string) (*arrow.Schema, error) {
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

func GetTableSchema(name string) (schema *arrow.Schema, ok bool) {
	schema, ok = tables[strings.ToLower(name)]
	return
}

func GetShowSelectColumns(name string, start int) (columns []string) {
	schema, ok := GetTableSchema(name)
	if !ok {
		return
	}
	fields := schema.Fields()
	return lo.Map(fields[start:], func(item arrow.Field, index int) string {
		return item.Name
	})
}

var (
	envSchema = arrow.NewSchema([]arrow.Field{
		{Name: "instance", Type: arrow.BinaryTypes.String},
		{Name: "key", Type: arrow.BinaryTypes.String},
		{Name: "value", Type: arrow.BinaryTypes.String},
		{Name: "default", Type: arrow.BinaryTypes.String},
	}, nil)
	masterSchema = arrow.NewSchema([]arrow.Field{
		{Name: "host_ip", Type: arrow.BinaryTypes.String},
		{Name: "host_name", Type: arrow.BinaryTypes.String},
		{Name: "http", Type: arrow.PrimitiveTypes.Int32},
		{Name: "version", Type: arrow.BinaryTypes.String},
		{Name: "online_time", Type: arrow.FixedWidthTypes.Timestamp_ns},
		{Name: "elect_time", Type: arrow.FixedWidthTypes.Timestamp_ns},
	}, nil)
	brokerSchema = arrow.NewSchema([]arrow.Field{
		{Name: "host_ip", Type: arrow.BinaryTypes.String},
		{Name: "host_name", Type: arrow.BinaryTypes.String},
		{Name: "version", Type: arrow.BinaryTypes.String},
		{Name: "online_time", Type: arrow.FixedWidthTypes.Timestamp_ns},
		{Name: "uptime", Type: arrow.FixedWidthTypes.Duration_ns},
		{Name: "grpc", Type: arrow.PrimitiveTypes.Int32},
		{Name: "http", Type: arrow.PrimitiveTypes.Int32},
	}, nil)
	storageSchema = arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int32},
		{Name: "host_ip", Type: arrow.BinaryTypes.String},
		{Name: "host_name", Type: arrow.BinaryTypes.String},
		{Name: "version", Type: arrow.BinaryTypes.String},
		{Name: "online_time", Type: arrow.FixedWidthTypes.Timestamp_ns},
		{Name: "uptime", Type: arrow.FixedWidthTypes.Duration_ns},
		{Name: "grpc", Type: arrow.PrimitiveTypes.Int32},
		{Name: "http", Type: arrow.PrimitiveTypes.Int32},
	}, nil)
	replicationSchema = arrow.NewSchema([]arrow.Field{
		{Name: "table_schema", Type: arrow.BinaryTypes.String},
		{Name: "node", Type: arrow.BinaryTypes.String},
		{Name: "shard", Type: arrow.PrimitiveTypes.Int32},
		{Name: "family", Type: arrow.BinaryTypes.String},
		{Name: "leader", Type: arrow.PrimitiveTypes.Int32},
		{Name: "replicator", Type: arrow.BinaryTypes.String},
		{Name: "type", Type: arrow.BinaryTypes.String},
		{Name: "append", Type: arrow.PrimitiveTypes.Int32},
		{Name: "consume", Type: arrow.PrimitiveTypes.Int32},
		{Name: "ack", Type: arrow.PrimitiveTypes.Int32},
		{Name: "pending", Type: arrow.PrimitiveTypes.Int32},
		{Name: "state", Type: arrow.BinaryTypes.String},
		{Name: "error", Type: arrow.BinaryTypes.String},
	}, nil)
	memoryDatabaseSchema = arrow.NewSchema([]arrow.Field{
		{Name: "table_schema", Type: arrow.BinaryTypes.String},
		{Name: "node", Type: arrow.BinaryTypes.String},
		{Name: "shard", Type: arrow.PrimitiveTypes.Int32},
		{Name: "family", Type: arrow.BinaryTypes.String},
		{Name: "state", Type: arrow.BinaryTypes.String},
		{Name: "uptime", Type: arrow.FixedWidthTypes.Duration_ns},
		{Name: "mem_size", Type: arrow.PrimitiveTypes.Int32},
		{Name: "num_of_series", Type: arrow.PrimitiveTypes.Int32},
	}, nil)
	enginesSchema = arrow.NewSchema([]arrow.Field{
		{Name: "engine", Type: arrow.BinaryTypes.String},  // metric/log/trace
		{Name: "support", Type: arrow.BinaryTypes.String}, // default/yes/no/disabled
	}, nil)
	schemataSchema = arrow.NewSchema([]arrow.Field{
		{Name: "schema_name", Type: arrow.BinaryTypes.String},
		{Name: "engine", Type: arrow.BinaryTypes.String},
		{Name: "statement", Type: arrow.BinaryTypes.String},
	}, nil)
	namespacesSchema = arrow.NewSchema([]arrow.Field{
		{Name: "table_schema", Type: arrow.BinaryTypes.String},
		{Name: "namespace", Type: arrow.BinaryTypes.String},
	}, nil)
	tableNamesSchema = arrow.NewSchema([]arrow.Field{
		{Name: "table_schema", Type: arrow.BinaryTypes.String},
		{Name: "namespace", Type: arrow.BinaryTypes.String},
		{Name: "table_name", Type: arrow.BinaryTypes.String},
	}, nil)
	columnsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "table_schema", Type: arrow.BinaryTypes.String},
		{Name: "namespace", Type: arrow.BinaryTypes.String},
		{Name: "table_name", Type: arrow.BinaryTypes.String},
		{Name: "column_name", Type: arrow.BinaryTypes.String},
		{Name: "data_type", Type: arrow.BinaryTypes.String},
	}, nil)
	metricsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "role", Type: arrow.BinaryTypes.String},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "tags", Type: arrow.BinaryTypes.String},
		{Name: "field_name", Type: arrow.BinaryTypes.String},
		{Name: "field_type", Type: arrow.BinaryTypes.String},
		{Name: "field_value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)
	metadataTypesSchema = arrow.NewSchema([]arrow.Field{
		{Name: "role", Type: arrow.BinaryTypes.String},
		{Name: "type", Type: arrow.BinaryTypes.String},
		{Name: "comment", Type: arrow.BinaryTypes.String},
	}, nil)
	metadatasSchema = arrow.NewSchema([]arrow.Field{
		{Name: "role", Type: arrow.BinaryTypes.String},
		{Name: "type", Type: arrow.BinaryTypes.String},
		{Name: "source", Type: arrow.BinaryTypes.String},
		{Name: "data", Type: arrow.BinaryTypes.String},
	}, nil)
	functionsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "template", Type: arrow.BinaryTypes.String},
	}, nil)
	snippetsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "template", Type: arrow.BinaryTypes.String},
	}, nil)
	streamingsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "statement", Type: arrow.BinaryTypes.String},
	}, nil)
	streamingJobsSchema = arrow.NewSchema([]arrow.Field{
		{Name: "streaming", Type: arrow.BinaryTypes.String},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "statement", Type: arrow.BinaryTypes.String},
	}, nil)

	// tables represents the schema of tables.
	tables = map[string]*arrow.Schema{
		constants.TableEnv:             envSchema,
		constants.TableMaster:          masterSchema,
		constants.TableBrokers:         brokerSchema,
		constants.TableStorages:        storageSchema,
		constants.TableReplications:    replicationSchema,
		constants.TableMemoryDatabases: memoryDatabaseSchema,
		constants.TableEngines:         enginesSchema,
		constants.TableSchemata:        schemataSchema,
		constants.TableMetrics:         metricsSchema,
		constants.TableNamespaces:      namespacesSchema,
		constants.TableTableNames:      tableNamesSchema,
		constants.TableColumns:         columnsSchema,
		constants.TableMetadataTypes:   metadataTypesSchema,
		constants.TableMetadatas:       metadatasSchema,
		constants.TableFunctions:       functionsSchema,
		constants.TableSnippets:        snippetsSchema,
		constants.TableStreamings:      streamingsSchema,
		constants.TableStreamingJobs:   streamingJobsSchema,
	}
)
