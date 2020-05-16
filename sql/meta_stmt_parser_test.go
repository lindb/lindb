package sql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestMetaStmt_validation(t *testing.T) {
	queryStmt := newMetaStmtParser(stmt.TagKey)
	// case 1: stmt err
	queryStmt.err = fmt.Errorf("err")
	s, err := queryStmt.build()
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestMetaStmt_ShowDatabases(t *testing.T) {
	sql := "show databases"
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Database, query.Type)
}

func TestMetaStmt_ShowNamespace(t *testing.T) {
	sql := "show namespaces"
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Namespace, query.Type)

	sql = "show namespaces where namespace='abc' limit 10"
	q, err = Parse(sql)
	query = q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Namespace, query.Type)
	assert.Equal(t, "abc", query.Prefix)
	assert.Equal(t, 10, query.Limit)
}

func TestMetaStmt_ShowMeasurement(t *testing.T) {
	sql := "show measurements"
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Metric, query.Type)

	sql = "show measurements on 'ns' where measurement='abc' limit 10"
	q, err = Parse(sql)
	query = q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Metric, query.Type)
	assert.Equal(t, "abc", query.Prefix)
	assert.Equal(t, "ns", query.Namespace)
	assert.Equal(t, 10, query.Limit)
}

func TestMetaStmt_ShowFields(t *testing.T) {
	sql := "show fields on 'ns' from 'cpu' "
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.Field, query.Type)
	assert.Equal(t, "cpu", query.MetricName)
	assert.Equal(t, "ns", query.Namespace)
}

func TestMetaStmt_ShowTagKeys(t *testing.T) {
	sql := "show tag keys on 'ns' from 'cpu' "
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.TagKey, query.Type)
	assert.Equal(t, "cpu", query.MetricName)
	assert.Equal(t, "ns", query.Namespace)
}

func TestMetaStmt_ShowTagValues(t *testing.T) {
	sql := "show tag values on 'ns' from 'cpu' with key = 'key1' where key1='value1' and key2='value2' limit 10"
	q, err := Parse(sql)
	query := q.(*stmt.Metadata)
	assert.Nil(t, err)
	assert.Equal(t, stmt.TagValue, query.Type)
	assert.Equal(t, "cpu", query.MetricName)
	assert.Equal(t, "ns", query.Namespace)
	assert.Equal(t, "key1", query.TagKey)
	assert.Equal(t, 10, query.Limit)
	expr := query.Condition.(*stmt.BinaryExpr)
	assert.Equal(t,
		stmt.BinaryExpr{
			Left:     &stmt.EqualsExpr{Key: "key1", Value: "value1"},
			Operator: stmt.AND,
			Right:    &stmt.EqualsExpr{Key: "key2", Value: "value2"},
		}, *expr)
}
