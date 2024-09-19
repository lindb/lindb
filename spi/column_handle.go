package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/sql/tree"
)

func init() {
	// register json encoder/decoder for column handle
	jsoniter.RegisterTypeEncoder("spi.ColumnHandle", &encoding.JSONEncoder[ColumnHandle]{})
	jsoniter.RegisterTypeDecoder("spi.ColumnHandle", &encoding.JSONDecoder[ColumnHandle]{})
}

type ColumnHandle interface{}

type ColumnAssignment struct {
	Handler ColumnHandle `json:"handler"`
	Column  string       `json:"column"`
}

type ColumnAggregation struct {
	Column      string
	AggFuncName tree.FuncName
}
