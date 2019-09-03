package tsdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

func TestDataFamily_BaseTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	timeRange := &timeutil.TimeRange{
		Start: 10,
		End:   50,
	}
	dataFamily := newDataFamily(int64(10000), timeRange, family)
	assert.Equal(t, timeRange, dataFamily.TimeRange())
	assert.Equal(t, int64(10000), dataFamily.Interval())
}

func TestDataFamily_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	dataFamily := newDataFamily(int64(1000), &timeutil.TimeRange{
		Start: 10,
		End:   50,
	}, family)

	//TODO need impl scan test logic
	dataFamily.Scan(&series.ScanContext{})
}
