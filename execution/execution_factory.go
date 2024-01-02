package execution

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/tree"
)

var (
	executionFactories = make(map[models.StatementType]ExecutionFactory)
)

func RegisterExecutionFactory(statementType models.StatementType, factory ExecutionFactory) {
	executionFactories[statementType] = factory
}

func GetExecutionFactory(statementType models.StatementType) ExecutionFactory {
	factory := executionFactories[statementType]
	return factory
}

type ExecutionFactory interface {
	CreateExecution(session *Session, statement *tree.PreparedStatement) Execution
}

type DataDefinitionExecutionFactory struct {
	deps *Deps
}

func NewDataDefinitionExecutionFactory(deps *Deps) ExecutionFactory {
	return &DataDefinitionExecutionFactory{
		deps: deps,
	}
}

func (f *DataDefinitionExecutionFactory) CreateExecution(session *Session, statement *tree.PreparedStatement) Execution {
	var task DataDefinitionTask
	switch sType := statement.Statement.(type) {
	case *tree.CreateDatabase:
		task = NewCreateDatabaseTask(f.deps, sType)
	}
	return NewDataDefinitionExecution(task)
}

type QueryExecutionFactory struct {
	deps *Deps
}

func NewQueryExecutionFactory(deps *Deps) ExecutionFactory {
	return &QueryExecutionFactory{
		deps: deps,
	}
}

func (f *QueryExecutionFactory) CreateExecution(session *Session, statement *tree.PreparedStatement) Execution {
	return NewQueryExecution(session, f.deps, statement)
}
