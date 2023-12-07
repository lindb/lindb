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

import "github.com/lindb/lindb/query/context"

// namespaceSuggest represents namespace suggest operator.
type namespaceSuggest struct {
	ctx *context.LeafMetadataContext
}

// NewNamespaceSuggest creates a namespaceSuggest instance.
func NewNamespaceSuggest(ctx *context.LeafMetadataContext) Operator {
	return &namespaceSuggest{
		ctx: ctx,
	}
}

// Execute returns namespace list by given prefix.
func (op *namespaceSuggest) Execute() error {
	req := op.ctx.Request
	limit := op.ctx.Limit
	rs, err := op.ctx.Database.MetaDB().SuggestNamespace(req.Prefix, limit)
	if err != nil {
		return err
	}
	op.ctx.ResultSet = rs
	return nil
}

// Identifier returns identifier string value of namespace suggest operator.
func (op *namespaceSuggest) Identifier() string {
	return "Namespace Suggest"
}
