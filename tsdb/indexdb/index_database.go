package indexdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

// for testing
var (
	newReader = invertedindex.NewReader
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
	tagKeysLength := len(tagKeys)
	tagKeyIDs := make([]uint32, tagKeysLength)
	// get tag key ids
	for idx, tagKey := range tagKeys {
		//TODO need opt, plan has got tag key ids
		tagKeyID, err := db.idGetter.GetTagKeyID(metricID, tagKey)
		if err != nil {
			return nil, err
		}
		tagKeyIDs[idx] = tagKeyID
	}

	snapShot := db.invertedIndexFamily.GetSnapshot()
	//FIXME need close snapshot after query completed
	//defer snapShot.Close()

	var err error
	defer func() {
		if err != nil {
			snapShot.Close()
		}
	}()
	var readers []table.Reader
	gCtx := query.NewGroupContext(tagKeysLength)
	for idx, tagKeyID := range tagKeyIDs {
		readers, err = snapShot.FindReaders(tagKeyID)
		if err != nil {
			return nil, err
		}
		reader := newReader(readers)
		tagValuesEntrySet := query.NewTagValuesEntrySet()
		gCtx.SetTagValuesEntrySet(idx, tagValuesEntrySet)
		err1 := reader.WalkTagValues(
			tagKeyID,
			"",
			func(tagValue []byte, it invertedindex.TagValueIterator) bool {
				for it.HasNext() {
					if it.DataVersion() == version {
						bitmapData := it.Next()
						seriesBitmap := roaring.New()
						err = seriesBitmap.UnmarshalBinary(bitmapData)
						if err != nil {
							return false
						}
						tagValuesEntrySet.AddTagValue(string(tagValue), seriesBitmap)
					}
				}
				return true
			})
		if err1 != nil {
			err = err1
		}
		if err != nil {
			return nil, err
		}
	}
	return gCtx, nil
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
