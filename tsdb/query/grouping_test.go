package query

import (
	"encoding/binary"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
)

func TestGroupingContext_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := series.NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]uint32{1}, map[uint32][]series.GroupingScanner{1: {scanner}})
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(roaring.BitmapOf(1, 2, 3, 10).GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	result := ctx.BuildGroup(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Len(t, result, 2)
	tagValueIDs := make([]byte, 4)
	binary.LittleEndian.PutUint32(tagValueIDs[0:], 10)
	seriesIDs := result[string(tagValueIDs)]
	assert.Equal(t, []uint16{1, 10}, seriesIDs)
	binary.LittleEndian.PutUint32(tagValueIDs[0:], 20)
	seriesIDs = result[string(tagValueIDs)]
	assert.Equal(t, []uint16{2}, seriesIDs)

	scanner.EXPECT().GetSeriesAndTagValue(uint16(2)).
		Return(roaring.BitmapOf(1, 2).GetContainerAtIndex(0), []uint32{30, 10})
	_ = ctx.BuildGroup(2, roaring.BitmapOf(1, 2).GetContainerAtIndex(0))
	// container not found
	scanner.EXPECT().GetSeriesAndTagValue(uint16(3)).Return(nil, nil)
	_ = ctx.BuildGroup(3, roaring.BitmapOf(1, 2).GetContainerAtIndex(0))
	// case: get group by tag value ids
	groupByTagValueIDs := ctx.GetGroupByTagValueIDs()
	assert.Len(t, groupByTagValueIDs, 1)
	assert.EqualValues(t, roaring.BitmapOf(10, 20, 30).ToArray(), groupByTagValueIDs[0].ToArray())
}
