package query

import (
	"context"
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pool"
	"github.com/lindb/lindb/series"
)

var Pool = pool.NewPool(200, 10)

// scanWorker represents dispatch the event of scanner
type scanWorker struct {
	hasGroupBy bool
	metricID   uint32
	tagKeys    []string

	metaGetter series.MetaGetter
	groupAgg   aggregation.GroupingAggregator
	resultCh   chan *series.TimeSeriesEvent

	events chan series.ScanEvent
	wait   sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

// createScanWorker creates scan worker dispatcher event to aggregate worker
func createScanWorker(ctx context.Context, metricID uint32,
	groupByTagKeys []string, metaGetter series.MetaGetter, groupedAgg aggregation.GroupingAggregator,
	resultCh chan *series.TimeSeriesEvent) series.ScanWorker {
	c, cancel := context.WithCancel(ctx)
	worker := &scanWorker{
		resultCh:   resultCh,
		metricID:   metricID,
		tagKeys:    groupByTagKeys,
		hasGroupBy: len(groupByTagKeys) > 0,
		metaGetter: metaGetter,
		events:     make(chan series.ScanEvent),
		groupAgg:   groupedAgg,
		ctx:        c,
		cancel:     cancel,
	}

	//FIXME add goroutine pool or add timeout/if need goroutine?
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("process scan event panic", logger.Any("err", err), logger.Stack())
			}
			worker.cancel()
		}()
		worker.process()
	}()

	return worker
}

// Emit emits the field event of spec series id
func (s *scanWorker) Emit(event series.ScanEvent) {
	if event == nil {
		return
	}
	s.wait.Add(1)
	Pool.JobQueue <- func() {
		if event.Scan() {
			s.events <- event
		} else {
			s.wait.Done()
		}
	}
}

// Close waits handle finish, then closes the event chan
func (s *scanWorker) Close() {
	go func() {
		s.wait.Wait()
		s.cancel()
	}()

	s.waitComplete()
	// if no group by tag keys, need send result after scan completed
	resultSet := s.groupAgg.ResultSet()
	if len(resultSet) > 0 {
		s.resultCh <- &series.TimeSeriesEvent{
			SeriesList: resultSet,
		}
	}
}

// process consumes event from chan, then handles the event
func (s *scanWorker) process() {
	for {
		select {
		case event := <-s.events:
			resultSet := event.ResultSet()
			if resultSet != nil {
				agg, ok := resultSet.(aggregation.FieldAggregates)
				if ok {
					s.groupAgg.Aggregate(agg.ResultSet(nil))
				}
			}
			event.Release()
			s.wait.Done()
		case <-s.ctx.Done():
			return
		}
	}
}

// waitComplete waits worker process complete
func (s *scanWorker) waitComplete() {
	// if ctx isn't done need wait, else timeout do not wait
	<-s.ctx.Done()
}
