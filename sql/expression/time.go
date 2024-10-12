package expression

import (
	"fmt"
	"time"

	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/spi/types"
)

type addSubDateFuncFactory struct{}

func (fct *addSubDateFuncFactory) NewFunc(args []Expression) Func {
	return &addSubDateFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

type addSubDateFunc struct {
	baseFunc
}

func (n *addSubDateFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(ctx, row)
	format := timeutil.DataTimeFormat2
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	duration, _, _ := n.args[1].EvalDuration(ctx, row)
	return time.UnixMilli(timestamp).Add(duration), false, nil
}

type nowFuncFactory struct{}

func (fct *nowFuncFactory) NewFunc(args []Expression) Func {
	return &nowFunc{}
}

type nowFunc struct {
	baseFunc
}

func (n *nowFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	fmt.Println(ctx)
	return ctx.CurrentTime(), false, nil
}

type strToDateFuncFactory struct{}

func (fct *strToDateFuncFactory) NewFunc(args []Expression) Func {
	return &strToDateFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

type strToDateFunc struct {
	baseFunc
}

func (n *strToDateFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(ctx, row)
	format, _, _ := n.args[1].EvalString(ctx, row)
	switch format {
	case "YYYYMMDD HH:mm:ss":
		format = timeutil.DataTimeFormat1
	case "YYYY-MM-DD HH:mm:ss":
		format = timeutil.DataTimeFormat2
	case "YYYY/MM/DD HH:mm:ss":
		format = timeutil.DataTimeFormat3
	case "YYYYMMDDHHmmss":
		format = timeutil.DataTimeFormat4
	}
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	return time.UnixMilli(timestamp), false, nil
}
