package version

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestSnapshot_FindReaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	cache := table.NewMockCache(ctrl)
	snapshot := newSnapshot("test", v, cache)
	gomock.InOrder(
		cache.EXPECT().GetReader("test", Table(int64(10))).Return(nil, fmt.Errorf("err")),
		cache.EXPECT().GetReader("test", Table(int64(11))).Return(table.NewMockReader(ctrl), nil),
	)
	_, err := snapshot.GetReader(int64(10))
	assert.NotNil(t, err)
	reader, err := snapshot.GetReader(int64(11))
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, reader)

	v.addFile(0, NewFileMeta(int64(10), 1, 30, 30))
	readers, _ := snapshot.FindReaders(uint32(80))
	assert.Equal(t, 0, len(readers))
	cache.EXPECT().GetReader("test", Table(int64(10))).Return(table.NewMockReader(ctrl), nil)
	readers, _ = snapshot.FindReaders(uint32(20))
	assert.Equal(t, 1, len(readers))

	cache.EXPECT().GetReader("test", Table(int64(10))).Return(nil, fmt.Errorf("err"))
	readers, err = snapshot.FindReaders(uint32(20))
	assert.NotNil(t, err)
	assert.Nil(t, readers)
}
