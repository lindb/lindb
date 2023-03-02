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
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
)

// seriesLimit represents series limit operator.
type seriesLimit struct {
	executeCtx *flow.ShardExecuteContext
	shard      tsdb.Shard
}

// NewSeriesLimit creates a seriesLimit instance.
func NewSeriesLimit(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &seriesLimit{
		executeCtx: executeCtx,
		shard:      shard,
	}
}

// Execute executes series limit.
func (op *seriesLimit) Execute() error {
	numOfSeries := op.executeCtx.SeriesIDsAfterFiltering.GetCardinality()
	if numOfSeries == 0 {
		return nil
	}
	limit := op.shard.Database().GetLimits()
	if numOfSeries > uint64(limit.MaxSeriesPerQuery) {
		return constants.ErrTooManySeriesFound
	}
	return nil
}

// Identifier returns identifier value of series limit operator.
func (op *seriesLimit) Identifier() string {
	return "Series Limit"
}
