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
	"github.com/lindb/lindb/tsdb/tblstore"
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
		newVersionBlockIterator = tblstore.NewVersionBlockIterator
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
	rs, err := dataFamily.Filter(uint32(10), nil, series.NewVersion(), nil)
	assert.Error(t, err)
	assert.Nil(t, rs)

	// test find kv readers nil
	snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
	rs, err = dataFamily.Filter(uint32(10), nil, series.NewVersion(), nil)
	assert.NoError(t, err)
	assert.Nil(t, rs)

	// test not find in reader
	reader := table.NewMockReader(ctrl)
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	reader.EXPECT().Get(gomock.Any()).Return(nil)
	rs, err = dataFamily.Filter(uint32(10), nil, series.NewVersion(), nil)
	assert.NoError(t, err)
	assert.Nil(t, rs)

	// test new version block err
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3})
	// create version block iterator err
	newVersionBlockIterator = func(block []byte) (iterator tblstore.VersionBlockIterator, e error) {
		return nil, fmt.Errorf("err")
	}
	rs, err = dataFamily.Filter(uint32(10), nil, series.NewVersion(), nil)
	assert.Error(t, err)
	assert.Nil(t, rs)

	// test normal case
	snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
	blockIt := tblstore.NewMockVersionBlockIterator(ctrl)
	newVersionBlockIterator = func(block []byte) (iterator tblstore.VersionBlockIterator, e error) {
		return blockIt, nil
	}
	reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3})
	blockIt.EXPECT().HasNext().Return(false)
	_, err = dataFamily.Filter(uint32(10), nil, series.NewVersion(), nil)
	assert.NoError(t, err)
}
