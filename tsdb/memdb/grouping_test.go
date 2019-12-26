package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series"
)

func TestGroupingContext_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStore := newMetricStore()
	ms := mStore.(*metricStore)
	gCtx := series.NewMockGroupingContext(ctrl)
	gCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(nil)

	ctx := &groupingContext{
		ms:   ms,
		gCtx: gCtx,
	}
	ctx.BuildGroup(uint16(10), nil)
}
