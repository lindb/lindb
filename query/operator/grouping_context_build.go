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

package operator

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
)

// groupingContextBuild represents grouping context build operator.
type groupingContextBuild struct {
	executeCtx *flow.ShardExecuteContext
	shard      tsdb.Shard
}

// NewGroupingContextBuild creates a groupingContextBuild instance.
func NewGroupingContextBuild(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &groupingContextBuild{
		executeCtx: executeCtx,
		shard:      shard,
	}
}

// Execute executes grouping context build based on series ids after tag filtering.
func (op *groupingContextBuild) Execute() error {
	if op.executeCtx.IsSeriesIDsEmpty() {
		return nil
	}
	return op.shard.IndexDB().GetGroupingContext(op.executeCtx)
}

// Identifier returns identifier string value of grouping context build operator.
func (op *groupingContextBuild) Identifier() string {
	return "Grouping Context Build"
}
