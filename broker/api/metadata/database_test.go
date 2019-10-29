package metadata

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/service"
)

func TestDatabaseAPI_ListDatabaseNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	databaseService := service.NewMockDatabaseService(ctrl)
	api := NewDatabaseAPI(databaseService)

	databaseService.EXPECT().List().Return(nil, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database/names",
		HandlerFunc:    api.ListDatabaseNames,
		ExpectHTTPCode: 200,
		RequestBody:    []string{},
	})

	databaseService.EXPECT().List().Return(
		[]*models.Database{
			{Name: "test1"},
			{Name: "test2"},
		},
		nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database/names",
		HandlerFunc:    api.ListDatabaseNames,
		ExpectHTTPCode: 200,
		RequestBody:    []string{"test1", "test2"},
	})

	databaseService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database/names",
		HandlerFunc:    api.ListDatabaseNames,
		ExpectHTTPCode: 500,
	})
}
