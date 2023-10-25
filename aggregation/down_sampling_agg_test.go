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

package aggregation

import (
	"math"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

func Test_fillInfBlock(t *testing.T) {
	for i := infBlockSize - 10; i < infBlockSize+100; i++ {
		assertBlockInf(t, i)
	}
	assert.Len(t, getFloat64Slice(1), 1)
}

func Test_getFloat64Slice(t *testing.T) {
	defer func() {
		float64Pool = sync.Pool{}
	}()
	float64Pool = sync.Pool{
		New: func() any {
			return nil
		},
	}
	assert.Len(t, getFloat64Slice(10), 10)
	float64Pool = sync.Pool{
		New: func() any {
			f := make([]float64, 1)
			return &f
		},
	}
	assert.Len(t, getFloat64Slice(5), 5)
	float64Pool = sync.Pool{
		New: func() any {
			f := make([]float64, 10)
			return &f
		},
	}
	assert.Len(t, getFloat64Slice(5), 5)
}

func assertBlockInf(t *testing.T, size int) {
	sl := getFloat64Slice(size)
	fillInfBlock(sl)
	for _, v := range sl {
		assert.True(t, math.IsInf(v, 1))
	}
	assert.Len(t, sl, size)
	putFloat64Slice(&sl)
}

func TestDownSampling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case 1: no data
	found := 0
	getter := encoding.NewMockTSDValueGetter(ctrl)
	getter.EXPECT().GetValue(gomock.Any()).Return(0.0, false)
	DownSampling(timeutil.SlotRange{}, timeutil.SlotRange{}, 1, 0, getter, nil)
	assert.Equal(t, 0, found)
	// case 2: out of target range
	getter.EXPECT().GetValue(gomock.Any()).Return(1.0, true)
	DownSampling(timeutil.SlotRange{}, timeutil.SlotRange{Start: 10}, 1, 0, getter, nil)
	assert.Equal(t, 0, found)
	getter.EXPECT().GetValue(gomock.Any()).Return(1.0, true)
	DownSampling(timeutil.SlotRange{Start: 10, End: 11}, timeutil.SlotRange{}, 1, 0, getter, nil)
	assert.Equal(t, 0, found)
	// case 3: find data
	getter.EXPECT().GetValue(gomock.Any()).Return(1.0, true)
	DownSampling(timeutil.SlotRange{}, timeutil.SlotRange{}, 1, 0, getter, func(_ int, _ float64) {
		found++
	})
	assert.Equal(t, 1, found)
}
