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
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// useStmtParser represents use statement parser.
type useStmtParser struct {
	name string
}

// newUseStmtParse creates a use statement parser.
func newUseStmtParse() *useStmtParser {
	return &useStmtParser{}
}

// visitName visits when production name expression is entered.
func (s *useStmtParser) visitName(ctx grammar.IIdentContext) {
	s.name = ctx.GetText()
}

// build returns the state statement.
func (s *useStmtParser) build() (stmt.Statement, error) {
	return &stmt.Use{Name: s.name}, nil
}
