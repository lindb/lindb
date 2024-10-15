package infoschema

import (
	"fmt"
	"strings"

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
		Columns: []types.ColumnMetadata{},
	}
	memoryDatabaseSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{},
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
	tablesSchema = &types.TableSchema{
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
			{Name: "metrics_name", DataType: types.DTString},
			{Name: "instance", DataType: types.DTString},
			{Name: "value", DataType: types.DTFloat},
		},
	}

	// tables represents the schema of tables.
	tables = map[string]*types.TableSchema{
		constants.TableMaster:         masterSchema,
		constants.TableBroker:         brokerSchema,
		constants.TableStorage:        storageSchema,
		constants.TableReplication:    replicationSchema,
		constants.TableMemoryDatabase: memoryDatabaseSchema,
		constants.TableEngines:        enginesSchema,
		constants.TableSchemata:       schemtatSchema,
		constants.TableMetrics:        metricsSchema,
		constants.TableNamespaces:     namespacesSchema,
		constants.TableTableNames:     tablesSchema,
		constants.TableColumns:        columnsSchema,
	}
)
