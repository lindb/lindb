package parallel

import (
	"context"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./result_merger.go -destination=./result_merger_mock.go -package=parallel

// ResultMerger represents a merger which merges the task response and aggregates the result
type ResultMerger interface {
	// merge merges the task response and aggregates the result
	merge(resp *pb.TaskResponse)

	close()
}

type seriesAgg struct {
	tags       map[string]string
	aggregator aggregation.SegmentAggregator
}

type resultMerger struct {
	query     *stmt.Query
	resultSet chan *series.TimeSeriesEvent

	aggregates map[string]seriesAgg
	events     chan *pb.TaskResponse

	closed chan struct{}
	ctx    context.Context

	err error
}

// newResultMerger create a result merger
func newResultMerger(ctx context.Context, query *stmt.Query, resultSet chan *series.TimeSeriesEvent) ResultMerger {
	merger := &resultMerger{
		query:      query,
		resultSet:  resultSet,
		aggregates: make(map[string]seriesAgg),
		events:     make(chan *pb.TaskResponse),
		closed:     make(chan struct{}),
		ctx:        ctx,
	}
	go func() {
		defer close(merger.closed)
		merger.process()
	}()
	return merger
}

// merge merges and aggregates the result
func (m *resultMerger) merge(resp *pb.TaskResponse) {
	m.events <- resp
}

func (m *resultMerger) close() {
	close(m.events)
	// waiting process completed
	<-m.closed
	// send result set
	if m.err != nil {
		m.resultSet <- &series.TimeSeriesEvent{Err: m.err}
	} else {
		// send all series data
		for _, agg := range m.aggregates {
			m.resultSet <- &series.TimeSeriesEvent{Series: agg.aggregator.Iterator(agg.tags)}
		}
	}
}

func (m *resultMerger) process() {
	for {
		select {
		case event, ok := <-m.events:
			if !ok {
				return
			}
			// if handle event fail, return
			if !m.handleEvent(event) {
				return
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *resultMerger) handleEvent(resp *pb.TaskResponse) bool {
	data := resp.Payload
	ts := &pb.TimeSeries{}
	err := ts.Unmarshal(data)
	if err != nil {
		m.err = err
		return false
	}
	// if no field data, ignore this response
	if len(ts.Fields) == 0 {
		return true
	}
	// 1. prepare series tags
	tagsStr := constants.EmptyGroupTagsStr
	tags := constants.EmptyGroupTags
	if m.query.HasGroupBy() {
		tags := ts.Tags
		// if time series hasn't tags of response, ignore this response
		if len(tags) == 0 {
			return true
		}
		tagsStr = models.TagsAsString(tags)
	}
	// 2. get series aggregator
	var agg seriesAgg
	ok := false
	agg, ok = m.aggregates[tagsStr]
	if !ok {
		agg = seriesAgg{
			tags:       tags,
			aggregator: aggregation.NewSeriesSegmentAggregator(m.query.Interval, &m.query.TimeRange),
		}
		m.aggregates[tagsStr] = agg
	}
	// 3. aggregate field data
	for name, data := range ts.Fields {
		it := series.NewFieldIterator(name, data)
		agg.aggregator.Aggregate(it)
	}
	return true
}
