package indexdb

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/indextbl"
	"github.com/lindb/lindb/tsdb/series"

	art "github.com/plar/go-adaptive-radix-tree"
)

const (
	// reserved for multi nameSpaces
	defaultNSID = 0
)

var (
	indexDBLogger   = logger.GetLogger("tsdb/indexdb")
	once4IndexDb    sync.Once
	indexDBInstance *indexDatabase
)

// tagKeyAndID holds the relation of tagKey and its ID
type tagKeyAndID struct {
	tagKey   string
	tagKeyID uint32
}

// fieldIDAndType holds the relation of field-id and field-type
type fieldIDAndType struct {
	fieldName string
	fieldID   uint16
	fieldType field.Type
}

// todo: @codingcrush, flush

// indexDatabase implements IndexDatabase
type indexDatabase struct {
	recoverError     error        // last error during recovery
	metricIDSequence uint32       // counter from 1
	tagKeyIDSequence uint32       // counter from 1
	rwMux            sync.RWMutex // readwrite lock for art-tree and map
	tree             artTreeINTF
	// unflushed generated id
	youngMetricNameIDs map[string]uint32           // metricName -> metricID
	youngTagKeyIDs     map[uint32][]tagKeyAndID    // metricID -> tagKey + tagKeyID
	youngFieldIDs      map[uint32][]fieldIDAndType // metricID -> fieldName + fieldType
	// index reaader
	nameIDsReader indextbl.MetricsNameIDReader
	metaReader    indextbl.MetricsMetaReader
	seriesReader  indextbl.SeriesIndexReader
}

// NewIndexDatabase returns a new IndexDatabase
func NewIndexDatabase(nameIDIndexSnapShot kv.Snapshot, metaIndexSnapShot kv.Snapshot,
	seriesIndexSnapShot kv.Snapshot) (IndexDatabase, error) {

	once4IndexDb.Do(func() {
		indexDBInstance = &indexDatabase{
			tree:               newArtTree(),
			youngMetricNameIDs: make(map[string]uint32),
			youngTagKeyIDs:     make(map[uint32][]tagKeyAndID),
			youngFieldIDs:      make(map[uint32][]fieldIDAndType),
			nameIDsReader:      indextbl.NewMetricsNameIDReader(nameIDIndexSnapShot),
			metaReader:         indextbl.NewMetricsMetaReader(metaIndexSnapShot),
			seriesReader:       indextbl.NewSeriesIndexReader(seriesIndexSnapShot)}
		indexDBInstance.recover()
	})
	return indexDBInstance, indexDBInstance.recoverError
}

// recover loads metric-names and metricIDs from the index file and build the tree
func (db *indexDatabase) recover() {
	db.rwMux.Lock()
	defer db.rwMux.Unlock()

	var err error
	data, metricIDSeq, tagIDSeq, ok := db.nameIDsReader.ReadMetricNS(defaultNSID)
	if ok {
		db.metricIDSequence = metricIDSeq
		db.tagKeyIDSequence = tagIDSeq
		for _, d := range data {
			err = db.tree.UnmarshalBinary(d)
			if db.recoverError == nil {
				db.recoverError = err
			}
		}
	}
}

// GenMetricID generates ID(uint32) from metricName
func (db *indexDatabase) GenMetricID(metricName string) uint32 {
	metricID, err := db.GetMetricID(metricName)
	if err == nil {
		return metricID
	}
	db.rwMux.Lock()
	defer db.rwMux.Unlock()
	// double check
	metricID, ok := db.youngMetricNameIDs[metricName]
	if ok {
		return metricID
	}
	newMetricID := atomic.AddUint32(&db.metricIDSequence, 1)
	db.youngMetricNameIDs[metricName] = newMetricID
	return newMetricID
}

// GenTagID generates tagID(uint32) from metricName and tagKey
func (db *indexDatabase) GenTagID(metricID uint32, tagKey string) uint32 {
	// case1: tagKeyID exist in memory
	db.rwMux.RLock()
	tagKeyID, ok := db.getTagIDInMem(metricID, tagKey)
	if ok {
		db.rwMux.RUnlock()
		return tagKeyID
	}
	// case2: tagKeyID exist on disk
	tagKeyID, ok = db.metaReader.ReadTagID(metricID, tagKey)
	if ok {
		db.rwMux.RUnlock()
		return tagKeyID
	}
	db.rwMux.RUnlock()
	// case3: double check
	db.rwMux.Lock()
	defer db.rwMux.Unlock()
	tagKeyID, ok = db.getTagIDInMem(metricID, tagKey)
	if ok {
		return tagKeyID
	}
	// case4: tagKeyID not exist, create a new one
	newTagKeyID := atomic.AddUint32(&db.tagKeyIDSequence, 1)
	tagKeyAndIDList, ok := db.youngTagKeyIDs[metricID]
	newItem := tagKeyAndID{tagKeyID: newTagKeyID, tagKey: tagKey}
	if ok {
		tagKeyAndIDList = append(tagKeyAndIDList, newItem)
	} else {
		tagKeyAndIDList = []tagKeyAndID{newItem}
	}
	db.youngTagKeyIDs[metricID] = tagKeyAndIDList
	return newTagKeyID
}

func (db *indexDatabase) getTagIDInMem(metricID uint32, tagKey string) (uint32, bool) {
	tagKeyAndIDList, ok := db.youngTagKeyIDs[metricID]
	if ok {
		for _, tagKeyWithIDItem := range tagKeyAndIDList {
			if tagKeyWithIDItem.tagKey == tagKey {
				return tagKeyWithIDItem.tagKeyID, true
			}
		}
	}
	return 0, false
}

func (db *indexDatabase) getFieldIDInMem(metricID uint32, fieldName string) (uint16, field.Type, bool) {
	fieldIDAndTypeList, ok := db.youngFieldIDs[metricID]
	if ok {
		for _, item := range fieldIDAndTypeList {
			if item.fieldName == fieldName {
				return item.fieldID, item.fieldType, true
			}
		}
	}
	return 0, 0, false
}

func (db *indexDatabase) getMaxFieldIDInMem(metricID uint32) uint16 {
	fieldIDAndTypeList, ok := db.youngFieldIDs[metricID]
	if ok {
		return fieldIDAndTypeList[len(fieldIDAndTypeList)-1].fieldID
	}
	return 0
}

// GenFieldID returns field ID(uint16), return error when fields-count exceed the limitation
// or type mis-match with the old one
func (db *indexDatabase) GenFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint16, error) {
	// find from memory
	db.rwMux.RLock()
	fID, fType, ok := db.getFieldIDInMem(metricID, fieldName)
	if ok {
		defer db.rwMux.RUnlock()
		if fType == fieldType {
			return fID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// find from disk
	fID, fType, err := db.GetFieldID(metricID, fieldName)
	db.rwMux.RUnlock()
	if err == nil {
		if fType == fieldType {
			return fID, nil
		}
		return 0, series.ErrWrongFieldType
	}

	db.rwMux.Lock()
	defer db.rwMux.Unlock()
	// double check
	_, _, ok = db.getFieldIDInMem(metricID, fieldName)
	if ok {
		return db.GenFieldID(metricID, fieldName, fieldType)
	}
	// create a new fieldID
	maxFieldIDOnDisk := db.metaReader.ReadMaxFieldID(metricID)
	maxFieldIDInMem := db.getMaxFieldIDInMem(metricID)
	maxFieldID := uint16(math.Max(float64(maxFieldIDInMem), float64(maxFieldIDOnDisk)))

	if maxFieldID >= constants.TStoreMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}

	fieldIDAndTypeList, ok := db.youngFieldIDs[metricID]
	newItem := fieldIDAndType{fieldID: maxFieldID + 1, fieldType: fieldType, fieldName: fieldName}
	if ok {
		fieldIDAndTypeList = append(fieldIDAndTypeList, newItem)
	} else {
		fieldIDAndTypeList = []fieldIDAndType{newItem}
	}
	db.youngFieldIDs[metricID] = fieldIDAndTypeList
	return newItem.fieldID, nil
}

// GetMetricID returns metric ID(uint32), if not exist return ErrMetaDataNotExist error
func (db *indexDatabase) GetMetricID(metricName string) (uint32, error) {
	db.rwMux.RLock()
	defer db.rwMux.RUnlock()
	// read memory
	metricID, ok := db.youngMetricNameIDs[metricName]
	if ok {
		return metricID, nil
	}
	val, ok := db.tree.Search(art.Key(metricName))
	if ok {
		return val.(uint32), nil
	}
	return 0, series.ErrMetaDataNotExist
}

// GetFieldID returns field ID(uint16), if not exist return ErrMetaDataNotExist error
func (db *indexDatabase) GetFieldID(metricID uint32, fieldName string) (
	fieldID uint16, fieldType field.Type, err error) {

	var ok bool
	fieldID, fieldType, ok = db.metaReader.ReadFieldID(metricID, fieldName)
	if !ok {
		return 0, 0, series.ErrMetaDataNotExist
	}
	return fieldID, fieldType, nil
}

// GetTagValues get tag values corresponding with the tagKeys
func (db *indexDatabase) GetTagValues(metricID uint32, tagKeys []string, version int64) (
	tagValues [][]string, err error) {
	return db.seriesReader.GetTagValues(metricID, tagKeys, version)
}

// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
func (db *indexDatabase) FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	return db.seriesReader.FindSeriesIDsByExpr(metricID, expr, timeRange)
}

// GetSeriesIDsForTag get series ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsForTag(metricID uint32, tagKey string,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	return db.seriesReader.GetSeriesIDsForTag(metricID, tagKey, timeRange)
}
