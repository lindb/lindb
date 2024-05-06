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

	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetricLoader_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockMetricReader(ctrl)
	r.EXPECT().GetTimeRange().Return(timeutil.SlotRange{}).MaxTimes(2)

	ctx := &flow.DataLoadContext{
		SeriesIDHighKey: 0,
		ShardExecuteCtx: &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				Query: &stmt.Query{},
			},
		},
	}
	var seriesOffsets *encoding.FixedOffsetDecoder
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "series id not found",
			prepare: func() {
				ctx.LowSeriesIDsContainer = roaring.BitmapOf(1).GetContainer(0)
			},
		},
		{
			name: "found series block err",
			prepare: func() {
				ctx.LowSeriesIDsContainer = roaring.BitmapOf(10).GetContainer(0)
				seriesOffsets = encoding.NewFixedOffsetDecoder()
				encoder := encoding.NewFixedOffsetEncoder(true)
				encoder.Add(100)
				data := encoder.MarshalBinary()
				_, _ = seriesOffsets.Unmarshal(data)
			},
		},
		{
			name: "found series block successfully",
			prepare: func() {
				ctx.LowSeriesIDsContainer = roaring.BitmapOf(10).GetContainer(0)
				seriesOffsets = encoding.NewFixedOffsetDecoder()
				encoder := encoding.NewFixedOffsetEncoder(true)
				encoder.Add(0)
				data := encoder.MarshalBinary()
				_, _ = seriesOffsets.Unmarshal(data)

				r.EXPECT().readSeriesData(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				seriesOffsets = nil
			}()
			if tt.prepare != nil {
				tt.prepare()
			}

			s := newMetricLoader(r, nil, roaring.BitmapOf(10).GetContainer(0), seriesOffsets)
			ctx.Grouping()
			s.Load(ctx)
		})
	}
}
