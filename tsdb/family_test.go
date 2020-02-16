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
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
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
	assert.NotNil(t, dataFamily.Family())
}

func TestDataFamily_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		newReaderFunc = metricsdata.NewReader
		newFilterFunc = metricsdata.NewFilter
	}()

	family := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	timeRange := timeutil.TimeRange{
		Start: 10,
		End:   50,
	}
	dataFamily := newDataFamily(timeutil.Interval(timeutil.OneSecond*10), timeRange, family)

	// test find kv readers err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	rs, err := dataFamily.Filter(uint32(10), nil, nil, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, rs)

	// case 1: find kv readers nil
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	rs, err = dataFamily.Filter(uint32(10), nil, nil, timeutil.TimeRange{})
	assert.NoError(t, err)
	assert.Nil(t, rs)

	// case 2: not find in reader
	reader := table.NewMockReader(ctrl)
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	reader.EXPECT().Get(gomock.Any()).Return(nil, false)
	rs, err = dataFamily.Filter(uint32(10), nil, nil, timeutil.TimeRange{})
	assert.NoError(t, err)
	assert.Nil(t, rs)

	// case 3: new metric reader err
	newReaderFunc = func(buf []byte) (reader metricsdata.Reader, err error) {
		return nil, fmt.Errorf("err")
	}
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, true)
	rs, err = dataFamily.Filter(uint32(10), nil, nil, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, rs)

	// case 4: normal case
	newReaderFunc = func(buf []byte) (reader metricsdata.Reader, err error) {
		return nil, nil
	}
	filter := metricsdata.NewMockFilter(ctrl)
	newFilterFunc = func(familyTime int64, snapshot version.Snapshot, readers []metricsdata.Reader) metricsdata.Filter {
		return filter
	}
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, true)
	filter.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err = dataFamily.Filter(uint32(10), nil, nil, timeutil.TimeRange{})
	assert.NoError(t, err)
}
