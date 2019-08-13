package indextbl

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/series"
)

// todo: @codingcrush, implementation

//go:generate mockgen -source ./series_reader.go -destination=./series_reader_mock.go -package indextbl

// SeriesIndexReader reads tag k/v info from the kv table
type SeriesIndexReader interface {
	series.Filter
	series.MetadataGetter
}

// seriesIndexReader implements SeriesIndexReader
type seriesIndexReader struct {
	snapshot kv.Snapshot
}

// NewSeriesIndexReader returns a new SeriesIndexReader
func NewSeriesIndexReader(snapshot kv.Snapshot) SeriesIndexReader {
	return &seriesIndexReader{snapshot: snapshot}
}

// GetTagValues returns tag values by tag keys and spec version for metric level
func (r *seriesIndexReader) GetTagValues(metricID uint32, tagKeys []string, version int64) (
	tagValues [][]string, err error) {
	return nil, nil
}

// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
func (r *seriesIndexReader) FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	return nil, nil
}

// GetSeriesIDsForTag get series ids for spec metric's tag key
func (r *seriesIndexReader) GetSeriesIDsForTag(metricID uint32, tagKey string,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	return nil, nil
}
