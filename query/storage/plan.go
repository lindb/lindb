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

package storagequery

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

var (
	errEmptySelectList = errors.New("select item list is empty")
)

// storageExecutePlan represents a storage level execute plan for data search,
// such as plan down sampling and aggregation specification.
type storageExecutePlan struct {
	ctx *executeContext

	fields map[field.ID]*aggregation.Aggregator

	err error
}

// newStorageExecutePlan creates a storage execute plan
func newStorageExecutePlan(ctx *executeContext) *storageExecutePlan {
	return &storageExecutePlan{
		ctx:    ctx,
		fields: make(map[field.ID]*aggregation.Aggregator),
	}
}

// Plan plans the query language, generates the execute plan for storage query
func (p *storageExecutePlan) Plan() error {
	// metric name => id, like table name
	query := p.ctx.storageExecuteCtx.Query
	metricID, err := p.ctx.getMetadata().MetadataDatabase().GetMetricID(query.Namespace, query.MetricName)
	if err != nil {
		return err
	}

	p.ctx.storageExecuteCtx.MetricID = metricID

	if err := p.groupBy(); err != nil {
		return err
	}
	if err := p.selectList(); err != nil {
		return err
	}
	if p.err != nil {
		return p.err
	}

	p.buildField()

	return nil
}

// groupBy parses group by tag keys
func (p *storageExecutePlan) groupBy() error {
	groupBy := p.ctx.storageExecuteCtx.Query.GroupBy
	lengthOfGroupByTagKeys := len(groupBy)
	if lengthOfGroupByTagKeys == 0 {
		return nil
	}
	p.ctx.storageExecuteCtx.GroupByTags = make(tag.Metas, lengthOfGroupByTagKeys)
	p.ctx.storageExecuteCtx.GroupByTagKeyIDs = make([]tag.KeyID, lengthOfGroupByTagKeys)
	queryStmt := p.ctx.storageExecuteCtx.Query
	for idx, tagKey := range groupBy {
		tagKeyID, err := p.ctx.getMetadata().MetadataDatabase().GetTagKeyID(queryStmt.Namespace, queryStmt.MetricName, tagKey)
		if err != nil {
			return err
		}
		p.ctx.storageExecuteCtx.GroupByTags[idx] = tag.Meta{Key: tagKey, ID: tagKeyID}
		p.ctx.storageExecuteCtx.GroupByTagKeyIDs[idx] = tagKeyID
	}

	// need cache found grouping tag value id
	p.ctx.storageExecuteCtx.GroupingTagValueIDs = make([]*roaring.Bitmap, lengthOfGroupByTagKeys)

	return nil
}

// getDownSamplingAggSpecs returns the down sampling aggregate specs.
func (p *storageExecutePlan) buildField() {
	lengthOfFields := len(p.fields)
	p.ctx.storageExecuteCtx.Fields = make(field.Metas, lengthOfFields)

	idx := 0
	for fieldID := range p.fields {
		f := p.fields[fieldID]
		p.ctx.storageExecuteCtx.Fields[idx] = field.Meta{
			ID:   fieldID,
			Type: f.DownSampling.GetFieldType(),
			Name: f.DownSampling.FieldName(),
		}
		idx++
	}
	p.ctx.storageExecuteCtx.SortFields()

	// after sort filed, build aggregation spec
	p.ctx.storageExecuteCtx.DownSamplingSpecs = make(aggregation.AggregatorSpecs, lengthOfFields)
	p.ctx.storageExecuteCtx.AggregatorSpecs = make(aggregation.AggregatorSpecs, lengthOfFields)
	for fieldIdx, fieldMeta := range p.ctx.storageExecuteCtx.Fields {
		f := p.fields[fieldMeta.ID]
		p.ctx.storageExecuteCtx.DownSamplingSpecs[fieldIdx] = f.DownSampling
		p.ctx.storageExecuteCtx.AggregatorSpecs[fieldIdx] = f.Aggregator
	}
}

// selectList plans the select list from down sampling aggregation specification
func (p *storageExecutePlan) selectList() error {
	selectItems := p.ctx.storageExecuteCtx.Query.SelectItems
	if len(selectItems) == 0 {
		return errEmptySelectList
	}

	for _, selectItem := range selectItems {
		if p.err != nil {
			return p.err
		}
		p.field(nil, selectItem)
	}
	return nil
}

// field plans the field expr from select list
func (p *storageExecutePlan) field(parentFunc *stmt.CallExpr, expr stmt.Expr) {
	if p.err != nil {
		return
	}
	switch e := expr.(type) {
	case *stmt.SelectItem:
		p.field(nil, e.Expr)
	case *stmt.CallExpr:
		if e.FuncType == function.Quantile {
			p.planHistogramFields(e)
			return
		}
		for _, param := range e.Params {
			p.field(e, param)
		}
	case *stmt.ParenExpr:
		p.field(nil, e.Expr)
	case *stmt.BinaryExpr:
		p.field(nil, e.Left)
		p.field(nil, e.Right)
	case *stmt.FieldExpr:
		queryStmt := p.ctx.storageExecuteCtx.Query
		fieldMeta, err := p.ctx.getMetadata().
			MetadataDatabase().GetField(queryStmt.Namespace, queryStmt.MetricName, field.Name(e.Name))
		if err != nil {
			p.err = err
			return
		}

		fieldType := fieldMeta.Type
		fieldID := fieldMeta.ID
		aggregator, exist := p.fields[fieldID]
		if !exist {
			aggregator = &aggregation.Aggregator{}
			aggregator.DownSampling = aggregation.NewAggregatorSpec(field.Name(e.Name), fieldType)
			aggregator.Aggregator = aggregation.NewAggregatorSpec(field.Name(e.Name), fieldType)
			p.fields[fieldID] = aggregator
		}

		var funcType function.FuncType
		// tests if it has func with field
		if parentFunc == nil {
			// if not using field default down sampling func
			funcType = fieldType.DownSamplingFunc()
			if funcType == function.Unknown {
				p.err = fmt.Errorf("cannot get default down sampling func for filed type[%s]", fieldType)
				return
			}
			aggregator.Aggregator.AddFunctionType(funcType)
		} else {
			// using input, and check func is supported
			if !fieldType.IsFuncSupported(parentFunc.FuncType) {
				p.err = fmt.Errorf("field type[%s] not support function[%s]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
			// TODO ignore down sampling func?
			aggregator.Aggregator.AddFunctionType(parentFunc.FuncType)
		}
		aggregator.DownSampling.AddFunctionType(funcType)
	}
}

func (p *storageExecutePlan) planHistogramFields(e *stmt.CallExpr) {
	if len(e.Params) != 1 {
		p.err = fmt.Errorf("qunantile params more than one")
		return
	}
	if v, err := strconv.ParseFloat(e.Params[0].Rewrite(), 64); err != nil {
		p.err = fmt.Errorf("quantile param: %s is not float", e.Params[0].Rewrite())
		return
	} else if v <= 0 || v >= 1 {
		p.err = fmt.Errorf("quantile param: %f is illegal", v)
		return
	}
	queryStmt := p.ctx.storageExecuteCtx.Query
	fieldMetas, err := p.ctx.getMetadata().MetadataDatabase().GetAllHistogramFields(queryStmt.Namespace, queryStmt.MetricName)
	if err != nil {
		p.err = err
		return
	}
	for _, fieldMeta := range fieldMetas {
		aggregator, exist := p.fields[fieldMeta.ID]
		if !exist {
			aggregator = &aggregation.Aggregator{}
			aggregator.DownSampling = aggregation.NewAggregatorSpec(fieldMeta.Name, fieldMeta.Type)
			aggregator.Aggregator = aggregation.NewAggregatorSpec(fieldMeta.Name, fieldMeta.Type)
			p.fields[fieldMeta.ID] = aggregator
		}
		aggregator.Aggregator.AddFunctionType(function.Sum)
		aggregator.DownSampling.AddFunctionType(function.Sum)
	}
}
