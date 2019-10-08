package aggregation

//import (
//	"github.com/golang/mock/gomock"
//	"github.com/lindb/lindb/aggregation/function"
//	"github.com/lindb/lindb/series"
//	"github.com/lindb/lindb/series/field"
//	"github.com/lindb/lindb/sql"
//	"github.com/lindb/lindb/sql/stmt"
//	"github.com/stretchr/testify/assert"
//	"math"
//	"testing"
//)
//
//func TestExpression_prepare(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	timeSeries := mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//		"f2": field.MaxField,
//	})
//
//	query, _ := sql.Parse("select f1,f2 from cpu")
//	expression := NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet := expression.ResultSet()
//	assert.Equal(t, 2, len(resultSet))
//	assert.NotNil(t, resultSet["f1"])
//	assert.NotNil(t, resultSet["f2"])
//
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{})
//	expression = NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//
//	timeSeries1 := series.NewMockIterator(ctrl)
//	timeSeries1.EXPECT().HasNext().Return(true)
//	it3 := series.NewMockFieldIterator(ctrl)
//	timeSeries1.EXPECT().Next().Return(it3)
//	timeSeries1.EXPECT().HasNext().Return(false)
//	expression = NewExpression(timeSeries1, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//
//	expression = NewExpression(nil, 10, []stmt.Expr{})
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//}
//
//func TestExpression_Paren(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	timeSeries := mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//		"f2": field.MaxField,
//		"f3": field.MaxField,
//	})
//
//	query, _ := sql.Parse("select (f1+f2)*f3 as f from cpu")
//	expression := NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet := expression.ResultSet()
//	assert.Equal(t, 1, len(resultSet))
//
//	value := resultSet["f"]
//	assert.Equal(t, 1, value.Size())
//	assert.Equal(t, 2.42, math.Floor(value.GetValue(4)*100)/100)
//}
//
//func TestExpression_BinaryEval(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	timeSeries := mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//		"f2": field.MaxField,
//	})
//
//	query, _ := sql.Parse("select f1+f2 from cpu")
//	expression := NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet := expression.ResultSet()
//	assert.Equal(t, 1, len(resultSet))
//
//	value := resultSet["f1+f2"]
//	assert.Equal(t, 1, value.Size())
//	assert.Equal(t, 2.2, value.GetValue(4))
//
//	it1 := mockSingleIterator(ctrl, "f1", field.SumField)
//	it2 := mockSingleIterator(ctrl, "f2", field.MaxField)
//	it3 := mockSingleIterator(ctrl, "f3", field.MaxField)
//
//	timeSeries1 := series.NewMockIterator(ctrl)
//	timeSeries1.EXPECT().HasNext().Return(true)
//	timeSeries1.EXPECT().Next().Return(it1)
//	timeSeries1.EXPECT().HasNext().Return(true)
//	timeSeries1.EXPECT().Next().Return(it2)
//	timeSeries1.EXPECT().HasNext().Return(true)
//	timeSeries1.EXPECT().Next().Return(it3)
//	timeSeries1.EXPECT().HasNext().Return(false)
//	query, _ = sql.Parse("select f1+f2*f3 from cpu")
//	expression = NewExpression(timeSeries1, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 1, len(resultSet))
//
//	value = resultSet["f1+f2*f3"]
//	assert.Equal(t, 1, value.Size())
//	assert.Equal(t, 2.31, math.Floor(value.GetValue(4)*100)/100)
//
//	// right is nil, return nil
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//	})
//	expression = NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//
//	// left is nil, return nil
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{
//		"f2": field.MaxField,
//	})
//	expression = NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//
//	// binary operator not accept, return nil
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//		"f2": field.MaxField,
//	})
//	expression = NewExpression(timeSeries, 10, []stmt.Expr{&stmt.SelectItem{Expr: &stmt.BinaryExpr{
//		Left:     &stmt.FieldExpr{Name: "f1"},
//		Operator: stmt.AND,
//		Right:    &stmt.FieldExpr{Name: "f2"},
//	}}})
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//}
//
//func TestExpression_FuncCall_Sum(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	timeSeries := mockTimeSeries(ctrl, map[string]field.Type{
//		"f1": field.SumField,
//	})
//
//	query, _ := sql.Parse("select sum(f1) from cpu")
//	expression := NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet := expression.ResultSet()
//	assert.Equal(t, 1, len(resultSet))
//
//	value := resultSet["sum(f1)"]
//	assert.Equal(t, 1, value.Size())
//	assert.Equal(t, 1.1, value.GetValue(4))
//
//	// return nil
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{
//		"f2": field.SumField,
//	})
//	expression = NewExpression(timeSeries, 10, query.SelectItems)
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//
//	timeSeries = mockTimeSeries(ctrl, map[string]field.Type{
//		"f2": field.SumField,
//	})
//	expression = NewExpression(timeSeries, 10, []stmt.Expr{&stmt.SelectItem{Expr: &stmt.CallExpr{
//		FuncType: function.Sum,
//	}}})
//	expression.Eval()
//	resultSet = expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//}
//
//func TestExpression_NotSupport_Expr(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	timeSeries := mockTimeSeries(ctrl)
//	expression := NewExpression(timeSeries, 10, []stmt.Expr{&stmt.EqualsExpr{}})
//	expression.Eval()
//	resultSet := expression.ResultSet()
//	assert.Equal(t, 0, len(resultSet))
//}
