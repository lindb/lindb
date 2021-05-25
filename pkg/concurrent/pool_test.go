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
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_Pool_Submit(t *testing.T) {
	grNum := runtime.NumGoroutine()
	pool := NewPool("test", 2, time.Second*5)
	// num. of pool + 1 dispatcher, workers has not been spawned
	assert.Equal(t, grNum+1, runtime.NumGoroutine())

	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.Submit(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(100)
	<-finished
	assert.True(t, grNum+2+1 <= runtime.NumGoroutine())
	pool.Stop()
	pool.Stop()
	// reject all task
	go do(100)
	<-finished
	assert.Equal(t, int32(100), c.Load())
}

func Test_Pool_Statistics(t *testing.T) {
	p := NewPool("test", 0, time.Millisecond*100)
	s := p.Statistics()
	assert.Zero(t, s.AliveWorkers)
	assert.Zero(t, s.CreatedWorkers)
	assert.Zero(t, s.KilledWorkers)
	assert.Zero(t, s.ConsumedTasks)
	for i := 0; i < 10; i++ {
		p.SubmitAndWait(nil)
		p.SubmitAndWait(func() {
		})
	}
	s = p.Statistics()
	assert.Equal(t, 1, s.AliveWorkers)
	assert.Equal(t, 1, s.CreatedWorkers)
	assert.Equal(t, 0, s.KilledWorkers)
	assert.Equal(t, 10, s.ConsumedTasks)

	time.Sleep(time.Second)
	p.Stop()
	s = p.Statistics()
	assert.Equal(t, 0, s.AliveWorkers)
	assert.Equal(t, 1, s.CreatedWorkers)
	assert.Equal(t, 1, s.KilledWorkers)
	assert.Equal(t, 10, s.ConsumedTasks)
}
