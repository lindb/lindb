package query

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

func TestMetadataStorageExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)
	metricMetaSuggester := series.NewMockMetricMetaSuggester(ctrl)
	db.EXPECT().MetricMetaSuggester().Return(metricMetaSuggester).AnyTimes()

	// suggest metric name
	exec := newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.Metric,
	})
	metricMetaSuggester.EXPECT().SuggestMetrics(gomock.Any(), gomock.Any()).Return([]string{"a"})
	result, err := exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// suggest tag keys
	exec = newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.TagKey,
	})
	metricMetaSuggester.EXPECT().SuggestTagKeys(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"})
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// suggest tag values
	exec = newMetadataStorageExecutor(db, []int32{1, 2}, &stmt.Metadata{
		Type: stmt.TagValue,
	})
	shard := tsdb.NewMockShard(ctrl)
	db.EXPECT().GetShard(int32(1)).Return(nil, false)
	db.EXPECT().GetShard(int32(2)).Return(shard, true)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB)
	indexDB.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"})
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
}
