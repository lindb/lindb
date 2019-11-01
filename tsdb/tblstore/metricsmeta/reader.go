package metricsmeta

import (
	"strings"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package metricsmeta

// Reader reads metric meta info from the kv table
type Reader interface {
	// ReadTagKeyID read TagKeyID by metricID and tagKey
	ReadTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, ok bool)
	// ReadMaxFieldID return the max field-id of this metric, return 0 if not exist
	ReadMaxFieldID(metricID uint32) (maxFieldID uint16)
	// ReadFieldID read fieldID and fieldType from metricID and fieldName
	ReadFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type, ok bool)
	// SuggestTagKeys returns suggestion of tagKeys by prefix
	SuggestTagKeys(metricID uint32, tagKeyPrefix string, limit int) []string
}

// reader implements Reader
type reader struct {
	readers []table.Reader
	sr      *stream.Reader
}

// NewReader returns a new MetricsMetaReader
func NewReader(readers []table.Reader) Reader {
	return &reader{
		readers: readers,
		sr:      stream.NewReader(nil)}
}

// ReadTagKeyID read tagKeyID by metricID and tagKey
func (r *reader) ReadTagKeyID(
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
		itr := newTagMetaIterator(tagMetaBlock)
		for itr.HasNext() {
			tagMeta := itr.Next()
			if tagMeta.Key == tagKey {
				return tagMeta.ID, true
			}
		}
	}
	return 0, false
}

// readMetasBlock reads the tagMeta and FieldMeta blocks from binary by metricID
func (r *reader) readMetasBlock(
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
func (r *reader) ReadMaxFieldID(
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
func (r *reader) ReadFieldID(
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
func (r *reader) SuggestTagKeys(
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
		itr := newTagMetaIterator(tagMetaBlock)
		for itr.HasNext() {
			// read tagKey
			if limit <= len(collectedTagKeys) {
				return collectedTagKeys
			}
			tagMeta := itr.Next()
			if strings.HasPrefix(tagMeta.Key, tagKeyPrefix) {
				collectedTagKeys = append(collectedTagKeys, tagMeta.Key)
			}
		}
	}
	return collectedTagKeys
}

type tagMetaIterator struct {
	sr       *stream.Reader
	tagKey   string
	tagKeyID uint32
}

func newTagMetaIterator(block []byte) *tagMetaIterator {
	return &tagMetaIterator{sr: stream.NewReader(block)}
}
func (ti *tagMetaIterator) HasNext() bool {
	tagKeyLen := ti.sr.ReadByte()
	ti.tagKey = string(ti.sr.ReadSlice(int(tagKeyLen)))
	ti.tagKeyID = ti.sr.ReadUint32()
	return ti.sr.Error() == nil
}

func (ti *tagMetaIterator) Next() (
	tagMeta tag.Meta,
) {
	return tag.Meta{Key: ti.tagKey, ID: ti.tagKeyID}
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
