package version

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	strm "github.com/lindb/lindb/pkg/stream"
)

func TestEditLogCodec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLog := NewMockLog(ctrl)
	RegisterLogType(1000, func() Log {
		return mockLog
	})
	defer func() {
		delete(newLogFuncMap, 1000)
	}()

	editLog := NewEditLog(1)
	assert.True(t, editLog.IsEmpty())

	newFile := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(NewDeleteFile(1, 123))

	v, err := editLog.marshal()

	assert.Nil(t, err)
	assert.True(t, len(v) > 0)

	editLog2 := NewEditLog(1)
	err2 := editLog2.unmarshal(v)
	assert.Nil(t, err2)

	assert.Equal(t, editLog, editLog2)

	editLog = NewEditLog(1)
	editLog.Add(mockLog)
	mockLog.EXPECT().Encode().Return(nil, fmt.Errorf("err"))
	_, err = editLog.marshal()
	assert.NotNil(t, err)
}

func TestEditLog_Unmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLog := NewMockLog(ctrl)
	RegisterLogType(1000, func() Log {
		return mockLog
	})
	defer func() {
		delete(newLogFuncMap, 1000)
	}()

	stream := strm.NewBufferWriter(nil)
	// write family id
	stream.PutVarint32(int32(1))
	// write num of logs
	stream.PutUvarint64(uint64(2))
	stream.PutUvarint32(uint32(10000))
	value, _ := stream.Bytes()
	editLog := NewEditLog(1)
	err := editLog.unmarshal(value)
	assert.NotNil(t, err)

	stream = strm.NewBufferWriter(nil)
	// write family id
	stream.PutVarint32(int32(1))
	// write num of logs
	stream.PutUvarint64(uint64(2))

	stream.PutVarint32(int32(1000))
	stream.PutUvarint32(uint32(3))
	stream.PutBytes([]byte("123"))
	value, _ = stream.Bytes()
	mockLog.EXPECT().Decode([]byte("123")).Return(fmt.Errorf("err"))
	err = editLog.unmarshal(value)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestEditLog_apply(t *testing.T) {
	initVersionSetTestData()
	defer destroyVersionTestData()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	familyVersion := vs.CreateFamilyVersion("family", 1)
	editLog := NewEditLog(1)
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	editLog.Add(newFile)
	version := newVersion(1, familyVersion)
	editLog.apply(version)

	assert.Equal(t, 1, len(version.getAllFiles()), "cannot add file into version")
	//delete file
	editLog2 := NewEditLog(1)
	editLog2.Add(NewDeleteFile(1, 12))
	editLog2.Add(NewNextFileNumber(int64(120)))
	editLog2.apply(version)
	assert.Equal(t, 2, len(editLog2.GetLogs()))
	assert.Equal(t, 0, len(version.getAllFiles()), "cannot delete file from version")
}

func TestEditLog_applyVersionSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vs := NewMockStoreVersionSet(ctrl)
	vs.EXPECT().setNextFileNumberWithoutLock(int64(120))

	editLog := NewEditLog(1)
	mockLog := NewMockLog(ctrl)
	editLog.Add(mockLog)
	editLog.Add(NewNextFileNumber(int64(120)))
	editLog.applyVersionSet(vs)
}
