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

package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xlab/treeprint"

	"github.com/lindb/lindb/pkg/ltoml"
)

// SeriesStats represents the stats for series.
type SeriesStats struct {
	NumOfSeries uint64 `json:"numOfSeries"`
}

// OperatorStats represents the stats of operator.
type OperatorStats struct {
	Identifier string      `json:"identifier"`
	Start      int64       `json:"start"`
	End        int64       `json:"end"`
	Cost       int64       `json:"cost"`
	Stats      interface{} `json:"stats,omitempty"`
	ErrMsg     string      `json:"errMsg,omitempty"`
}

// StageStats represents the stats of stage.
type StageStats struct {
	Identifier string           `json:"identifier"`
	Start      int64            `json:"start"`
	End        int64            `json:"end"`
	Cost       int64            `json:"cost"`
	State      string           `json:"state"`
	ErrMsg     string           `json:"errMsg"`
	Async      bool             `json:"async"`
	Operators  []*OperatorStats `json:"operators,omitempty"`

	Children []*StageStats `json:"children"`
}

// NodeStats represents query stats of node.
type NodeStats struct {
	Node       string        `json:"node"`
	WaitCost   int64         `json:"waitCost,omitempty"` // wait intermediate or leaf response duration
	WaitStart  int64         `json:"waitStart,omitempty"`
	WaitEnd    int64         `json:"waitEnd,omitempty"`
	NetPayload int64         `json:"netPayload,omitempty"`
	TotalCost  int64         `json:"totalCost"`
	Start      int64         `json:"start"`
	End        int64         `json:"end"`
	Stages     []*StageStats `json:"stages,omitempty"`

	Children []*NodeStats `json:"children,omitempty"`
}

// Stats represents the time stats
type Stats struct {
	TotalCost int64 `json:"totalCost"`
	Min       int64 `json:"min"`
	Max       int64 `json:"max"`
	Count     int   `json:"count"`
	Series    int   `json:"series,omitempty"`
}

// ToTable returns the result of query as table if it has value, else return empty string.
func (s *NodeStats) ToTable() (rows int, tableStr string) {
	// 1. set headers
	headers := table.Row{}
	headers = append(headers, "Query Plan")
	result := NewTableFormatter()
	result.AppendHeader(headers)
	// fix calc row width
	treeprint.EdgeTypeLink = "!"
	treeprint.EdgeTypeMid = "^^"
	treeprint.EdgeTypeEnd = "~~"
	treeprint.IndentSize = 2
	tree := treeprint.NewWithRoot(nodeTitle(s))

	for _, stage := range s.Stages {
		stageToTable(tree, stage)
	}
	for _, child := range s.Children {
		nodeToTable(tree, child)
	}
	str := strings.TrimSuffix(tree.String(), "\n")
	result.AppendRow(table.Row{str})
	rs := result.Render()
	rs = strings.ReplaceAll(rs, "!", "│")
	rs = strings.ReplaceAll(rs, "^^", "├─")
	rs = strings.ReplaceAll(rs, "~~", "└─")
	return 1, rs
}

// nodeToTable returns node info.
func nodeToTable(tree treeprint.Tree, node *NodeStats) {
	sub := tree.AddBranch(nodeTitle(node))
	for _, stage := range node.Stages {
		stageToTable(sub, stage)
	}
	for _, child := range node.Children {
		nodeToTable(sub, child)
	}
}

// nodeTitle returns the title of node.
func nodeTitle(node *NodeStats) string {
	costs := []string{fmt.Sprintf("Cost: %s", time.Duration(node.TotalCost))}
	if node.WaitStart > 0 {
		costs = append(costs, fmt.Sprintf("Wait: %s", time.Duration(node.WaitCost)))
	}
	if node.NetPayload > 0 {
		costs = append(costs, fmt.Sprintf("Network: %s", ltoml.Size(node.NetPayload)))
	}
	return fmt.Sprintf("%s: [%s]",
		node.Node, strings.Join(costs, ", "),
	)
}

// stageToTable builds stage stats as table.
func stageToTable(tree treeprint.Tree, stage *StageStats) {
	stageNode := tree.AddBranch(fmt.Sprintf("Stage(%s), [Cost:%s]", stage.Identifier, time.Duration(stage.Cost)))
	for _, op := range stage.Operators {
		stageNode.AddNode(fmt.Sprintf("Operator(%s), [Cost:%s]", op.Identifier, time.Duration(op.Cost)))
	}
	for _, child := range stage.Children {
		stageToTable(stageNode, child)
	}
}
