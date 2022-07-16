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
	"fmt"
	"strconv"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

// metadataLookup represents metadata lookup operator.
type metadataLookup struct {
	database   tsdb.Database
	metadata   metadb.MetadataDatabase
	executeCtx *flow.StorageExecuteContext

	fields map[field.ID]*aggregation.Aggregator

	err error
}

// NewMetadataLookup creates a metadataLookup instance.
func NewMetadataLookup(executeCtx *flow.StorageExecuteContext, database tsdb.Database) Operator {
	return &metadataLookup{
		database:   database,
		metadata:   database.Metadata().MetadataDatabase(),
		executeCtx: executeCtx,
		fields:     make(map[field.ID]*aggregation.Aggregator),
	}
}

// Execute executes metadata(metric/field/grouping keys) lookup.
func (op *metadataLookup) Execute() error {
	// metric name => id, like table name
	query := op.executeCtx.Query
	metricID, err := op.metadata.GetMetricID(query.Namespace, query.MetricName)
	if err != nil {
		return err
	}

	op.executeCtx.MetricID = metricID
	op.executeCtx.TagKeys = make(map[string]tag.KeyID)

	if err := op.groupBy(); err != nil {
		return err
	}
	if err := op.selectList(); err != nil {
		return err
	}

	op.buildField()
	return nil
}

// groupBy parses group by tag keys
func (op *metadataLookup) groupBy() error {
	groupBy := op.executeCtx.Query.GroupBy
	lengthOfGroupByTagKeys := len(groupBy)
	if lengthOfGroupByTagKeys == 0 {
		return nil
	}
	op.executeCtx.GroupByTags = make(tag.Metas, lengthOfGroupByTagKeys)
	op.executeCtx.GroupByTagKeyIDs = make([]tag.KeyID, lengthOfGroupByTagKeys)

	queryStmt := op.executeCtx.Query
	metadata := op.metadata
	for idx, tagKey := range groupBy {
		tagKeyID, err := metadata.GetTagKeyID(queryStmt.Namespace, queryStmt.MetricName, tagKey)
		if err != nil {
			return err
		}
		op.executeCtx.GroupByTags[idx] = tag.Meta{Key: tagKey, ID: tagKeyID}
		op.executeCtx.GroupByTagKeyIDs[idx] = tagKeyID
		// cache tag keys in context
		op.executeCtx.TagKeys[tagKey] = tagKeyID
	}

	// init grouping tag value collection, need cache found grouping tag value id
	op.executeCtx.GroupingTagValueIDs = make([]*roaring.Bitmap, lengthOfGroupByTagKeys)

	return nil
}

// getDownSamplingAggSpecs returns the down sampling aggregate specs.
func (op *metadataLookup) buildField() {
	lengthOfFields := len(op.fields)
	op.executeCtx.Fields = make(field.Metas, lengthOfFields)

	idx := 0
	for fieldID := range op.fields {
		f := op.fields[fieldID]
		op.executeCtx.Fields[idx] = field.Meta{
			ID:   fieldID,
			Type: f.DownSampling.GetFieldType(),
			Name: f.DownSampling.FieldName(),
		}
		idx++
	}
	// first sort field by field id
	op.executeCtx.SortFields()
	// after sort filed, build aggregation spec
	op.executeCtx.DownSamplingSpecs = make(aggregation.AggregatorSpecs, lengthOfFields)
	op.executeCtx.AggregatorSpecs = make(aggregation.AggregatorSpecs, lengthOfFields)
	for fieldIdx, fieldMeta := range op.executeCtx.Fields {
		f := op.fields[fieldMeta.ID]
		op.executeCtx.DownSamplingSpecs[fieldIdx] = f.DownSampling
		op.executeCtx.AggregatorSpecs[fieldIdx] = f.Aggregator
	}
}

// selectList plans the select list from down sampling aggregation specification
func (op *metadataLookup) selectList() error {
	selectItems := op.executeCtx.Query.SelectItems
	if len(selectItems) == 0 {
		return constants.ErrEmptySelectList
	}

	for _, selectItem := range selectItems {
		op.field(nil, selectItem)
		if op.err != nil {
			return op.err
		}
	}
	return nil
}

// field plans the field expr from select list
func (op *metadataLookup) field(parentFunc *stmt.CallExpr, expr stmt.Expr) {
	if op.err != nil {
		return
	}
	switch e := expr.(type) {
	case *stmt.SelectItem:
		op.field(nil, e.Expr)
	case *stmt.CallExpr:
		if e.FuncType == function.Quantile {
			op.planHistogramFields(e)
			return
		}
		for _, param := range e.Params {
			op.field(e, param)
		}
	case *stmt.ParenExpr:
		op.field(nil, e.Expr)
	case *stmt.BinaryExpr:
		op.field(nil, e.Left)
		op.field(nil, e.Right)
	case *stmt.FieldExpr:
		queryStmt := op.executeCtx.Query
		fieldMeta, err := op.metadata.GetField(queryStmt.Namespace, queryStmt.MetricName, field.Name(e.Name))
		if err != nil {
			op.err = err
			return
		}

		fieldType := fieldMeta.Type
		fieldID := fieldMeta.ID
		aggregator, exist := op.fields[fieldID]
		if !exist {
			aggregator = &aggregation.Aggregator{}
			aggregator.DownSampling = aggregation.NewAggregatorSpec(field.Name(e.Name), fieldType)
			aggregator.Aggregator = aggregation.NewAggregatorSpec(field.Name(e.Name), fieldType)
			op.fields[fieldID] = aggregator
		}

		var funcType function.FuncType
		// tests if it has func with field
		if parentFunc == nil {
			// if not using field default down sampling func
			funcType = fieldType.DownSamplingFunc()
			if funcType == function.Unknown {
				op.err = fmt.Errorf("cannot get default down sampling func for filed type[%s]", fieldType)
				return
			}
			aggregator.Aggregator.AddFunctionType(funcType)
		} else {
			// using input, and check func is supported
			if !fieldType.IsFuncSupported(parentFunc.FuncType) {
				op.err = fmt.Errorf("field type[%s] not support function[%s]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
			// TODO ignore down sampling func?
			aggregator.Aggregator.AddFunctionType(parentFunc.FuncType)
		}
		aggregator.DownSampling.AddFunctionType(funcType)
	}
}

func (op *metadataLookup) planHistogramFields(e *stmt.CallExpr) {
	if len(e.Params) != 1 {
		op.err = fmt.Errorf("qunantile params more than one")
		return
	}
	if v, err := strconv.ParseFloat(e.Params[0].Rewrite(), 64); err != nil {
		op.err = fmt.Errorf("quantile param: %s is not float", e.Params[0].Rewrite())
		return
	} else if v <= 0 || v >= 1 {
		op.err = fmt.Errorf("quantile param: %f is illegal", v)
		return
	}
	queryStmt := op.executeCtx.Query
	fieldMetas, err := op.metadata.GetAllHistogramFields(queryStmt.Namespace, queryStmt.MetricName)
	if err != nil {
		op.err = err
		return
	}
	for _, fieldMeta := range fieldMetas {
		aggregator, exist := op.fields[fieldMeta.ID]
		if !exist {
			aggregator = &aggregation.Aggregator{}
			aggregator.DownSampling = aggregation.NewAggregatorSpec(fieldMeta.Name, fieldMeta.Type)
			aggregator.Aggregator = aggregation.NewAggregatorSpec(fieldMeta.Name, fieldMeta.Type)
			op.fields[fieldMeta.ID] = aggregator
		}
		aggregator.Aggregator.AddFunctionType(function.Sum)
		aggregator.DownSampling.AddFunctionType(function.Sum)
	}
}
