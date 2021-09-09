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

func Test_Min(t *testing.T) {
	m1 := newMin("min")
	var wg sync.WaitGroup
	for range [10]struct{}{} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				m1.Update(10)
				m1.Update(20)
				m1.Update(10)
				m1.Update(9)
				m1.Update(21)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, float64(9), m1.Get())
}
