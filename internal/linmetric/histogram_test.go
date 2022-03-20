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

package linmetric

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Histogram(t *testing.T) {
	dh := NewHistogram()
	defer func() {
		dh.Update(func() {
		})
	}()

	concurrentDo(
		func() {
			dh.WithLinearBuckets(time.Millisecond, time.Second*5, 100)
			dh.WithExponentBuckets(time.Second, time.Second*5, 5)
		})
	concurrentDo(
		func() {
			dh.UpdateMilliseconds(2000)                      // bucket2
			dh.UpdateMilliseconds(5001)                      // bucket3
			dh.UpdateSeconds(6)                              // bucket3
			dh.UpdateDuration(time.Millisecond * 1700)       // bucket1
			dh.UpdateDuration(time.Millisecond * 1710)       // bucket2
			dh.UpdateDuration(time.Millisecond * 2920)       // bucket2
			dh.UpdateDuration(time.Millisecond * 2950)       // bucket3
			dh.UpdateDuration(time.Millisecond * 4990)       // bucket2
			dh.UpdateDuration(time.Second * 6)               // bucket3
			dh.UpdateSince(time.Now().Add(time.Second))      // drop
			dh.UpdateSince(time.Now().Add(-1 * time.Second)) // bucket0
		})
	assert.InDeltaSlice(t, []float64{100, 100, 300, 200, 300}, dh.bkts.values, 0.01)
}

func concurrentDo(f func()) {
	var wg sync.WaitGroup
	for range [100]struct{}{} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}
