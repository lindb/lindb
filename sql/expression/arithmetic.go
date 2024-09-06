package expression

import (
	"fmt"

	"github.com/lindb/lindb/spi"
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

func (f *arithmeticPlusFunc) EvalInt(row spi.Row) (val int64, err error) {
	lv, _ := f.args[0].EvalInt(row)
	rv, _ := f.args[1].EvalInt(row)
	fmt.Println("plus.....")
	return lv + rv, nil
}

func (f *arithmeticPlusFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return 0, nil
}

func (f *arithmeticPlusFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return nil, nil
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

func (f *arithmeticMinusFunc) EvalInt(row spi.Row) (val int64, err error) {
	lv, _ := f.args[0].EvalInt(row)
	rv, _ := f.args[1].EvalInt(row)
	fmt.Println("minus.....")
	return lv - rv, nil
}

func (f *arithmeticMinusFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return 0, nil
}

func (f *arithmeticMinusFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return nil, nil
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

func (f *arithmeticMulFunc) EvalInt(row spi.Row) (val int64, err error) {
	lv, _ := f.args[0].EvalInt(row)
	rv, _ := f.args[1].EvalInt(row)
	fmt.Println("mul.....")
	return lv * rv, nil
}

func (f *arithmeticMulFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return 0, nil
}

func (f *arithmeticMulFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return
}

type arithmeticDivFuncFactory struct{}

func (fct *arithmeticDivFuncFactory) NewFunc(args []Expression) Func {
	return &arithmeticDivFunc{
		baseFunc: baseFunc{args: args},
	}
}

type arithmeticDivFunc struct{ baseFunc }

func (f *arithmeticDivFunc) EvalInt(row spi.Row) (val int64, err error) {
	lv, _ := f.args[0].EvalInt(row)
	rv, _ := f.args[1].EvalInt(row)
	fmt.Println("div.....")
	return lv / rv, nil
}

func (f *arithmeticDivFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return
}

func (f *arithmeticDivFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return
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

func (f *arithmeticModFunc) EvalInt(row spi.Row) (val int64, err error) {
	lv, _ := f.args[0].EvalInt(row)
	rv, _ := f.args[1].EvalInt(row)

	fmt.Println("mod.....")
	return lv % rv, nil
}

func (f *arithmeticModFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return
}

func (f *arithmeticModFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return
}
