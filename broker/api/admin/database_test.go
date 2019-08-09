package admin

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/service"
)

func TestDatabaseAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	databaseService := service.NewMockDatabaseService(ctrl)

	api := NewDatabaseAPI(databaseService)

	db := models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "test",
				NumOfShard:    12,
				ReplicaFactor: 3,
				Engine:        option.EngineOption{Interval: "10s"},
			},
		},
	}
	// create success
	databaseService.EXPECT().Save(gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    db,
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 204,
	})
	// create err
	databaseService.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("err"))
	db.Name = ""
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    db,
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 500,
	})
	db.Clusters = append(db.Clusters, models.DatabaseCluster{Engine: option.EngineOption{Interval: "aa"}})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    db,
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 500,
	})

	// get success
	databaseService.EXPECT().Get(gomock.Any()).Return(&db, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database?name=test",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 200,
		ExpectResponse: db,
	})
	// no database name
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})
	databaseService.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database?name=test",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})
}
