package tagkeymeta

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVFlusher := kv.NewMockFlusher(ctrl)
	flusher := NewFlusher(mockKVFlusher)
	assert.NotNil(t, flusher)

	// mock commit error
	mockKVFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, flusher.Commit())

	// mock commit ok
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := flusher.FlushTagKeyID(333, 100)
	assert.Nil(t, err)
}

func TestFlushTagKeyID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVFlusher := kv.NewMockFlusher(ctrl)
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	flusher := NewFlusher(mockKVFlusher)

	// flush tagValue1
	flusher.EnsureSize((1 << 8) * (1 << 8))
	for x := 1; x < 1<<8; x++ {
		for y := 1; y < 1<<8; y++ {
			flusher.FlushTagValue([]byte(fmt.Sprintf("192.168.%d.%d", x, y)), uint32(x*y))
		}
	}
	// flush tagKeyID
	assert.Nil(t, flusher.FlushTagKeyID(1, 10))
}
