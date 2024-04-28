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

package index

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

const testMetricDir = "./metric_index_segment"

func TestNewMetricIndexSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	var (
		metaDir  = path.Join(name, "meta")
		indexDir = path.Join(name, "index")
	)

	cases := []struct {
		segment int
	}{
		{
			segment: 202403,
		},
		{
			segment: 202404,
		},
	}

	// make segment dir
	for _, c := range cases {
		segment := strconv.Itoa(c.segment)
		err := mkDirFunc(path.Join(indexDir, segment))
		assert.NoError(t, err)
	}

	// add some noise
	for _, noise := range []string{"20240102test", "20230102test"} {
		err := mkDirFunc(path.Join(indexDir, noise))
		assert.NoError(t, err)
	}

	metaDB, err := NewMetricMetaDatabase(metaDir)
	assert.NoError(t, err)

	indexSegment, err := NewMetricIndexSegment(indexDir, metaDB)
	assert.NoError(t, err)

	// assert type
	segment := indexSegment.(*metricIndexSegment)

	// assert number of indexDB
	assert.Equal(t, len(cases), len(segment.indices))

	for _, c := range cases {
		index := segment.indices[c.segment]
		assert.NotNil(t, index)
	}
}

func TestMetricIndexSegment_GetOrCreateIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	var (
		metaDir  = path.Join(name, "meta")
		indexDir = path.Join(name, "index")
	)

	metaDB, err := NewMetricMetaDatabase(metaDir)
	assert.NoError(t, err)

	indexSegment, err := NewMetricIndexSegment(indexDir, metaDB)
	assert.NoError(t, err)

	// assert type
	segment := indexSegment.(*metricIndexSegment)

	segments := []int{202403, 202404}
	for _, s := range segments {
		tm, err := time.Parse("200601", strconv.Itoa(s))
		assert.NoError(t, err)
		index, err := segment.GetOrCreateIndex(tm.UnixMilli())
		assert.NoError(t, err)
		assert.NotNil(t, index)
		later := tm.UnixMilli() + 15*24*time.Hour.Milliseconds()
		index, err = segment.GetOrCreateIndex(later)
		assert.NoError(t, err)
		assert.NotNil(t, index)
	}

	assert.Equal(t, len(segments), len(segment.indices))
}

func TestMetricIndexSegment_GetSeriesIDsByTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	var (
		metaDir  = path.Join(name, "meta")
		indexDir = path.Join(name, "index")
	)

	metaDB, err := NewMetricMetaDatabase(metaDir)
	assert.NoError(t, err)

	segment, err := NewMetricIndexSegment(indexDir, metaDB)
	assert.NoError(t, err)

	tm, err := time.Parse("200601", strconv.Itoa(202404))
	assert.NoError(t, err)
	index, err := segment.GetOrCreateIndex(tm.UnixMilli())
	assert.NoError(t, err)

	var timeRange = timeutil.TimeRange{
		Start: tm.UnixMilli(),
		End:   tm.UnixMilli() + 15*24*time.Hour.Milliseconds(),
	}

	cases := []struct {
		namespace  []byte
		metricName []byte
	}{
		{
			namespace:  []byte("n"),
			metricName: []byte("m"),
		},
	}

	for _, c := range cases {
		metricID, err := metaDB.genMetricID(c.namespace, c.metricName)
		assert.NoError(t, err)
		assert.Equal(t, metricID, metric.ID(0))

		tags := tag.Tags{tag.Tag{Key: []byte("idc"), Value: []byte("sh")}}
		seriesID := uint32(0)
		ch := make(chan struct{})

		mNotifier := &MetaNotifier{
			Namespace:  string(c.namespace),
			MetricName: string(c.metricName),
			MetricID:   metricID,
			TagHash:    1,
			Tags:       tags,
			Callback: func(id uint32, err error) {
				assert.NoError(t, err)
				seriesID = id
				ch <- struct{}{}
			},
		}

		index.Notify(mNotifier)
		<-ch
		assert.Equal(t, seriesID, uint32(0))
		index.PrepareFlush()
		_ = index.Flush()
		time.Sleep(120 * time.Millisecond)

		var (
			tagKeyID   tag.KeyID
			tagValueID uint32
		)

		series, err := index.GetSeriesIDsForMetric(metricID)
		assert.NoError(t, err)
		series2, err := segment.GetSeriesIDsForMetric(metricID, timeRange)
		assert.NoError(t, err)
		assert.NotNil(t, series)
		assert.NotNil(t, series2)
		assert.False(t, series.IsEmpty())
		assert.False(t, series2.IsEmpty())
		assert.Equal(t, series, series2)

		series, err = index.GetSeriesIDsForTag(tagKeyID)
		assert.NoError(t, err)
		series2, err = segment.GetSeriesIDsForTag(tagKeyID, timeRange)
		assert.NoError(t, err)
		assert.NotNil(t, series)
		assert.NotNil(t, series2)
		assert.False(t, series.IsEmpty())
		assert.False(t, series2.IsEmpty())
		assert.Equal(t, series, series2)

		tagValueIDs := roaring.New()
		tagValueIDs.Add(tagValueID)
		series, err = index.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
		assert.NoError(t, err)
		series2, err = segment.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs, timeRange)
		assert.NoError(t, err)
		assert.NotNil(t, series)
		assert.NotNil(t, series2)
		assert.False(t, series.IsEmpty())
		assert.False(t, series2.IsEmpty())
		assert.Equal(t, series, series2)

		seriesIDsAfterFiltering := roaring.New()
		seriesIDsAfterFiltering.Add(seriesID)
		shardExecuteContext := &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				GroupByTagKeyIDs: []tag.KeyID{tagKeyID},
				Query: &stmt.Query{
					TimeRange: timeRange,
				},
			},
			SeriesIDsAfterFiltering: seriesIDsAfterFiltering,
		}
		scanner, err := index.GetGroupingContext(shardExecuteContext)
		assert.NoError(t, err)
		assert.NotNil(t, scanner)
		err = segment.GetGroupingContext(shardExecuteContext)
		assert.NoError(t, err)

		assert.NoError(t, segment.Close())
	}
}

func TestMetricIndexSegment_Flush_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		newIndexDBFunc = NewMetricIndexDatabase
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	indexDB := NewMockMetricIndexDatabase(ctrl)
	newIndexDBFunc = func(dir string, metaDB MetricMetaDatabase) (MetricIndexDatabase, error) {
		return indexDB, nil
	}
	metaDB := NewMockMetricMetaDatabase(ctrl)
	indexSegment, err := NewMetricIndexSegment(name, metaDB)
	assert.NoError(t, err)
	indexDB.EXPECT().Notify(gomock.Any()).Do(func(n Notifier) {
		mn := n.(*FlushNotifier)
		mn.Callback(fmt.Errorf("err"))
	})
	_, err = indexSegment.GetOrCreateIndex(time.Now().UnixMilli())
	assert.NoError(t, err)
	assert.Error(t, indexSegment.Flush())
}

func TestMetricIndexSegment_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		newIndexDBFunc = NewMetricIndexDatabase
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	indexDB := NewMockMetricIndexDatabase(ctrl)
	newIndexDBFunc = func(dir string, metaDB MetricMetaDatabase) (MetricIndexDatabase, error) {
		return indexDB, nil
	}
	metaDB := NewMockMetricMetaDatabase(ctrl)
	indexSegment, err := NewMetricIndexSegment(name, metaDB)
	assert.NoError(t, err)
	indexDB.EXPECT().Notify(gomock.Any()).Do(func(n Notifier) {
		mn := n.(*FlushNotifier)
		mn.Callback(nil)
	})
	_, err = indexSegment.GetOrCreateIndex(time.Now().UnixMilli())
	assert.NoError(t, err)
	assert.NoError(t, indexSegment.Flush())
}

func TestMetricIndexSegment_Close_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		newIndexDBFunc = NewMetricIndexDatabase
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	indexDB := NewMockMetricIndexDatabase(ctrl)
	newIndexDBFunc = func(dir string, metaDB MetricMetaDatabase) (MetricIndexDatabase, error) {
		return indexDB, nil
	}
	metaDB := NewMockMetricMetaDatabase(ctrl)
	indexSegment, err := NewMetricIndexSegment(name, metaDB)
	assert.NoError(t, err)
	indexDB.EXPECT().Close().Return(fmt.Errorf("err"))
	_, err = indexSegment.GetOrCreateIndex(time.Now().UnixMilli())
	assert.NoError(t, err)
	assert.Error(t, indexSegment.Close())
}

func TestMetricIndexSegment_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := testMetricDir
	defer func() {
		newIndexDBFunc = NewMetricIndexDatabase
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	indexDB := NewMockMetricIndexDatabase(ctrl)
	newIndexDBFunc = func(dir string, metaDB MetricMetaDatabase) (MetricIndexDatabase, error) {
		return indexDB, nil
	}
	metaDB := NewMockMetricMetaDatabase(ctrl)
	indexSegment, err := NewMetricIndexSegment(name, metaDB)
	assert.NoError(t, err)
	indexDB.EXPECT().Close().Return(nil)
	_, err = indexSegment.GetOrCreateIndex(time.Now().UnixMilli())
	assert.NoError(t, err)
	assert.NoError(t, indexSegment.Close())
}

func TestMetricIndexSegment_getSegment(t *testing.T) {
	index := metricIndexSegment{}
	now := time.Now()
	later := now.AddDate(0, 1, 0)
	start, end := index.getSegment(timeutil.TimeRange{
		Start: now.UnixMilli(),
		End:   later.UnixMilli(),
	})
	s := now.Format("200601")
	e := later.Format("200601")
	gotStart, err := strconv.Atoi(s)
	assert.NoError(t, err)
	gotEnd, err := strconv.Atoi(e)
	assert.NoError(t, err)
	assert.Equal(t, start, gotStart)
	assert.Equal(t, end, gotEnd)
}
