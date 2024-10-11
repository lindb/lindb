package expression

import (
	"errors"
	"fmt"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
)

type arithmeticPlusFuncFactory struct{}

func (fct *arithmeticPlusFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticPlusFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticPlusFunc struct {
	baseFunc
}

func (f *arithmeticPlusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("plus int.....")
	return lv + rv, false, nil
}

func (f *arithmeticPlusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("plus float.....")
	return
}

func (f *arithmeticPlusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("plus time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv + rv
	})
}

func (f *arithmeticPlusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalTime(row)
	rv, _, _ := f.args[1].EvalDuration(row)
	val = lv.Add(rv)
	return
}

func evalTimeSeries(row types.Row, args []Expression, math func(lv, rv float64) float64) (val *types.TimeSeries, isNull bool, err error) {
	l, lIsNull, err := args[0].EvalTimeSeries(row)
	if err != nil {
		return nil, false, err
	}
	r, rIsNull, err := args[1].EvalTimeSeries(row)
	if err != nil {
		return nil, false, err
	}
	if lIsNull {
		return r, rIsNull, nil
	}
	if rIsNull {
		return l, lIsNull, nil
	}
	fmt.Println(l)
	fmt.Println(r)
	// check num. of points whether match
	if !l.IsSingleValue() && !r.IsSingleValue() && l.Size() != r.Size() {
		return nil, true, errors.New("num. of points not match")
	}
	var result *types.TimeSeries
	if !l.IsSingleValue() {
		result = types.NewTimeSeries(l.TimeRange, timeutil.Interval(l.Interval))
	} else if !r.IsSingleValue() {
		result = types.NewTimeSeries(r.TimeRange, timeutil.Interval(r.Interval))
	} else {
		result = types.NewTimeSeriesWithSingleValue(0)
	}
	for i := 0; i < result.Size(); i++ {
		result.Put(i, math(l.Get(i), r.Get(i)))
	}
	return result, false, nil
}

type arithmeticMinusFuncFactory struct{}

func (fct *arithmeticMinusFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticMinusFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticMinusFunc struct {
	baseFunc
}

func (f *arithmeticMinusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("minus int.....")
	return lv - rv, false, nil
}

func (f *arithmeticMinusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("minus float.....")
	return
}

func (f *arithmeticMinusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("minus time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv - rv
	})
}

func (f *arithmeticMinusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalTime(row)
	rv, _, _ := f.args[1].EvalDuration(row)
	val = lv.Add(-rv)
	return
}

type arithmeticMulFuncFactory struct{}

func (fct *arithmeticMulFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticMulFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticMulFunc struct {
	baseFunc
}

func (f *arithmeticMulFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("mul int.....")
	return lv * rv, false, nil
}

func (f *arithmeticMulFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("mul float.....")
	return
}

func (f *arithmeticMulFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("mul time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv * rv
	})
}

type arithmeticDivFuncFactory struct{}

func (fct *arithmeticDivFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticDivFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticDivFunc struct{ baseFunc }

func (f *arithmeticDivFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("div int.....")
	return lv / rv, false, nil
}

func (f *arithmeticDivFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("div float.....")
	return
}

func (f *arithmeticDivFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("div time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		if rv == 0 {
			return 0
		}
		return lv / rv
	})
}

type arithmeticModFuncFactory struct{}

func (fct *arithmeticModFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticModFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticModFunc struct {
	baseFunc
}

func (f *arithmeticModFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)

	fmt.Println("mod int.....")
	return lv % rv, false, nil
}

func (f *arithmeticModFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("mod float.....")
	return
}

func (f *arithmeticModFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("mod time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		if rv == 0 {
			return 0
		}
		return float64(int64(lv) % int64(rv))
	})
}
