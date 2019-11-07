package tsdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

func TestDataFamily_BaseTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	timeRange := timeutil.TimeRange{
		Start: 10,
		End:   50,
	}
	dataFamily := newDataFamily(timeutil.Interval(timeutil.OneSecond*10), timeRange, family)
	assert.Equal(t, timeRange, dataFamily.TimeRange())
	assert.Equal(t, int64(10000), dataFamily.Interval())
}

func TestDataFamily_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	dataFamily := newDataFamily(timeutil.Interval(timeutil.OneSecond*10), timeutil.TimeRange{
		Start: 10,
		End:   50,
	}, family)

	mockSnapShot := version.NewMockSnapshot(ctrl)
	mockSnapShot.EXPECT().Close().Return().AnyTimes()
	family.EXPECT().GetSnapshot().Return(mockSnapShot).AnyTimes()

	// WithFindReadersError
	mockSnapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("error"))
	dataFamily.Scan(&series.ScanContext{})

	// WithFindReadersOK
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	mockSnapShot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{mockReader}, nil)
	dataFamily.Scan(&series.ScanContext{})
}
