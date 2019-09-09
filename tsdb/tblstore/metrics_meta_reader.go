package tblstore

import (
	"strings"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/tsdb/field"
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
			if theTagKey == tagKey && theTagKeyID != 0 {
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
	if block == nil {
		return nil, nil
	}
	r.sr.Reset(block)
	// read length of tagMeta
	keyMetaLength := r.sr.ReadUvarint64()
	startOfTagMeta := r.sr.Position()
	// jump to end of tagMeta block
	r.sr.ShiftAt(uint32(keyMetaLength))
	endOfTagMeta := r.sr.Position()
	// block size too small
	if r.sr.Error() != nil {
		return nil, nil
	}
	tagMetaBlock = block[startOfTagMeta:endOfTagMeta]
	// read length of fieldMeta
	fieldMetaLen := r.sr.ReadUvarint64()
	startOfFieldMeta := r.sr.Position()
	r.sr.ShiftAt(uint32(fieldMetaLen))
	endOfFieldMeta := r.sr.Position()
	// failing assertion: the remaining block is field block
	if r.sr.Error() != nil || !r.sr.Empty() {
		return nil, nil
	}
	return tagMetaBlock, block[startOfFieldMeta:endOfFieldMeta]
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
	itr := newFieldIDIterator(fieldMetaBlock)
	for itr.HasNext() {
		_, _, fieldID := itr.Next()
		maxFieldID = fieldID
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
	var thisFieldName string
	for _, reader := range r.readers {
		_, fieldMetaBlock := r.readMetasBlock(reader.Get(metricID))
		if fieldMetaBlock == nil {
			continue
		}
		itr := newFieldIDIterator(fieldMetaBlock)
		for itr.HasNext() {
			thisFieldName, fieldType, fieldID = itr.Next()
			if thisFieldName == fieldName && fieldID != 0 && fieldType != 0 {
				ok = true
				return
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
	sr *stream.Reader
}

func newTagKeyIDIterator(block []byte) *tagKeyIDIterator {
	return &tagKeyIDIterator{sr: stream.NewReader(block)}
}
func (ti *tagKeyIDIterator) HasNext() bool { return !ti.sr.Empty() && ti.sr.Error() == nil }
func (ti *tagKeyIDIterator) Next() (
	tagKey string,
	tagKeyID uint32,
) {
	tagKeyLen := ti.sr.ReadByte()
	tagKey = string(ti.sr.ReadBytes(int(tagKeyLen)))
	tagKeyID = ti.sr.ReadUint32()
	return
}

type fieldIDIterator struct {
	sr *stream.Reader
}

func newFieldIDIterator(block []byte) *fieldIDIterator {
	return &fieldIDIterator{sr: stream.NewReader(block)}
}
func (fi *fieldIDIterator) HasNext() bool { return !fi.sr.Empty() && fi.sr.Error() == nil }
func (fi *fieldIDIterator) Next() (
	fieldName string,
	fieldType field.Type,
	fieldID uint16,
) {
	// read field-name
	thisFieldNameLen := fi.sr.ReadByte()
	fieldName = string(fi.sr.ReadBytes(int(thisFieldNameLen)))
	// read field-type
	fieldType = field.Type(fi.sr.ReadByte())
	// read field-ID
	fieldID = fi.sr.ReadUint16()
	return
}
