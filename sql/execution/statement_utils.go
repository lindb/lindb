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

package execution

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/tree"
)

var statementTypes = make(map[reflect.Type]models.StatementType)

func init() {
	// DDL
	statementTypes[reflect.TypeOf(&tree.CreateDatabase{})] = models.DataDefinition
	statementTypes[reflect.TypeOf(&tree.DropDatabase{})] = models.DataDefinition
	// DML
	statementTypes[reflect.TypeOf(&tree.Query{})] = models.Select
	// Explain
	statementTypes[reflect.TypeOf(&tree.Explain{})] = models.Select
	// Show replication/memory databases/namespaces/table names/columns
	statementTypes[reflect.TypeOf(&tree.Show{})] = models.Select
}

func GetStatementType(statement tree.Statement) models.StatementType {
	if statementType, ok := statementTypes[reflect.TypeOf(statement)]; ok {
		return statementType
	}
	panic(fmt.Sprintf("unknown statement type for '%s'", reflect.TypeOf(statement)))
}
