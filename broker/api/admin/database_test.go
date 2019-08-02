package admin

import (
	"net/http"
	"testing"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"

	"gopkg.in/check.v1"
)

type testDatabaseAPISuite struct {
	mock.RepoTestSuite
}

func TestDatabaseAPI(t *testing.T) {
	check.Suite(&testDatabaseAPISuite{})
	test = t
	check.TestingT(t)
}

func (ts *testDatabaseAPISuite) TestGetDatabase(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	api := NewDatabaseAPI(service.NewDatabaseService(repo))

	db := models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "test",
				NumOfShard:    12,
				ReplicaFactor: 3,
			},
		},
	}
	//create success
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    db,
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 204,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    models.Database{},
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 500,
	})

	// get success
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database?name=test",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 200,
		ExpectResponse: db,
	})
	// no database name
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})
	// wrong database name
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/database?name=test2",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 500,
	})
}
