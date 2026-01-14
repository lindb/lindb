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

package template

import (
	"html/template"
	"slices"
	"strings"
	"text/template/parse"
)

// Parse parses the given template string and returns the template if it contains action nodes.
func Parse(name, tpl string) (*template.Template, error) {
	// check if contains template delimiters
	if !strings.Contains(tpl, "{{") {
		return nil, nil
	}

	t, err := template.New(name).Parse(tpl)
	if err != nil {
		return nil, nil
	}
	// check if syntax tree contains non-text nodes
	if t.Tree == nil || t.Tree.Root == nil {
		return nil, nil
	}

	// iterate over the syntax tree to determine if there are non-text nodes
	if hasNonTextNode(t.Tree.Root) {
		return t, nil
	}
	return nil, nil
}

// hashNonTextNode checks if the parse tree contains any non-text nodes.
func hasNonTextNode(node parse.Node) bool {
	switch n := node.(type) {
	case *parse.ListNode:
		if slices.ContainsFunc(n.Nodes, hasNonTextNode) {
			return true
		}
	case *parse.TextNode:
		// plain text, not considered an "action"
		return false

	// has action nodes
	case *parse.ActionNode, // {{.}}、{{call ...}} etc.
		*parse.RangeNode,    // {{range ...}} ... {{end}}
		*parse.IfNode,       // {{if ...}} ... {{end}}
		*parse.WithNode,     // {{with ...}} ... {{end}}
		*parse.TemplateNode, // {{template "name" .}}
		*parse.VariableNode, // $x := ...
		*parse.CommandNode,
		*parse.PipeNode: // pipe
		return true
	// other node type (not expected in typical templates)
	default:
		return false
	}

	return false
}
