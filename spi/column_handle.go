package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	// register json encoder/decoder for column handle
	jsoniter.RegisterTypeEncoder("spi.ColumnHandle", &encoding.JSONEncoder[ColumnHandle]{})
	jsoniter.RegisterTypeDecoder("spi.ColumnHandle", &encoding.JSONDecoder[ColumnHandle]{})
}

type ColumnHandle interface{}
