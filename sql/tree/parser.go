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
	parser := grammar.NewSQLParser(tokens)
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&MyErrorListener{})
	// first, try parsing with potentially faster SLL mode
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	// TODO: fail to LL mode
	parseTree := parser.Statement()

	visitor := NewAstVisitor(idAllocator)
	node := visitor.Visit(parseTree)
	if node != nil {
		stmt = node.(Statement)
	}
	return stmt, nil
}
