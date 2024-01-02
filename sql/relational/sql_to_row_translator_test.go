package relational

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/sql/tree"
)

func TestSQLToExpression_Translate(t *testing.T) {
	translator := &SQLToRowExpressionTranslator{}
	result := translator.Translate(&tree.StringLiteral{Value: "test"})
	fmt.Println(result)
}
