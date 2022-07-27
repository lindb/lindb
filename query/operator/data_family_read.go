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

// dataFamilyRead represents data family filtering operator based on series ids.
type dataFamilyRead struct {
	executeCtx *flow.ShardExecuteContext
	family     tsdb.DataFamily
}

// NewDataFamilyRead creates a dataFamilyRead instance.
func NewDataFamilyRead(executeCtx *flow.ShardExecuteContext, family tsdb.DataFamily) Operator {
	return &dataFamilyRead{
		executeCtx: executeCtx,
		family:     family,
	}
}

// Execute executes data family(file/memory) based on series ids, then add result set into time segment context.
func (op *dataFamilyRead) Execute() error {
	family := op.family
	resultSet, err := family.Filter(op.executeCtx)
	if err != nil {
		return err
	}
	for _, rs := range resultSet {
		op.executeCtx.TimeSegmentContext.AddFilterResultSet(family.Interval(), rs)
	}
	return nil
}
