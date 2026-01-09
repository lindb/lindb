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

package plan

import (
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type Symbol struct {
	Name     string              `json:"name"`
	DataType types.DataType      `json:"datatype"`
	Hidden   bool                `json:"hidden,omitempty"`
	AggType  types.AggregateType `json:"agg_type,omitempty"`
}

func (s *Symbol) ToSymbolReference() *tree.SymbolReference {
	return &tree.SymbolReference{
		Name:     s.Name,
		DataType: s.DataType,
		Hidden:   s.Hidden,
		AggType:  s.AggType,
	}
}

func SymbolFrom(expression tree.Expression) *Symbol {
	if symbolRef, ok := expression.(*tree.SymbolReference); ok {
		return &Symbol{
			Name:     symbolRef.Name,
			DataType: symbolRef.DataType,
			Hidden:   symbolRef.Hidden,
			AggType:  symbolRef.AggType,
		}
	}
	panic(fmt.Sprintf("new symbol with unexpected expression: %s, type:%T",
		tree.FormatExpression(expression), expression))
}

func (s *Symbol) String() string {
	var h string
	var agg string
	if s.Hidden {
		h = "!" // mark symbol hidden
	}
	if s.AggType != types.ATUnknown {
		agg = fmt.Sprintf("@%s", s.AggType.String())
	}
	return fmt.Sprintf("%s%s:%s%s", h, s.Name, s.DataType, agg)
}
