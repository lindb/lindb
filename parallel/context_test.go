package parallel

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
)

func TestBrokerExecuteContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expression := aggregation.NewMockExpression(ctrl)

	query, err := sql.Parse("select f from cpu")
	assert.NoError(t, err)
	query.Interval = 10 * timeutil.OneSecond

	ctx := NewBrokerExecuteContext(query)
	brokerCtx := ctx.(*brokerExecuteContext)
	brokerCtx.expression = expression
	ctx.RetainTask(10)
	assert.NotNil(t, brokerCtx.expression)
	assert.NotNil(t, ctx.ResultCh())
	it := series.NewMockGroupedIterator(ctrl)
	it.EXPECT().Tags().Return(nil)
	expression.EXPECT().Eval(gomock.Any())
	values := collections.NewFloatArray(10)
	values.SetValue(1, 10.0)
	expression.EXPECT().ResultSet().Return(map[string]collections.FloatArray{"test": nil, "f": values})
	expression.EXPECT().Reset()
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{it},
	})
	ctx.Emit(&series.TimeSeriesEvent{
		Err: fmt.Errorf("err"),
	})

	rs, err := ctx.ResultSet()
	ctx.Complete(nil)
	ctx.Complete(fmt.Errorf("err"))
	assert.Error(t, err)
	assert.NotNil(t, rs.Series[0].Fields["f"])
}

func TestStorageExecuteContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := pb.NewMockTaskService_HandleServer(ctrl)

	ctx := newStorageExecutorContext(context.TODO(), &pb.TaskRequest{
		JobID:        10,
		ParentTaskID: "task_1",
	}, stream)
	assert.NotNil(t, ctx)

	stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))

	ctx.RetainTask(2)
	ctx.Emit(&series.TimeSeriesEvent{
		Err: fmt.Errorf("err"),
	})
	ctx.Complete(nil)
	ctx.Complete(fmt.Errorf("err"))
	ctx.Emit(nil)

	// test normal case
	ctx = newStorageExecutorContext(context.TODO(), &pb.TaskRequest{
		JobID:        10,
		ParentTaskID: "task_1",
	}, stream)
	ctx.RetainTask(1)
	gIt := series.NewMockGroupedIterator(ctrl)
	it := series.NewMockIterator(ctrl)
	fIt := series.NewMockFieldIterator(ctrl)
	gomock.InOrder(
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(it),
		it.EXPECT().FieldType().Return(field.SumField),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("ff")), //err
		gIt.EXPECT().HasNext().Return(false),
	)
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{gIt},
	})
	gomock.InOrder(
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(nil), // empty
		gIt.EXPECT().HasNext().Return(false),
	)
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{gIt},
	})
	gomock.InOrder(
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(it),
		it.EXPECT().FieldType().Return(field.SumField),
		it.EXPECT().HasNext().Return(true),
		it.EXPECT().Next().Return(int64(10), fIt),
		fIt.EXPECT().MarshalBinary().Return([]byte{1, 1, 1}, nil), //normal
		it.EXPECT().HasNext().Return(false),
		it.EXPECT().FieldName().Return("f"),
		gIt.EXPECT().HasNext().Return(false),
		gIt.EXPECT().Tags().Return(nil),
	)
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{gIt},
	})

	stream.EXPECT().Send(gomock.Any()).Return(nil)
	ctx.Complete(nil)
}
