package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

func TestTagKVEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer func() {
		entrySetToIDSetFunc = entrySetToIDSet
	}()

	entries := TagKVEntries{}
	assert.Equal(t, 0, entries.TagValuesCount())

	tagKVEntry := NewMockTagKVEntrySetINTF(ctrl)
	tagKVEntry.EXPECT().TagValuesCount().Return(10)
	entries = TagKVEntries{tagKVEntry}
	assert.Equal(t, 10, entries.TagValuesCount())

	entrySetToIDSetFunc = func(entrySet TagKVEntrySetINTF, timeRange timeutil.TimeRange, offsets []int) (idSet *series.MultiVerSeriesIDSet, err error) {
		return nil, fmt.Errorf("err")
	}
	idSet, err := entries.GetSeriesIDs(timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, idSet)

	entrySetToIDSetFunc = func(entrySet TagKVEntrySetINTF, timeRange timeutil.TimeRange, offsets []int) (idSet *series.MultiVerSeriesIDSet, err error) {
		return series.NewMultiVerSeriesIDSet(), nil
	}
	idSet, err = entries.GetSeriesIDs(timeutil.TimeRange{})
	assert.NoError(t, err)
	assert.NotNil(t, idSet)
}
