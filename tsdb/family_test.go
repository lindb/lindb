package tsdb

import (
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/series"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDataFamily_BaseTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	dataFamily := newDataFamily(int64(1000), family)
	assert.Equal(t, int64(1000), dataFamily.BaseTime())
}

func TestDataFamily_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	dataFamily := newDataFamily(int64(1000), family)

	//TODO need impl scan test logic
	assert.Nil(t, dataFamily.Scan(series.ScanContext{}))
}
