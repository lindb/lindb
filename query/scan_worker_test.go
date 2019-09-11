package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series"
)

func TestScanWorker_Emit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aggWorker := NewMockaggregateWorker(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), nil, nil, aggWorker)
	gomock.InOrder(
		aggWorker.EXPECT().emit(gomock.Any()),
		aggWorker.EXPECT().close(),
		aggWorker.EXPECT().sendResult(gomock.Any()),
	)
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FieldIt:         nil,
		FamilyStartTime: 10,
	})
	w := worker.(*scanWorker)
	worker.Emit(nil)

	w.Close()
}

func TestScanWorker_Emit_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aggWorker := NewMockaggregateWorker(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), nil, nil, aggWorker)
	aggWorker.EXPECT().close()
	aggWorker.EXPECT().sendResult(gomock.Any())
	worker.Close()
	worker.Complete(uint32(10))
}

func TestScanWorker_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aggWorker := NewMockaggregateWorker(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), nil, nil, aggWorker)
	gomock.InOrder(
		aggWorker.EXPECT().emit(gomock.Any()),
		aggWorker.EXPECT().close(),
		aggWorker.EXPECT().sendResult(gomock.Any()),
	)
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FieldIt:         nil,
		FamilyStartTime: 10,
	})
	worker.Complete(uint32(10))

	worker.Close()
}

func TestScanWorker_GroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aggWorker := NewMockaggregateWorker(ctrl)
	metaGetter := series.NewMockMetaGetter(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), []string{"host", "disk"}, metaGetter, aggWorker)
	gomock.InOrder(
		aggWorker.EXPECT().emit(gomock.Any()),
		metaGetter.EXPECT().GetTagValues(uint32(10), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),
		aggWorker.EXPECT().emit(gomock.Any()),
		aggWorker.EXPECT().sendResult(gomock.Any()),
		aggWorker.EXPECT().emit(gomock.Any()),
		metaGetter.EXPECT().GetTagValues(uint32(10), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),
		aggWorker.EXPECT().sendResult(gomock.Any()),
		aggWorker.EXPECT().close(),
	)
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FamilyStartTime: 10,
	})
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FamilyStartTime: 11,
	})
	worker.Complete(uint32(10))
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(11),
		Completed:       false,
		FamilyStartTime: 10,
	})
	worker.Complete(uint32(11))
	worker.Close()

	// test panic
	worker = createScanWorker(context.TODO(), uint32(10), []string{"host", "disk"}, nil, aggWorker)
	gomock.InOrder(
		aggWorker.EXPECT().emit(gomock.Any()),
		aggWorker.EXPECT().close().AnyTimes(),
	)
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FamilyStartTime: 10,
	})
	worker.Close()

	// test get group by tag values err
	worker = createScanWorker(context.TODO(), uint32(10), []string{"host", "disk"}, metaGetter, aggWorker)
	gomock.InOrder(
		aggWorker.EXPECT().emit(gomock.Any()),
		metaGetter.EXPECT().GetTagValues(uint32(10), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")),
		aggWorker.EXPECT().close().AnyTimes(),
	)
	worker.Emit(&series.FieldEvent{
		Version:         series.Version(10),
		SeriesID:        uint32(10),
		Completed:       false,
		FamilyStartTime: 10,
	})
	worker.Close()
}
