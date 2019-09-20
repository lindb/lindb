package query

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetricAPI_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	api := NewMetricAPI(nil, nil, executorFactory, nil)

	ch := make(chan *series.TimeSeriesEvent)

	time.AfterFunc(10*time.Millisecond, func() {
		it := series.NewMockGroupedIterator(ctrl)
		it.EXPECT().Tags().Return(nil)
		it.EXPECT().HasNext().Return(true)
		pIt := series.NewMockPrimitiveIterator(ctrl)
		pIt.EXPECT().HasNext().Return(true)
		pIt.EXPECT().Next().Return(10, 10.0)
		pIt.EXPECT().HasNext().Return(false)
		fIt := series.NewMockFieldIterator(ctrl)
		fIt.EXPECT().SegmentStartTime().Return(int64(10))
		fIt.EXPECT().HasNext().Return(true)
		fIt.EXPECT().FieldMeta().Return(field.Meta{Name: "f1"})
		fIt.EXPECT().Next().Return(pIt)
		fIt.EXPECT().HasNext().Return(false)
		it.EXPECT().Next().Return(fIt)
		it.EXPECT().HasNext().Return(false)
		ch <- &series.TimeSeriesEvent{
			Series: it,
		}
		close(ch)
	})

	exec := parallel.NewMockExecutor(ctrl)
	executorFactory.EXPECT().
		NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	exec.EXPECT().Execute().Return(ch)
	exec.EXPECT().Error().Return(nil).MaxTimes(2)
	exec.EXPECT().Statement().Return(&stmt.Query{MetricName: "test.metric.name", Interval: 10 * timeutil.OneSecond})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 200,
	})
}

func TestNewMetricAPI_Search_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	api := NewMetricAPI(nil, nil, executorFactory, nil)

	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})

	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})

	exec := parallel.NewMockExecutor(ctrl)
	executorFactory.EXPECT().
		NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	exec.EXPECT().Execute().Return(nil)
	exec.EXPECT().Error().Return(fmt.Errorf("err")).MaxTimes(2)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})

	ch := make(chan *series.TimeSeriesEvent)

	time.AfterFunc(10*time.Millisecond, func() {
		ch <- &series.TimeSeriesEvent{
			Err: fmt.Errorf("err"),
		}
		close(ch)
	})

	executorFactory.EXPECT().
		NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	exec.EXPECT().Execute().Return(ch)
	exec.EXPECT().Error().Return(nil).MaxTimes(2)
	exec.EXPECT().Statement().Return(&stmt.Query{MetricName: "test.metric.name", Interval: 10 * timeutil.OneSecond})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})
}
