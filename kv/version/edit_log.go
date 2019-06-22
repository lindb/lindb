package version

import (
	"fmt"
	"reflect"

	"github.com/eleme/lindb/pkg/logger"
	strm "github.com/eleme/lindb/pkg/stream"
)

// StoreFamilyID is store level edit log,
// actually store family is not actual family just store store level edit log for metadata.
const StoreFamilyID = 0

// EditLog contains all metadata edit log
type EditLog struct {
	logs     []Log
	familyID int
}

// NewEditLog new EditLog instance
func NewEditLog(familyID int) *EditLog {
	return &EditLog{
		familyID: familyID,
	}
}

// Add adds edit log into log list
func (el *EditLog) Add(log Log) {
	el.logs = append(el.logs, log)
}

// IsEmpty returns edit logs is empty or not.
func (el *EditLog) IsEmpty() bool {
	return len(el.logs) == 0
}

// marshal encodes edit log to binary data
func (el *EditLog) marshal() ([]byte, error) {
	stream := strm.BinaryWriter()
	// write family id
	stream.PutInt32(int32(el.familyID))
	// write num of logs
	stream.PutUvarint64(uint64(len(el.logs)))
	// write detail log data
	for _, log := range el.logs {
		logType := logTypes[reflect.TypeOf(log)]
		stream.PutInt32(logType)
		value, err := log.Encode()
		if err != nil {
			return nil, fmt.Errorf("edit logs encode error: %s", err)
		}
		stream.PutUvarint32(uint32(len(value))) // write log bytes length
		stream.PutBytes(value)                  // write log bytes data
	}
	return stream.Bytes()
}

// unmarshal create an edit log from its seriealized in buf
func (el *EditLog) unmarshal(buf []byte) error {
	stream := strm.BinaryReader(buf)
	el.familyID = int(stream.ReadInt32())
	// read num of logs
	count := stream.ReadUvarint64()
	// read detail log data
	for ; count > 0; count-- {
		logType := stream.ReadInt32()
		fn, ok := newLogFuncMap[logType]
		if !ok {
			return fmt.Errorf("cannot get log type new func, type is:[%d]", logType)
		}
		l := fn()
		length := int(stream.ReadUvarint32())
		logData := stream.ReadBytes(length)
		if err := l.Decode(logData); err != nil {
			return fmt.Errorf("unmarshal log data error, type is:[%d],error:%s", logType, err)
		}
		el.Add(l)
	}
	return stream.Error()
}

// apply family edit logs into version metadata
func (el *EditLog) apply(version *Version) {
	for _, log := range el.logs {
		log.apply(version)

		if v, ok := log.(StoreLog); ok {
			// if log is store log, need to apply version set
			v.applyVersionSet(version.fv.versionSet)
		}
	}
}

// apply store edit logs into version set
func (el *EditLog) applyVersionSet(versionSet *StoreVersionSet) {
	l := logger.GetLogger()
	for _, log := range el.logs {
		switch v := log.(type) {
		case StoreLog:
			v.applyVersionSet(versionSet)
		default:
			l.Warn("cannot apply family edit log to version set")
		}
	}
}
