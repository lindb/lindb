package query

import (
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb"
)

// scanWorker represents dispatch the event of scanner
type scanWorker struct {
	hasGroupBy bool
	metricID   uint32
	tagKeys    []string

	metaGetter series.MetaGetter
	groupAgg   aggregation.GroupingAggregator

	executorPool *tsdb.ExecutorPool

	ctx     parallel.ExecuteContext
	pending atomic.Int32

	done atomic.Bool

	mutex sync.Mutex
}

// createScanWorker creates scan worker dispatcher event to aggregate worker
func createScanWorker(
	ctx parallel.ExecuteContext,
	metricID uint32,
	groupByTagKeys []string,
	metaGetter series.MetaGetter,
	groupedAgg aggregation.GroupingAggregator,
	executorPool *tsdb.ExecutorPool,
) series.ScanWorker {
	worker := &scanWorker{
		metricID:     metricID,
		executorPool: executorPool,
		tagKeys:      groupByTagKeys,
		hasGroupBy:   len(groupByTagKeys) > 0,
		metaGetter:   metaGetter,
		groupAgg:     groupedAgg,
		ctx:          ctx,
	}
	return worker
}

// Emit emits the field event of spec series id
func (s *scanWorker) Emit(event series.ScanEvent) {
	if event == nil {
		return
	}
	s.pending.Inc()
	s.executorPool.Scanners.Execute(func() {
		if event.Scan() {
			s.executorPool.Mergers.Execute(func() {
				defer s.complete()

				resultSet := event.ResultSet()
				if resultSet != nil {
					agg, ok := resultSet.(aggregation.FieldAggregates)
					if ok {
						s.mutex.Lock()
						s.groupAgg.Aggregate(agg.ResultSet(nil))
						s.mutex.Unlock()
					}
				}
				event.Release()
			})
		} else {
			s.complete()
		}
	})
}

// Close marks scan worker can be done
func (s *scanWorker) Close() {
	s.done.Store(true)
}

// complete completes the worker if all pending events is done
func (s *scanWorker) complete() {
	pending := s.pending.Dec()
	if pending == 0 && s.done.Load() {
		resultSet := s.groupAgg.ResultSet()
		if len(resultSet) > 0 {
			s.ctx.Emit(&series.TimeSeriesEvent{
				SeriesList: resultSet,
			})
		}
		// complete the scan task
		s.ctx.Complete(nil)
	}
}
