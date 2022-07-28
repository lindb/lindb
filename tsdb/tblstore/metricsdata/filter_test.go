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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestFileFilterResultSet_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockMetricReader(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	rs := newFileFilterResultSet(1, nil, reader, snapshot)
	reader.EXPECT().Load(gomock.Any())
	rs.Load(&flow.DataLoadContext{})
	assert.Equal(t, int64(1), rs.FamilyTime())
	reader.EXPECT().GetTimeRange().Return(timeutil.SlotRange{
		Start: 10,
		End:   20,
	})
	assert.Equal(t, timeutil.SlotRange{
		Start: 10,
		End:   20,
	}, rs.SlotRange())
	snapshot.EXPECT().Close()
	rs.Close()
}

func TestMetricsDataFilter_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockMetricReader(ctrl)
	filter := NewFilter(10, nil, []MetricReader{reader})

	// case 1: field not found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	rs, err := filter.Filter(roaring.BitmapOf(1, 2, 3), field.Metas{{ID: 1}, {ID: 30}})
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, rs)
	// case 2: series ids found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	reader.EXPECT().GetSeriesIDs().Return(roaring.BitmapOf(10, 200))
	rs, err = filter.Filter(roaring.BitmapOf(1, 2, 3), field.Metas{{ID: 2}, {ID: 30}})
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Nil(t, rs)
	// case 3: data found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	reader.EXPECT().GetSeriesIDs().Return(roaring.BitmapOf(10, 200))
	rs, err = filter.Filter(roaring.BitmapOf(1, 200, 3), field.Metas{{ID: 2}, {ID: 30}})
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.EqualValues(t, roaring.BitmapOf(200).ToArray(), rs[0].SeriesIDs().ToArray())
	reader.EXPECT().Path().Return("1.sst")
	assert.Equal(t, "1.sst", rs[0].Identifier())
}
