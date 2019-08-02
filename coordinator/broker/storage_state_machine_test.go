package broker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

type testStorageStateMachineSuite struct {
	mock.RepoTestSuite
}

func TestStorageStateMachine(t *testing.T) {
	check.Suite(&testStorageStateMachineSuite{})
	check.TestingT(t)
}

func (ts *testStorageStateMachineSuite) TestStorageState(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/coordinator/storage/state",
		Endpoints: ts.Cluster.Endpoints,
	})

	stateMachine, err := NewStorageStateMachine(context.TODO(), repo)
	if err != nil {
		c.Fatal(err)
	}
	clusters := stateMachine.List()
	c.Assert(0, check.Equals, len(clusters))

	storageState := models.NewStorageState()
	storageState.Name = "LinDB_Storage"
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})
	srv := service.NewStorageStateService(repo)
	err = srv.Save("Test_LinDB", storageState)
	if err != nil {
		c.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	clusters = stateMachine.List()
	c.Assert(1, check.Equals, len(clusters))
	c.Assert(*storageState, check.DeepEquals, *clusters[0])

	_ = stateMachine.Close()
	clusters = stateMachine.List()
	c.Assert(0, check.Equals, len(clusters))

	// re-open state exist
	repo, _ = state.NewRepo(state.Config{
		Namespace: "/coordinator/storage/state",
		Endpoints: ts.Cluster.Endpoints,
	})
	stateMachine, err = NewStorageStateMachine(context.TODO(), repo)
	if err != nil {
		c.Fatal(err)
	}
	clusters = stateMachine.List()
	c.Assert(1, check.Equals, len(clusters))
	c.Assert(*storageState, check.DeepEquals, *clusters[0])
}

func (ts *testStorageStateMachineSuite) TestStorageState_Fail(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/coordinator/storage/state/fail",
		Endpoints: ts.Cluster.Endpoints,
	})

	stateMachine, err := NewStorageStateMachine(context.TODO(), repo)
	if err != nil {
		c.Fatal(err)
	}
	clusters := stateMachine.List()
	c.Assert(0, check.Equals, len(clusters))

	storageState := models.NewStorageState()
	storageState.Name = "Test_LinDB"
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})
	srv := service.NewStorageStateService(repo)
	err = srv.Save("Test_LinDB", storageState)
	if err != nil {
		c.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	clusters = stateMachine.List()
	c.Assert(1, check.Equals, len(clusters))
	c.Assert(*storageState, check.DeepEquals, *clusters[0])

	// wrong data
	_ = repo.Put(context.TODO(), pathutil.GetStorageClusterStatePath("test"), []byte("ddd"))
	time.Sleep(200 * time.Millisecond)
	clusters = stateMachine.List()
	c.Assert(1, check.Equals, len(clusters))

	_ = repo.Delete(context.TODO(), pathutil.GetStorageClusterStatePath("Test_LinDB"))
	time.Sleep(200 * time.Millisecond)
	clusters = stateMachine.List()
	c.Assert(0, check.Equals, len(clusters))

	data, _ := json.Marshal(models.NewStorageState())
	_ = repo.Put(context.TODO(), pathutil.GetStorageClusterStatePath("test1"), data)
	time.Sleep(200 * time.Millisecond)
	clusters = stateMachine.List()
	c.Assert(0, check.Equals, len(clusters))
}
