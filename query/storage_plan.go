package query

import (
	"fmt"
	"sort"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
)

// storageExecutePlan represents a storage level execute plan for data search,
// such as plan down sampling and aggregation specification.
type storageExecutePlan struct {
	query    *stmt.Query
	idGetter diskdb.IDGetter

	fields map[uint16]*aggregation.AggregatorSpec

	metricID uint32

	err error
}

// newStorageExecutePlan creates a storage execute plan
func newStorageExecutePlan(index diskdb.IDGetter, query *stmt.Query) Plan {
	return &storageExecutePlan{
		idGetter: index,
		query:    query,
		fields:   make(map[uint16]*aggregation.AggregatorSpec),
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

	if err := p.selectList(); err != nil {
		return err
	}

	if p.err != nil {
		return p.err
	}
	if len(p.fields) == 0 {
		return fmt.Errorf("field cannot be empty for select list")
	}

	return nil
}

// getFieldIDs returns sorted slice of field ids
func (p *storageExecutePlan) getFieldIDs() []uint16 {
	var result []uint16
	for fieldID := range p.fields {
		result = append(result, fieldID)
	}
	// sort field ids
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

// selectList plans the select list from down sampling aggregation specification
func (p *storageExecutePlan) selectList() error {
	selectItems := p.query.SelectItems
	if len(selectItems) == 0 {
		return fmt.Errorf("select item list is empty")
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
				p.err = fmt.Errorf("cannot get default down sampling func for filed type[%d]", fieldType)
				return
			}
		} else {
			// using use input, and check func is supported
			if !aggregation.IsSupportFunc(fieldType, parentFunc.FuncType) {
				//TODO format error msg
				p.err = fmt.Errorf("field type[%d] not supprot function[%d]", fieldType, parentFunc.FuncType)
				return
			}
			funcType = parentFunc.FuncType
		}
		downSampling, exist := p.fields[fieldID]
		if !exist {
			downSampling = aggregation.NewAggregatorSpec(fieldID, e.Name, fieldType)
			p.fields[fieldID] = downSampling
		}
		downSampling.AddFunctionType(funcType)
	}
}
