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
	"strconv"
	"strings"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/tree"
)

type PlanNodeIDAllocator struct {
	next PlanNodeID
}

func NewPlanNodeIDAllocator() *PlanNodeIDAllocator {
	return &PlanNodeIDAllocator{}
}

func (a *PlanNodeIDAllocator) Next() PlanNodeID {
	a.next++
	return a.next
}

type SymbolAllocator struct {
	analyzerContext *analyzer.AnalyzerContext

	symbols map[string]struct{}
	mapping map[tree.NodeID]*Symbol
	next    int
}

func NewSymbolAllocator(analyzerContext *analyzer.AnalyzerContext) *SymbolAllocator {
	return &SymbolAllocator{
		analyzerContext: analyzerContext,
		symbols:         make(map[string]struct{}),
		mapping:         make(map[tree.NodeID]*Symbol),
	}
}

func (a *SymbolAllocator) FromExpression(expression tree.Expression, dataType types.DataType) *Symbol {
	if symbol, ok := a.mapping[expression.GetID()]; ok {
		return symbol
	}

	fmt.Printf("new symbol=%T\n", expression)
	nameHint := "expr"
	var hidden bool
	switch expr := expression.(type) {
	case *tree.Identifier:
		nameHint = expr.Value
	case *tree.SymbolReference:
		nameHint = expr.Name
		hidden = expr.Hidden
	case *tree.FunctionCall:
		if expr.RefField != nil {
			// FIXME: func call,not use ref field name
			nameHint = expr.RefField.Name
			dataType = expr.RefField.DataType
		} else {
			nameHint = string(expr.Name)
		}
	}
	symbol := a.NewSymbol(nameHint, dataType, hidden)
	a.mapping[expression.GetID()] = symbol
	return symbol
}

func (a *SymbolAllocator) FromSymbol(symbolHint *Symbol, dataType types.DataType, hidden bool) *Symbol {
	return a.NewSymbol(symbolHint.Name, dataType, hidden)
}

func (a *SymbolAllocator) NewSymbol(nameHint string, dataType types.DataType, hidden bool) *Symbol {
	nameHint = cleanNameHint(nameHint)

	// TODO: modify for?
	_, exist := a.symbols[nameHint]
	if exist {
		nameHint = fmt.Sprintf("%s_%d", nameHint, a.next)
		a.next++
	}
	a.symbols[nameHint] = struct{}{}
	fmt.Printf("...................................nameHint=%v, symbols=%v\n", nameHint, a.symbols)

	return &Symbol{
		Name:     nameHint,
		DataType: dataType,
		Hidden:   hidden,
	}
}

func cleanNameHint(nameHint string) string {
	index := strings.LastIndex(nameHint, "_")
	if index > 0 {
		tail := nameHint[index+1:]
		// only strip if tail is numeric or _ is the last character
		if _, err := strconv.Atoi(tail); err == nil || index == len(nameHint)-1 {
			nameHint = nameHint[:index]
		}
	}

	if nameHint == "" {
		nameHint = "col"
	}

	return nameHint
}
