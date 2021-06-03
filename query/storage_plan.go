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

package query

import (
	"errors"
	"fmt"
	"sort"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

var (
	errEmptySelectList = errors.New("select item list is empty")
)

// storageExecutePlan represents a storage level execute plan for data search,
// such as plan down sampling and aggregation specification.
type storageExecutePlan struct {
	namespace string
	query     *stmt.Query
	metadata  metadb.Metadata

	fieldMetas field.Metas

	metricID    uint32
	fields      map[field.ID]*aggregation.Aggregator
	groupByTags []tag.Meta

	err error
}

// newStorageExecutePlan creates a storage execute plan
func newStorageExecutePlan(namespace string, metadata metadb.Metadata, query *stmt.Query) Plan {
	return &storageExecutePlan{
		namespace: namespace,
		metadata:  metadata,
		query:     query,
		fields:    make(map[field.ID]*aggregation.Aggregator),
	}
}

// Plan plans the query language, generates the execute plan for storage query
func (p *storageExecutePlan) Plan() error {
	// metric name => id, like table name
	metricID, err := p.metadata.MetadataDatabase().GetMetricID(p.namespace, p.query.MetricName)
	if err != nil {
		return err
	}
	p.metricID = metricID
	if err := p.groupBy(); err != nil {
		return err
	}
	if err := p.selectList(); err != nil {
		return err
	}
	if p.err != nil {
		return p.err
	}
	p.fieldMetas = make(field.Metas, len(p.fields))
	idx := 0
	for fieldID := range p.fields {
		f := p.fields[fieldID]
		p.fieldMetas[idx] = field.Meta{
			ID:   fieldID,
			Type: f.DownSampling.GetFieldType(),
			Name: f.DownSampling.FieldName(),
		}
		idx++
	}
	// sort field ids
	sort.Slice(p.fieldMetas, func(i, j int) bool {
		return p.fieldMetas[i].ID < p.fieldMetas[j].ID
	})

	return nil
}

// groupByKeyIDs returns group by tag key ids
func (p *storageExecutePlan) groupByKeyIDs() []tag.Meta {
	return p.groupByTags
}

// groupBy parses group by tag keys
func (p *storageExecutePlan) groupBy() error {
	groupByTags := len(p.query.GroupBy)
	if groupByTags == 0 {
		return nil
	}
	p.groupByTags = make([]tag.Meta, groupByTags)

	for idx, tagKey := range p.query.GroupBy {
		tagKeyID, err := p.metadata.MetadataDatabase().GetTagKeyID(p.namespace, p.query.MetricName, tagKey)
		if err != nil {
			return err
		}
		p.groupByTags[idx] = tag.Meta{
			Key: tagKey,
			ID:  tagKeyID,
		}
	}
	return nil
}

// getDownSamplingAggSpecs returns the down sampling aggregate specs.
func (p *storageExecutePlan) getDownSamplingAggSpecs() aggregation.AggregatorSpecs {
	result := make(aggregation.AggregatorSpecs, len(p.fieldMetas))
	for idx, f := range p.fieldMetas {
		result[idx] = p.fields[f.ID].DownSampling
	}
	return result
}

// getAggregatorSpecs returns aggregator specs for group by.
func (p *storageExecutePlan) getAggregatorSpecs() aggregation.AggregatorSpecs {
	result := make(aggregation.AggregatorSpecs, len(p.fieldMetas))
	for idx, f := range p.fieldMetas {
		result[idx] = p.fields[f.ID].Aggregator
	}
	return result
}

// getFields returns sorted of field.Metas.
func (p *storageExecutePlan) getFields() field.Metas {
	return p.fieldMetas
}

// selectList plans the select list from down sampling aggregation specification
func (p *storageExecutePlan) selectList() error {
	selectItems := p.query.SelectItems
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
		for _, param := range e.Params {
			p.field(e, param)
		}
	case *stmt.ParenExpr:
		p.field(nil, e.Expr)
	case *stmt.BinaryExpr:
		p.field(nil, e.Left)
		p.field(nil, e.Right)
	case *stmt.FieldExpr:
		fieldMeta, err := p.metadata.MetadataDatabase().GetField(p.namespace, p.query.MetricName, field.Name(e.Name))
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
		// tests if has func with field
		if parentFunc == nil {
			// if not using field default down sampling func
			funcType = fieldType.DownSamplingFunc()
			if funcType == function.Unknown {
				p.err = fmt.Errorf("cannot get default down sampling func for filed type[%s]", fieldType)
				return
			}
			aggregator.Aggregator.AddFunctionType(funcType)
		} else {
			// using use input, and check func is supported
			if !fieldType.IsFuncSupported(parentFunc.FuncType) {
				p.err = fmt.Errorf("field type[%s] not supprot function[%s]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
			//TODO ignore down sampling func?
			aggregator.Aggregator.AddFunctionType(parentFunc.FuncType)
		}
		aggregator.DownSampling.AddFunctionType(funcType)
	}
}
