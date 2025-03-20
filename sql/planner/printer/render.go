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

package printer

import (
	"regexp"
	"strings"

	"github.com/lindb/lindb/pkg/strutil"
)

type Render interface {
	Render(paln *PlanRepresentation) string
}

type TextRender struct {
	level int
}

func NewTextRender(level int) Render {
	return &TextRender{
		level: level,
	}
}

func (r *TextRender) Render(plan *PlanRepresentation) string {
	root := plan.getRoot()
	sb := &strings.Builder{}
	hasChildren := hasChildren(plan, root)

	return strings.TrimSuffix(r.writeTextOutput(sb, plan, strutil.NewIndent(r.level, hasChildren), root), "\n")
}

func (r *TextRender) writeTextOutput(
	output *strings.Builder, plan *PlanRepresentation,
	indent *strutil.Indent, node *NodeRepresentation,
) string {
	output.WriteString(indent.NodeIndent())
	output.WriteString(node.getName())
	var kvs []string
	for key, value := range node.descriptor {
		kvs = append(kvs, key+" = "+value)
	}
	output.WriteString("[" + strings.Join(kvs, ", ") + "]")
	output.WriteString("\n")

	output.WriteString(indentMultilineString("Layout: "+formatSymbols(node.outputs), indent.DetailIndent()))
	output.WriteString("\n")

	if len(node.details) > 0 {
		details := strings.Join(node.details, "\n")
		details = indentMultilineString(details, indent.DetailIndent())
		output.WriteString(details)
		if !strings.HasSuffix(details, "\n") {
			output.WriteString("\n")
		}
	}

	// process children
	childrenIDs := node.children
	for i, childID := range childrenIDs {
		child := plan.getNode(childID)
		if child != nil {
			r.writeTextOutput(output, plan,
				indent.ForChild(i == len(childrenIDs)-1, hasChildren(plan, child)), child)
		}
	}

	return output.String()
}

func indentMultilineString(str, indent string) string {
	m1 := regexp.MustCompile("(?m)^")
	return m1.ReplaceAllString(str, indent)
}

func hasChildren(plan *PlanRepresentation, node *NodeRepresentation) bool {
	for _, childID := range node.children {
		child := plan.getNode(childID)
		if child != nil {
			return true
		}
	}
	return false
}
