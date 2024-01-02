package tree

import (
	"errors"

	"github.com/antlr4-go/antlr/v4"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/grammar"
)

var log = logger.GetLogger("SQL", "Parser")

var walker = antlr.ParseTreeWalkerDefault

type SQLParser struct{}

func GetParser() *SQLParser {
	return &SQLParser{}
}

func (p *SQLParser) CreateStatement(sql string, idAllocator *NodeIDAllocator) (stmt Statement, err error) {
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
	// lexer.AddErrorListener(errorHandle)

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := grammar.NewSQLParser(tokens)
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	// parser.AddErrorListener(errorHandle)
	// first, try parsing with potentially faster SLL mode
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	// TODO: fail to LL mode
	parseTree := parser.Statement()

	visitor := NewAstVisitor(idAllocator)
	node := visitor.Visit(parseTree)
	if node != nil {
		stmt = node.(Statement)
	}
	return
}
