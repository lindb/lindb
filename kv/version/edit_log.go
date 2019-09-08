package version

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
)

// StoreFamilyID is store level edit log,
// actually store family is not actual family just store store level edit log for metadata.
const StoreFamilyID = -99999999

// EditLog contains all metadata edit log
type EditLog struct {
	logs     []Log
	familyID int
	logger   *logger.Logger
}

// NewEditLog new EditLog instance
func NewEditLog(familyID int) *EditLog {
	return &EditLog{
		familyID: familyID,
		logger:   logger.GetLogger("kv", "EditLog"),
	}
}

// GetLogs return the logs under edit log
func (el *EditLog) GetLogs() []Log {
	return el.logs
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
	sw := stream.NewBufferWriter(nil)
	// write family id
	sw.PutVarint32(int32(el.familyID))
	// write num of logs
	sw.PutUvarint64(uint64(len(el.logs)))
	// write detail log data
	for _, log := range el.logs {
		logType := logTypes[reflect.TypeOf(log)]
		sw.PutVarint32(logType)
		value, err := log.Encode()
		if err != nil {
			return nil, fmt.Errorf("edit logs encode error: %s", err)
		}
		sw.PutUvarint32(uint32(len(value))) // write log bytes length
		sw.PutBytes(value)                  // write log bytes data
	}
	return sw.Bytes()
}

// unmarshal create an edit log from its serialized in buf
func (el *EditLog) unmarshal(buf []byte) error {
	reader := stream.NewReader(buf)
	el.familyID = int(reader.ReadVarint32())
	// read num of logs
	count := reader.ReadUvarint64()
	// read detail log data
	for ; count > 0; count-- {
		logType := reader.ReadVarint32()
		fn, ok := newLogFuncMap[logType]
		if !ok {
			return fmt.Errorf("cannot get log type new func, type is:[%d]", logType)
		}
		l := fn()
		length := int(reader.ReadUvarint32())
		logData := reader.ReadBytes(length)
		if err := l.Decode(logData); err != nil {
			return fmt.Errorf("unmarshal log data error, type is:[%d],error:%s", logType, err)
		}
		el.Add(l)
	}
	return reader.Error()
}

// apply family edit logs into version metadata
func (el *EditLog) apply(version *Version) {
	for _, log := range el.logs {
		log.apply(version)

		if v, ok := log.(StoreLog); ok {
			// if log is store log, need to apply version set
			v.applyVersionSet(version.fv.GetVersionSet())
		}
	}
}

// apply store edit logs into version set
func (el *EditLog) applyVersionSet(versionSet StoreVersionSet) {
	for _, log := range el.logs {
		switch v := log.(type) {
		case StoreLog:
			v.applyVersionSet(versionSet)
		default:
			el.logger.Warn("cannot apply family edit log to version set")
		}
	}
}
