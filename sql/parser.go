// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package sql

import (
	"errors"
	"strings"
	"sync"

	antlr "github.com/antlr/antlr4/runtime/Go/antlr/v4"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/grammar"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	getSQLParserFunc = getSQLParser
)

var log = logger.GetLogger("SQL", "Parser")

var errorHandle = &errorListener{}

var walker = antlr.ParseTreeWalkerDefault

// Parse parses sql using the grammar of LinDB query language
func Parse(sql string) (stmt stmtpkg.Statement, err error) {
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

	sql = strings.ReplaceAll(sql, `\"`, `"`)
	input := antlr.NewInputStream(sql)

	lexer := getSQLLexer(input)
	defer putSQLLexer(lexer)

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	parser := getSQLParserFunc(tokens)
	defer putSQLParser(parser)

	ctx := parser.Statement()

	// create sql listener
	sqlListener := listener{}

	walker.Walk(&sqlListener, ctx)

	stmt, err = sqlListener.statement()
	return stmt, err
}

var (
	lexerPool  sync.Pool
	parserPool sync.Pool
)

func getSQLLexer(input *antlr.InputStream) *grammar.SQLLexer {
	lexer := lexerPool.Get()
	if lexer == nil {
		lexer := grammar.NewSQLLexer(input)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(errorHandle)
		return lexer
	}
	l := lexer.(*grammar.SQLLexer)
	l.SetInputStream(input)
	return l
}

func putSQLLexer(l *grammar.SQLLexer) {
	lexerPool.Put(l)
}

// getSQLParser picks a cached parser from the pool
func getSQLParser(tokenStream *antlr.CommonTokenStream) *grammar.SQLParser {
	parser := parserPool.Get()
	if parser == nil {
		parser := grammar.NewSQLParser(tokenStream)
		parser.BuildParseTrees = true
		parser.RemoveErrorListeners()
		parser.AddErrorListener(errorHandle)
		return parser
	}
	p := parser.(*grammar.SQLParser)
	p.SetInputStream(tokenStream)
	return p
}

func putSQLParser(p *grammar.SQLParser) {
	parserPool.Put(p)
}
