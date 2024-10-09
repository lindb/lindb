package utils

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

func SingleValueQuery(columnName, value string) *tree.Query {
	page := types.NewPage()
	column := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{Name: columnName, DataType: types.DTString}, column)
	column.AppendString(value)
	return &tree.Query{
		BaseNode: tree.BaseNode{
			ID: 1,
		},
		QueryBody: &tree.QuerySpecification{
			BaseNode: tree.BaseNode{
				ID: 2,
			},
			Select: &tree.Select{
				SelectItems: []tree.SelectItem{&tree.AllColumns{}},
			},
			From: &tree.Values{
				BaseNode: tree.BaseNode{
					ID: 3,
				},
				Rows: page,
			},
		},
	}
}
