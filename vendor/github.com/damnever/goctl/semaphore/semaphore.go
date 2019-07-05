// Package semaphore provides semaphore implementations.
package semaphore

import (
	"context"
	"errors"
)

var (
	// ErrOpMismatch is returned when you call Release before Acquire,
	// or Release after Acquire failed.
	ErrOpMismatch = errors.New("operation mismatch: Release called without a successful Acquire")
)

// Semaphore is a semaphore.
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a new Semaphore.
func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, limit),
	}
}

// Acquire acquires a semaphore, blocks until ctx done.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.sem <- struct{}{}:
		return nil
	}
}

// Release releases a semaphore.
func (s *Semaphore) Release() error {
	select {
	case <-s.sem:
		return nil
	default:
		return ErrOpMismatch
	}
}
