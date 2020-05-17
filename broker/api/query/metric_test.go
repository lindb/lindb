package query

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/series"
)

func TestMetricAPI_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	brokerExecutor := parallel.NewMockBrokerExecutor(ctrl)
	executeCtx := parallel.NewMockBrokerExecuteContext(ctrl)
	brokerExecutor.EXPECT().ExecuteContext().Return(executeCtx)
	brokerExecutor.EXPECT().Execute()

	executorFactory.EXPECT().NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(brokerExecutor)

	api := NewMetricAPI(nil, nil, nil, executorFactory, nil)

	ch := make(chan *series.TimeSeriesEvent)

	executeCtx.EXPECT().ResultCh().Return(ch)
	executeCtx.EXPECT().Emit(gomock.Any())
	executeCtx.EXPECT().ResultSet().Return(&models.ResultSet{}, nil)

	time.AfterFunc(100*time.Millisecond, func() {
		ch <- nil
		close(ch)
	})

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
	api := NewMetricAPI(nil, nil, nil, executorFactory, nil)

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

	brokerExecutor := parallel.NewMockBrokerExecutor(ctrl)
	executeCtx := parallel.NewMockBrokerExecuteContext(ctrl)
	brokerExecutor.EXPECT().ExecuteContext().Return(executeCtx)
	brokerExecutor.EXPECT().Execute()

	executorFactory.EXPECT().NewBrokerExecutor(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(brokerExecutor)

	ch := make(chan *series.TimeSeriesEvent)

	executeCtx.EXPECT().ResultCh().Return(ch)
	executeCtx.EXPECT().ResultSet().Return(&models.ResultSet{}, fmt.Errorf("err"))

	time.AfterFunc(100*time.Millisecond, func() {
		close(ch)
	})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})
}
