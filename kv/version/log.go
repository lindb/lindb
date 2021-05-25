// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package version

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./log.go -destination=./log_mock.go -package=version

// LogType represents the version edit log type
type LogType int

// Defines all version edit log types,
// Rollup/Reference file is spec version edit log for metric data rollup with difference target time interval.
const (
	NewFileLog LogType = iota + 1
	DeleteFileLog
	NextFileNumberLog
	NewRollupFileLog
	DeleteRollupFileLog
	NewReferenceFileLog
	DeleteReferenceFileLog
)

func init() {
	// register new file
	RegisterLogType(NewFileLog, func() Log {
		return &newFile{}
	})
	// register delete file
	RegisterLogType(DeleteFileLog, func() Log {
		return &deleteFile{}
	})
	// register new file number
	RegisterLogType(NextFileNumberLog, func() Log {
		return &nextFileNumber{}
	})
	// register new rollup file
	RegisterLogType(NewRollupFileLog, func() Log {
		return &newRollupFile{}
	})
	// register delete rollup file
	RegisterLogType(DeleteRollupFileLog, func() Log {
		return &deleteRollupFile{}
	})
	// register reference file
	RegisterLogType(NewReferenceFileLog, func() Log {
		return &newReferenceFile{}
	})
	// register delete reference
	RegisterLogType(DeleteReferenceFileLog, func() Log {
		return &deleteReferenceFile{}
	})
}

// NewLogFunc creates specific edit log instance
type NewLogFunc func() Log

// Stores log create function => log type mapping
var newLogFuncMap = make(map[LogType]NewLogFunc)
var logTypes = make(map[reflect.Type]LogType)

// RegisterLogType registers edit log type when system init,
// if has duplicate log type, system need panic and exit.
func RegisterLogType(logType LogType, fn NewLogFunc) {
	if _, ok := newLogFuncMap[logType]; ok {
		panic(fmt.Sprintf("log type already registered: %d", logType))
	}
	newLogFuncMap[logType] = fn

	// register log type
	log := fn()
	logTypes[reflect.TypeOf(log)] = logType
}

// Log represents metadata edit log for family level
type Log interface {
	// Encode writes log from binary, if error return err
	Encode() ([]byte, error)
	// Decode reads log from binary, if error return err
	Decode(v []byte) error
	// apply applies edit log to family's current version
	apply(version Version)
}

// StoreLog represents metadata edit log store level
type StoreLog interface {
	Log
	// applyVersionSet applies edit to store version set
	applyVersionSet(versionSet StoreVersionSet)
}

// newFile represents version edit log for adding new file into metadata
type newFile struct {
	level int32
	file  *FileMeta
}

// CreateNewFile creates NewFile instance for add new file
func CreateNewFile(level int32, file *FileMeta) Log {
	return &newFile{
		level: level,
		file:  file,
	}
}

// Encode writes new file data to binary, if error return err
func (n *newFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint32(n.level)                        // level
	writer.PutVarint64(n.file.GetFileNumber().Int64()) // file number
	writer.PutUvarint32(n.file.GetMinKey())            // min key
	writer.PutUvarint32(n.file.GetMaxKey())            // max key
	writer.PutVarint32(n.file.GetFileSize())           // file size
	return writer.Bytes()
}

// Decode reads new file from binary, if error return err
func (n *newFile) Decode(v []byte) error {
	reader := stream.NewReader(v)
	// read level
	n.level = reader.ReadVarint32()
	// read file meta
	n.file = NewFileMeta(table.FileNumber(reader.ReadVarint64()),
		reader.ReadUvarint32(), reader.ReadUvarint32(), reader.ReadVarint32())
	// if error, return it
	return reader.Error()
}

// String returns string value of new file log
func (n *newFile) String() string {
	return fmt.Sprintf("addFile:{level:%d,file:%s}", n.level, n.file)
}

// apply applies new file edit log to version
func (n *newFile) apply(version Version) {
	version.AddFile(int(n.level), n.file)
}

// deleteFile represents version edit log for deleting file from metadata
type deleteFile struct {
	level      int32
	fileNumber table.FileNumber
}

// NewDeleteFile creates DeleteFile instance
func NewDeleteFile(level int32, fileNumber table.FileNumber) Log {
	return &deleteFile{
		level:      level,
		fileNumber: fileNumber,
	}
}

// Encode writes delete file data into binary
func (d *deleteFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint32(d.level)
	writer.PutVarint64(d.fileNumber.Int64())
	return writer.Bytes()
}

// Decode reads delete file data from binary
func (d *deleteFile) Decode(v []byte) error {
	reader := stream.NewReader(v)

	d.level = reader.ReadVarint32()
	d.fileNumber = table.FileNumber(reader.ReadVarint64())

	return reader.Error()
}

// String returns string value of delete file log
func (d *deleteFile) String() string {
	return fmt.Sprintf("deleteFile:{level:%d,fileNumber:%d}", d.level, d.fileNumber)
}

// apply applies remove file from version
func (d *deleteFile) apply(version Version) {
	version.DeleteFile(int(d.level), d.fileNumber)
}

// nextFileNumber represent version edit log for next file number for metadata
type nextFileNumber struct {
	fileNumber table.FileNumber
}

// NewNextFileNumber creates NextFileNumber instance
func NewNextFileNumber(fileNumber table.FileNumber) Log {
	return &nextFileNumber{
		fileNumber: fileNumber,
	}
}

// Encode writes next file number data into binary
func (n *nextFileNumber) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint64(n.fileNumber.Int64())
	return writer.Bytes()
}

// Decode reads next file number data from binary
func (n *nextFileNumber) Decode(v []byte) error {
	reader := stream.NewReader(v)

	n.fileNumber = table.FileNumber(reader.ReadVarint64())
	return reader.Error()
}

// apply does nothing for next file number
func (n *nextFileNumber) apply(version Version) {
	// do nothing
}

// String returns string value of file number log
func (n *nextFileNumber) String() string {
	return fmt.Sprintf("fileNumber:%d", n.fileNumber)
}

// applyVersionSet applies edit to store version set
func (n *nextFileNumber) applyVersionSet(versionSet StoreVersionSet) {
	versionSet.setNextFileNumberWithoutLock(n.fileNumber)
}

// newRollupFile represent version edit log for new rollup file for rollup job
type newRollupFile struct {
	fileNumber table.FileNumber  // file number
	interval   timeutil.Interval // target time interval
}

// CreateNewRollupFile creates a new rollup file
func CreateNewRollupFile(fileNumber table.FileNumber, interval timeutil.Interval) Log {
	return &newRollupFile{
		fileNumber: fileNumber,
		interval:   interval,
	}
}

// Encode writes new rollup file data into binary
func (n *newRollupFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint64(n.fileNumber.Int64())
	writer.PutVarint64(n.interval.Int64())
	return writer.Bytes()
}

// Decode reads new rollup file data from binary
func (n *newRollupFile) Decode(v []byte) error {
	reader := stream.NewReader(v)
	n.fileNumber = table.FileNumber(reader.ReadVarint64())
	n.interval = timeutil.Interval(reader.ReadVarint64())
	return reader.Error()
}

// String returns string value of add rollup file log
func (n *newRollupFile) String() string {
	return fmt.Sprintf("addRollup:{fileNumber:%d,interval:%d}", n.fileNumber, n.interval)
}

// apply applies new rollup file edit log to version
func (n *newRollupFile) apply(version Version) {
	version.AddRollupFile(n.fileNumber, n.interval)
}

// deleteRollupFile represent version edit log for delete rollup file for rollup job
type deleteRollupFile struct {
	fileNumber table.FileNumber // file number
}

// CreateDeleteRollupFile creates a remove rollup file
func CreateDeleteRollupFile(fileNumber table.FileNumber) Log {
	return &deleteRollupFile{
		fileNumber: fileNumber,
	}
}

// Encode writes remove rollup file data into binary
func (d *deleteRollupFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint64(d.fileNumber.Int64())
	return writer.Bytes()
}

// Decode reads remove rollup file data from binary
func (d *deleteRollupFile) Decode(v []byte) error {
	reader := stream.NewReader(v)
	d.fileNumber = table.FileNumber(reader.ReadVarint64())
	return reader.Error()
}

// String returns string value of delete rollup file log
func (d *deleteRollupFile) String() string {
	return fmt.Sprintf("deleteRollup:{fileNumber:%d}", d.fileNumber)
}

// apply applies remove rollup file edit log to version
func (d *deleteRollupFile) apply(version Version) {
	version.DeleteRollupFile(d.fileNumber)
}

// newReferenceFile represent version edit log for new reference file for rollup job
type newReferenceFile struct {
	familyID   FamilyID         // source family id
	fileNumber table.FileNumber // source file number
}

// CreateNewReferenceFile creates a new reference file
func CreateNewReferenceFile(familyID FamilyID, fileNumber table.FileNumber) Log {
	return &newReferenceFile{
		fileNumber: fileNumber,
		familyID:   familyID,
	}
}

// Encode writes new reference file data into binary
func (n *newReferenceFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)
	writer.PutVarint64(n.fileNumber.Int64())
	writer.PutVarint32(n.familyID.Int32())
	return writer.Bytes()
}

// Decode reads new reference file data from binary
func (n *newReferenceFile) Decode(v []byte) error {
	reader := stream.NewReader(v)
	n.fileNumber = table.FileNumber(reader.ReadVarint64())
	n.familyID = FamilyID(reader.ReadVarint32())
	return reader.Error()
}

// String returns string value of add reference file log
func (n *newReferenceFile) String() string {
	return fmt.Sprintf("addRefFile:{familyID:%d,fileNumber:%d}", n.familyID, n.fileNumber)
}

// apply applies new reference file edit log to version
func (n *newReferenceFile) apply(version Version) {
	version.AddReferenceFile(n.familyID, n.fileNumber)
}

// deleteReferenceFile represent version edit log for remove reference file for rollup job
type deleteReferenceFile struct {
	familyID   FamilyID         // source family id
	fileNumber table.FileNumber // source file number
}

// CreateDeleteReferenceFile creates a delete reference file
func CreateDeleteReferenceFile(familyID FamilyID, fileNumber table.FileNumber) Log {
	return &deleteReferenceFile{
		fileNumber: fileNumber,
		familyID:   familyID,
	}
}

// Encode writes remove reference file data into binary
func (n *deleteReferenceFile) Encode() ([]byte, error) {
	writer := stream.NewBufferWriter(nil)
	writer.PutVarint64(n.fileNumber.Int64())
	writer.PutVarint32(n.familyID.Int32())
	return writer.Bytes()
}

// Decode reads delete reference file data from binary
func (n *deleteReferenceFile) Decode(v []byte) error {
	reader := stream.NewReader(v)
	n.fileNumber = table.FileNumber(reader.ReadVarint64())
	n.familyID = FamilyID(reader.ReadVarint32())
	return reader.Error()
}

// String returns string value of delete reference file log
func (n *deleteReferenceFile) String() string {
	return fmt.Sprintf("deleteRefFile:{familyID:%d,fileNumber:%d}", n.familyID, n.fileNumber)
}

// apply applies remove reference file edit log to version
func (n *deleteReferenceFile) apply(version Version) {
	version.DeleteReferenceFile(n.familyID, n.fileNumber)
}
