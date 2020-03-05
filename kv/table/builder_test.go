package table

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
)

const (
	testKVPath = "test_builder"
)

var bitmapMarshal = encoding.BitmapMarshal

func TestFileNumber_Int64(t *testing.T) {
	assert.Equal(t, int64(10), FileNumber(10).Int64())
}

func TestStoreBuilder_magicNumber(t *testing.T) {
	code := []byte("eleme-ci")
	assert.Len(t, code, 8)
	assert.Equal(t, magicNumberOffsetFile, binary.LittleEndian.Uint64(code))
}

func TestStoreBuilder_BuildStore(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	var builder, err = NewStoreBuilder(10, testKVPath+"/000010.sst")
	defer func() {
		_ = os.RemoveAll(testKVPath)
		_ = builder.Close()
	}()

	assert.Nil(t, err)

	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	// reject for duplicate key
	err = builder.Add(1, []byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), builder.Count())

	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	assert.Equal(t, uint32(1), builder.MinKey())
	assert.Equal(t, uint32(10), builder.MaxKey())
	assert.Equal(t, FileNumber(10), builder.FileNumber())
	assert.True(t, builder.Size() > 0)
}

func TestStoreBuilder_Build_Err(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	ctrl := gomock.NewController(t)
	defer func() {
		newBufioWriterFunc = bufioutil.NewBufioWriter
		encoding.BitmapMarshal = bitmapMarshal
		_ = os.Remove(testKVPath)
		ctrl.Finish()
	}()
	writer := bufioutil.NewMockBufioWriter(ctrl)
	newBufioWriterFunc = func(fileName string) (bufioutil.BufioWriter, error) {
		return writer, nil
	}
	builder, err := NewStoreBuilder(10, testKVPath+"/000200.sst")
	assert.NoError(t, err)
	writer.EXPECT().Size().Return(int64(10)).AnyTimes()

	// case 1: write value err
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = builder.Add(10, []byte{1, 2, 3})
	assert.Error(t, err)
	// case 2: close empty keys
	err = builder.Close()
	assert.Equal(t, ErrEmptyKeys, err)
	// case 3: close write offset err
	writer.EXPECT().Write([]byte{1, 2, 3}).Return(10, nil)
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = builder.Add(10, []byte{1, 2, 3})
	assert.NoError(t, err)
	err = builder.Close()
	assert.Error(t, err)
	// case 4: bitmap marshal err
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	writer.EXPECT().Write(gomock.Any()).Return(10, nil)
	err = builder.Close()
	assert.Error(t, err)
	// case 5: write keys err
	encoding.BitmapMarshal = bitmapMarshal
	writer.EXPECT().Write(gomock.Any()).Return(10, nil)              // write offset
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err")) // write keys
	err = builder.Close()
	assert.Error(t, err)
	// case 6: write footer err
	writer.EXPECT().Write(gomock.Any()).Return(10, nil).MaxTimes(2)  // write offset/keys
	writer.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err")) // write footer
	err = builder.Close()
	assert.Error(t, err)
	// case 6: new builder err
	newBufioWriterFunc = func(fileName string) (bufioutil.BufioWriter, error) {
		return nil, fmt.Errorf("err")
	}
	builder, err = NewStoreBuilder(10, testKVPath+"/000200.sst")
	assert.Error(t, err)
	assert.Nil(t, builder)
}

func TestStoreBuilder_Abandon(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()
	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	assert.NoError(t, err)
	_ = builder.Add(1, []byte("test"))
	err = builder.Abandon()
	assert.NoError(t, err)
}
