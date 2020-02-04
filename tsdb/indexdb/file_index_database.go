package indexdb

//
//import (
//	"github.com/lindb/roaring"
//
//	"github.com/lindb/lindb/constants"
//	"github.com/lindb/lindb/kv"
//	"github.com/lindb/lindb/kv/table"
//	"github.com/lindb/lindb/pkg/timeutil"
//	"github.com/lindb/lindb/series"
//	"github.com/lindb/lindb/sql/stmt"
//	"github.com/lindb/lindb/tsdb/metadb"
//	"github.com/lindb/lindb/tsdb/query"
//	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
//)
//
//// for testing
//var (
//	newReader = invertedindex.NewReader
//)
//
//// fileIndexDatabase implements IndexDatabase
//type fileIndexDatabase struct {
//	idGetter            metadb.IDGetter
//	invertedIndexFamily kv.Family
//}
//
//// NewFileIndexDatabase returns a new IndexDatabase
//func NewFileIndexDatabase(
//	idGetter metadb.IDGetter,
//	invertedIndexFamily kv.Family,
//) FileIndexDatabase {
//	return &fileIndexDatabase{
//		idGetter:            idGetter,
//		invertedIndexFamily: invertedIndexFamily,
//	}
//}
//
//// SuggestTagValues returns suggestions from given tag key id and prefix of tagValue
//func (db *fileIndexDatabase) SuggestTagValues(
//	tagKeyID uint32,
//	tagValuePrefix string,
//	limit int,
//) []string {
//	if limit <= 0 {
//		return nil
//	}
//	if limit > constants.MaxSuggestions {
//		limit = constants.MaxSuggestions
//	}
//	snapShot := db.invertedIndexFamily.GetSnapshot()
//	defer snapShot.Close()
//
//	readers, err := snapShot.FindReaders(tagKeyID)
//	if err != nil {
//		return nil
//	}
//	return invertedindex.NewReader(readers).SuggestTagValues(tagKeyID, tagValuePrefix, limit)
//}
//
//func (db *fileIndexDatabase) GetGroupingContext(tagKeyIDs []uint32,
//	version series.Version,
//) (series.GroupingContext, error) {
//	snapShot := db.invertedIndexFamily.GetSnapshot()
//	//FIXME need close snapshot after query completed
//	//defer snapShot.Close()
//
//	var err error
//	defer func() {
//		if err != nil {
//			snapShot.Close()
//		}
//	}()
//	var readers []table.Reader
//	gCtx := query.NewGroupContext(len(tagKeyIDs))
//	for idx, tagKeyID := range tagKeyIDs {
//		readers, err = snapShot.FindReaders(tagKeyID)
//		if err != nil {
//			return nil, err
//		}
//		reader := newReader(readers)
//		tagValuesEntrySet := query.NewTagValuesEntrySet()
//		gCtx.SetTagValuesEntrySet(idx, tagValuesEntrySet)
//		err1 := reader.WalkTagValues(
//			tagKeyID,
//			"",
//			func(tagValue []byte, it invertedindex.TagValueIterator) bool {
//				for it.HasNext() {
//					if it.DataVersion() == version {
//						bitmapData := it.Next()
//						seriesBitmap := roaring.New()
//						err = seriesBitmap.UnmarshalBinary(bitmapData)
//						if err != nil {
//							return false
//						}
//						tagValuesEntrySet.AddTagValue(string(tagValue), seriesBitmap)
//					}
//				}
//				return true
//			})
//		if err1 != nil {
//			err = err1
//		}
//		if err != nil {
//			return nil, err
//		}
//	}
//	return gCtx, nil
//}
//
//// FindSeriesIDsByExpr finds series ids by tag filter expr for tag key id
//func (db *fileIndexDatabase) FindSeriesIDsByExpr(
//	tagKeyID uint32,
//	expr stmt.TagFilter,
//	timeRange timeutil.TimeRange,
//) (
//	*series.MultiVerSeriesIDSet,
//	error,
//) {
//	snapShot := db.invertedIndexFamily.GetSnapshot()
//	defer snapShot.Close()
//
//	readers, err := snapShot.FindReaders(tagKeyID)
//	if err != nil {
//		return nil, err
//	}
//	if len(readers) == 0 {
//		return nil, series.ErrNotFound
//	}
//	return invertedindex.NewReader(readers).FindSeriesIDsByExprForTagKeyID(tagKeyID, expr, timeRange)
//}
//
//// GetSeriesIDsForTag get series ids for spec metric's tag key
//func (db *fileIndexDatabase) GetSeriesIDsForTag(
//	tagKeyID uint32,
//	timeRange timeutil.TimeRange,
//) (
//	*series.MultiVerSeriesIDSet,
//	error,
//) {
//	snapShot := db.invertedIndexFamily.GetSnapshot()
//	defer snapShot.Close()
//
//	readers, err := snapShot.FindReaders(tagKeyID)
//	if err != nil {
//		return nil, err
//	}
//	return invertedindex.NewReader(readers).GetSeriesIDsForTagKeyID(tagKeyID, timeRange)
//}
