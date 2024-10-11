package expression

import (
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

func EvalTime(expression tree.Expression) (time.Time, error) {
	expr := Rewrite(&RewriteContext{}, expression)
	val, _, err := expr.EvalTime(types.EmptyRow)
	return val, err
}
