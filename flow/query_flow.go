package flow

import (
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/concurrent"
)

//go:generate mockgen -source=./query_flow.go -destination=./query_flow_mock.go -package=flow

// StorageQueryFlow represents the storage query engine execute flow
type StorageQueryFlow interface {
	// Prepare prepares the query flow, builds the flow execute context based on down sampling aggregator specs
	Prepare(downSamplingSpecs aggregation.AggregatorSpecs)
	// Filtering does the filtering task
	Filtering(task concurrent.Task)
	// Grouping does the grouping task
	Grouping(task concurrent.Task)
	// Scanner does the scan task
	Scanner(task concurrent.Task)
	// Reduce reduces the down sampling aggregator's result
	Reduce(tags string, agg aggregation.FieldAggregates)
	// ReduceTagValues reduces the group by tag values
	ReduceTagValues(tagKeyIndex int, tagValues map[uint32]string)
	// GetAggregator gets the down sampling filed aggregator
	GetAggregator() (agg aggregation.FieldAggregates)
	// Complete completes the query flow with error
	Complete(err error)
}
