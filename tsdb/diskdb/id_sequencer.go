package diskdb

import (
	"math"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore"

	art "github.com/plar/go-adaptive-radix-tree"
	"go.uber.org/atomic"
)

const (
	// reserved for multi nameSpaces
	defaultNSID = 0
)

var (
	once4IDSequencer     sync.Once
	idSequencerSingleton *idSequencer
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

// idSequencer implements IDSequencer
type idSequencer struct {
	metricIDSequence *atomic.Uint32 // counter from 1
	tagKeyIDSequence *atomic.Uint32 // counter from 1
	rwMux            sync.RWMutex   // readwrite lock for art-tree and map
	tree             art.Tree
	// unflushed generated id
	youngMetricNameIDs map[string]uint32           // metricName -> metricID
	youngTagKeyIDs     map[uint32][]tagKeyAndID    // metricID -> tagKey + tagKeyID
	youngFieldIDs      map[uint32][]fieldIDAndType // metricID -> fieldName + fieldType
	// family files for id-generating
	nameIDsFamily kv.Family
	metaFamily    kv.Family
}

// NewIDSequencer returns a new IDSequencer
func NewIDSequencer(nameIDsFamily, metaFamily kv.Family) IDSequencer {
	once4IDSequencer.Do(func() {
		idSequencerSingleton = &idSequencer{
			metricIDSequence:   atomic.NewUint32(0),
			tagKeyIDSequence:   atomic.NewUint32(0),
			tree:               art.New(),
			youngMetricNameIDs: make(map[string]uint32),
			youngTagKeyIDs:     make(map[uint32][]tagKeyAndID),
			youngFieldIDs:      make(map[uint32][]fieldIDAndType),
			nameIDsFamily:      nameIDsFamily,
			metaFamily:         metaFamily}
	})
	return idSequencerSingleton
}

// Recover loads metric-names and metricIDs from the index file to build the tree
func (seq *idSequencer) Recover() error {
	snapShot := seq.nameIDsFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(defaultNSID)
	if err != nil {
		return err
	}
	seq.rwMux.Lock()
	defer seq.rwMux.Unlock()

	nameIDReader := tblstore.NewMetricsNameIDReader(readers)
	data, metricIDSeq, tagKeyIDSeq, ok := nameIDReader.ReadMetricNS(defaultNSID)
	if ok {
		seq.metricIDSequence.Store(metricIDSeq)
		seq.tagKeyIDSequence.Store(tagKeyIDSeq)
		for _, d := range data {
			if err = nameIDReader.UnmarshalBinaryToART(seq.tree, d); err != nil {
				return err
			}
		}
	}
	return nil
}

// SuggestMetrics returns suggestions of metricNames from a given prefix,
func (seq *idSequencer) SuggestMetrics(prefix string, limit int) (suggestions []string) {
	if limit <= 0 {
		return nil
	}
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	suggestions = make([]string, 128)[:0]

	seq.rwMux.RLock()
	defer seq.rwMux.RUnlock()

	seq.tree.ForEachPrefix(art.Key(prefix), func(node art.Node) (cont bool) {
		if len(suggestions) >= limit {
			return false
		}
		suggestions = append(suggestions, string(node.Key()))
		return true
	})
	return suggestions
}

// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
func (seq *idSequencer) SuggestTagKeys(metricName, tagKeyPrefix string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	metricID, err := seq.GetMetricID(metricName)
	if err != nil {
		return nil
	}
	snapShot := seq.metaFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(defaultNSID)
	if err != nil {
		return nil
	}
	metaReader := tblstore.NewMetricsMetaReader(readers)
	return metaReader.SuggestTagKeys(metricID, tagKeyPrefix, limit)
}

// GenMetricID generates ID(uint32) from metricName
func (seq *idSequencer) GenMetricID(metricName string) uint32 {
	metricID, err := seq.GetMetricID(metricName)
	if err == nil {
		return metricID
	}
	seq.rwMux.Lock()
	defer seq.rwMux.Unlock()
	// double check
	metricID, ok := seq.youngMetricNameIDs[metricName]
	if ok {
		return metricID
	}
	newMetricID := seq.metricIDSequence.Add(1)
	seq.youngMetricNameIDs[metricName] = newMetricID
	return newMetricID
}

// GenTagKeyID generates tagKeyID(uint32) from metricName and tagKey
func (seq *idSequencer) GenTagKeyID(
	metricID uint32,
	tagKey string,
) (tagKeyID uint32) {
	// case1, 2: check if it is in memory or on disk
	tagKeyID, err := seq.GetTagKeyID(metricID, tagKey)
	if err == nil {
		return tagKeyID
	}
	// case3: double check
	seq.rwMux.Lock()
	defer seq.rwMux.Unlock()
	tagKeyID, ok := seq.getTagKeyIDInMem(metricID, tagKey)
	if ok {
		return tagKeyID
	}
	// case4: tagKeyID not exist, create a new one
	newTagKeyID := seq.tagKeyIDSequence.Add(1)
	tagKeyAndIDList, ok := seq.youngTagKeyIDs[metricID]
	newItem := tagKeyAndID{tagKeyID: newTagKeyID, tagKey: tagKey}
	if ok {
		tagKeyAndIDList = append(tagKeyAndIDList, newItem)
	} else {
		tagKeyAndIDList = []tagKeyAndID{newItem}
	}
	seq.youngTagKeyIDs[metricID] = tagKeyAndIDList
	return newTagKeyID
}

func (seq *idSequencer) getTagKeyIDInMem(
	metricID uint32,
	tagKey string,
) (
	tagKeyID uint32,
	ok bool,
) {
	tagKeyAndIDList, ok := seq.youngTagKeyIDs[metricID]
	if ok {
		for _, tagKeyWithIDItem := range tagKeyAndIDList {
			if tagKeyWithIDItem.tagKey == tagKey {
				return tagKeyWithIDItem.tagKeyID, true
			}
		}
	}
	return 0, false
}

// GetTagKeyID returns tag ID(uint32), return ErrNotFound if not exist
func (seq *idSequencer) GetTagKeyID(metricID uint32, tagKey string) (tagID uint32, err error) {
	// case1: tagKeyID exist in memory
	seq.rwMux.RLock()
	defer seq.rwMux.RUnlock()
	tagKeyID, ok := seq.getTagKeyIDInMem(metricID, tagKey)
	if ok {
		return tagKeyID, nil
	}
	// case2: tagKeyID exist on disk
	snapShot := seq.metaFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(metricID)
	if err != nil {
		return 0, err
	}
	return seq.readTagKeyID(tblstore.NewMetricsMetaReader(readers), metricID, tagKey)
}

// readTagKeyID reads the tagKeyID from reader.
func (seq *idSequencer) readTagKeyID(
	reader tblstore.MetricsMetaReader,
	metricID uint32,
	tagKey string,
) (
	tagKeyID uint32,
	err error,
) {
	tagKeyID, ok := reader.ReadTagKeyID(metricID, tagKey)
	if ok {
		return tagKeyID, nil
	}
	return 0, series.ErrNotFound
}

func (seq *idSequencer) getFieldIDInMem(
	metricID uint32,
	fieldName string,
) (
	uint16,
	field.Type,
	bool,
) {
	fieldIDAndTypeList, ok := seq.youngFieldIDs[metricID]
	if ok {
		for _, item := range fieldIDAndTypeList {
			if item.fieldName == fieldName {
				return item.fieldID, item.fieldType, true
			}
		}
	}
	return 0, 0, false
}

func (seq *idSequencer) getMaxFieldIDInMem(metricID uint32) uint16 {
	fieldIDAndTypeList, ok := seq.youngFieldIDs[metricID]
	if ok {
		return fieldIDAndTypeList[len(fieldIDAndTypeList)-1].fieldID
	}
	return 0
}

// GenFieldID returns field ID(uint16), return error when fields-count exceed the limitation
// or type mis-match with the old one
func (seq *idSequencer) GenFieldID(
	metricID uint32,
	fieldName string,
	fieldType field.Type,
) (
	uint16,
	error,
) {
	// find from memory
	seq.rwMux.RLock()
	fID, fType, ok := seq.getFieldIDInMem(metricID, fieldName)
	if ok {
		defer seq.rwMux.RUnlock()
		if fType == fieldType {
			return fID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	seq.rwMux.RUnlock()

	snapShot := seq.metaFamily.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(metricID)
	if err != nil {
		return 0, err
	}
	metaReader := tblstore.NewMetricsMetaReader(readers)
	return seq.genFieldID(metaReader, metricID, fieldName, fieldType)
}

// genFieldID generate fieldID from reader.
func (seq *idSequencer) genFieldID(
	reader tblstore.MetricsMetaReader,
	metricID uint32,
	fieldName string,
	fieldType field.Type,
) (
	uint16,
	error,
) {

	// find from disk
	fID, fType, err := seq.readFieldID(reader, metricID, fieldName)
	if err == nil {
		if fType == fieldType {
			return fID, nil
		}
		return 0, series.ErrWrongFieldType
	}

	seq.rwMux.Lock()
	// double check
	_, _, ok := seq.getFieldIDInMem(metricID, fieldName)
	if ok {
		seq.rwMux.Unlock()
		return seq.GenFieldID(metricID, fieldName, fieldType)
	}
	defer seq.rwMux.Unlock()
	// create a new fieldID
	maxFieldIDOnDisk := reader.ReadMaxFieldID(metricID)
	maxFieldIDInMem := seq.getMaxFieldIDInMem(metricID)
	maxFieldID := uint16(math.Max(float64(maxFieldIDInMem), float64(maxFieldIDOnDisk)))

	if maxFieldID >= constants.TStoreMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}

	fieldIDAndTypeList, ok := seq.youngFieldIDs[metricID]
	newItem := fieldIDAndType{fieldID: maxFieldID + 1, fieldType: fieldType, fieldName: fieldName}
	if ok {
		fieldIDAndTypeList = append(fieldIDAndTypeList, newItem)
	} else {
		fieldIDAndTypeList = []fieldIDAndType{newItem}
	}
	seq.youngFieldIDs[metricID] = fieldIDAndTypeList
	return newItem.fieldID, nil

}

// GetMetricID returns metric ID(uint32), if not exist return ErrMetaDataNotExist error
func (seq *idSequencer) GetMetricID(metricName string) (uint32, error) {
	seq.rwMux.RLock()
	defer seq.rwMux.RUnlock()
	// read memory
	metricID, ok := seq.youngMetricNameIDs[metricName]
	if ok {
		return metricID, nil
	}
	val, ok := seq.tree.Search(art.Key(metricName))
	if ok {
		return val.(uint32), nil
	}
	return 0, series.ErrNotFound
}

// GetFieldID returns field ID(uint16), if not exist return ErrMetaDataNotExist error
func (seq *idSequencer) GetFieldID(
	metricID uint32,
	fieldName string,
) (
	fieldID uint16,
	fieldType field.Type,
	err error,
) {

	seq.rwMux.RLock()
	fID, fType, ok := seq.getFieldIDInMem(metricID, fieldName)
	if ok {
		seq.rwMux.RUnlock()
		return fID, fType, nil
	}
	seq.rwMux.RUnlock()

	snapShot := seq.metaFamily.GetSnapshot()
	defer snapShot.Close()
	readers, err := snapShot.FindReaders(metricID)
	if err != nil {
		return 0, 0, err
	}
	return seq.readFieldID(tblstore.NewMetricsMetaReader(readers), metricID, fieldName)
}

// readFieldID read fieldID from the reader
func (seq *idSequencer) readFieldID(
	reader tblstore.MetricsMetaReader,
	metricID uint32,
	fieldName string,
) (
	fieldID uint16,
	fieldType field.Type,
	err error,
) {
	var ok bool
	fieldID, fieldType, ok = reader.ReadFieldID(metricID, fieldName)
	if !ok {
		return 0, 0, series.ErrNotFound
	}
	return fieldID, fieldType, nil
}

// FlushNameIDs flushes metricName and metricID to family
func (seq *idSequencer) FlushNameIDs() error {
	kvFlusher := seq.nameIDsFamily.NewFlusher()
	return seq.flushNameIDsTo(tblstore.NewMetricsNameIDFlusher(kvFlusher))
}

// flushNameIDsTo flushes metricName and metricID to flusher
func (seq *idSequencer) flushNameIDsTo(flusher tblstore.MetricsNameIDFlusher) error {
	seq.rwMux.Lock()
	unflushed := seq.youngMetricNameIDs
	seq.youngMetricNameIDs = make(map[string]uint32)
	for metricName, metricID := range unflushed {
		seq.tree.Insert([]byte(metricName), metricID)
	}
	seq.rwMux.Unlock()

	for metricName, metricID := range unflushed {
		flusher.FlushNameID(metricName, metricID)
	}
	if err := flusher.FlushMetricsNS(defaultNSID,
		seq.metricIDSequence.Load(),
		seq.tagKeyIDSequence.Load()); err != nil {
		return err
	}
	return flusher.Commit()
}

// FlushMetricsMeta flushes tagKey, tagKeyId, fieldName, fieldID to family
func (seq *idSequencer) FlushMetricsMeta() error {
	kvFlusher := seq.metaFamily.NewFlusher()
	return seq.flushMetricsMetaTo(tblstore.NewMetricsMetaFlusher(kvFlusher))
}

// flushMetricsMetaTo flushes tagKey, tagKeyId, fieldName, fieldID to flusher
func (seq *idSequencer) flushMetricsMetaTo(flusher tblstore.MetricsMetaFlusher) error {
	seq.rwMux.Lock()
	unflushedTagKeys := seq.youngTagKeyIDs
	unflushedFields := seq.youngFieldIDs
	seq.youngTagKeyIDs = make(map[uint32][]tagKeyAndID)
	seq.youngFieldIDs = make(map[uint32][]fieldIDAndType)
	seq.rwMux.Unlock()

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
	return flusher.Commit()
}
