package sql

import (
	"errors"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/sql/grammar"
	"github.com/eleme/lindb/sql/stmt"
)

var log = logger.GetLogger("sql/parser")
var errorHandle = &errorListener{}
var walker = antlr.ParseTreeWalkerDefault

// Parse parses sql using the grammar of LinDB query language
func Parse(sql string) (stmt *stmt.Query, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic when sql parse")
			}
			log.Error("parse sql", logger.String("sql", sql), logger.Error(err), logger.Stack())
			stmt = nil
		}
	}()

	input := antlr.NewInputStream(sql)

	lexer := grammar.NewSQLLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorHandle)

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	parser := grammar.NewSQLParser(tokens)
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorHandle)

	ctx := parser.Statement()

	// create sql listener
	listener := listener{}

	walker.Walk(&listener, ctx)

	stmt, err = listener.statement()
	return stmt, err
}
