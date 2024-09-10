package model

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
)

type Schema struct {
	Columns   []spi.ColumnMetadata `json:"columns,omitempty"`
	Partition []models.Partition   `json:"partitions,omitempty"`
}

type ResultSet struct {
	Schema *Schema `json:"schema,omitempty"`
	Rows   [][]any `json:"rows,omitempty"`
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		Schema: &Schema{},
	}
}
