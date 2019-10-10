package aggregation

import (
	"github.com/lindb/lindb/aggregation/fields"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./expression.go -destination=./expression_mock.go -package=aggregation

// Expression represents expression eval like math calc, function call etc.
type Expression interface {
	// Eval evaluates the select item's expression
	Eval(timeSeries series.GroupedIterator)
	// ResultSet returns the eval result
	ResultSet() map[string]collections.FloatArray
	// Reset resets the expression context for reusing
	Reset()
}

// expression implement Expression interface, operator as below:
// 1. prepare field store based on time series iterator
// 2. eval the expression
// 3. build result set
type expression struct {
	pointCount  int
	interval    int64
	timeRange   timeutil.TimeRange
	selectItems []stmt.Expr

	fieldStore map[string]fields.Field
	resultSet  map[string]collections.FloatArray
}

// NewExpression creates an expression
func NewExpression(timeRange timeutil.TimeRange, interval int64, selectItems []stmt.Expr) Expression {
	return &expression{
		pointCount:  timeutil.CalPointCount(timeRange.Start, timeRange.End, interval),
		interval:    interval,
		timeRange:   timeRange,
		selectItems: selectItems,
		fieldStore:  make(map[string]fields.Field),
		resultSet:   make(map[string]collections.FloatArray),
	}
}

// Eval evaluates the select item's expression
func (e *expression) Eval(timeSeries series.GroupedIterator) {
	if len(e.selectItems) == 0 {
		return
	}
	// prepare expression context
	e.prepare(timeSeries)

	if len(e.fieldStore) == 0 {
		return
	}

	for _, selectItem := range e.selectItems {
		values := e.eval(nil, selectItem)
		if len(values) != 0 {
			item, ok := selectItem.(*stmt.SelectItem)
			if ok && len(item.Alias) > 0 {
				e.resultSet[item.Alias] = values[0]
			} else {
				e.resultSet[item.Rewrite()] = values[0]
			}
		}
	}
}

// ResultSet returns the eval result
func (e *expression) ResultSet() map[string]collections.FloatArray {
	return e.resultSet
}

// prepare prepares the field store
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

// eval evaluates the expression
func (e *expression) eval(parentFunc *stmt.CallExpr, expr stmt.Expr) []collections.FloatArray {
	switch ex := expr.(type) {
	case *stmt.SelectItem:
		return e.eval(nil, ex.Expr)
	case *stmt.CallExpr:
		return e.funcCall(ex)
	case *stmt.ParenExpr:
		return e.eval(nil, ex.Expr)
	case *stmt.BinaryExpr:
		return e.binaryEval(ex)
	case *stmt.FieldExpr:
		fieldName := ex.Name
		fieldValues, ok := e.fieldStore[fieldName]
		if !ok {
			return nil
		}

		// tests if has func with field
		if parentFunc == nil {
			return fieldValues.GetDefaultValues()
		}
		return fieldValues.GetValues(parentFunc.FuncType)
	default:
		return nil
	}
}

// funcCall calls the function
func (e *expression) funcCall(expr *stmt.CallExpr) []collections.FloatArray {
	var params []collections.FloatArray
	for _, param := range expr.Params {
		paramValues := e.eval(expr, param)
		if len(paramValues) != 1 {
			return nil
		}
		params = append(params, paramValues[0])
	}
	result := function.FuncCall(expr.FuncType, params...)
	if result == nil {
		return nil
	}
	return []collections.FloatArray{result}
}

// binaryEval evaluates binary operator
func (e *expression) binaryEval(expr *stmt.BinaryExpr) []collections.FloatArray {
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
		return []collections.FloatArray{result}
	}

	return nil
}

// Reset resets the expression context for reusing
func (e *expression) Reset() {
	for _, f := range e.fieldStore {
		f.Reset()
	}
	e.resultSet = make(map[string]collections.FloatArray)
}
