package infoschema

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

func init() {
	spi.RegisterGetTableSchemaFn(spi.InfoSchema, func(db, ns, table string) (*types.TableSchema, error) {
		return GetTableSchema(table), nil
	})
}

func InitInfoSchema(metadataMgr meta.MetadataManager) {
	spi.RegisterSplitSourceProvider(&TableHandle{}, NewSplitSourceProvider(metadataMgr))
	spi.RegisterPageSourceProvider(&TableHandle{}, NewPageSourceProvider(NewReader(metadataMgr)))
}

var (
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

	// tables represents the schema of tables.
	tables = map[string]*types.TableSchema{
		constants.TableSchemata: schemtatSchema,
		"tables":                tablesSchema,
		"namespaces":            namespacesSchema,
		"columns":               columnsSchema,
	}
)

func GetTableSchema(name string) *types.TableSchema {
	return tables[name]
}
