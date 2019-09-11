package query

import (
	"context"
	"fmt"
	"sync"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
)

// scanWorker represents dispatch the event of scanner
type scanWorker struct {
	hasGroupBy bool
	metricID   uint32
	//seriesID  uint32
	tagValues []string
	tagKeys   []string

	metaGetter series.MetaGetter
	aggWorker  aggregateWorker

	events chan *series.FieldEvent
	wait   sync.WaitGroup
	done   chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
}

// createScanWorker creates scan worker dispatcher event to aggregate worker
func createScanWorker(ctx context.Context, metricID uint32,
	groupByTagKeys []string, metaGetter series.MetaGetter, aggWorker aggregateWorker) series.ScanWorker {
	c, cancel := context.WithCancel(ctx)
	worker := &scanWorker{
		metricID:   metricID,
		tagKeys:    groupByTagKeys,
		hasGroupBy: len(groupByTagKeys) > 0,
		metaGetter: metaGetter,
		events:     make(chan *series.FieldEvent, 1),
		done:       make(chan struct{}),
		aggWorker:  aggWorker,
		ctx:        c,
		cancel:     cancel,
	}

	//FIXME add goroutine pool or add timeout/if need goroutine?
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("process scan event panic", logger.Any("err", err), logger.Stack())
			}
			close(worker.done)
		}()
		worker.process()
	}()
	return worker
}

// Emit emits the field event of spec series id
func (s *scanWorker) Emit(event *series.FieldEvent) {
	if event == nil {
		return
	}
	s.wait.Add(1)
	s.events <- event
}

// Complete completes current series id scan
func (s *scanWorker) Complete(seriesID uint32) {
	s.Emit(&series.FieldEvent{
		SeriesID:  seriesID,
		Completed: true,
	})
}

// Close waits handle finish, then closes the event chan
func (s *scanWorker) Close() {
	go func() {
		s.wait.Wait()
		s.cancel()
	}()

	s.waitComplete()

	s.aggWorker.close()
	// if no group by tag keys, need send result after scan completed
	if !s.hasGroupBy {
		s.aggWorker.sendResult(nil)
	}
}

// process consumes event from chan, then handles the event
func (s *scanWorker) process() {
	for {
		select {
		case event := <-s.events:
			err := s.handleEvent(event)
			s.wait.Done()
			if err != nil {
				log.Error("handle event error", logger.Error(err))
				return
			}
		case <-s.ctx.Done():
			log.Warn("scan worker timeout")
			return
		}
	}
}

// handleEvent emits event to aggregate worker,
// if query has group by tag keys, searches tag values by group by tag keys
func (s *scanWorker) handleEvent(event *series.FieldEvent) error {
	// test current series if scan complete
	if event.Completed {
		// todo
		if s.hasGroupBy {
			s.aggWorker.sendResult(nil)
			s.tagValues = nil
		}
		return nil
	}
	s.aggWorker.emit(event)

	if s.hasGroupBy && len(s.tagValues) == 0 {
		if err := s.getGroupByTagValues(event.Version, event.SeriesID); err != nil {
			return err
		}
	}
	return nil
}

// getGroupByTagValues gets group by tag values
func (s *scanWorker) getGroupByTagValues(version series.Version, seriesID uint32) error {
	_, err := s.metaGetter.GetTagValues(s.metricID, s.tagKeys, version, roaring.BitmapOf(seriesID))
	fmt.Println(version)
	fmt.Println(seriesID)
	if err != nil {
		return err
	}
	//TODO set group by tag values
	s.tagValues = []string{"a", "b"}
	return nil
}

// waitComplete waits worker process complete
func (s *scanWorker) waitComplete() {
	// if ctx isn't done need wait, else timeout do not wait
	select {
	case <-s.done:
	case <-s.ctx.Done():
	}
}
