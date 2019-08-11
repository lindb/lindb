package cluser

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestMasterAPI_GetMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	master := coordinator.NewMockMaster(ctrl)

	api := NewMasterAPI(master)

	m := models.Master{ElectTime: timeutil.Now(), Node: models.Node{IP: "1.1.1.1", Port: 8000}}
	// get success
	master.EXPECT().GetMaster().Return(&m)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/master",
		HandlerFunc:    api.GetMaster,
		ExpectHTTPCode: 200,
		ExpectResponse: &m,
	})

	master.EXPECT().GetMaster().Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/master",
		HandlerFunc:    api.GetMaster,
		ExpectHTTPCode: 404,
	})
}
