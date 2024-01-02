package value

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("value.Block", &encoding.JSONEncoder[Block]{})
	jsoniter.RegisterTypeDecoder("value.Block", &encoding.JSONDecoder[Block]{})

	encoding.RegisterNodeType(TimeSeries{})
}
