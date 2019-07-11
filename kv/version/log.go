package version

import (
	"fmt"
	"reflect"

	strm "github.com/eleme/lindb/pkg/stream"
)

func init() {
	// register new file
	RegisterLogType(1, func() Log {
		return &NewFile{}
	})
	// register delete file
	RegisterLogType(2, func() Log {
		return &DeleteFile{}
	})
	// register new file number
	RegisterLogType(3, func() Log {
		return &NextFileNumber{}
	})
}

// NewLogFunc create specific edit log instance
type NewLogFunc func() Log

var newLogFuncMap = make(map[int32]NewLogFunc)
var logTypes = make(map[reflect.Type]int32)

// RegisterLogType register edit log type when system init,
// if has duplicate log type, system need panic and exit.
func RegisterLogType(logType int32, fn NewLogFunc) {
	if _, ok := newLogFuncMap[logType]; ok {
		panic(fmt.Sprintf("log type already registered: %d", logType))
	}
	newLogFuncMap[logType] = fn

	// register log type
	log := fn()
	logTypes[reflect.TypeOf(log)] = logType
}

// Log metadata edit log for family level
type Log interface {
	// Encode write log from binary, if error return err
	Encode() ([]byte, error)
	// Decode reads log from binary, if error return err
	Decode(v []byte) error
	// apply edit log to family's current version
	apply(version *Version)
}

// StoreLog metadata dit log store level
type StoreLog interface {
	Log
	// applyVersionSet apply edito to store version set
	applyVersionSet(versionSet *StoreVersionSet)
}

// NewFile add new file into metadata
type NewFile struct {
	level int32
	file  *FileMeta
}

// CreateNewFile new NewFile instance for add new file
func CreateNewFile(level int32, file *FileMeta) *NewFile {
	return &NewFile{
		level: level,
		file:  file,
	}
}

// Encode writes new file data to binary, if error return err
func (n *NewFile) Encode() ([]byte, error) {
	var stream = strm.BinaryWriter()

	stream.PutVarint32(n.level)                // level
	stream.PutVarint64(n.file.GetFileNumber()) // file number
	stream.PutUvarint32(n.file.GetMinKey())    // min key
	stream.PutUvarint32(n.file.GetMaxKey())    // max key
	stream.PutVarint32(n.file.GetFileSize())   // file size

	return stream.Bytes()
}

// Decode reads new file from binary, if error return err
func (n *NewFile) Decode(v []byte) error {
	var stream = strm.BinaryReader(v)
	// read level
	n.level = stream.ReadVarint32()
	// read file meta
	n.file = NewFileMeta(stream.ReadVarint64(), stream.ReadUvarint32(), stream.ReadUvarint32(), stream.ReadVarint32())
	// if error, return it
	return stream.Error()
}

// Apply new file edit log to version
func (n *NewFile) apply(version *Version) {
	version.addFile(int(n.level), n.file)
}

// DeleteFile remove file from metadata
type DeleteFile struct {
	level      int32
	fileNumber int64
}

// NewDeleteFile create DeleteFile instance
func NewDeleteFile(level int32, fileNumber int64) *DeleteFile {
	return &DeleteFile{
		level:      level,
		fileNumber: fileNumber,
	}
}

// Encode writes delete file data into binary
func (d *DeleteFile) Encode() ([]byte, error) {
	var stream = strm.BinaryWriter()

	stream.PutVarint32(d.level)
	stream.PutVarint64(d.fileNumber)

	return stream.Bytes()
}

// Decode reads delete file data from binary
func (d *DeleteFile) Decode(v []byte) error {
	var stream = strm.BinaryReader(v)

	d.level = stream.ReadVarint32()
	d.fileNumber = stream.ReadVarint64()

	return stream.Error()
}

// Apply removes file from version
func (d *DeleteFile) apply(version *Version) {
	version.deleteFile(int(d.level), d.fileNumber)
}

// NextFileNumber set next file number for metadata
type NextFileNumber struct {
	fileNumber int64
}

// NewNextFileNumber creates NextFileNumber instance
func NewNextFileNumber(fileNumber int64) *NextFileNumber {
	return &NextFileNumber{
		fileNumber: fileNumber,
	}
}

// Encode writes next file number data into binary
func (n *NextFileNumber) Encode() ([]byte, error) {
	var stream = strm.BinaryWriter()
	stream.PutVarint64(n.fileNumber)
	return stream.Bytes()
}

// Decode reads next file number data from binary
func (n *NextFileNumber) Decode(v []byte) error {
	var stream = strm.BinaryReader(v)
	n.fileNumber = stream.ReadVarint64()
	return stream.Error()
}

// Apply do nothing for next file number
func (n *NextFileNumber) apply(version *Version) {
	// do nothing
}

//applyVersionSet applies edito to store version set
func (n *NextFileNumber) applyVersionSet(versionSet *StoreVersionSet) {
	versionSet.setNextFileNumberWithoutLock(n.fileNumber)
}
