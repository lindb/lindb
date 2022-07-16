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
	"github.com/lindb/lindb/query/context"
)

// leafReduce represents aggregate down sampling result set operator.
type leafReduce struct {
	leafExecuteCtx *context.LeafExecuteContext
	executeCtx     *flow.DataLoadContext
}

// NewLeafReduce creates a leafReduce instance.
func NewLeafReduce(leafExecuteCtx *context.LeafExecuteContext, executeCtx *flow.DataLoadContext) Operator {
	return &leafReduce{
		leafExecuteCtx: leafExecuteCtx,
		executeCtx:     executeCtx,
	}
}

// Execute executes aggregate down sampling result set after all down sampling operators completed.
func (op *leafReduce) Execute() error {
	if op.executeCtx.PendingDataLoadTasks.Load() == 0 {
		// after load, need to reduce the aggregator's result to query flow.
		op.executeCtx.Reduce(op.leafExecuteCtx.ReduceCtx.Reduce)
	}
	return nil
}
