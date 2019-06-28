package sql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ShowStats(t *testing.T) {
	sql := "show stats"
	statement := sqlPlan.Plan(sql).GetStatement().GetShowStats()
	assert.Equal(t, "", statement.GetModule())

	sql = "show stats for database"
	statement = sqlPlan.Plan(sql).GetStatement().GetShowStats()
	assert.Equal(t, "database", statement.GetModule())
}

func Test_ShowDatabases(t *testing.T) {
	assert.True(t, true, sqlPlan.Plan("show databases").GetStatement().GetShowDatabases())
	assert.True(t, true, sqlPlan.Plan("SHOW databases").GetStatement().GetShowDatabases())
}

func Test_ShowNode(t *testing.T) {
	assert.True(t, true, sqlPlan.Plan("show node").GetStatement().GetShowNode())
}

func Test_ShowQueries(t *testing.T) {
	sql := "show queries"
	assert.True(t, true, sqlPlan.Plan(sql).GetStatement().GetShowQueries())

	sql = "SHOW queries"
	assert.True(t, true, sqlPlan.Plan(sql).GetStatement().GetShowQueries())
}

func Test_KillQuery(t *testing.T) {
	sql := "show measurements"
	statement := sqlPlan.Plan(sql).GetStatement().GetShowMetric()
	assert.Equal(t, "", statement.GetName())
	assert.Equal(t, int32(50), statement.GetLimit())

	sql = "show measurements with measurement = abc limit 100"
	statement = sqlPlan.Plan(sql).GetStatement().GetShowMetric()
	assert.Equal(t, "abc", statement.GetName())
	assert.Equal(t, int32(100), statement.GetLimit())
}

func Test_ShowFieldKeys(t *testing.T) {
	sql := "show field keys from cpu limit 100"
	statement := sqlPlan.Plan(sql).GetStatement().GetShowFieldKeys()
	assert.Equal(t, "cpu", statement.GetMeasurement())
	assert.Equal(t, int32(100), statement.GetLimit())
}

func Test_ShowTagKeys(t *testing.T) {
	sql := "show tag keys from cpu limit 100"
	statement := sqlPlan.Plan(sql).GetStatement().GetShowTagKeys()
	assert.Equal(t, "cpu", statement.GetMeasurement())
	assert.Equal(t, int32(100), statement.GetLimit())
}

func Test_ShowTagValues(t *testing.T) {
	sql := "show tag values from cpu with key = host limit 100"
	statement := sqlPlan.Plan(sql).statement.GetShowTagValues()
	assert.Equal(t, "cpu", statement.GetMeasurement())
	assert.Equal(t, int32(100), statement.GetLimit())
	assert.Equal(t, "host", statement.GetTagKey())

	sql = "show tag values from 'cpu' with key = 'host' limit 100"
	statement = sqlPlan.Plan(sql).statement.GetShowTagValues()
	assert.Equal(t, "cpu", statement.GetMeasurement())
	assert.Equal(t, int32(100), statement.GetLimit())
	assert.Equal(t, "host", statement.GetTagKey())

	sql = "show tag values from 'cpu' with key = 'host' where value = 'host1*' limit 100"
	statement = sqlPlan.Plan(sql).GetStatement().GetShowTagValues()
	assert.Equal(t, "cpu", statement.GetMeasurement())
	assert.Equal(t, int32(100), statement.GetLimit())
	assert.Equal(t, "host", statement.GetTagKey())
	assert.Equal(t, "host1", statement.GetTagValue())
}
