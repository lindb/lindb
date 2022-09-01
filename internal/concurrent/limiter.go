// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package concurrent

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/lindb/lindb/metrics"
)

var ErrConcurrencyLimiterTimeout = errors.New("reaches the max concurrency for writing")

type Limiter struct {
	ctx     context.Context
	timeout time.Duration
	tokens  chan struct{}

	statistics *metrics.LimitStatistics
}

// NewLimiter creates a limiter based of buffer channel.
// It limits the concurrency for writing.
func NewLimiter(ctx context.Context, maxConcurrency int, timeout time.Duration, statistics *metrics.LimitStatistics) *Limiter {
	return &Limiter{
		ctx:        ctx,
		timeout:    timeout,
		tokens:     make(chan struct{}, maxConcurrency),
		statistics: statistics,
	}
}

func (l *Limiter) Do(f func() error) error {
	select {
	case l.tokens <- struct{}{}:
		err := f()
		l.statistics.Processed.Incr()
		<-l.tokens
		return err
	default:
		// tokens are taken, so waits one to be free
	}
	l.statistics.Throttles.Incr()

	timer := acquireTimer(l.timeout)
	select {
	case l.tokens <- struct{}{}:
		releaseTimer(timer)
		err := f()
		l.statistics.Processed.Incr()
		<-l.tokens
		return err
	case <-l.ctx.Done():
		return nil
	case <-timer.C:
		releaseTimer(timer)
		l.statistics.Timeouts.Incr()
		return ErrConcurrencyLimiterTimeout
	}
}

var timerPool sync.Pool

func acquireTimer(d time.Duration) *time.Timer {
	item := timerPool.Get()
	if item == nil {
		return time.NewTimer(d)
	}
	t := item.(*time.Timer)
	// timer.Reset returns if timer is active.
	if t.Reset(d) {
		return time.NewTimer(d)
	}
	return t
}

func releaseTimer(t *time.Timer) {
	if !t.Stop() {
		return
	}
	timerPool.Put(t)
}
