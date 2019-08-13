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
func (r *metricsNameIDReader) ReadMetricNS(nsID uint32) (data [][]byte, metricIDSeq, tagIDSeq uint32, ok bool) {
	for _, reader := range r.snapshot.Readers() {
		block := reader.Get(nsID)
		if len(block) < metricNameIDSequenceSize {
			continue
		}
		seqOffset := len(block) - metricNameIDSequenceSize
		data = append(data, block[:seqOffset])
		ok = true
		metricIDSeq = binary.BigEndian.Uint32(block[seqOffset : seqOffset+metricIDSize])
		tagIDSeq = binary.BigEndian.Uint32(block[seqOffset+metricIDSize:])
	}
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
func (r *metricsMetaReader) ReadTagID(metricID uint32, tagKey string) (tagID uint32, ok bool) {
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
				break
			}
			tagID = binary.BigEndian.Uint32(thisBinaryTagID)
			if thisTagKey == tagKey && tagID != 0 {
				return tagID, true
			}
			if sr.Error() != nil {
				break
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

// ReadMaxFieldID return the max field-id of this metric
func (r *metricsMetaReader) ReadMaxFieldID(metricID uint32) (maxFieldID uint16) {
	readers := r.snapshot.Readers()
	if len(readers) == 0 {
		return 0
	}
	_, fieldMeta := r.readMetasBlock(readers[len(readers)-1], metricID)
	if fieldMeta == nil {
		return 0
	}
	sr := stream.BinaryReader(fieldMeta)
	for !sr.Empty() {
		thisFieldNameLen := sr.ReadByte()
		// read field-name
		sr.ReadBytes(int(thisFieldNameLen))
		// read field-type
		sr.ReadByte()
		// read field-ID binary
		thisBinaryFieldID := sr.ReadBytes(fieldIDSize)
		// data corruption
		if len(thisBinaryFieldID) != fieldIDSize {
			break
		}
		maxFieldID = binary.BigEndian.Uint16(thisBinaryFieldID)
		if sr.Error() != nil {
			break
		}
	}
	return
}

// ReadFieldID read fieldID and fieldType from metricID and fieldName
func (r *metricsMetaReader) ReadFieldID(metricID uint32, fieldName string) (
	fieldID uint16, fieldType field.Type, ok bool) {

	for _, reader := range r.snapshot.Readers() {
		_, fieldMeta := r.readMetasBlock(reader, metricID)
		if fieldMeta == nil {
			continue
		}
		sr := stream.BinaryReader(fieldMeta)
		for !sr.Empty() {
			// read field-name
			thisFieldNameLen := sr.ReadByte()
			thisFieldName := string(sr.ReadBytes(int(thisFieldNameLen)))
			// read field-type
			fieldType = field.Type(sr.ReadByte())
			// read field-ID binary
			thisBinaryFieldID := sr.ReadBytes(fieldIDSize)
			// data corruption
			if len(thisBinaryFieldID) != fieldIDSize {
				break
			}
			fieldID = binary.BigEndian.Uint16(thisBinaryFieldID)
			if thisFieldName == fieldName && fieldID != 0 && fieldType != 0 {
				ok = true
				return
			}
			if sr.Error() != nil {
				break
			}
		}
	}
	return 0, field.Type(0), false
}
