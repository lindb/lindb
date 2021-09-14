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

	"github.com/lindb/lindb/internal/linmetric"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_Pool_Submit(t *testing.T) {
	grNum := runtime.NumGoroutine()
	pool := NewPool("test", 2, time.Second*5, linmetric.NewScope("1"))
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
	pool.Stop()
	pool.Stop()
	// reject all task
	go do(100)
	<-finished
	assert.Equal(t, int32(100), c.Load())
}

func Test_Pool_Statistics(t *testing.T) {
	p := NewPool("test", 0, time.Millisecond*100, linmetric.NewScope("2"))
	wp := p.(*workerPool)

	for i := 0; i < 10; i++ {
		p.SubmitAndWait(nil)
		p.SubmitAndWait(func() {
		})
	}
	assert.Equal(t, float64(1), wp.workersAlive.Get())

	time.Sleep(time.Second)
	p.Stop()
	assert.Equal(t, float64(0), wp.workersAlive.Get())
}
