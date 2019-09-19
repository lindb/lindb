package tblstore

import (
	"strings"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./metrics_meta_reader.go -destination=./metrics_meta_reader_mock.go -package tblstore

// MetricsMetaReader reads metric meta info from the kv table
type MetricsMetaReader interface {
	// ReadTagKeyID read TagKeyID by metricID and tagKey
	ReadTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, ok bool)
	// ReadMaxFieldID return the max field-id of this metric, return 0 if not exist
	ReadMaxFieldID(metricID uint32) (maxFieldID uint16)
	// ReadFieldID read fieldID and fieldType from metricID and fieldName
	ReadFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type, ok bool)
	// SuggestTagKeys returns suggestion of tagKeys by prefix
	SuggestTagKeys(metricID uint32, tagKeyPrefix string, limit int) []string
}

// metricsMetaReader implements MetricsMetaReader
type metricsMetaReader struct {
	readers []table.Reader
	sr      *stream.Reader
}

// NewMetricsMetaReader returns a new MetricsMetaReader
func NewMetricsMetaReader(readers []table.Reader) MetricsMetaReader {
	return &metricsMetaReader{
		readers: readers,
		sr:      stream.NewReader(nil)}
}

// ReadTagKeyID read tagKeyID by metricID and tagKey
func (r *metricsMetaReader) ReadTagKeyID(
	metricID uint32,
	tagKey string,
) (
	tagKeyID uint32,
	ok bool,
) {
	for _, reader := range r.readers {
		tagMetaBlock, _ := r.readMetasBlock(reader.Get(metricID))
		if tagMetaBlock == nil {
			continue
		}
		itr := newTagKeyIDIterator(tagMetaBlock)
		for itr.HasNext() {
			theTagKey, theTagKeyID := itr.Next()
			if theTagKey == tagKey {
				return theTagKeyID, true
			}
		}
	}
	return 0, false
}

// readMetasBlock reads the tagMeta and FieldMeta blocks from binary by metricID
func (r *metricsMetaReader) readMetasBlock(
	block []byte,
) (
	tagMetaBlock []byte,
	fieldMetaBlock []byte,
) {
	if len(block) <= 4 { // posOfFieldMeta
		return nil, nil
	}
	// read pos of field-meta
	r.sr.Reset(block)
	r.sr.ReadSlice(len(block) - 4)
	posOfFieldMetaPos := int(r.sr.ReadUint32())
	// read tag-meta and field-meta
	r.sr.SeekStart()
	tagMetaBlock = r.sr.ReadSlice(posOfFieldMetaPos)
	fieldMetaBlock = r.sr.ReadSlice(len(block) - posOfFieldMetaPos - 4)
	// failing assertion: the remaining block is field block
	_ = r.sr.ReadSlice(4)
	if r.sr.Error() != nil || !r.sr.Empty() {
		return nil, nil
	}
	return tagMetaBlock, fieldMetaBlock
}

// ReadMaxFieldID return the max field-id of this metric
func (r *metricsMetaReader) ReadMaxFieldID(
	metricID uint32,
) (maxFieldID uint16) {
	if len(r.readers) == 0 {
		return 0
	}
	_, fieldMetaBlock := r.readMetasBlock(r.readers[len(r.readers)-1].Get(metricID))
	if fieldMetaBlock == nil {
		return 0
	}
	itr := newFieldMetaIterator(fieldMetaBlock)
	for itr.HasNext() {
		meta := itr.Next()
		maxFieldID = meta.ID
	}
	return
}

// ReadFieldID read fieldID and fieldType from metricID and fieldName
func (r *metricsMetaReader) ReadFieldID(
	metricID uint32,
	fieldName string,
) (
	fieldID uint16,
	fieldType field.Type,
	ok bool,
) {
	for _, reader := range r.readers {
		_, fieldMetaBlock := r.readMetasBlock(reader.Get(metricID))
		if fieldMetaBlock == nil {
			continue
		}
		itr := newFieldMetaIterator(fieldMetaBlock)
		for itr.HasNext() {
			fieldMeta := itr.Next()
			if fieldMeta.Name == fieldName {
				return fieldMeta.ID, fieldMeta.Type, true
			}
		}
	}
	return 0, field.Type(0), false
}

// SuggestTagKeys returns suggestion of tagKeys by prefix
func (r *metricsMetaReader) SuggestTagKeys(
	metricID uint32,
	tagKeyPrefix string,
	limit int,
) []string {
	var collectedTagKeys []string
	for _, reader := range r.readers {
		tagMetaBlock, _ := r.readMetasBlock(reader.Get(metricID))
		if tagMetaBlock == nil {
			continue
		}
		itr := newTagKeyIDIterator(tagMetaBlock)
		for itr.HasNext() {
			// read tagKey
			if limit <= len(collectedTagKeys) {
				return collectedTagKeys
			}
			theTagKey, _ := itr.Next()
			if strings.HasPrefix(theTagKey, tagKeyPrefix) {
				collectedTagKeys = append(collectedTagKeys, theTagKey)
			}
		}
	}
	return collectedTagKeys
}

type tagKeyIDIterator struct {
	sr       *stream.Reader
	tagKey   string
	tagKeyID uint32
}

func newTagKeyIDIterator(block []byte) *tagKeyIDIterator {
	return &tagKeyIDIterator{sr: stream.NewReader(block)}
}
func (ti *tagKeyIDIterator) HasNext() bool {
	tagKeyLen := ti.sr.ReadByte()
	ti.tagKey = string(ti.sr.ReadSlice(int(tagKeyLen)))
	ti.tagKeyID = ti.sr.ReadUint32()
	return ti.sr.Error() == nil
}

func (ti *tagKeyIDIterator) Next() (
	tagKey string,
	tagKeyID uint32,
) {
	return ti.tagKey, ti.tagKeyID
}

type fieldMetaIterator struct {
	sr   *stream.Reader
	meta field.Meta
}

func newFieldMetaIterator(block []byte) *fieldMetaIterator {
	return &fieldMetaIterator{sr: stream.NewReader(block)}
}
func (fi *fieldMetaIterator) HasNext() bool {
	var meta field.Meta
	// read field-ID
	meta.ID = fi.sr.ReadUint16()
	// read field-type
	meta.Type = field.Type(fi.sr.ReadByte())
	// read field-name
	fieldNameLen := fi.sr.ReadUvarint64()
	meta.Name = string(fi.sr.ReadSlice(int(fieldNameLen)))
	fi.meta = meta
	return fi.sr.Error() == nil
}

func (fi *fieldMetaIterator) Next() (
	meta field.Meta,
) {
	return fi.meta
}
