package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/planner/plan"
)

type RelationPlan struct {
	Root          plan.PlanNode
	Scope         *analyzer.Scope
	OutContext    *TranslationMap
	FieldMappings []*plan.Symbol
}

func (r *RelationPlan) getSymbol(fieldIdx int) *plan.Symbol {
	if fieldIdx < 0 || fieldIdx >= len(r.FieldMappings) {
		panic(fmt.Sprintf("no field->symbol mapping for field %d", fieldIdx))
	}
	return r.FieldMappings[fieldIdx]
}
