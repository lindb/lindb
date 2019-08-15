package metric

import (
	"fmt"
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
		URL:            "/metric/sum?cluster=dal",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	cm.EXPECT().GetChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	ch := replication.NewMockChannel(ctrl)
	cm.EXPECT().GetChannel(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch, nil).AnyTimes()
	ch.EXPECT().Write(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	ch.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 200,
	})

}
