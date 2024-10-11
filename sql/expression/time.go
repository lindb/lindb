package expression

import (
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

func (n *addSubDateFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(row)
	format := timeutil.DataTimeFormat2
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	duration, _, _ := n.args[1].EvalDuration(row)
	return time.UnixMilli(timestamp).Add(duration), false, nil
}

type nowFuncFactory struct{}

func (fct *nowFuncFactory) NewFunc(args []Expression) Func {
	return &nowFunc{}
}

type nowFunc struct {
	baseFunc
}

func (n *nowFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	return time.Now(), false, nil
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

func (n *strToDateFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(row)
	format, _, _ := n.args[1].EvalString(row)
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
