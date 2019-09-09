package tblstore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_MetricsNameIDFlusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("write failure")).AnyTimes()
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit failure")).AnyTimes()

	nameIDFlusher := NewMetricsNameIDFlusher(mockFlusher)
	for nsID := 0; nsID < 2; nsID++ {
		for i := 0; i < 10000; i++ {
			nameIDFlusher.FlushNameID(strconv.Itoa(i), uint32(i))
		}
		assert.NotNil(t, nameIDFlusher.FlushMetricsNS(uint32(nsID), 1, 2))
	}
	assert.NotNil(t, nameIDFlusher.Commit())
}
