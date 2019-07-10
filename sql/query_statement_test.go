package sql

import (
	"github.com/eleme/lindb/pkg/proto"

	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_QueryPlan(t *testing.T) {
	sql := "select f from test where a='a1' or b='b1' and time>'20190410 00:00:00' " +
		"and time<'20190410 10:00:00' order by f desc"
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, 2, len(query.ConditionAggregators))
	assert.Equal(t, proto.LogicOperator_OR, query.Condition.Operator)
	assert.Equal(t, 0, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
	assert.Equal(t, int32(20), query.Limit)
	assert.Equal(t, true, query.OrderBy.GetDesc())
}

func Test_QueryPlan1(t *testing.T) {
	sql := "select rate(max(d_avg(a), 'name' = '1*')),b,max(c),d_max(e),d_avg(a)" +
		" from test group by 'type',time(1m)"
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(60000), query.Interval)
	assert.Equal(t, 1, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, []string{"type"}, query.GroupByExpr.GroupBy)
	assert.Equal(t, 5, len(query.FieldExprList))
}

func Test_QueryPlan2(t *testing.T) {
	sql := "select f from test where a='a1' and (e='ee1' or (f='f1' and  " +
		"(g='g1' or g='g2' ))) and (aa='a1' and bb='b1') and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 6, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}

func Test_QueryPlan3(t *testing.T) {
	sql := "select f from test where a='a1' or b='b1' and time>'20190410 00:00:00'" +
		" and time<'20190410 10:00:00'"
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, proto.LogicOperator_OR, query.Condition.Operator)
	assert.Equal(t, 0, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}

func Test_QueryPlan4(t *testing.T) {
	sql := "select f from test where (a='a1' or b='b1') and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 1, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}

func Test_QueryPlan5(t *testing.T) {
	sql := "select f from test where a='a1' and (e='e1' or (f='f1' and f='f2') )" +
		" and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, "f", query.FieldExprList[0].Alias)
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, "f", query.ConditionAggregators[0].Field)
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 4, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}

func Test_QueryPlan6(t *testing.T) {
	sql := "select f from test where a='a1' and (e='e1' or (f='f1' and " +
		" (g='g1' or g='g2' )) ) and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, "f", query.FieldExprList[0].Alias)
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, "f", query.ConditionAggregators[0].Field)
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 4, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}

func Test_QueryPlan7(t *testing.T) {
	sql := "select f from test where a='a1' and (ee='e1' and (e='e1' or f='f1') and" +
		" (g='g1' or g='g2' )) and (aa='a1' or bb='b1') and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, "f", query.FieldExprList[0].Alias)
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, "f", query.ConditionAggregators[0].Field)
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 4, len(query.Condition.TagFilters))
	assert.Equal(t, 2, len(query.Condition.Condition))
}

func Test_QueryPlan8(t *testing.T) {
	sql := "select f from test where a='a1' and (ee='ee1' and ((e='e1' or f='f1') or " +
		"(g='g1' or g='g2' ))) and (aa='a1' or bb='b1') and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, "f", query.FieldExprList[0].Alias)
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, "f", query.ConditionAggregators[0].Field)
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 4, len(query.Condition.TagFilters))
	assert.Equal(t, 2, len(query.Condition.Condition))
}

func Test_QueryPlan9(t *testing.T) {
	sql := "select f from test where a='a1' and (e='ee1' or (f='f1' and  " +
		"(g='g1' or g='g2' ))) and (aa='a1' and bb='b1') and b='b1' and d='d1' and c='c1' "
	query := sqlParser.Parser(sql).stmt.build()
	assert.Equal(t, "test", query.Measurement)
	assert.Equal(t, int64(0), query.Interval)
	assert.Equal(t, 0, len(query.GroupByExpr.GroupBy))
	assert.Equal(t, 1, len(query.FieldExprList))
	assert.Equal(t, "f", query.FieldExprList[0].Alias)
	assert.Equal(t, 1, len(query.ConditionAggregators))
	assert.Equal(t, "f", query.ConditionAggregators[0].Field)
	assert.Equal(t, proto.LogicOperator_AND, query.Condition.Operator)
	assert.Equal(t, 6, len(query.Condition.TagFilters))
	assert.Equal(t, 1, len(query.Condition.Condition))
}
