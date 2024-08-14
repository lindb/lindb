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

import "github.com/lindb/lindb/flow"

// groupingTagsLookup represents grouping tags lookup operator.
type groupingTagsLookup struct {
	executeCtx *flow.DataLoadContext
}

// NewGroupingTagsLookup creates a groupingTagsLookup instance.
func NewGroupingTagsLookup(executeCtx *flow.DataLoadContext) Operator {
	return &groupingTagsLookup{
		executeCtx: executeCtx,
	}
}

// Execute executes grouping tag value ids lookup, if it hasn't grouping tag key returns no grouping.
func (op *groupingTagsLookup) Execute() error {
	op.executeCtx.Grouping()
	if op.executeCtx.ShardExecuteCtx.GroupingContext != nil {
		// lookup grouping tags, grouped series: tags => series IDs(based on low series ids)
		op.executeCtx.ShardExecuteCtx.GroupingContext.BuildGroup(op.executeCtx)
	} else {
		op.executeCtx.PrepareAggregatorWithoutGrouping()
	}
	return nil
}

// Identifier returns identifier string value of grouping tags lookup operator.
func (op *groupingTagsLookup) Identifier() string {
	return "Grouping Tags Lookup"
}
