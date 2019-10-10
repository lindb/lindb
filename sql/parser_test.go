package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
)

func Test_SQL_Parse(t *testing.T) {
	query, err := Parse("select f+100 from cpu")
	assert.NoError(t, err)
	_ = encoding.JSONMarshal(&query)
	//query1 := stmt.Query{}
	//err = encoding.JSONUnmarshal(data, &query1)
	//assert.NoError(t, err)
}

func BenchmarkSQLParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("select f from cpu " +
			"where (ip in ('1.1.1.1','2.2.2.2')" +
			" and region='sh') and (path='/data' or path='/home')")
	}
}
