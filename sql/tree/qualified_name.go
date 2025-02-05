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

package tree

import (
	"strings"
)

type QualifiedName struct {
	Prefix        *QualifiedName
	Name          string
	Suffix        string
	OriginalParts []*Identifier
	Parts         []string
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
	return start >= 0 // FIXME:
}
