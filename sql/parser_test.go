package sql

import (
	"fmt"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

func Test_SQL_Parse(t *testing.T) {
	query, err := Parse("select f+100 from cpu")
	assert.NoError(t, err)
	data := encoding.JSONMarshal(&query)
	query1 := stmt.Query{}
	err = encoding.JSONUnmarshal(data, &query1)
	assert.NoError(t, err)
}

func Test_Meta_SQL_Parse(t *testing.T) {
	query, err := Parse("show databases")
	assert.NoError(t, err)
	_, ok := query.(*stmt.Metadata)
	assert.True(t, ok)
}

func TestParse_panic(t *testing.T) {
	defer func() {
		newSQLParserFunc = grammar.NewSQLParser
	}()
	newSQLParserFunc = func(input antlr.TokenStream) *grammar.SQLParser {
		panic(fmt.Errorf("err"))
	}
	_, err := Parse("select f+100 from cpu")
	assert.Error(t, err)

	newSQLParserFunc = func(input antlr.TokenStream) *grammar.SQLParser {
		panic(123)
	}
	_, err = Parse("select f+100 from cpu")
	assert.Error(t, err)
}

func BenchmarkSQLParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("select f from cpu " +
			"where (ip in ('1.1.1.1','2.2.2.2')" +
			" and region='sh') and (path='/data' or path='/home')")
	}
}
