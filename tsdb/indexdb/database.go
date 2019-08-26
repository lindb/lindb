package indexdb

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/series"
	"github.com/lindb/lindb/tsdb/tblstore"

	art "github.com/plar/go-adaptive-radix-tree"
)

const (
	// reserved for multi nameSpaces
	defaultNSID = 0
)

var (
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

// indexDatabase implements IndexDatabase
type indexDatabase struct {
	metricIDSequence uint32       // counter from 1
	tagKeyIDSequence uint32       // counter from 1
	rwMux            sync.RWMutex // readwrite lock for art-tree and map
	tree             artTreeINTF
	// unflushed generated id
	youngMetricNameIDs map[string]uint32           // metricName -> metricID
	youngTagKeyIDs     map[uint32][]tagKeyAndID    // metricID -> tagKey + tagKeyID
	youngFieldIDs      map[uint32][]fieldIDAndType // metricID -> fieldName + fieldType
	// index reader
	metaReader          tblstore.MetricsMetaReader
	invertedIndexReader tblstore.InvertedIndexReader
	forwardIndexReader  tblstore.ForwardIndexReader
}

// NewIndexDatabase returns a new IndexDatabase
func NewIndexDatabase(metaIndexSnapShot kv.Snapshot, seriesIndexSnapShot kv.Snapshot) IndexDatabase {
	once4IndexDb.Do(func() {
		indexDBInstance = &indexDatabase{
			tree:                newArtTree(),
			youngMetricNameIDs:  make(map[string]uint32),
			youngTagKeyIDs:      make(map[uint32][]tagKeyAndID),
			youngFieldIDs:       make(map[uint32][]fieldIDAndType),
			metaReader:          tblstore.NewMetricsMetaReader(metaIndexSnapShot),
			invertedIndexReader: tblstore.NewInvertedIndexReader(seriesIndexSnapShot)}
	})
	return indexDBInstance
}

// Recover loads metric-names and metricIDs from the index file to build the tree
func (db *indexDatabase) Recover(nameIDsReader tblstore.MetricsNameIDReader) error {
	db.rwMux.Lock()
	defer db.rwMux.Unlock()

	var err error
	data, metricIDSeq, tagIDSeq, ok := nameIDsReader.ReadMetricNS(defaultNSID)
	if ok {
		db.metricIDSequence = metricIDSeq
		db.tagKeyIDSequence = tagIDSeq
		for _, d := range data {
			if err = db.tree.UnmarshalBinary(d); err != nil {
				return err
			}
		}
	}
	return nil
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
	// case1, 2: check if it is in memory or on disk
	tagID, err := db.GetTagID(metricID, tagKey)
	if err == nil {
		return tagID
	}
	// case3: double check
	db.rwMux.Lock()
	defer db.rwMux.Unlock()
	tagKeyID, ok := db.getTagIDInMem(metricID, tagKey)
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

// GetTagID returns tag ID(uint32), return ErrNotFound if not exist
func (db *indexDatabase) GetTagID(metricID uint32, tagKey string) (tagID uint32, err error) {
	// case1: tagKeyID exist in memory
	db.rwMux.RLock()
	defer db.rwMux.RUnlock()
	tagKeyID, ok := db.getTagIDInMem(metricID, tagKey)
	if ok {
		return tagKeyID, nil
	}
	// case2: tagKeyID exist on disk
	tagKeyID, ok = db.metaReader.ReadTagID(metricID, tagKey)
	if ok {
		return tagKeyID, nil
	}
	return 0, series.ErrNotFound
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
	return 0, series.ErrNotFound
}

// GetFieldID returns field ID(uint16), if not exist return ErrMetaDataNotExist error
func (db *indexDatabase) GetFieldID(metricID uint32, fieldName string) (
	fieldID uint16, fieldType field.Type, err error) {

	var ok bool
	fieldID, fieldType, ok = db.metaReader.ReadFieldID(metricID, fieldName)
	if !ok {
		return 0, 0, series.ErrNotFound
	}
	return fieldID, fieldType, nil
}

// GetTagValues get tag values corresponding with the tagKeys
func (db *indexDatabase) GetTagValues(metricID uint32, tagKeys []string, version uint32) (
	tagValues [][]string, err error) {
	return db.forwardIndexReader.GetTagValues(metricID, tagKeys, version)
}

// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
func (db *indexDatabase) FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	tagID, err := db.GetTagID(metricID, expr.TagKey())
	if err != nil {
		return nil, err
	}
	return db.invertedIndexReader.FindSeriesIDsByExprForTagID(tagID, expr, timeRange)
}

// GetSeriesIDsForTag get series ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsForTag(metricID uint32, tagKey string,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	tagID, err := db.GetTagID(metricID, tagKey)
	if err != nil {
		return nil, err
	}
	return db.invertedIndexReader.GetSeriesIDsForTagID(tagID, timeRange)
}

// FlushNameIDsTo flushes metricName and metricID to flusher
func (db *indexDatabase) FlushNameIDsTo(flusher tblstore.MetricsNameIDFlusher) error {
	db.rwMux.Lock()
	unflushed := db.youngMetricNameIDs
	db.youngMetricNameIDs = make(map[string]uint32)
	for metricName, metricID := range unflushed {
		db.tree.Insert([]byte(metricName), metricID)
	}
	db.rwMux.Unlock()

	compressor := newNameIDCompressor()
	for metricName, metricID := range unflushed {
		compressor.AddNameID(metricName, metricID)
	}
	data, err := compressor.Close()
	if err != nil {
		return err
	}
	return flusher.FlushMetricsNS(defaultNSID, data,
		atomic.LoadUint32(&db.metricIDSequence),
		atomic.LoadUint32(&db.tagKeyIDSequence))
}

// FlushMetricsMetaTo flushes tagKey, tagKeyId, fieldName, fieldID to flusher
func (db *indexDatabase) FlushMetricsMetaTo(flusher tblstore.MetricsMetaFlusher) error {
	db.rwMux.Lock()
	unflushedTagKeys := db.youngTagKeyIDs
	unflushedFields := db.youngFieldIDs
	db.youngTagKeyIDs = make(map[uint32][]tagKeyAndID)
	db.youngFieldIDs = make(map[uint32][]fieldIDAndType)
	db.rwMux.Unlock()

	// union of metricID
	metricIDs := make(map[uint32]struct{})
	for metricID := range unflushedTagKeys {
		metricIDs[metricID] = struct{}{}
	}
	for metricID := range unflushedFields {
		metricIDs[metricID] = struct{}{}
	}
	// flush process
	for metricID := range metricIDs {
		items1, ok := unflushedTagKeys[metricID]
		if ok {
			for _, item := range items1 {
				flusher.FlushTagKeyID(item.tagKey, item.tagKeyID)
			}
		}
		items2, ok := unflushedFields[metricID]
		if ok {
			for _, item := range items2 {
				flusher.FlushFieldID(item.fieldName, item.fieldType, item.fieldID)
			}
		}
		if err := flusher.FlushMetricMeta(metricID); err != nil {
			return err
		}
	}
	return nil
}
