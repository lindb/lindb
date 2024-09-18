package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	// register json encoder/decoder for table handle
	jsoniter.RegisterTypeEncoder("spi.TableHandle", &encoding.JSONEncoder[TableHandle]{})
	jsoniter.RegisterTypeDecoder("spi.TableHandle", &encoding.JSONDecoder[TableHandle]{})
}

type TableKind int

const (
	MetricTable TableKind = iota + 1
)

// TableHandle represents a table handle that connect the storage engine.
type TableHandle interface {
	// String returns table info, format: ${database}:${namespace}:${tableName}
	String() string
}
