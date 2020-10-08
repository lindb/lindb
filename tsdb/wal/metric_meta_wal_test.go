package wal

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
	"github.com/lindb/lindb/series/field"
)

var testMetaWALPath = "metaWAL"
var ns = "ns"

func TestNewMetricMetaWAL(t *testing.T) {
	defer func() {
		mkDirFunc = fileutil.MkDirIfNotExist

		_ = fileutil.RemoveDir(testMetaWALPath)
	}()

	// case 1: make path err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	wal, err := NewMetricMetaWAL(testMetaWALPath)
	assert.Error(t, err)
	assert.Nil(t, wal)
	mkDirFunc = fileutil.MkDirIfNotExist
	// case 2: create wal success
	wal, err = NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.False(t, wal.NeedRecovery())

	err = wal.Sync()
	assert.NoError(t, err)

	err = wal.Close()
	assert.NoError(t, err)

	wal, err = NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	assert.True(t, wal.NeedRecovery())

	err = wal.Close()
	assert.NoError(t, err)
}

func TestMetricMetaWAL_Append(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testMetaWALPath)
	}()
	mockAppendData(t)
}

func TestMetricMetaWAL_Append_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testMetaWALPath)

		ctrl.Finish()
	}()
	fct := page.NewMockFactory(ctrl)
	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err")).MaxTimes(3)
	wal, err := NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)

	wal1 := wal.(*metricMetaWAL)
	wal1.base.walFactory = fct
	wal1.base.pageSize = 1

	assert.Error(t, wal.AppendTagKey(1, 1, "tagKey"))
	assert.Error(t, wal.AppendField(1, 1, "f", field.SumField))
	assert.Error(t, wal.AppendMetric(ns, "metric", 1))

	err = wal.Close()
	assert.NoError(t, err)
}

func TestMetricMetaWAL_Recovery(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testMetaWALPath)
	}()
	mockAppendData(t)

	metaWAL, err := NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, metaWAL)
	count := 0
	metaWAL.Recovery(func(namespace, metricName string, metricID uint32) error {
		if namespace == ns && metricName == "metric-2" && metricID == 2 {
			count++
			return nil
		} else if namespace == ns && metricName == "metric-1" && metricID == 1 {
			count++
			return nil
		}
		count++
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		if metricID == 1 && fID == field.ID(1) && fType == field.SumField && fieldName == "f-1" {
			count++
			return nil
		} else if metricID == 2 && fID == field.ID(2) && fType == field.HistogramField && fieldName == "f-2" {
			count++
			return nil
		}
		count++
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		if metricID == 1 && tagKeyID == 1 && tagKey == "tagKey-1" {
			count++
			return nil
		} else if metricID == 2 && tagKeyID == 2 && tagKey == "tagKey-2" {
			count++
			return nil
		}
		count++
		return nil
	}, func() error {
		count++
		return nil
	})
	assert.Equal(t, 7, count)
	assert.False(t, metaWAL.NeedRecovery())

	err = metaWAL.Close()
	assert.NoError(t, err)
}

func TestMetricMetaWAL_Recovery_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testMetaWALPath)

		ctrl.Finish()
	}()
	mockAppendData(t)
	wal, err := NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
	// case 1: commit err
	wal.Recovery(func(namespace, metricName string, metricID uint32) error {
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		return nil
	}, func() error {
		return fmt.Errorf("err")
	})
	assert.True(t, wal.NeedRecovery())
	// case 2: metric recovery err
	wal.Recovery(func(namespace, metricName string, metricID uint32) error {
		return fmt.Errorf("err")
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		return nil
	}, func() error {
		return fmt.Errorf("err")
	})
	assert.True(t, wal.NeedRecovery())
	// case 3: field recovery err
	wal.Recovery(func(namespace, metricName string, metricID uint32) error {
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		return fmt.Errorf("err")
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		return nil
	}, func() error {
		return fmt.Errorf("err")
	})
	assert.True(t, wal.NeedRecovery())
	// case 4: tag key recovery err
	wal.Recovery(func(namespace, metricName string, metricID uint32) error {
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		return fmt.Errorf("err")
	}, func() error {
		return fmt.Errorf("err")
	})
	assert.True(t, wal.NeedRecovery())

	// case 5: release err
	wal1 := wal.(*metricMetaWAL)
	fct := page.NewMockFactory(ctrl)
	mockPage := page.NewMockMappedPage(ctrl)
	wal1.base.walFactory = fct
	fct.EXPECT().GetPage(int64(0)).Return(nil, false)
	fct.EXPECT().GetPage(int64(1)).Return(mockPage, true)
	mockPage.EXPECT().ReadUint8(0).Return(uint8(0))
	fct.EXPECT().ReleasePage(gomock.Any()).Return(fmt.Errorf("err"))
	wal.Recovery(func(namespace, metricName string, metricID uint32) error {
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		return nil
	}, func() error {
		return nil
	})
	assert.False(t, wal.NeedRecovery())

	err = wal.Close()
	assert.NoError(t, err)
}

func mockAppendData(t *testing.T) {
	wal, err := NewMetricMetaWAL(testMetaWALPath)
	assert.NoError(t, err)
	assert.NotNil(t, wal)

	assert.NoError(t, wal.AppendTagKey(1, 1, "tagKey-1"))
	assert.NoError(t, wal.AppendField(1, 1, "f-1", field.SumField))
	assert.NoError(t, wal.AppendMetric(ns, "metric-1", 1))
	assert.NoError(t, wal.AppendField(2, 2, "f-2", field.HistogramField))
	assert.NoError(t, wal.AppendTagKey(2, 2, "tagKey-2"))
	assert.NoError(t, wal.AppendMetric(ns, "metric-2", 2))

	err = wal.Close()
	assert.NoError(t, err)
}
