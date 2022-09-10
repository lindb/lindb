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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

func Test_Limiter(t *testing.T) {
	limiter := NewLimiter(
		context.TODO(),
		10,
		500*time.Millisecond,
		metrics.NewLimitStatistics("test", linmetric.BrokerRegistry),
	)
	var (
		wg          sync.WaitGroup
		atomicError atomic.Error
	)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := limiter.Do(func() error {
				time.Sleep(time.Millisecond)
				return nil
			}); err != nil {
				atomicError.Store(err)
			}
		}()
	}
	wg.Wait()
	assert.NoError(t, atomicError.Load())
}

func Test_Limiter_Timeout(t *testing.T) {
	limiter := NewLimiter(
		context.TODO(),
		1,
		10*time.Millisecond,
		metrics.NewLimitStatistics("test", linmetric.BrokerRegistry),
	)
	var (
		wg          sync.WaitGroup
		atomicError atomic.Error
	)
	limiter.tokens <- struct{}{} // put one
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := limiter.Do(func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		}); err != nil {
			atomicError.Store(err)
		}
	}()
	wg.Wait()
	assert.Equal(t, ErrConcurrencyLimiterTimeout, atomicError.Load())
}

func Test_Limiter_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel() // cancel contexxt
	limiter := NewLimiter(
		ctx,
		1,
		10*time.Millisecond,
		metrics.NewLimitStatistics("test", linmetric.BrokerRegistry),
	)
	var (
		wg          sync.WaitGroup
		atomicError atomic.Error
	)
	limiter.tokens <- struct{}{} // put one
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := limiter.Do(func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}); err != nil {
			atomicError.Store(err)
		}
	}()
	wg.Wait()
	assert.NoError(t, atomicError.Load())
}

func TestTimerPool(t *testing.T) {
	defer func() {
		timerPool = sync.Pool{}
	}()
	t.Run("test put failure", func(_ *testing.T) {
		timerPool = sync.Pool{}
		for i := 0; i < 100; i++ {
			timer := acquireTimer(time.Minute)
			assert.NotNil(t, timer)
			timer.Stop()
			releaseTimer(timer)
		}
	})
	t.Run("timer active, need new one", func(_ *testing.T) {
		timerPool = sync.Pool{}
		timerPool.Put(time.NewTimer(time.Minute))
		timer := acquireTimer(time.Minute)
		assert.NotNil(t, timer)
	})
}

func Benchmark_TimerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tm := acquireTimer(time.Microsecond)
		<-tm.C
		releaseTimer(tm)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func Benchmark_Timer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tm := time.NewTimer(time.Microsecond)
		<-tm.C
	}
	b.StopTimer()
	b.ReportAllocs()
}
