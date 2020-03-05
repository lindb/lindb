package version

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source=./edit_log.go -destination=./edit_log_mock.go -package=version

// StoreFamilyID is store level edit log,
// actually store family is not actual family just store store level edit log for metadata.
const StoreFamilyID = -99999999

// EditLog represents the version metadata edit log
type EditLog interface {
	// FamilyID return family id
	FamilyID() FamilyID
	// Add adds edit log into log list
	Add(log Log)
	// GetLogs return the logs under edit log
	GetLogs() []Log
	// IsEmpty returns edit logs is empty or not.
	IsEmpty() bool
	// marshal encodes edit log to binary data
	marshal() ([]byte, error)
	// unmarshal create an edit log from its serialized in buf
	unmarshal(buf []byte) error
	// apply applies family edit logs into version metadata
	apply(version Version)
	// applyVersionSet applies store edit logs into version set
	applyVersionSet(versionSet StoreVersionSet)
}

// editLog contains all metadata edit log
type editLog struct {
	logs     []Log
	familyID FamilyID
}

// NewEditLog new editLog instance
func NewEditLog(familyID FamilyID) EditLog {
	return &editLog{
		familyID: familyID,
	}
}

// newEmptyEditLog create empty edit log without family id for unmarshal
func newEmptyEditLog() EditLog {
	return &editLog{}
}

// FamilyID return family id
func (el *editLog) FamilyID() FamilyID {
	return el.familyID
}

// GetLogs return the logs under edit log
func (el *editLog) GetLogs() []Log {
	return el.logs
}

// Add adds edit log into log list
func (el *editLog) Add(log Log) {
	el.logs = append(el.logs, log)
}

// IsEmpty returns edit logs is empty or not.
func (el *editLog) IsEmpty() bool {
	return len(el.logs) == 0
}

// marshal encodes edit log to binary data
func (el *editLog) marshal() ([]byte, error) {
	sw := stream.NewBufferWriter(nil)
	// write family id
	sw.PutVarint32(int32(el.familyID))
	// write num of logs
	sw.PutUvarint64(uint64(len(el.logs)))
	// write detail log data
	for _, log := range el.logs {
		logType := logTypes[reflect.TypeOf(log)]
		sw.PutVarint32(int32(logType))
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
func (el *editLog) unmarshal(buf []byte) error {
	reader := stream.NewReader(buf)
	el.familyID = FamilyID(reader.ReadVarint32())
	// read num of logs
	count := reader.ReadUvarint64()
	// read detail log data
	for ; count > 0; count-- {
		logType := reader.ReadVarint32()
		fn, ok := newLogFuncMap[LogType(logType)]
		if !ok {
			return fmt.Errorf("cannot get log type new func, type is:[%d]", logType)
		}
		l := fn()
		length := int(reader.ReadUvarint32())
		logData := reader.ReadSlice(length)
		if err := l.Decode(logData); err != nil {
			return fmt.Errorf("unmarshal log data error, type is:[%d],error:%s", logType, err)
		}
		el.Add(l)
	}
	return reader.Error()
}

// apply applies family edit logs into version metadata
func (el *editLog) apply(version Version) {
	for _, log := range el.logs {
		log.apply(version)

		if v, ok := log.(StoreLog); ok {
			// if log is store log, need to apply version set
			v.applyVersionSet(version.GetFamilyVersion().GetVersionSet())
		}
	}
}

// applyVersionSet applies store edit logs into version set
func (el *editLog) applyVersionSet(versionSet StoreVersionSet) {
	for _, log := range el.logs {
		switch v := log.(type) {
		case StoreLog:
			v.applyVersionSet(versionSet)
		default:
			versionLogger.Warn("cannot apply family edit log to version set")
		}
	}
}
