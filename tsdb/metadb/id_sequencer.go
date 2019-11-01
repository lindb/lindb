package metadb

import (
	"math"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/tblstore"

	art "github.com/plar/go-adaptive-radix-tree"
	"go.uber.org/atomic"
)

const (
	// reserved for multi nameSpaces
	defaultNSID = 0
)

// todo: @codingcrush fix it, not singleton
var (
	once4IDSequencer     sync.Once
	idSequencerSingleton *idSequencer
)

// idSequencer implements IDSequencer
type idSequencer struct {
	metricIDSequence atomic.Uint32 // counter from 1
	tagKeyIDSequence atomic.Uint32 // counter from 1
	rwMux            sync.RWMutex  // readwrite lock for art-tree and map
	tree             art.Tree
	// unflushed generated id
	newNameIDs    map[string]uint32       // metricName -> metricID
	newTagMetas   map[uint32][]tag.Meta   // metricID -> tagKey + tagKeyID
	newFieldMetas map[uint32][]field.Meta // metricID -> fieldName + fieldType
	// family files for id-generating
	nameIDsFamily kv.Family
	metaFamily    kv.Family
}

// NewIDSequencer returns a new IDSequencer
func NewIDSequencer(nameIDsFamily, metaFamily kv.Family) IDSequencer {
	once4IDSequencer.Do(func() {
		idSequencerSingleton = &idSequencer{
			metricIDSequence: *atomic.NewUint32(0),
			tagKeyIDSequence: *atomic.NewUint32(0),
			tree:             art.New(),
			newNameIDs:       make(map[string]uint32),
			newTagMetas:      make(map[uint32][]tag.Meta),
			newFieldMetas:    make(map[uint32][]field.Meta),
			nameIDsFamily:    nameIDsFamily,
			metaFamily:       metaFamily}
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
	metricID, ok := seq.newNameIDs[metricName]
	if ok {
		return metricID
	}
	newMetricID := seq.metricIDSequence.Add(1)
	seq.newNameIDs[metricName] = newMetricID
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
	tagMetas, ok := seq.newTagMetas[metricID]
	newTagMeta := tag.Meta{ID: newTagKeyID, Key: tagKey}
	if ok {
		tagMetas = append(tagMetas, newTagMeta)
	} else {
		tagMetas = []tag.Meta{newTagMeta}
	}
	seq.newTagMetas[metricID] = tagMetas
	return newTagKeyID
}

func (seq *idSequencer) getTagKeyIDInMem(
	metricID uint32,
	tagKey string,
) (
	tagKeyID uint32,
	ok bool,
) {
	tagMetas, ok := seq.newTagMetas[metricID]
	if ok {
		for _, tagMeta := range tagMetas {
			if tagMeta.Key == tagKey {
				return tagMeta.ID, true
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
	fieldMetas, ok := seq.newFieldMetas[metricID]
	if ok {
		for _, fieldMeta := range fieldMetas {
			if fieldMeta.Name == fieldName {
				return fieldMeta.ID, fieldMeta.Type, true
			}
		}
	}
	return 0, 0, false
}

func (seq *idSequencer) getMaxFieldIDInMem(metricID uint32) uint16 {
	fieldMetas, ok := seq.newFieldMetas[metricID]
	if ok {
		return fieldMetas[len(fieldMetas)-1].ID
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

	metaList, ok := seq.newFieldMetas[metricID]
	newItem := field.Meta{ID: maxFieldID + 1, Type: fieldType, Name: fieldName}
	if ok {
		metaList = append(metaList, newItem)
	} else {
		metaList = []field.Meta{newItem}
	}
	seq.newFieldMetas[metricID] = metaList
	return newItem.ID, nil

}

// GetMetricID returns metric ID(uint32), if not exist return ErrMetaDataNotExist error
func (seq *idSequencer) GetMetricID(metricName string) (uint32, error) {
	seq.rwMux.RLock()
	defer seq.rwMux.RUnlock()
	// read memory
	metricID, ok := seq.newNameIDs[metricName]
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
	unflushed := seq.newNameIDs
	seq.newNameIDs = make(map[string]uint32)
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
	metricIDs := make(map[uint32]struct{})
	emptyTagMetas := make(map[uint32][]tag.Meta)
	emptyFieldMetas := make(map[uint32][]field.Meta)

	seq.rwMux.Lock()
	defer seq.rwMux.Unlock()
	// union of metricID
	for metricID := range seq.newTagMetas {
		metricIDs[metricID] = struct{}{}
	}
	for metricID := range seq.newFieldMetas {
		metricIDs[metricID] = struct{}{}
	}
	// flush process
	for metricID := range metricIDs {
		tagMetas, ok := seq.newTagMetas[metricID]
		if ok {
			for _, tagMeta := range tagMetas {
				flusher.FlushTagMeta(tagMeta)
			}
		}
		fieldMetas, ok := seq.newFieldMetas[metricID]
		if ok {
			for _, fieldMeta := range fieldMetas {
				flusher.FlushFieldMeta(fieldMeta)
			}
		}
		if err := flusher.FlushMetricMeta(metricID); err != nil {
			return err
		}
	}
	// replace it only on success
	seq.newTagMetas = emptyTagMetas
	seq.newFieldMetas = emptyFieldMetas
	return flusher.Commit()
}
