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

	"github.com/stretchr/testify/assert"
)

func Test_Counter(t *testing.T) {
	c1 := NewCounter("count")
	var wg sync.WaitGroup
	for range [10]struct{}{} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				c1.Add(2)
				c1.Sub(1)
				c1.Incr()
				c1.Decr()
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, float64(100), c1.Get())
	assert.Equal(t, float64(100), c1.getAndReset())
	// reset
	assert.Equal(t, float64(0), c1.Get())
}
