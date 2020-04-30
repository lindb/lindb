package wal

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
)

var testSeriesWALPath = "seriesWAL"

func TestNewSeriesWAL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		mkDirFunc = fileutil.MkDirIfNotExist
		newPageFactoryFunc = page.NewFactory

		_ = fileutil.RemoveDir(testSeriesWALPath)

		ctrl.Finish()
	}()

	// case 1: make path err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	wal, err := NewSeriesWAL(testSeriesWALPath)
	assert.Error(t, err)
	assert.Nil(t, wal)
	mkDirFunc = fileutil.MkDirIfNotExist
	// case 2: new page factory err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.Error(t, err)
	assert.Nil(t, wal)

	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}
	// case 3: AcquirePage err
	fct.EXPECT().Close().Return(fmt.Errorf("err"))
	fct.EXPECT().GetPageIDs().Return([]int64{19, 20, 21})
	fct.EXPECT().AcquirePage(int64(22)).Return(nil, fmt.Errorf("err"))
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.Error(t, err)
	assert.Nil(t, wal)
	// case 4: init wal success with re-open
	fct.EXPECT().GetPageIDs().Return([]int64{19, 20, 21})
	fct.EXPECT().AcquirePage(int64(22)).Return(nil, nil)
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.True(t, wal.NeedRecovery())
	// case 5: init wal success with empty
	fct.EXPECT().GetPageIDs().Return(nil)
	fct.EXPECT().AcquirePage(int64(1)).Return(nil, nil)
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.False(t, wal.NeedRecovery())
}

func TestSeriesWAL_Append(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newPageFactoryFunc = page.NewFactory
		_ = fileutil.RemoveDir(testSeriesWALPath)

		ctrl.Finish()
	}()
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}
	mockPage := page.NewMockMappedPage(ctrl)
	fct.EXPECT().GetPageIDs().Return(nil)
	fct.EXPECT().AcquirePage(int64(1)).Return(mockPage, nil)
	wal, err := NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	wal1 := wal.(*seriesWAL)
	wal1.base.pageSize = 32
	// case 1: put series id
	gomock.InOrder(
		mockPage.EXPECT().PutUint32(uint32(10), 0),
		mockPage.EXPECT().PutUint64(uint64(20), 4),
		mockPage.EXPECT().PutUint32(uint32(100), 12),
	)
	err = wal.Append(10, 20, 100)
	assert.NoError(t, err)
	// case 2: put series id
	gomock.InOrder(
		mockPage.EXPECT().PutUint32(uint32(110), 16),
		mockPage.EXPECT().PutUint64(uint64(210), 20),
		mockPage.EXPECT().PutUint32(uint32(1100), 28),
	)
	err = wal.Append(110, 210, 1100)
	assert.NoError(t, err)
	// case 3: create new data page err
	gomock.InOrder(
		mockPage.EXPECT().Sync().Return(fmt.Errorf("err")),
		fct.EXPECT().AcquirePage(wal1.base.pageIndex.Load()+1).Return(nil, fmt.Errorf("err")),
	)
	err = wal.Append(10, 20, 100)
	assert.Error(t, err)
	// case 4: create new data page success, then write new series data
	gomock.InOrder(
		mockPage.EXPECT().Sync().Return(fmt.Errorf("err")),
		fct.EXPECT().AcquirePage(wal1.base.pageIndex.Load()+1).Return(mockPage, nil),
		mockPage.EXPECT().PutUint32(uint32(10), 0),
		mockPage.EXPECT().PutUint64(uint64(20), 4),
		mockPage.EXPECT().PutUint32(uint32(100), 12),
	)
	err = wal.Append(10, 20, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), wal1.base.pageIndex.Load())
}

func TestSeriesWAL_Recovery(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testSeriesWALPath)
	}()
	wal, err := NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	err = wal.Append(10, 20, 100)
	assert.NoError(t, err)
	err = wal.Append(10, 210, 1100)
	assert.NoError(t, err)
	assert.False(t, wal.NeedRecovery())
	err = wal.Close()
	assert.NoError(t, err)
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.True(t, wal.NeedRecovery())
	count := 0
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		if metricID == 10 && tagsHash == 20 && seriesID == 100 {
			count++
			return nil
		} else if metricID == 10 && tagsHash == 210 && seriesID == 1100 {
			count++
			return nil
		}
		return fmt.Errorf("err")
	}, func() error {
		count++
		return nil
	})
	assert.Equal(t, 3, count)
	assert.False(t, wal.NeedRecovery())
	err = wal.Close()
	assert.NoError(t, err)
	// case: re-open
	wal, err = NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.True(t, wal.NeedRecovery())
	// empty data page
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		return fmt.Errorf("err")
	}, func() error {
		return nil
	})
	assert.False(t, wal.NeedRecovery())
	err = wal.Close()
	assert.NoError(t, err)
}

func TestSeriesWAL_Recovery_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newPageFactoryFunc = page.NewFactory
		_ = fileutil.RemoveDir(testSeriesWALPath)

		ctrl.Finish()
	}()
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}
	mockPage := page.NewMockMappedPage(ctrl)
	fct.EXPECT().GetPageIDs().Return(nil)
	fct.EXPECT().AcquirePage(int64(1)).Return(mockPage, nil)
	wal, err := NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	wal1 := wal.(*seriesWAL)
	wal1.base.commitPageIndex.Store(10)
	wal1.base.pageIndex.Store(11)
	// case 1: get nil page by page id
	fct.EXPECT().GetPage(int64(10)).Return(nil, false)
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		return fmt.Errorf("err")
	}, func() error {
		return fmt.Errorf("err")
	})
	// case 2: metric id = 0
	fct.EXPECT().GetPage(int64(10)).Return(mockPage, true).AnyTimes()
	mockPage.EXPECT().ReadUint32(0).Return(uint32(0))
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		return fmt.Errorf("err")
	}, func() error {
		return fmt.Errorf("err")
	})
	// case 3: recovery err
	mockPage.EXPECT().ReadUint32(0).Return(uint32(10))
	mockPage.EXPECT().ReadUint64(4).Return(uint64(10))
	mockPage.EXPECT().ReadUint32(12).Return(uint32(10))
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		return fmt.Errorf("err")
	}, func() error {
		return fmt.Errorf("err")
	})
	// case 4: release page err
	mockPage.EXPECT().ReadUint32(0).Return(uint32(0))
	fct.EXPECT().ReleasePage(int64(10)).Return(fmt.Errorf("err"))
	wal.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		return fmt.Errorf("err")
	}, func() error {
		return nil
	})

	assert.False(t, wal.NeedRecovery())
}

func TestSeriesWAL_Close(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testSeriesWALPath)
	}()
	wal, err := NewSeriesWAL(testSeriesWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.NoError(t, wal.Sync())
	assert.NoError(t, wal.Close())
}
