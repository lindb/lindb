package state

import (
	"context"
	"net/http"
	"testing"
	"time"

	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

type testStorageAPISuite struct {
	mock.RepoTestSuite
}

var test *testing.T

func TestDatabaseAPI(t *testing.T) {
	check.Suite(&testStorageAPISuite{})
	test = t
	check.TestingT(t)
}

func (ts *testStorageAPISuite) TestStorageState(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/api/storage/state",
		Endpoints: ts.Cluster.Endpoints,
	})

	stateMachine, err := broker.NewStorageStateMachine(context.TODO(), repo)
	if err != nil {
		c.Fatal(err)
	}
	api := NewStorageAPI(stateMachine)

	storageState := models.NewStorageState()
	storageState.Name = "LinDB_Storage"
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})
	srv := service.NewStorageStateService(repo)
	_ = srv.Save("Test_LinDB", storageState)
	time.Sleep(100 * time.Millisecond)

	// get success
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/state",
		HandlerFunc:    api.ListStorageCluster,
		ExpectHTTPCode: 200,
		ExpectResponse: []*models.StorageState{storageState},
	})
}
