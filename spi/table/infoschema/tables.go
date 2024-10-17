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
	spi.RegisterSplitSourceProvider(&TableHandle{}, NewSplitSourceProvider(metadataMgr))
	spi.RegisterPageSourceProvider(&TableHandle{}, NewPageSourceProvider(NewReader(metadataMgr)))
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
	masterSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTString},
			{Name: "elect_time", DataType: types.DTString},
		},
	}
	brokerSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTString},
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
			{Name: "online_time", DataType: types.DTString},
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
			{Name: "type", DataType: types.DTString},
			{Name: "value", DataType: types.DTFloat},
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
	// tables represents the schema of tables.
	tables = map[string]*types.TableSchema{
		constants.TableMaster:          masterSchema,
		constants.TableBroker:          brokerSchema,
		constants.TableStorage:         storageSchema,
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
	}
)
