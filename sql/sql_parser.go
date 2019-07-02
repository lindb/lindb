package sql

import (
	"github.com/eleme/lindb/pkg/proto"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"sync"

	parser "github.com/eleme/lindb/sql/grammar"
)

var (
	instance *LinDbSQLParser
	once     sync.Once
)

type LinDbSQLParser struct {
	lexer    *parser.SQLLexer
	parser   *parser.SQLParser
	walker   *antlr.ParseTreeWalker
	listener *Listener
}

// InitSQLPlan init lindb sql antlr4 engine
func (l *LinDbSQLParser) InitSQLParser() {
	input := antlr.NewInputStream("")

	// create sql lexer
	l.lexer = parser.NewSQLLexer(input)

	// create sql token stream
	stream := antlr.NewCommonTokenStream(l.lexer, antlr.TokenDefaultChannel)
	l.lexer.RemoveErrorListeners()

	// create the sql parser
	p := parser.NewSQLParser(stream)
	p.BuildParseTrees = true
	l.parser = p
	l.parser.RemoveErrorListeners()

	// create sql listener
	listener := Listener{}
	listener.InitSQLListener()

	// finally create default walk tree
	l.walker = antlr.ParseTreeWalkerDefault
	l.listener = &listener
}

// Plan antlr4 parse lindb sql
func (l *LinDbSQLParser) Parser(sql string) *Listener {
	input := antlr.NewInputStream(sql)
	l.lexer = parser.NewSQLLexer(input)
	tokens := antlr.NewCommonTokenStream(l.lexer, antlr.TokenDefaultChannel)
	l.parser = parser.NewSQLParser(tokens)
	ctx := l.parser.Statement()
	l.walker.Walk(l.listener, ctx)
	return l.listener
}

func (l *LinDbSQLParser) ParserTemp(sql string) *proto.Stmt {
	errorListener := sqlErrorListener{}
	is := antlr.NewInputStream(sql)
	lexer := parser.NewSQLLexer(is)
	l.lexer = lexer
	l.lexer.RemoveErrorListeners()
	l.lexer.AddErrorListener(&errorListener)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	l.parser = parser.NewSQLParser(tokens)
	l.parser.RemoveErrorListeners()
	l.parser.AddErrorListener(&errorListener)
	l.walker = antlr.ParseTreeWalkerDefault
	ctx := parser.NewEmptyStatementContext()
	listener := Listener{}
	listener.InitSQLListener()
	l.walker.Walk(&listener, ctx)
	return listener.GetStatement()
}

func GetInstance() *LinDbSQLParser {
	once.Do(func() {
		instance = new(LinDbSQLParser)
		instance.InitSQLParser()
	})
	return instance
}
