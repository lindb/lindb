package metric

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/replication"
)

func TestWriteAPI_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewWriteAPI(cm)
	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	cm.EXPECT().Write(gomock.Any()).Return(errors.New("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal&c=1",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	cm.EXPECT().Write(gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal&c=1",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 200,
	})

}
