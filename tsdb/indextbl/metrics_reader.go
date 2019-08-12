package indextbl

import (
	"bytes"
	"encoding/binary"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./metrics_reader.go -destination=./metrics_reader_mock.go -package indextbl

const (
	metricNameIDSequenceSize = 4 + // metricID sequence
		4 // tagID sequence
	metricIDSize = 4 // uint32
	fieldIDSize  = 2 // uint16
	tagIDSize    = 4 // uint32
)

// MetricsNameIDReader reads metricNameID info from the kv table
type MetricsNameIDReader interface {
	// ReadMetricNS read metricNameID data by the namespace-id
	ReadMetricNS(nsID uint32) (data []byte, metricIDSeq, tagIDSeq uint32)
}

// MetricsMetaReader reads metric meta info from the kv table
type MetricsMetaReader interface {
	// ReadTagID read tagIDs by metricID and tagKey
	ReadTagID(metricID uint32, tagKey string) (tagID uint32)
	// ReadFieldID read fieldID and fieldType from metricID and fieldName
	ReadFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type)
}

// metricsNameIDReader implements MetricsNameIDReader
type metricsNameIDReader struct {
	snapshot kv.Snapshot
}

// NewMetricsNameIDReader returns a new MetricsNameIDReader
func NewMetricsNameIDReader(snapshot kv.Snapshot) MetricsNameIDReader {
	return &metricsNameIDReader{snapshot: snapshot}
}

// ReadMetricNS read metricNameID data by the namespace-id
func (r *metricsNameIDReader) ReadMetricNS(nsID uint32) (data []byte, metricIDSeq, tagIDSeq uint32) {
	var buffer bytes.Buffer
	for _, reader := range r.snapshot.Readers() {
		block := reader.Get(nsID)
		if len(block) < metricNameIDSequenceSize {
			continue
		}
		seqOffset := len(block) - metricNameIDSequenceSize
		buffer.Write(block[:seqOffset])
		metricIDSeq = binary.BigEndian.Uint32(block[seqOffset : seqOffset+metricIDSize])
		tagIDSeq = binary.BigEndian.Uint32(block[seqOffset+metricIDSize:])
	}
	data = buffer.Bytes()
	return
}

// metricsMetaReader implements MetricsMetaReader
type metricsMetaReader struct {
	snapshot kv.Snapshot
}

// NewMetricsMetaReader returns a new MetricsMetaReader
func NewMetricsMetaReader(snapshot kv.Snapshot) MetricsMetaReader {
	return &metricsMetaReader{snapshot: snapshot}
}

// ReadTagID read tagIDs by metricID and tagKey
func (r *metricsMetaReader) ReadTagID(metricID uint32, tagKey string) (tagID uint32) {
	for _, reader := range r.snapshot.Readers() {
		tagMeta, _ := r.readMetasBlock(reader, metricID)
		if tagMeta == nil {
			continue
		}
		sr := stream.BinaryReader(tagMeta)
		for !sr.Empty() {
			tagKeyLen := sr.ReadByte()
			thisTagKey := string(sr.ReadBytes(int(tagKeyLen)))
			thisBinaryTagID := sr.ReadBytes(tagIDSize)
			if len(thisBinaryTagID) != tagIDSize {
				continue
			}
			tagID = binary.BigEndian.Uint32(thisBinaryTagID)
			if thisTagKey == tagKey && tagID != 0 {
				return tagID
			}
			if sr.Error() != nil {
				break
			}
		}
	}
	return 0
}

// readTagFieldBlock reads the tagMeta and FieldMeta blocks from binary by metricID
func (r *metricsMetaReader) readMetasBlock(reader table.Reader, metricID uint32) (tagMeta []byte, fieldMeta []byte) {
	block := reader.Get(metricID)
	if block == nil {
		return nil, nil
	}
	// validation of tagMeta
	keyMetaLength, err := binary.ReadUvarint(bytes.NewBuffer(block))
	// read block failure
	if err != nil {
		return nil, nil
	}
	sizeOfKeyLen := bufioutil.GetUVariantLength(keyMetaLength)
	keyMetaEndPos := sizeOfKeyLen + int(keyMetaLength)
	// block size too small
	if len(block) < keyMetaEndPos {
		return nil, nil
	}
	tagMeta = block[sizeOfKeyLen:keyMetaEndPos]
	// validation of fieldMeta
	remainingBlock := block[keyMetaEndPos:]
	fieldLength, err := binary.ReadUvarint(bytes.NewBuffer(remainingBlock))
	if err != nil {
		return nil, nil
	}
	// failing assertion: the remaining block is field block
	sizeOfFieldLen := bufioutil.GetUVariantLength(fieldLength)
	if len(remainingBlock) != sizeOfFieldLen+int(fieldLength) {
		return nil, nil
	}
	return tagMeta, remainingBlock[sizeOfFieldLen:]
}

// ReadFieldID read fieldID and fieldType from metricID and fieldName
func (r *metricsMetaReader) ReadFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type) {
	for _, reader := range r.snapshot.Readers() {
		_, fieldMeta := r.readMetasBlock(reader, metricID)
		if fieldMeta == nil {
			continue
		}
		sr := stream.BinaryReader(fieldMeta)
		for !sr.Empty() {
			thisFieldNameLen := sr.ReadByte()
			thisFieldName := string(sr.ReadBytes(int(thisFieldNameLen)))
			fieldType = field.Type(sr.ReadByte())
			thisBinaryFieldID := sr.ReadBytes(fieldIDSize)
			if len(thisBinaryFieldID) != fieldIDSize {
				continue
			}
			fieldID = binary.BigEndian.Uint16(thisBinaryFieldID)
			if thisFieldName == fieldName && fieldID != 0 && fieldType != 0 {
				return
			}
			if sr.Error() != nil {
				break
			}
		}
	}
	return 0, field.Type(0)
}
