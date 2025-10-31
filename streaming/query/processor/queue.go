package processor

import (
	"context"
)

type Queue struct {
	eventCh chan any
}

func NewQueue(eventCh chan any) *Queue {
	return &Queue{
		eventCh: eventCh,
	}
}

func (q *Queue) Produce(event any) {
	if event == nil {
		return
	}
	q.eventCh <- event
}

func (q *Queue) Consume(ctx context.Context) (any, bool) {
	select {
	case err := <-ctx.Done():
		panic(err)
	case page, ok := <-q.eventCh:
		return page, ok
	}
}

func (q *Queue) GetInbound() chan any {
	return q.eventCh
}

func (q *Queue) Close() {
	close(q.eventCh)
}
