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

package memdb

import (
	"fmt"
	"testing"

	"github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/metric"
)

func TestIndexDatabase_GetOrCreateTimeSeriesIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := NewIndexDatabase(nil, nil)

	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}

	row := protoToStorageRow(m)
	tsIndex := idx.GetOrCreateTimeSeriesIndex(row)
	tsIndex.GenMemTimeSeriesID(row.TagsHash(), idx.GenMemSeriesID)
	assert.NotNil(t, tsIndex)

	tsIndex1 := idx.(*indexDatabase).getOrCreateTimeSeriesIndex(row.NameHash())
	assert.NotNil(t, tsIndex1)
	assert.Equal(t, tsIndex, tsIndex1)

	// not found
	tsIndex, ok := idx.GetTimeSeriesIndex(100)
	assert.Nil(t, tsIndex)
	assert.False(t, ok)

	db := NewMockMemoryDatabase(ctrl)
	db.EXPECT().CreatedTime().Return(timeutil.Now())
	db.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(1, 2, 3))
	tsIndex1.IndexTimeSeries(100, 1)
	assert.NotZero(t, tsIndex1.NumOfSeries())
	idx.Cleanup(db)
	assert.NotZero(t, tsIndex1.NumOfSeries())

	db.EXPECT().CreatedTime().Return(timeutil.Now())
	db.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(3))
	tsIndex1.ExpireTimeSeriesIDs(roaring.BitmapOf(1, 0), timeutil.Now()-4*timeutil.OneHour)
	assert.NotZero(t, tsIndex1.NumOfSeries())
	idx.Cleanup(db)
	assert.Zero(t, tsIndex1.NumOfSeries())

	idx.Close()
}

func TestIndexDatabase_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	indexDB.EXPECT().PrepareFlush()
	indexDB.EXPECT().Flush().Return(nil)
	ch := make(chan error)
	idx := NewIndexDatabase(nil, indexDB)
	idx.Notify(&FlushEvent{
		Callback: func(err error) {
			ch <- err
		},
	})

	err := <-ch
	assert.NoError(t, err)
}

func TestIndexDatabase_handleRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	memMetaDB := NewMetadataDatabase(&models.DatabaseConfig{}, metaDB)
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	memIndexDB := NewIndexDatabase(memMetaDB, indexDB)

	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}

	row := protoToStorageRow(m)
	row.Add(100)

	metaDB.EXPECT().GenMetricID([]byte("ns"), []byte("test1")).Return(metric.ID(0), fmt.Errorf("err"))
	memIndexDB.(*indexDatabase).handleRow(row)

	metaDB.EXPECT().GenMetricID([]byte("ns"), []byte("test1")).Return(metric.ID(0), nil)
	indexDB.EXPECT().GenSeriesID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
	memIndexDB.(*indexDatabase).handleRow(row)
}
