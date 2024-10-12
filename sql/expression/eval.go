package expression

import (
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

func EvalTime(ctx EvalContext, expression tree.Expression) (time.Time, error) {
	expr := Rewrite(&RewriteContext{}, expression)
	val, _, err := expr.EvalTime(ctx, types.EmptyRow)
	return val, err
}
