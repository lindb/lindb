package tree

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("tree.Expression", &encoding.JSONEncoder[Expression]{})
	jsoniter.RegisterTypeDecoder("tree.Expression", &encoding.JSONDecoder[Expression]{})

	encoding.RegisterNodeType(ComparisonExpression{})
	encoding.RegisterNodeType(StringLiteral{})
	encoding.RegisterNodeType(Identifier{})
	encoding.RegisterNodeType(Cast{})
}
