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

package tree

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/grammar"
)

var log = logger.GetLogger("SQL", "Parser")

type Parser struct{}

func GetParser() *Parser {
	return &Parser{}
}

type MyErrorListener struct {
	*antlr.DefaultErrorListener
	errors []string
}

// SyntaxError captures ANTLR syntax errors so CreateStatement can surface them
// instead of silently producing a partial (incorrect) parse result.
func (l *MyErrorListener) SyntaxError(_ antlr.Recognizer, _ any, line, col int, msg string, _ antlr.RecognitionException) {
	l.errors = append(l.errors, fmt.Sprintf("line %d:%d %s", line, col, msg))
}

func (p *Parser) CreateStatement(sql string, idAllocator *NodeIDAllocator) (stmt Statement, err error) {
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

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	newParser := func() (*grammar.SQLParser, *MyErrorListener) {
		listener := &MyErrorListener{}
		tokens.Reset()
		p := grammar.NewSQLParser(tokens)
		p.BuildParseTrees = true
		p.RemoveErrorListeners()
		p.AddErrorListener(listener)
		return p, listener
	}

	// Stage 1: attempt fast SLL prediction mode.
	sllParser, sllListener := newParser()
	sllParser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	parseTree := sllParser.Statement()

	if len(sllListener.errors) > 0 {
		return nil, fmt.Errorf("sql syntax error: %s", strings.Join(sllListener.errors, "; "))
	}

	visitor := NewAstVisitor(idAllocator, input)
	node := visitor.Visit(parseTree)
	if node != nil {
		stmt = node.(Statement)
	}
	return stmt, nil
}
