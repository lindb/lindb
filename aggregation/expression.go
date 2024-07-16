// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package aggregation

import (
	"strconv"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/aggregation/fields"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./expression.go -destination=./expression_mock.go -package=aggregation

var log = logger.GetLogger("Aggregation", "Expression")

// Expression represents Expression eval like math calc, function call etc.
// 1. prepare field store based on time series iterator
// 2. eval the Expression
// 3. build result set
type Expression interface {
	// Eval evaluates the select item's expression.
	Eval(timeSeries series.GroupedIterator)
	// ResultSet returns the eval result, returns field name(alias) => series data.
	ResultSet() map[string]*collections.FloatArray
	// Reset resets the Expression context for reusing.
	Reset()
}

// expression implements Expression interface.
type expression struct {
	fieldStore  map[field.Name]fields.Field
	resultSet   map[string]*collections.FloatArray
	selectItems []stmt.Expr
	timeRange   timeutil.TimeRange
	pointCount  int
	interval    int64
}

// NewExpression creates an Expression instance.
func NewExpression(timeRange timeutil.TimeRange, interval int64, selectItems []stmt.Expr) Expression {
	return &expression{
		pointCount:  timeutil.CalPointCount(timeRange.Start, timeRange.End, interval) + 1,
		interval:    interval,
		timeRange:   timeRange,
		selectItems: selectItems,
		fieldStore:  make(map[field.Name]fields.Field),
		resultSet:   make(map[string]*collections.FloatArray),
	}
}

// Eval evaluates the select item's Expression
func (e *expression) Eval(timeSeries series.GroupedIterator) {
	if len(e.selectItems) == 0 {
		return
	}
	// prepare Expression context
	e.prepare(timeSeries)

	if len(e.fieldStore) == 0 {
		return
	}

	for _, selectItem := range e.selectItems {
		values := e.eval(nil, selectItem)
		if len(values) != 0 {
			if item, ok := selectItem.(*stmt.SelectItem); ok && item.Alias != "" {
				e.resultSet[item.Alias] = values[0]
			} else {
				e.resultSet[item.Rewrite()] = values[0]
			}
		}
	}
}

// ResultSet returns the eval result, returns field name(alias) => series data.
func (e *expression) ResultSet() map[string]*collections.FloatArray {
	return e.resultSet
}

// prepare the field store.
func (e *expression) prepare(timeSeries series.GroupedIterator) {
	if timeSeries == nil {
		return
	}
	for timeSeries.HasNext() {
		fieldSeries := timeSeries.Next()
		fieldName := fieldSeries.FieldName()
		fieldType := fieldSeries.FieldType()
		f := fields.NewDynamicField(fieldType, e.timeRange.Start, e.interval, e.pointCount)
		e.fieldStore[fieldName] = f
		f.SetValue(fieldSeries)
	}
}

// eval evaluates the Expression.
func (e *expression) eval(parentFunc *stmt.CallExpr, expr stmt.Expr) []*collections.FloatArray {
	switch ex := expr.(type) {
	case *stmt.SelectItem:
		return e.eval(nil, ex.Expr)
	case *stmt.CallExpr:
		switch ex.FuncType {
		case function.Quantile:
			return e.quantile(ex)
		default:
			return e.funcCall(ex)
		}
	case *stmt.ParenExpr:
		return e.eval(nil, ex.Expr)
	case *stmt.BinaryExpr:
		return e.binaryEval(ex)
	case *stmt.NumberLiteral:
		values := collections.NewFloatArray(e.pointCount)
		for i := 0; i < e.pointCount; i++ {
			values.SetValue(i, ex.Val)
		}
		values.SetSingle(true)
		return []*collections.FloatArray{values}
	case *stmt.FieldExpr:
		fieldName := ex.Name
		if fieldValues, ok := e.fieldStore[field.Name(fieldName)]; ok {
			// tests if it has func with field
			if parentFunc == nil {
				return fieldValues.GetDefaultValues()
			}
			// get field data by function type
			return fieldValues.GetValues(parentFunc.FuncType)
		}
		return nil
	default:
		return nil
	}
}

func (e *expression) quantile(expr *stmt.CallExpr) []*collections.FloatArray {
	histogramFields := make(map[float64][]*collections.FloatArray)
	if len(expr.Params) != 1 {
		return nil
	}
	quantileValue, err := strconv.ParseFloat(expr.Params[0].Rewrite(), 64)
	if err != nil {
		return nil
	}
	for fieldName, df := range e.fieldStore {
		if df.Type() == field.HistogramField {
			var upperBound float64
			upperBound, err = metric.UpperBound(fieldName.String())
			if err != nil {
				continue
			}
			histogramFields[upperBound] = df.GetDefaultValues()
		}
	}
	if len(histogramFields) == 0 {
		return nil
	}
	array, err := function.QuantileCall(quantileValue, histogramFields)
	if err != nil {
		log.Warn("histogram quantile call error", logger.Error(err))
		return nil
	}
	return []*collections.FloatArray{array}
}

// funcCall calls the function
func (e *expression) funcCall(expr *stmt.CallExpr) []*collections.FloatArray {
	var params []*collections.FloatArray
	for _, param := range expr.Params {
		paramValues := e.eval(expr, param)
		if len(paramValues) == 0 {
			return nil
		}
		params = append(params, paramValues...)
	}
	var result *collections.FloatArray
	switch expr.FuncType {
	case function.Avg:
		result = function.AvgCall(params...)
	case function.Rate:
		result = function.RateCall(e.interval, params...)
	default:
		result = function.FuncCall(expr.FuncType, params...)
	}
	if result == nil {
		return nil
	}
	return []*collections.FloatArray{result}
}

// binaryEval evaluates binary operator
func (e *expression) binaryEval(expr *stmt.BinaryExpr) []*collections.FloatArray {
	binaryOP := expr.Operator
	if binaryOP == stmt.ADD || binaryOP == stmt.SUB || binaryOP == stmt.DIV || binaryOP == stmt.MUL {
		left := e.eval(nil, expr.Left)
		if len(left) != 1 {
			return nil
		}
		right := e.eval(nil, expr.Right)
		if len(right) != 1 {
			return nil
		}
		result := binaryEval(binaryOP, left[0], right[0])
		return []*collections.FloatArray{result}
	}

	return nil
}

// Reset resets the Expression context for reusing.
func (e *expression) Reset() {
	for _, f := range e.fieldStore {
		f.Reset()
	}
	e.resultSet = make(map[string]*collections.FloatArray)
}
