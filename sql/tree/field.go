package tree

import (
	"strings"

	"github.com/lindb/lindb/spi/types"
)

type Field struct {
	RelationAlias *QualifiedName
	Name          string
	DataType      types.DataType
}

func (f *Field) MatchesPrefix(prefix *QualifiedName) bool {
	return prefix == nil || (f.RelationAlias != nil && f.RelationAlias.HasSuffix(prefix))
}

func (f *Field) CanResolve(name *QualifiedName) bool {
	if f.Name == "" {
		return false
	}
	// TODO: need to know whether the qualified name and the name of this field were quoted
	return f.MatchesPrefix(name.Prefix) && strings.EqualFold(f.Name, name.Suffix)
}

func (f *Field) String() string {
	if f.RelationAlias == nil {
		return f.Name
	}
	return f.RelationAlias.Name + "." + f.Name
}
