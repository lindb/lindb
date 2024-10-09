package expression

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

func Eval(expression tree.Expression) (any, error) {
	fmt.Printf("eval expr=%T\n", expression)
	return nil, nil
}
