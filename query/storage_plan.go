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

	fieldIDs []field.ID

	metricID    uint32
	fields      map[field.ID]aggregation.AggregatorSpec
	groupByTags []tag.Meta

	err error
}

// newStorageExecutePlan creates a storage execute plan
func newStorageExecutePlan(namespace string, metadata metadb.Metadata, query *stmt.Query) Plan {
	return &storageExecutePlan{
		namespace: namespace,
		metadata:  metadata,
		query:     query,
		fields:    make(map[field.ID]aggregation.AggregatorSpec),
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
	p.fieldIDs = make([]field.ID, len(p.fields))
	idx := 0
	for fieldID := range p.fields {
		p.fieldIDs[idx] = fieldID
		idx++
	}
	// sort field ids
	sort.Slice(p.fieldIDs, func(i, j int) bool {
		return p.fieldIDs[i] < p.fieldIDs[j]
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

// getDownSamplingAggSpecs returns the down sampling aggregate specs
func (p *storageExecutePlan) getDownSamplingAggSpecs() aggregation.AggregatorSpecs {
	result := make(aggregation.AggregatorSpecs, len(p.fieldIDs))
	for idx, fieldID := range p.fieldIDs {
		result[idx] = p.fields[fieldID]
	}
	return result

}

// getFieldIDs returns sorted slice of field ids
func (p *storageExecutePlan) getFieldIDs() []field.ID {
	return p.fieldIDs
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
		var funcType function.FuncType
		// tests if has func with field
		if parentFunc == nil {
			// if not using field default down sampling func
			funcType = fieldType.DownSamplingFunc()
			if funcType == function.Unknown {
				p.err = fmt.Errorf("cannot get default down sampling func for filed type[%s]", fieldType)
				return
			}
		} else {
			// using use input, and check func is supported
			if !fieldType.IsFuncSupported(parentFunc.FuncType) {
				p.err = fmt.Errorf("field type[%s] not supprot function[%s]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
		}
		downSampling, exist := p.fields[fieldID]
		if !exist {
			downSampling = aggregation.NewDownSamplingSpec(field.Name(e.Name), fieldType)
			p.fields[fieldID] = downSampling
		}
		downSampling.AddFunctionType(funcType)
	}
}
