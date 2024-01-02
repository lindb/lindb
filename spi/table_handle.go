package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("spi.TableHandle", &encoding.JSONEncoder[TableHandle]{})
	jsoniter.RegisterTypeDecoder("spi.TableHandle", &encoding.JSONDecoder[TableHandle]{})
}

type TableHandle interface {
}
