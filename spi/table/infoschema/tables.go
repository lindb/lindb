package infoschema

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

func init() {
	spi.RegisterGetTableSchemaFn(spi.InfoSchema, func(db, ns, table string) (*types.TableSchema, error) {
		schema := GetTableSchema(table)
		if schema == nil {
			return nil, fmt.Errorf("information table schema not found: %s", table)
		}
		return schema, nil
	})
}

func InitInfoSchema(metadataMgr meta.MetadataManager) {
	spi.RegisterSplitSourceProvider(&TableHandle{}, NewSplitSourceProvider(metadataMgr))
	spi.RegisterPageSourceProvider(&TableHandle{}, NewPageSourceProvider(NewReader(metadataMgr)))
}

var (
	masterSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "host_ip", DataType: types.DTString},
			{Name: "host_name", DataType: types.DTString},
			{Name: "version", DataType: types.DTString},
			{Name: "online_time", DataType: types.DTInt},
			{Name: "elect_time", DataType: types.DTInt},
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
			{Name: "namespace_name", DataType: types.DTString},
		},
	}
	tablesSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "table_name", DataType: types.DTString},
		},
	}
	columnsSchema = &types.TableSchema{
		Columns: []types.ColumnMetadata{
			{Name: "column_name", DataType: types.DTString},
			{Name: "data_type", DataType: types.DTString},
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
		constants.TableMaster:   masterSchema,
		constants.TableSchemata: schemtatSchema,
		constants.TableMetrics:  metricsSchema,
		"tables":                tablesSchema,
		"namespaces":            namespacesSchema,
		"columns":               columnsSchema,
	}
)

func GetTableSchema(name string) *types.TableSchema {
	return tables[name]
}
