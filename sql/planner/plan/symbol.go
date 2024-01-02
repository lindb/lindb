package plan

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type Symbol struct {
	Name string `json:"name"`
}

func (s *Symbol) ToSymbolReference() *tree.SymbolReference {
	return &tree.SymbolReference{
		Name: s.Name,
	}
}

func SymbolFrom(expression tree.Expression) *Symbol {
	if symbolRef, ok := expression.(*tree.SymbolReference); ok {
		return &Symbol{Name: symbolRef.Name}
	}
	panic(fmt.Sprintf("new symbol with unexpected expression: %T", expression))
}
