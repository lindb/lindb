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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

var now, _ = timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
var familyTime, _ = timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")

func TestExpression_prepare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sumSeries := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	maxSeries := mockTimeSeries(ctrl, familyTime+timeutil.OneHour, "f2", field.MinField, field.Min)
	timeSeries := series.NewMockGroupedIterator(ctrl)

	// case 1: text express
	q, _ := sql.Parse("select f1,f2 from cpu")
	query := q.(*stmt.Query)
	expression := NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(sumSeries),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(maxSeries),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet := expression.ResultSet()
	assert.Equal(t, 2, len(resultSet))
	rs := resultSet["f1"]
	assert.Equal(t, 50.0, rs.GetValue(50-10))
	rs = resultSet["f2"]
	assert.Equal(t, 4.0, rs.GetValue(4+60-10))
	assert.Equal(t, 50.0, rs.GetValue(50+60-10))

	// case 2: test reset
	expression.Reset()
	expression.Eval(nil)
	resultSet = expression.ResultSet()
	rs = resultSet["f1"]
	assert.True(t, rs.IsEmpty())
	rs = resultSet["f2"]
	assert.True(t, rs.IsEmpty())

	// case 3: test new expression for nil eval
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	expression.Eval(nil)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))

	// case 4：test no data
	timeSeries1 := series.NewMockIterator(ctrl)
	timeSeries1.EXPECT().FieldName().Return(field.Name("f1"))
	timeSeries1.EXPECT().FieldType().Return(field.SumField)
	timeSeries1.EXPECT().HasNext().Return(false)
	timeSeries.EXPECT().HasNext().Return(true)
	timeSeries.EXPECT().Next().Return(timeSeries1)
	timeSeries.EXPECT().HasNext().Return(false)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))

	// case 5: test no match field
	sumSeries = mockTimeSeries(ctrl, familyTime, "f3", field.SumField, field.Sum)
	maxSeries = mockTimeSeries(ctrl, familyTime+timeutil.OneHour, "f4", field.MinField, field.Min)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(sumSeries),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(maxSeries),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	q, _ = sql.Parse("select f1,f2 from cpu")
	query = q.(*stmt.Query)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))
}

func TestExpression_Paren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	series1 := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	series2 := mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	series3 := mockTimeSeries(ctrl, familyTime, "f3", field.MinField, field.Min)
	timeSeries := series.NewMockGroupedIterator(ctrl)

	q, _ := sql.Parse("select (f1+f2)*f3 as f from cpu")
	query := q.(*stmt.Query)
	expression := NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series3),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet := expression.ResultSet()
	assert.Equal(t, 1, len(resultSet))

	value := resultSet["f"]
	assert.Equal(t, 1, value.Size())
	assert.Equal(t, (50.0+50.0)*50.0, value.GetValue(50-10))
}

func TestExpression_BinaryEval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	series1 := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	series2 := mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	timeSeries := series.NewMockGroupedIterator(ctrl)

	q, _ := sql.Parse("select (f1+f2)*100 as f from cpu")
	query := q.(*stmt.Query)
	expression := NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet := expression.ResultSet()
	value := resultSet["f"]
	assert.Equal(t, 1, value.Size())
	assert.Equal(t, (50.0+50.0)*100.0, value.GetValue(50-10))

	series1 = mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	series2 = mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	series3 := mockTimeSeries(ctrl, familyTime, "f3", field.MinField, field.Min)
	q, _ = sql.Parse("select f1+f2*f3 from cpu")
	query = q.(*stmt.Query)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series3),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 1, len(resultSet))
	value = resultSet["f1+f2*f3"]
	assert.Equal(t, 1, value.Size())
	assert.Equal(t, 50.0+50.0*50.0, value.GetValue(50-10))

	// right is nil, return nil
	series1 = mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Nil(t, resultSet["f1+f2*f3"])

	// left is nil, return nil
	series2 = mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Nil(t, resultSet["f1+f2*f3"])

	// binary operator not accept, return nil
	series1 = mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	series2 = mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, []stmt.Expr{&stmt.SelectItem{Expr: &stmt.BinaryExpr{
		Left:     &stmt.FieldExpr{Name: "f1"},
		Operator: stmt.AND,
		Right:    &stmt.FieldExpr{Name: "f2"},
	}}})
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))
}

func TestExpression_FuncCall_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	series1 := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	timeSeries := series.NewMockGroupedIterator(ctrl)

	q, _ := sql.Parse("select sum(f1) from cpu")
	query := q.(*stmt.Query)
	expression := NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet := expression.ResultSet()
	assert.Equal(t, 1, len(resultSet))

	value := resultSet["sum(f1)"]
	assert.Equal(t, 1, value.Size())
	assert.Equal(t, 50.0, value.GetValue(50-10))

	// return nil
	series1 = mockTimeSeries(ctrl, familyTime, "f2", field.SumField, field.Sum)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))

	series1 = mockTimeSeries(ctrl, familyTime, "f2", field.SumField, field.Sum)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, []stmt.Expr{&stmt.SelectItem{Expr: &stmt.CallExpr{
		FuncType: function.Sum,
	}}})
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))
}

func TestExpression_NotSupport_Expr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expression := NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, []stmt.Expr{})
	expression.Eval(nil)
	resultSet := expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))

	timeSeries := series.NewMockGroupedIterator(ctrl)
	series1 := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	expression = NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, []stmt.Expr{&stmt.EqualsExpr{}})
	gomock.InOrder(
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	expression.Eval(timeSeries)
	resultSet = expression.ResultSet()
	assert.Equal(t, 0, len(resultSet))
}
