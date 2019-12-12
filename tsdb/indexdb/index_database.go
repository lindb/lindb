package indexdb

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

// indexDatabase implements IndexDatabase
type indexDatabase struct {
	idGetter            metadb.IDGetter
	invertedIndexFamily kv.Family
}

// NewIndexDatabase returns a new IndexDatabase
func NewIndexDatabase(
	idGetter metadb.IDGetter,
	invertedIndexFamily kv.Family,
) IndexDatabase {
	return &indexDatabase{
		idGetter:            idGetter,
		invertedIndexFamily: invertedIndexFamily,
	}
}

// SuggestTagValues returns suggestions from given metricName, tagKey and prefix of tagValue
func (db *indexDatabase) SuggestTagValues(
	metricName string,
	tagKey string,
	tagValuePrefix string,
	limit int,
) []string {
	if limit <= 0 {
		return nil
	}
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	metricID, err := db.idGetter.GetMetricID(metricName)
	if err != nil {
		return nil
	}
	tagKeyID, err := db.idGetter.GetTagKeyID(metricID, tagKey)
	if err != nil {
		return nil
	}
	snapShot := db.invertedIndexFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(tagKeyID)
	if err != nil {
		return nil
	}
	return invertedindex.NewReader(readers).SuggestTagValues(tagKeyID, tagValuePrefix, limit)
}

func (db *indexDatabase) GetGroupingContext(metricID uint32, tagKeys []string,
	version series.Version,
) (series.GroupingContext, error) {
	//FIXME need impl
	//tagKeyIDs := make([]uint32, len(tagKeys))
	//// get tag key ids
	//for idx, tagKey := range tagKeys {
	//	//TODO need opt, plan has got tag key ids
	//	tagKeyID, err := db.idGetter.GetTagKeyID(metricID, tagKey)
	//	if err != nil {
	//		return nil, err
	//	}
	//	tagKeyIDs[idx] = tagKeyID
	//}
	return nil, nil
}

// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
func (db *indexDatabase) FindSeriesIDsByExpr(
	metricID uint32,
	expr stmt.TagFilter,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	tagKeyID, err := db.idGetter.GetTagKeyID(metricID, expr.TagKey())
	if err != nil {
		return nil, err
	}
	snapShot := db.invertedIndexFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(tagKeyID)
	if err != nil {
		return nil, err
	}
	if len(readers) == 0 {
		return nil, series.ErrNotFound
	}
	return invertedindex.NewReader(readers).FindSeriesIDsByExprForTagKeyID(tagKeyID, expr, timeRange)
}

// GetSeriesIDsForTag get series ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsForTag(
	metricID uint32,
	tagKey string,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	tagKeyID, err := db.idGetter.GetTagKeyID(metricID, tagKey)
	if err != nil {
		return nil, err
	}
	snapShot := db.invertedIndexFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(tagKeyID)
	if err != nil {
		return nil, err
	}
	return invertedindex.NewReader(readers).GetSeriesIDsForTagKeyID(tagKeyID, timeRange)
}
