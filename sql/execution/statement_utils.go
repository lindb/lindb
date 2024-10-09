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
	// DML
	statementTypes[reflect.TypeOf(&tree.Query{})] = models.Select
	// Explain
	statementTypes[reflect.TypeOf(&tree.Explain{})] = models.Select
}

func GetStatementType(statement tree.Statement) models.StatementType {
	if statementType, ok := statementTypes[reflect.TypeOf(statement)]; ok {
		return statementType
	}
	panic(fmt.Sprintf("unknown statement type for '%s'", reflect.TypeOf(statement)))
}
