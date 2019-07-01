package query

import (
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/query/aggregation"
)

type tsdbPlan struct {
	query models.Query
	// err error
}

func NewTSDBPlan(query models.Query) Plan {
	return &tsdbPlan{
		query: query,
	}
}

func (p *tsdbPlan) Plan() *aggregation.AggregatorStreamSpec {
	return nil
}
