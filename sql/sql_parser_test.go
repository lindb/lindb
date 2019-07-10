package sql

import (
	"github.com/eleme/lindb/pkg/proto"

	"github.com/stretchr/testify/assert"

	"testing"
)

var sqlParser = GetInstance()

func Test_Sql_Parser(t *testing.T) {
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
