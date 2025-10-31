package expression

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
)

type mapValuesFunc struct {
	baseFunc
}

func newMapValuesFunc(args []Expression) Func {
	return &mapValuesFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

func (n *mapValuesFunc) EvalMap(ctx EvalContext, row types.Row) (val map[string]string, isNull bool, err error) {
	// TODO: check error/now func
	tsStr, _, _ := n.args[0].EvalMap(ctx, row)
	// if err != nil {
	// 	return "kk", true, err
	// }
	duration, _, _ := n.args[1].EvalString(ctx, row)
	return lo.PickByKeys(tsStr, []string{duration}), true, nil
}
