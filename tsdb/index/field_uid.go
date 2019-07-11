package index

import (
	"bytes"
	"fmt"
	"sort"

	"go.uber.org/zap"

	"github.com/eleme/lindb/kv"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/stream"
	"github.com/eleme/lindb/pkg/util"
)

//FieldUID represents field unique under the metric name.
type FieldUID struct {
	metricID   uint32
	fieldMap   map[string]uint32 //key is field name
	sequenceID uint16            // sequence id
	family     kv.Family
	dbField    zap.Field
}

//FieldReader represents parses byte arrays for reading
type FieldReader struct {
	reader     *stream.ByteBufReader //byte buf reader
	seq        uint16                //field sequence id
	fieldCount int                   //field count
	fieldPos   int                   //field start position
}

//newFieldReader returns FieldReader
func newFieldReader(byteArray []byte) *FieldReader {
	bufReader := stream.NewBufReader(byteArray)
	sequence := uint16(bufReader.ReadUvarint64())
	count := int(bufReader.ReadUvarint64())

	return &FieldReader{
		reader:     bufReader,
		seq:        sequence,
		fieldCount: count,
		fieldPos:   bufReader.GetPosition(),
	}
}

//NewFieldUID creation requires kvFamily
func NewFieldUID(f kv.Family) *FieldUID {
	return &FieldUID{
		fieldMap: make(map[string]uint32),
		family:   f,
		dbField:  zap.String("db", "db"),
	}
}

//GetOrCreateFieldID  returns find the ID associated with a given field name and field type or create it.
func (f *FieldUID) GetOrCreateFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint32, error) {
	//get fieldID from disk
	fieldID := f.GetFieldID(metricID, fieldName)
	if NotFoundFieldID == fieldID {
		if f.metricID != metricID {
			err := f.Flush()
			if nil != err {
				logger.GetLogger("tsdb/index").Error("flush metric field error!", f.dbField, logger.Error(err))
				return NotFoundFieldID, err
			}
			f.metricID = metricID
			//clear all
			for k := range f.fieldMap {
				delete(f.fieldMap, k)
			}
			f.sequenceID = f.getLastFieldSequenceID(metricID)
		}

		f.sequenceID++
		f.fieldMap[fieldName] = util.ShortToInt(uint16(fieldType), f.sequenceID)
		return f.fieldMap[fieldName], nil
	}
	return fieldID, nil
}

//GetFields returns get all fields within the metric name
func (f *FieldUID) GetFields(metricID uint32, limit int16) map[string]struct{} {
	f.family.Lookup(metricID, func(bytes []byte) bool {
		fieldReader := newFieldReader(bytes)
		for i := 0; i < fieldReader.fieldCount; i++ {
			_, key := fieldReader.reader.ReadLenBytes()
			id := fieldReader.reader.ReadUvarint64()
			//todo
			fmt.Println("fieldName", key, " id:", id)
		}
		return true
	})
	return nil
}

//Flush represents forces a flush of in-memory data, and clear it
func (f *FieldUID) Flush() error {
	//flush fieldId to kv-store
	if len(f.fieldMap) > 0 {
		writer := stream.BinaryWriter()
		writer.PutUvarint64(uint64(f.sequenceID))
		writer.PutUvarint64(uint64(len(f.fieldMap)))

		fieldNames := f.getSortFieldNames()
		for _, fieldName := range fieldNames {
			writer.PutLenBytes([]byte(fieldName))
			writer.PutUvarint64(uint64(f.fieldMap[fieldName]))
		}

		by, err := writer.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode metric field error:", f.dbField, logger.Error(err))
			return err
		}

		//flusher
		flusher := f.family.NewFlusher()
		addError := flusher.Add(f.metricID, by)
		if nil != addError {
			logger.GetLogger("tsdb/index").Error("write metric field error!",
				f.dbField, zap.String("metricID", string(f.metricID)), logger.Error(addError))
			return addError
		}
		//commit
		commitError := flusher.Commit()
		if nil != commitError {
			logger.GetLogger("tsdb/index").Error("flush metric fieldId error!", f.dbField, logger.Error(commitError))
			return commitError
		}

		//clear in-memory data
		for k := range f.fieldMap {
			delete(f.fieldMap, k)
		}
	}
	return nil
}

//getLastFieldSequenceID returns last field sequence by metricID
func (f *FieldUID) getLastFieldSequenceID(metricID uint32) uint16 {
	seq := uint16(0)
	f.family.Lookup(metricID, func(bytes []byte) bool {
		fieldReader := newFieldReader(bytes)
		seq = fieldReader.seq
		return true
	})
	return seq
}

//GetFieldID returns get fieldID by fieldName within the metric name
func (f *FieldUID) GetFieldID(metricID uint32, field string) uint32 {
	var fieldID = NotFoundFieldID

	f.family.Lookup(metricID, func(byteArray []byte) bool {
		fieldReader := newFieldReader(byteArray)
		for i := 0; i < fieldReader.fieldCount; i++ {
			_, key := fieldReader.reader.ReadLenBytes()
			id := fieldReader.reader.ReadUvarint64()
			if bytes.Equal(key, []byte(field)) {
				fieldID = uint32(id)
				return true
			}
		}
		return false
	})
	return fieldID
}

//getSortFieldNames returns get the sorted field names
func (f *FieldUID) getSortFieldNames() []string {
	fieldNames := make([]string, len(f.fieldMap))
	var idx int
	for fieldName := range f.fieldMap {
		fieldNames[idx] = fieldName
		idx++
	}
	sort.Strings(fieldNames)
	return fieldNames
}
