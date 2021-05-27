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

package metricsdata

import (
	"testing"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
)

func TestMetricLoader_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockMetricReader(ctrl)
	r.EXPECT().GetTimeRange().Return(timeutil.SlotRange{}).MaxTimes(2)

	// case 1: series id not exist
	s := newMetricLoader(r, roaring.BitmapOf(10).GetContainer(0), nil)
	s.Load(1)
	// case 2: read series data
	r.EXPECT().readSeriesData(gomock.Any())
	encoder := encoding.NewFixedOffsetEncoder()
	encoder.Add(100)
	data := encoder.MarshalBinary()
	seriesOffsets := encoding.NewFixedOffsetDecoder(data)
	s = newMetricLoader(r, roaring.BitmapOf(10).GetContainer(0), seriesOffsets)
	s.Load(10)
}
