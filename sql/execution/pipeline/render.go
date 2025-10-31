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

package pipeline

import (
	"strings"

	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/execution/operator"
)

func renderText(root operator.Operator) string {
	sb := &strings.Builder{}
	return strings.TrimSuffix(writeTextOutput(sb, strutil.NewIndent(0, len(root.Children()) > 0), root), "\n")
}

func writeTextOutput(sb *strings.Builder, indent *strutil.Indent, node operator.Operator) string {
	sb.WriteString(indent.NodeIndent())
	sb.WriteString(node.String())
	sb.WriteString("\n")
	children := node.Children()
	for i, child := range children {
		writeTextOutput(sb, indent.ForChild(i == len(children)-1, len(child.Children()) > 0), child)
	}
	return sb.String()
}
