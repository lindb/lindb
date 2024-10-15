package spi

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi/types"
)

type ApplyAggregationResult struct {
	ColumnAssignments []*ColumnAssignment
}

type CreateTable func(db, ns, table string) TableHandle

type applyAggregation func(table TableHandle, tableMeta *types.TableMetadata,
	aggregations []ColumnAggregation) *ApplyAggregationResult

var (
	createTableFn      = make(map[DatasourceKind]CreateTable)
	getTableSchemaFn   = make(map[DatasourceKind]func(db, ns, table string) (*types.TableSchema, error))
	applyAggregationFn = make(map[DatasourceKind]applyAggregation)
)

func GetTableSchema(kind DatasourceKind, database, ns, table string) (*types.TableSchema, error) {
	return getTableSchemaFn[kind](database, ns, table)
}

func RegisterCreateTableFn(kind DatasourceKind, fn CreateTable) {
	createTableFn[kind] = fn
}

func RegisterGetTableSchemaFn(kind DatasourceKind, fn func(database, ns, table string) (*types.TableSchema, error)) {
	getTableSchemaFn[kind] = fn
}

func RegisterApplyAggregationFn(kind DatasourceKind, fn applyAggregation) {
	applyAggregationFn[kind] = fn
}

func ApplyAggregation(table TableHandle, tableMeta *types.TableMetadata, aggregations []ColumnAggregation) *ApplyAggregationResult {
	applyFn, ok := applyAggregationFn[table.Kind()]
	if !ok {
		panic(fmt.Sprintf("apply aggregation func not exist, kind: %v", table.Kind()))
	}
	return applyFn(table, tableMeta, aggregations)
}

type MetadataManager interface {
	GetTableMetadata(db, ns, table string) (*types.TableMetadata, error)
	GetTableHandle(db, ns, table string) TableHandle
}

type metadataManager struct {
	metadataMgr meta.MetadataManager
}

func NewMetadataManager(metadataMgr meta.MetadataManager) MetadataManager {
	return &metadataManager{
		metadataMgr: metadataMgr,
	}
}

func (mgr *metadataManager) GetTableHandle(db, ns, table string) TableHandle {
	fmt.Printf("get table handle %v\n", db)
	var kind DatasourceKind
	if db == constants.InformationSchema {
		kind = InfoSchema
	} else {
		_, ok := mgr.metadataMgr.GetDatabase(db)
		if !ok {
			panic(constants.ErrDatabaseNotExist)
		}
		// FIXME: get table kind by database
		kind = Metric
	}
	fn, ok := createTableFn[kind]
	if !ok {
		panic(fmt.Sprintf("create table handle func not exist, kind: %v", Metric))
	}
	return fn(db, ns, table)
}

func (mgr *metadataManager) GetTableMetadata(database, ns, table string) (*types.TableMetadata, error) {
	if database == constants.InformationSchema {
		schema, err := GetTableSchema(InfoSchema, database, ns, table)
		if err != nil {
			return nil, err
		}
		partitions, err := mgr.metadataMgr.GetPartitions(database, ns, table)
		if err != nil {
			return nil, err
		}
		return &types.TableMetadata{
			Schema:     schema,
			Partitions: partitions,
		}, nil
	}
	return mgr.metadataMgr.GetTableMetadata(database, ns, table)
}
