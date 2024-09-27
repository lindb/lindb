package types

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("types.Block", &encoding.JSONEncoder[Block]{})
	jsoniter.RegisterTypeDecoder("types.Block", &encoding.JSONDecoder[Block]{})

	encoding.RegisterNodeType(TimeSeries{})
	encoding.RegisterNodeType(String(""))
	encoding.RegisterNodeType(Float(0))
	encoding.RegisterNodeType(Int(0))
}
