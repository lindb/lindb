package admin

import (
	"net/http"
	"testing"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"

	check "gopkg.in/check.v1"
)

type testStorageClusterAPISuite struct {
	mock.RepoTestSuite
}

var test *testing.T

func TestStorageClusterAPI(t *testing.T) {
	check.Suite(&testStorageClusterAPISuite{})
	test = t
	check.TestingT(t)
}

func (ts *testStorageClusterAPISuite) TestStorageCluster(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	api := NewStorageClusterAPI(service.NewStorageClusterService(repo))

	cfg := models.StorageCluster{
		Name: "test1",
	}
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/stroage/cluster",
		RequestBody:    cfg,
		HandlerFunc:    api.Create,
		ExpectHTTPCode: 204,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/stroage/cluster",
		RequestBody:    models.StorageCluster{},
		HandlerFunc:    api.Create,
		ExpectHTTPCode: 500,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 200,
		ExpectResponse: cfg,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/stroage/cluster",
		HandlerFunc:    api.List,
		ExpectHTTPCode: 200,
		ExpectResponse: []models.StorageCluster{cfg},
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodDelete,
		URL:            "/stroage/cluster?name=test1",
		HandlerFunc:    api.DeleteByName,
		ExpectHTTPCode: 204,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster?name=test1",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 404,
	})
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/cluster?name=test19999",
		HandlerFunc:    api.GetByName,
		ExpectHTTPCode: 404,
	})
}
