package spi

import "github.com/lindb/lindb/models"

type TableMetadata struct {
	Schema     *TableSchema
	Partitions map[models.InternalNode][]int
}

type TableSchema struct {
	Columns []ColumnMetadata `json:"columns,omitempty"`
}

func NewTableSchema() *TableSchema {
	return &TableSchema{}
}

func (s *TableSchema) AddColumn(column ColumnMetadata) {
	s.Columns = append(s.Columns, column)
}

func (s *TableSchema) AddColumns(columns []ColumnMetadata) {
	s.Columns = append(s.Columns, columns...)
}
