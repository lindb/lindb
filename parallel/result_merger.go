package parallel

import (
	"context"

	"github.com/lindb/lindb/aggregation"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./result_merger.go -destination=./result_merger_mock.go -package=parallel

// ResultMerger represents a merger which merges the task response and aggregates the result
type ResultMerger interface {
	// merge merges the task response and aggregates the result
	merge(resp *pb.TaskResponse)

	close()
}

type resultMerger struct {
	resultSet chan *series.TimeSeriesEvent

	groupAgg aggregation.GroupingAggregator

	events chan *pb.TaskResponse

	closed chan struct{}
	ctx    context.Context

	err error
}

// newResultMerger create a result merger
func newResultMerger(ctx context.Context, groupAgg aggregation.GroupingAggregator, resultSet chan *series.TimeSeriesEvent) ResultMerger {
	merger := &resultMerger{
		resultSet: resultSet,
		groupAgg:  groupAgg,
		events:    make(chan *pb.TaskResponse),
		closed:    make(chan struct{}),
		ctx:       ctx,
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
		resultSet := m.groupAgg.ResultSet()
		if len(resultSet) > 0 {
			m.resultSet <- &series.TimeSeriesEvent{
				SeriesList: resultSet,
			}
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
	tsList := &pb.TimeSeriesList{}
	err := tsList.Unmarshal(data)
	if err != nil {
		m.err = err
		return false
	}
	for _, ts := range tsList.TimeSeriesList {
		// if no field data, ignore this response
		if len(ts.Fields) == 0 {
			return true
		}
		m.groupAgg.Aggregate(series.NewGroupedIterator(ts.Tags, ts.Fields))
	}
	return true
}
