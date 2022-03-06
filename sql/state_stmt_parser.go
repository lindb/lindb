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

import "github.com/lindb/lindb/sql/stmt"

// stateStmtParse represents show state statement parser.
type stateStmtParse struct {
	stateType stmt.StateType
}

// newStateStmtParse creates a show state statement parser.
func newStateStmtParse(stateType stmt.StateType) *stateStmtParse {
	return &stateStmtParse{
		stateType: stateType,
	}
}

// build returns the state statement.
func (s *stateStmtParse) build() (stmt.Statement, error) {
	return &stmt.State{Type: s.stateType}, nil
}
