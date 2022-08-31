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

// LeafNodeStats represents query stats in storage side
type LeafNodeStats struct {
	NetPayload int64         `json:"netPayload"`
	TotalCost  int64         `json:"totalCost"`
	Start      int64         `json:"start"`
	End        int64         `json:"end"`
	Stages     []*StageStats `json:"stages,omitempty"`
}

// Stats represents the time stats
type Stats struct {
	TotalCost int64 `json:"totalCost"`
	Min       int64 `json:"min"`
	Max       int64 `json:"max"`
	Count     int   `json:"count"`
	Series    int   `json:"series,omitempty"`
}

// QueryStats represents the query stats when need explain query flow stat
type QueryStats struct {
	Root         string                    `json:"root"`
	BrokerNodes  map[string]*QueryStats    `json:"brokerNodes,omitempty"`
	LeafNodes    map[string]*LeafNodeStats `json:"leafNodes,omitempty"`
	NetPayload   int64                     `json:"netPayload"`
	PlanCost     int64                     `json:"planCost,omitempty"`
	PlanStart    int64                     `json:"planStart,omitempty"`
	PlanEnd      int64                     `json:"planEnd,omitempty"`
	WaitCost     int64                     `json:"waitCost,omitempty"` // wait intermediate or leaf response duration
	WaitStart    int64                     `json:"waitStart,omitempty"`
	WaitEnd      int64                     `json:"waitEnd,omitempty"`
	ExpressCost  int64                     `json:"expressCost,omitempty"`
	ExpressStart int64                     `json:"expressStart,omitempty"`
	ExpressEnd   int64                     `json:"expressEnd,omitempty"`
	TotalCost    int64                     `json:"totalCost,omitempty"` // total query cost
	Start        int64                     `json:"start"`
	End          int64                     `json:"end"`
}

// NewQueryStats creates the query stats
func NewQueryStats() *QueryStats {
	return &QueryStats{
		BrokerNodes: make(map[string]*QueryStats),
		LeafNodes:   make(map[string]*LeafNodeStats),
	}
}

// MergeBrokerTaskStats merges intermediate task execution stats
func (s *QueryStats) MergeBrokerTaskStats(nodeID string, stats *QueryStats) {
	s.BrokerNodes[nodeID] = stats
}

// MergeLeafTaskStats merges leaf task execution stats
func (s *QueryStats) MergeLeafTaskStats(nodeID string, stats *LeafNodeStats) {
	s.LeafNodes[nodeID] = stats
}

// ToTable returns the result of query as table if it has value, else return empty string.
func (s *QueryStats) ToTable() (rows int, tableStr string) {
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
	tree := treeprint.NewWithRoot(fmt.Sprintf("Root(%s): [Cost:%s, Plan:%s, Wait:%s, Express: %s], Net Payload:%s",
		s.Root, time.Duration(s.TotalCost), time.Duration(s.PlanCost), time.Duration(s.WaitCost),
		time.Duration(s.ExpressCost), ltoml.Size(s.NetPayload),
	))

	for node, leaf := range s.LeafNodes {
		leafNode := tree.AddBranch(fmt.Sprintf("Leaf(%s): [Cost:%s], Net Payload:%s",
			node, time.Duration(leaf.TotalCost), ltoml.Size(leaf.NetPayload)))
		for _, stage := range leaf.Stages {
			s.stageToTable(leafNode, stage)
		}
		leafNode = tree.AddBranch(fmt.Sprintf("Leaf(%s): [Cost:%s], Net Payload:%s",
			node, time.Duration(leaf.TotalCost), ltoml.Size(leaf.NetPayload)))
		for _, stage := range leaf.Stages {
			s.stageToTable(leafNode, stage)
		}
	}
	str := strings.TrimSuffix(tree.String(), "\n")
	result.AppendRow(table.Row{str})
	rs := result.Render()
	rs = strings.ReplaceAll(rs, "!", "│")
	rs = strings.ReplaceAll(rs, "^^", "├─")
	rs = strings.ReplaceAll(rs, "~~", "└─")
	return 1, rs
}

// stageToTable builds stage stats as table.
func (s *QueryStats) stageToTable(tree treeprint.Tree, stage *StageStats) {
	stageNode := tree.AddBranch(fmt.Sprintf("Stage(%s), [Cost:%s]", stage.Identifier, time.Duration(stage.Cost)))
	for _, op := range stage.Operators {
		stageNode.AddNode(fmt.Sprintf("Operator(%s), [Cost:%s]", op.Identifier, time.Duration(op.Cost)))
	}
	for _, child := range stage.Children {
		s.stageToTable(stageNode, child)
	}
}
