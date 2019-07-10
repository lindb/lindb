package admin

import (
	"net/http"
	"testing"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/service"
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

	db := models.Database{Name: "test", NumOfShard: 1, ReplicaFactor: 1}
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
		RequestBody:    models.Database{NumOfShard: 1, ReplicaFactor: 1},
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 500,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    models.Database{Name: "test", ReplicaFactor: 1},
		HandlerFunc:    api.Save,
		ExpectHTTPCode: 500,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/database",
		RequestBody:    models.Database{Name: "test", NumOfShard: 1},
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
