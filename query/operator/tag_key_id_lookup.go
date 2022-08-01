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

// tagKeyIDLookup represents tag key id lookup operator.
type tagKeyIDLookup struct {
	ctx *context.LeafMetadataContext
}

// NewTagKeyIDLookup create a tagKeyIDLookup instance.
func NewTagKeyIDLookup(ctx *context.LeafMetadataContext) Operator {
	return &tagKeyIDLookup{
		ctx: ctx,
	}
}

// Execute finds tag key id by given namespace/metric/tag key.
func (op *tagKeyIDLookup) Execute() error {
	req := op.ctx.Request
	tagKeyID, err := op.ctx.Database.Metadata().MetadataDatabase().GetTagKeyID(req.Namespace, req.MetricName, req.TagKey)
	if err != nil {
		return err
	}
	op.ctx.TagKeyID = tagKeyID
	return nil
}

// Identifier returns identifier value of tag key lookup operator.
func (op *tagKeyIDLookup) Identifier() string {
	return "Tag Key Lookup"
}
