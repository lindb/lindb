package query

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/parallel"
)

func TestMetricAPI_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	api := NewMetricAPI(nil, nil, executorFactory)

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

	executorFactory.EXPECT().
		NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 200,
	})
}
