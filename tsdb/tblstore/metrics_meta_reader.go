package tblstore

import (
	"strings"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./metrics_meta_reader.go -destination=./metrics_meta_reader_mock.go -package tblstore

const (
	metricNameIDSequenceSize = 4 + // metricID sequence
		4 // tagID sequence
)

// MetricsNameIDReader reads metricNameID info from the kv table
type MetricsNameIDReader interface {
	// ReadMetricNS read metricNameID data by the namespace-id
	ReadMetricNS(nsID uint32) (data [][]byte, metricIDSeq, tagIDSeq uint32, ok bool)
}

// MetricsMetaReader reads metric meta info from the kv table
type MetricsMetaReader interface {
	// ReadTagID read tagIDs by metricID and tagKey
	ReadTagID(metricID uint32, tagKey string) (tagID uint32, ok bool)
	// ReadMaxFieldID return the max field-id of this metric, return 0 if not exist
	ReadMaxFieldID(metricID uint32) (maxFieldID uint16)
	// ReadFieldID read fieldID and fieldType from metricID and fieldName
	ReadFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type, ok bool)
	// SuggestTagKeys returns suggestion of tagKeys by prefix
	SuggestTagKeys(metricID uint32, tagKeyPrefix string, limit int) []string
}

// metricsNameIDReader implements MetricsNameIDReader
type metricsNameIDReader struct {
	snapshot version.Snapshot
}

// NewMetricsNameIDReader returns a new MetricsNameIDReader
func NewMetricsNameIDReader(snapshot version.Snapshot) MetricsNameIDReader {
	return &metricsNameIDReader{snapshot: snapshot}
}

// ReadMetricNS read metricNameID data by the namespace-id
func (r *metricsNameIDReader) ReadMetricNS(nsID uint32) (data [][]byte, metricIDSeq, tagIDSeq uint32, ok bool) {
	readers, err := r.snapshot.FindReaders(nsID)
	if err != nil {
		//todo need check
		return
	}
	for _, reader := range readers {
		block := reader.Get(nsID)
		if len(block) < metricNameIDSequenceSize {
			continue
		}
		sr := stream.NewReader(block)
		sr.ShiftAt(uint32(len(block) - metricNameIDSequenceSize))
		ok = true
		metricIDSeq = sr.ReadUint32()
		tagIDSeq = sr.ReadUint32()
	}
	return
}

// metricsMetaReader implements MetricsMetaReader
type metricsMetaReader struct {
	snapshot version.Snapshot
}

// NewMetricsMetaReader returns a new MetricsMetaReader
func NewMetricsMetaReader(snapshot version.Snapshot) MetricsMetaReader {
	return &metricsMetaReader{snapshot: snapshot}
}

// ReadTagID read tagIDs by metricID and tagKey
func (r *metricsMetaReader) ReadTagID(metricID uint32, tagKey string) (tagID uint32, ok bool) {
	readers, err := r.snapshot.FindReaders(metricID)
	if err != nil {
		//todo need check
		return 0, false
	}
	for _, reader := range readers {
		tagMeta, _ := r.readMetasBlock(reader, metricID)
		if tagMeta == nil {
			continue
		}
		sr := stream.NewReader(tagMeta)
		for !sr.Empty() && sr.Error() == nil {
			tagKeyLen := sr.ReadByte()
			thisTagKey := string(sr.ReadBytes(int(tagKeyLen)))
			tagID = sr.ReadUint32()
			if thisTagKey == tagKey && tagID != 0 {
				return tagID, true
			}
		}
	}
	return 0, false
}

// readTagFieldBlock reads the tagMeta and FieldMeta blocks from binary by metricID
func (r *metricsMetaReader) readMetasBlock(reader table.Reader, metricID uint32) (tagMeta []byte, fieldMeta []byte) {
	block := reader.Get(metricID)
	if block == nil {
		return nil, nil
	}
	sr := stream.NewReader(block)

	// read length of tagMeta
	keyMetaLength := sr.ReadUvarint64()
	startOfTagMeta := sr.Position()
	// jump to end of tagMeta block
	sr.ShiftAt(uint32(keyMetaLength))
	endOfTagMeta := sr.Position()
	// block size too small
	if sr.Error() != nil {
		return nil, nil
	}
	tagMeta = block[startOfTagMeta:endOfTagMeta]
	// read length of fieldMeta
	fieldMetaLen := sr.ReadUvarint64()
	startOfFieldMeta := sr.Position()
	sr.ShiftAt(uint32(fieldMetaLen))
	endOfFieldMeta := sr.Position()
	// failing assertion: the remaining block is field block
	if sr.Error() != nil || !sr.Empty() {
		return nil, nil
	}
	return tagMeta, block[startOfFieldMeta:endOfFieldMeta]
}

// ReadMaxFieldID return the max field-id of this metric
func (r *metricsMetaReader) ReadMaxFieldID(metricID uint32) (maxFieldID uint16) {
	readers, err := r.snapshot.FindReaders(metricID)
	//TODO need check
	if err != nil || len(readers) == 0 {
		return 0
	}
	_, fieldMeta := r.readMetasBlock(readers[len(readers)-1], metricID)
	if fieldMeta == nil {
		return 0
	}
	sr := stream.NewReader(fieldMeta)
	for !sr.Empty() {
		thisFieldNameLen := sr.ReadByte()
		// read field-name
		sr.ReadBytes(int(thisFieldNameLen))
		// read field-type
		sr.ReadByte()
		thisFieldID := sr.ReadUint16()
		if sr.Error() != nil {
			break
		}
		maxFieldID = thisFieldID
	}
	return
}

// ReadFieldID read fieldID and fieldType from metricID and fieldName
func (r *metricsMetaReader) ReadFieldID(metricID uint32, fieldName string) (
	fieldID uint16, fieldType field.Type, ok bool) {
	readers, err := r.snapshot.FindReaders(metricID)
	if err != nil {
		//TODO need check
		return fieldID, fieldType, ok
	}

	for _, reader := range readers {
		_, fieldMeta := r.readMetasBlock(reader, metricID)
		if fieldMeta == nil {
			continue
		}
		sr := stream.NewReader(fieldMeta)
		for !sr.Empty() && sr.Error() == nil {
			// read field-name
			thisFieldNameLen := sr.ReadByte()
			thisFieldName := string(sr.ReadBytes(int(thisFieldNameLen)))
			// read field-type
			fieldType = field.Type(sr.ReadByte())
			// data corruption
			fieldID = sr.ReadUint16()
			if thisFieldName == fieldName && fieldID != 0 && fieldType != 0 {
				ok = true
				return
			}
		}
	}
	return 0, field.Type(0), false
}

// SuggestTagKeys returns suggestion of tagKeys by prefix
func (r *metricsMetaReader) SuggestTagKeys(metricID uint32, tagKeyPrefix string, limit int) []string {
	var collectedTagKeys []string
	readers, err := r.snapshot.FindReaders(metricID)
	if err != nil {
		return nil
	}
	for _, reader := range readers {
		tagMeta, _ := r.readMetasBlock(reader, metricID)
		if tagMeta == nil {
			continue
		}
		sr := stream.NewReader(tagMeta)
		for !sr.Empty() && sr.Error() == nil {
			// read tagKey
			if limit <= len(collectedTagKeys) {
				return collectedTagKeys
			}
			tagKeyLen := sr.ReadByte()
			thisTagKey := string(sr.ReadBytes(int(tagKeyLen)))
			// readTagID
			_ = sr.ReadUint32()
			if strings.HasPrefix(thisTagKey, tagKeyPrefix) {
				collectedTagKeys = append(collectedTagKeys, thisTagKey)
			}
		}
	}
	return collectedTagKeys
}
