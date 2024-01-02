package tree

import (
	"strings"
)

type QualifiedName struct {
	OriginalParts []*Identifier
	Parts         []string
	Name          string
	Prefix        *QualifiedName
	Suffix        string
}

func NewQualifiedName(partsIdent []*Identifier) *QualifiedName {
	var parts []string
	for _, ident := range partsIdent {
		parts = append(parts, ident.Value)
	}
	var suffix string
	var prefix *QualifiedName
	if len(partsIdent) == 1 {
		suffix = parts[0]
	} else {
		prefix = NewQualifiedName(partsIdent[0 : len(partsIdent)-1])
		suffix = parts[len(parts)-1] // last
	}
	return &QualifiedName{
		OriginalParts: partsIdent,
		Parts:         parts,
		Name:          strings.Join(parts, "."),
		Prefix:        prefix,
		Suffix:        suffix,
	}
}

func (qn *QualifiedName) HasSuffix(suffix *QualifiedName) bool {
	if len(qn.Parts) < len(suffix.Parts) {
		return false
	}
	start := len(qn.Parts) - len(suffix.Parts)
	return start >= 0 //FIXME:
}
