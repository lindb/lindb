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

package rewrite

import (
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/tree"
)

type StatementRewrite struct {
	rewrites []interfaces.Rewrite
}

func NewStatementRewrite(rewrites []interfaces.Rewrite) interfaces.Rewrite {
	return &StatementRewrite{
		rewrites: rewrites,
	}
}

func (rw *StatementRewrite) Rewrite(statement tree.Statement) tree.Statement {
	for _, rewrite := range rw.rewrites {
		statement = rewrite.Rewrite(statement)
	}
	return statement
}
