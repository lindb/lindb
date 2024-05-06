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
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/query/context"
)

// fieldSuggest represents field suggest operator.
type fieldSuggest struct {
	ctx *context.LeafMetadataContext
}

// NewFieldSuggest creates a fieldSuggest operator.
func NewFieldSuggest(ctx *context.LeafMetadataContext) Operator {
	return &fieldSuggest{
		ctx: ctx,
	}
}

// Execute returns all fields by given metric.
func (op *fieldSuggest) Execute() error {
	req := op.ctx.Request
	metricID, err := op.ctx.Database.MetaDB().GetMetricID(req.Namespace, req.MetricName)
	if err != nil {
		return err
	}
	schema, err := op.ctx.Database.MetaDB().GetSchema(metricID)
	if err != nil {
		return err
	}
	if schema != nil {
		var result []string
		result = append(result, string(encoding.JSONMarshal(schema.Fields)))
		op.ctx.ResultSet = result
	}
	return nil
}

// Identifier returns identifier string value of field suggest operator.
func (op *fieldSuggest) Identifier() string {
	return "Field Suggest"
}
