package operator

import (
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ProjectionOperatorFactory struct {
	project      *plan.ProjectionNode
	sourceLayout []*plan.Symbol
}

func NewProjectionOperatorFactory(project *plan.ProjectionNode, sourceLayout []*plan.Symbol) OperatorFactory {
	return &ProjectionOperatorFactory{
		project:      project,
		sourceLayout: sourceLayout,
	}
}

// CreateOperator implements OperatorFactory.
func (fct *ProjectionOperatorFactory) CreateOperator() Operator {
	return NewProjectionOperator(fct.project, fct.sourceLayout)
}

type ProjectionOperator struct {
	project *plan.ProjectionNode
	page    *spi.Page // TODO: refact

	sourceLayout []*plan.Symbol
}

func NewProjectionOperator(project *plan.ProjectionNode, sourceLayout []*plan.Symbol) Operator {
	return &ProjectionOperator{project: project, sourceLayout: sourceLayout}
}

// AddInput implements Operator.
func (h *ProjectionOperator) AddInput(page *spi.Page) {
	h.page = page
}

// Finish implements Operator.
func (h *ProjectionOperator) Finish() {
}

// GetOutput implements Operator.
func (h *ProjectionOperator) GetOutput() *spi.Page {
	fmt.Println("projection operator...............................")
	it := h.page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		for _, assign := range h.project.Assignments {
			val, _ := expression.Rewrite(&expression.RewriteContext{
				SourceLayout: h.sourceLayout,
			}, assign.Expression).EvalInt(row)
			fmt.Printf("projection op result=%v\n", val)
		}
	}
	return h.page
}

// IsFinished implements Operator.
func (h *ProjectionOperator) IsFinished() bool {
	return true
}
