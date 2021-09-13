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
	"errors"
	"sync"
	"time"

	"github.com/lindb/lindb/internal/linmetric"
)

var ErrConcurrencyLimiterTimeout = errors.New("reaches the max concurrency for writing")

type Limiter struct {
	timeout time.Duration
	tokens  chan struct{}

	statistics struct {
		throttles *linmetric.BoundCounter // counter reaches the max-concurrency
		timeouts  *linmetric.BoundCounter // counter pending and then timeout
	}
}

// NewLimiter creates a limiter based of buffer channel.
// It limits the concurrency for writing.
func NewLimiter(maxConcurrency int, timeout time.Duration, scope linmetric.Scope) *Limiter {
	l := &Limiter{
		timeout: timeout,
		tokens:  make(chan struct{}, maxConcurrency),
	}
	l.statistics.throttles = scope.NewCounter("throttle_requests")
	l.statistics.timeouts = scope.NewCounter("timeout_requests")
	return l
}

func (l *Limiter) Do(f func() error) error {
	select {
	case l.tokens <- struct{}{}:
		err := f()
		<-l.tokens
		return err
	default:
		// tokens are taken, so waits one to be free
	}
	l.statistics.throttles.Incr()

	timer := acquireTimer(l.timeout)
	select {
	case l.tokens <- struct{}{}:
		releaseTimer(timer)
		err := f()
		<-l.tokens
		return err
	case <-timer.C:
		releaseTimer(timer)
		l.statistics.timeouts.Incr()
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
