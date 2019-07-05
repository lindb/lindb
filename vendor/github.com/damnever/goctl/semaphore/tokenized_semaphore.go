package semaphore

import (
	"context"
	"sync"
)

type tokenWrapper struct {
	token string
	doneC chan struct{}
}

// TokenizedSemaphore is a kind of semaphore which only allow
// the same token acquire once until it released.
type TokenizedSemaphore struct {
	l       sync.Mutex
	seq     int64
	count   int
	limit   int
	pending map[int64]tokenWrapper
	tokens  map[string]bool
}

// NewTokenizedSemaphore creates a new TokenizedSemaphore.
func NewTokenizedSemaphore(limit int) *TokenizedSemaphore {
	return &TokenizedSemaphore{
		seq:     0,
		count:   0,
		limit:   limit,
		pending: make(map[int64]tokenWrapper),
		tokens:  make(map[string]bool, limit),
	}
}

// Acquire acquires a semaphore with given token.
func (s *TokenizedSemaphore) Acquire(ctx context.Context, token string) error {
	s.l.Lock()
	if s.count < s.limit && !s.tokens[token] {
		s.tokens[token] = true
		s.count++
		s.l.Unlock()
		return nil
	}

	s.seq++ // Overflow?? no such thing..
	seq := s.seq
	td := tokenWrapper{
		token: token,
		doneC: make(chan struct{}),
	}
	s.pending[seq] = td
	s.l.Unlock()

	select {
	case <-ctx.Done():
		s.l.Lock()
		defer s.l.Unlock()
		select {
		case <-td.doneC: // Double check.
			return nil // Must let user to release it.
		default:
			delete(s.pending, seq)
		}
		return ctx.Err()
	case <-td.doneC:
		return nil
	}
}

// Release releases a semaphore with given token.
func (s *TokenizedSemaphore) Release(token string) error {
	s.l.Lock()
	defer s.l.Unlock()

	if !s.tokens[token] {
		return ErrOpMismatch
	}

	delete(s.tokens, token)
	s.count--

	// Resume pendings.
	for seq, td := range s.pending {
		if s.tokens[td.token] {
			continue
		}
		s.count++
		s.tokens[td.token] = true
		close(td.doneC)
		delete(s.pending, seq)
		return nil
	}
	return nil
}
