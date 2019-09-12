package query

import (
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./aggregate_worker.go -destination=./aggregate_worker_mock.go -package=query

// aggregateWorker represents
type aggregateWorker interface {
	// emit emits the field event of spec series
	emit(it *series.FieldEvent)
	// sendResult sends the current series aggregate result
	sendResult(tags map[string]string)
	// close closes the aggregate worker(chan)
	close()
}

type aggWorker struct {
	queryInterval  int64
	queryTimeRange *timeutil.TimeRange
	aggSpecs       map[string]*aggregation.AggregatorSpec

	aggregates map[int64]aggregation.SegmentAggregator

	itEvents chan *series.FieldEvent
	resultCh chan *series.TimeSeriesEvent

	err  error
	wait sync.WaitGroup
}

// createAggWorker creates series level aggregate worker
func createAggWorker(queryInterval int64, queryTimeRange *timeutil.TimeRange,
	appSpecs map[string]*aggregation.AggregatorSpec, resultCh chan *series.TimeSeriesEvent) aggregateWorker {
	worker := &aggWorker{
		queryTimeRange: queryTimeRange,
		queryInterval:  queryInterval,
		aggregates:     make(map[int64]aggregation.SegmentAggregator),
		itEvents:       make(chan *series.FieldEvent, 1),
		aggSpecs:       appSpecs,
		resultCh:       resultCh,
	}
	go worker.process()
	return worker
}

// sendResult sends series aggregate result
func (aw *aggWorker) sendResult(tags map[string]string) {
	aw.wait.Wait()
	for _, agg := range aw.aggregates {
		aw.resultCh <- &series.TimeSeriesEvent{
			Series: agg.Iterator(tags),
			Err:    aw.err,
		}
	}
}

// emit emits event to chan
func (aw *aggWorker) emit(it *series.FieldEvent) {
	aw.wait.Add(1)
	aw.itEvents <- it
}

// close closes the chan
func (aw *aggWorker) close() {
	close(aw.itEvents)
}

// process consumes the event from chan, then processes the event
func (aw *aggWorker) process() {
	for event := range aw.itEvents {
		familyStartTime := event.FamilyStartTime
		var agg aggregation.SegmentAggregator
		ok := false
		agg, ok = aw.aggregates[familyStartTime]
		if !ok {
			agg = aggregation.NewSegmentAggregator(aw.queryInterval, event.Interval,
				aw.queryTimeRange, event.FamilyStartTime, aw.aggSpecs)
			aw.aggregates[familyStartTime] = agg
		}
		agg.Aggregate(event.FieldIt)
		aw.wait.Done()
	}
}
