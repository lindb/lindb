package plan

import (
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type Symbol struct {
	Name string `json:"name"`
	// FIXME: remove suffix?
	Suffix   string         `json:"suffix"`
	DataType types.DataType `json:"Datatype"`
}

func (s *Symbol) ToSymbolReference() *tree.SymbolReference {
	return &tree.SymbolReference{
		Name:     s.Name,
		DataType: s.DataType,
	}
}

func SymbolFrom(expression tree.Expression) *Symbol {
	if symbolRef, ok := expression.(*tree.SymbolReference); ok {
		return &Symbol{Name: symbolRef.Name, DataType: symbolRef.DataType}
	}
	panic(fmt.Sprintf("new symbol with unexpected expression: %T", expression))
}

func (s *Symbol) String() string {
	if s.Suffix != "" {
		return fmt.Sprintf("%s$%s:%s", s.Name, s.Suffix, s.DataType)
	}
	return fmt.Sprintf("%s:%s", s.Name, s.DataType)
}
