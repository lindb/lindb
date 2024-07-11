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

package encoding

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DeltaBitPackingEncoder_Add(t *testing.T) {
	p := NewDeltaBitPackingEncoder()

	p.Add(1)
	p.Add(10)
	p.Add(1)
	for i := 0; i < 100; i++ {
		p.Add(100)
	}

	p.Add(200)

	b := p.Bytes()

	d := NewDeltaBitPackingDecoder(b)

	count := 0
	for d.HasNext() {
		_ = d.Next()
		count++
	}
	assert.Equal(t, 104, count)
}

func Test_DeltaBitPackingEncoder_Reset(t *testing.T) {
	p := NewDeltaBitPackingEncoder()
	for i := 0; i < 100; i++ {
		p.Add(100)
	}
	b1 := p.Bytes()
	p.Reset()
	for i := 0; i < 100; i++ {
		p.Add(100)
	}
	b2 := p.Bytes()
	assert.Equal(t, b1, b2)
}

func Test_DeltaBitPackingEncoder_Decoder(t *testing.T) {
	p := NewDeltaBitPackingEncoder()
	d := NewDeltaBitPackingDecoder(nil)

	for range [10]struct{}{} {
		p.Reset()
		list := getRandomList()
		for _, v := range list {
			p.Add(v)
		}
		b := p.Bytes()

		d.Reset(b)
		count := 0
		for d.HasNext() {
			value := d.Next()
			assert.Equal(t, list[count], value)
			count++
		}
		assert.Equal(t, count, len(list))
	}
}

func getRandomList() []int32 {
	var list []int32

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < math.MaxUint16; i++ {
		list = append(list, r.Int31n(math.MaxInt32))
	}
	return list
}
