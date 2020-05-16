package sql

import (
	"errors"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	newSQLParserFunc = grammar.NewSQLParser
)

var log = logger.GetLogger("sql", "Parser")
var errorHandle = &errorListener{}
var walker = antlr.ParseTreeWalkerDefault

// Parse parses sql using the grammar of LinDB query language
func Parse(sql string) (stmt stmt.Statement, err error) {
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

	parser := newSQLParserFunc(tokens)
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
