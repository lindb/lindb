package queue

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrEmpty indicates queue is empty.
	ErrEmpty = errors.New("queue is empty")
)

type getter chan interface{}

// Queue is a simple synchronized queue,
// just like a channel with infinite buffer size.
type Queue struct {
	l       sync.RWMutex
	items   *Ring
	getters *Ring
}

// NewQueue creates a new Queue.
func NewQueue() *Queue {
	return &Queue{
		items:   NewRing(),
		getters: NewRing(),
	}
}

// Get gets an item from queue, block until ctx done.
func (q *Queue) Get(ctx context.Context) (interface{}, error) {
	q.l.Lock()
	if q.items.Len() > 0 {
		item := q.items.Pop()
		q.l.Unlock()
		return item, nil
	}

	waiter := make(getter, 1) // Need buffer to prevent dead lock..
	// Using the same lock here, so there is no chance Get op
	// got starving even if items is not empty.
	q.getters.Append(waiter)
	q.l.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case item := <-waiter:
		return item, nil
	}
}

// GetNoWait gets an item from queue immediately,
// ErrEmpty returned if queue is empty.
func (q *Queue) GetNoWait() (interface{}, error) {
	q.l.Lock()
	defer q.l.Unlock()
	if q.items.Len() > 0 {
		return q.items.Pop(), nil
	}
	return nil, ErrEmpty
}

// Put puts an item into queue, always no wait.
func (q *Queue) Put(item interface{}) {
	q.l.Lock()
	q.items.Append(item)
	for q.getters.Len() > 0 && q.items.Len() > 0 {
		waiter := q.getters.Pop().(getter)
		waiter <- q.items.Pop()
	}
	q.l.Unlock()
}

// Len returns the number of items.
func (q *Queue) Len() int {
	q.l.RLock()
	n := q.items.Len()
	q.l.RUnlock()
	return n
}
