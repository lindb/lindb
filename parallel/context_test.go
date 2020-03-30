package parallel

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql"
)

func TestBrokerExecuteContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expression := aggregation.NewMockExpression(ctrl)

	query, err := sql.Parse("select f from cpu")
	assert.NoError(t, err)
	query.Interval = timeutil.Interval(10 * timeutil.OneSecond)

	ctx := NewBrokerExecuteContext(timeutil.NowNano(), query)
	brokerCtx := ctx.(*brokerExecuteContext)
	brokerCtx.expression = expression
	assert.NotNil(t, brokerCtx.expression)
	assert.NotNil(t, ctx.ResultCh())
	it := series.NewMockGroupedIterator(ctrl)
	expression.EXPECT().Eval(gomock.Any())
	values := collections.NewFloatArray(10)
	values.SetValue(1, 10.0)
	expression.EXPECT().ResultSet().Return(map[string]collections.FloatArray{"test": nil, "f": values})
	expression.EXPECT().Reset()
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{it},
	})
	ctx.Emit(&series.TimeSeriesEvent{
		Err:   fmt.Errorf("err"),
		Stats: &models.QueryStats{},
	})

	rs, err := ctx.ResultSet()
	ctx.Complete(nil)
	ctx.Complete(fmt.Errorf("err"))
	assert.Error(t, err)
	assert.NotNil(t, rs.Series[0].Fields["f"])
}

func TestBrokerExecuteContext_Emit_GroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expression := aggregation.NewMockExpression(ctrl)

	query, err := sql.Parse("select f from cpu group by host")
	assert.NoError(t, err)
	query.Interval = timeutil.Interval(10 * timeutil.OneSecond)

	ctx := NewBrokerExecuteContext(timeutil.NowNano(), query)
	brokerCtx := ctx.(*brokerExecuteContext)
	brokerCtx.expression = expression
	assert.NotNil(t, brokerCtx.expression)
	assert.NotNil(t, ctx.ResultCh())
	it := series.NewMockGroupedIterator(ctrl)
	it.EXPECT().Tags().Return("host")
	expression.EXPECT().Eval(gomock.Any())
	values := collections.NewFloatArray(10)
	values.SetValue(1, 10.0)
	expression.EXPECT().ResultSet().Return(map[string]collections.FloatArray{"test": nil, "f": values})
	expression.EXPECT().Reset()
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{it},
		Stats:      &models.QueryStats{},
	})
	ctx.Emit(&series.TimeSeriesEvent{
		Err: fmt.Errorf("err"),
	})

	rs, err := ctx.ResultSet()
	ctx.Complete(nil)
	ctx.Complete(fmt.Errorf("err"))
	assert.Error(t, err)
	assert.NotNil(t, rs.Series[0].Fields["f"])

	ctx = NewBrokerExecuteContext(timeutil.NowNano(), query)
	brokerCtx = ctx.(*brokerExecuteContext)
	brokerCtx.expression = expression
	assert.NotNil(t, brokerCtx.expression)
	assert.NotNil(t, ctx.ResultCh())
	it = series.NewMockGroupedIterator(ctrl)
	it.EXPECT().Tags().Return("")
	ctx.Emit(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{it},
	})
	rs, err = ctx.ResultSet()
	ctx.Complete(nil)
	assert.NoError(t, err)
	assert.Len(t, rs.Series, 0)
}

func TestBrokerExecuteContext_ResultSet(t *testing.T) {
	ctx := NewBrokerExecuteContext(timeutil.NowNano(), nil)
	ctx.Complete(fmt.Errorf("err"))
	rs, err := ctx.ResultSet()
	assert.Error(t, err)
	assert.NotNil(t, rs)
}
