package query

import (
	"errors"
	"fmt"
	"sort"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
)

var (
	errEmptySelectList = errors.New("select item list is empty")
)

// storageExecutePlan represents a storage level execute plan for data search,
// such as plan down sampling and aggregation specification.
type storageExecutePlan struct {
	query    *stmt.Query
	idGetter diskdb.IDGetter

	fieldIDs []uint16

	metricID       uint32
	fields         map[uint16]aggregation.AggregatorSpec
	groupByTagKeys map[string]uint32

	err error
}

// newStorageExecutePlan creates a storage execute plan
func newStorageExecutePlan(index diskdb.IDGetter, query *stmt.Query) Plan {
	return &storageExecutePlan{
		idGetter:       index,
		query:          query,
		fields:         make(map[uint16]aggregation.AggregatorSpec),
		groupByTagKeys: make(map[string]uint32),
	}
}

// Plan plans the query language, generates the execute plan for storage query
func (p *storageExecutePlan) Plan() error {
	// metric name => id, like table name
	metricID, err := p.idGetter.GetMetricID(p.query.MetricName)
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
	p.fieldIDs = make([]uint16, len(p.fields))
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

// hasGroupBy returns if query has group by tag keys
func (p *storageExecutePlan) hasGroupBy() bool {
	return len(p.query.GroupBy) > 0
}

// groupBy parses group by tag keys
func (p *storageExecutePlan) groupBy() error {
	if len(p.query.GroupBy) == 0 {
		return nil
	}

	for _, tagKey := range p.query.GroupBy {
		tagKeyID, err := p.idGetter.GetTagKeyID(p.metricID, tagKey)
		if err != nil {
			return err
		}
		p.groupByTagKeys[tagKey] = tagKeyID
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
func (p *storageExecutePlan) getFieldIDs() []uint16 {
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
		fieldID, fieldType, err := p.idGetter.GetFieldID(p.metricID, e.Name)
		if err != nil {
			p.err = err
			return
		}
		var funcType function.FuncType
		// tests if has func with field
		if parentFunc == nil {
			// if not using field default down sampling func
			funcType = aggregation.DownSamplingFunc(fieldType)
			if funcType == function.Unknown {
				p.err = fmt.Errorf("cannot get default down sampling func for filed type[%s]", fieldType)
				return
			}
		} else {
			// using use input, and check func is supported
			if !aggregation.IsSupportFunc(fieldType, parentFunc.FuncType) {
				p.err = fmt.Errorf("field type[%s] not supprot function[%s]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
		}
		downSampling, exist := p.fields[fieldID]
		if !exist {
			downSampling = aggregation.NewAggregatorSpec(e.Name, fieldType)
			p.fields[fieldID] = downSampling
		}
		downSampling.AddFunctionType(funcType)
	}
}
