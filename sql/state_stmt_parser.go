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
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// stateStmtParser represents show state statement parser.
type stateStmtParser struct {
	state *stmt.State
}

// newStateStmtParse creates a show state statement parser.
func newStateStmtParse(stateType stmt.StateType) *stateStmtParser {
	return &stateStmtParser{
		state: &stmt.State{Type: stateType},
	}
}

// visitDatabaseFilter visits database filter.
func (s *stateStmtParser) visitDatabaseFilter(ctx *grammar.DatabaseFilterContext) {
	s.state.Database = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitMetricList visits metric name list.
func (s *stateStmtParser) visitMetricList(ctx *grammar.MetricListContext) {
	names := ctx.AllIdent()
	for _, n := range names {
		s.state.MetricNames = append(s.state.MetricNames, strutil.GetStringValue(n.GetText()))
	}
}

// build returns the state statement.
func (s *stateStmtParser) build() (stmt.Statement, error) {
	return s.state, nil
}
