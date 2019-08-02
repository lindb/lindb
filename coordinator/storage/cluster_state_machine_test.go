package storage

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"

	check "gopkg.in/check.v1"
)

type testClusterStateMachineSuite struct {
	mock.RepoTestSuite
}

func TestClusterStateMachine(t *testing.T) {
	check.Suite(&testClusterStateMachineSuite{})
	check.TestingT(t)
}

func (ts *testClusterStateMachineSuite) TestDiscoveryFail(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "fail",
		Endpoints: ts.Cluster.Endpoints,
	})

	stateMachine, _ := NewClusterStateMachine(context.TODO(), repo)

	data, _ := json.Marshal(models.StorageCluster{
		Name: "test1",
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/1", data)
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/2", []byte("fail"))
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/3", []byte("{}"))
	time.Sleep(200 * time.Millisecond)
	c.Assert(0, check.Equals, len(stateMachine.GetAllCluster()))
	_ = stateMachine.Close()
}

func (ts *testClusterStateMachineSuite) TestDiscovery(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	stateMachine, _ := NewClusterStateMachine(context.TODO(), repo)

	storage1 := state.Config{
		Namespace: "storage1",
		Endpoints: ts.Cluster.Endpoints,
	}
	data1, _ := json.Marshal(models.StorageCluster{
		Name:   "storage1",
		Config: storage1,
	})
	storage2 := state.Config{
		Namespace: "storage2",
		Endpoints: ts.Cluster.Endpoints,
	}
	data2, _ := json.Marshal(models.StorageCluster{
		Name:   "storage2",
		Config: storage2,
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage1", data1)
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage2", data2)
	time.Sleep(200 * time.Millisecond)
	c.Assert(2, check.Equals, len(stateMachine.GetAllCluster()))
	cluster1 := stateMachine.GetCluster("storage1")
	cluster2 := stateMachine.GetCluster("storage2")
	c.Assert(cluster1, check.NotNil)
	c.Assert(cluster2, check.NotNil)

	repo2, _ := state.NewRepo(state.Config{
		Namespace: "storage1",
		Endpoints: ts.Cluster.Endpoints,
	})
	node1, _ := json.Marshal(models.ActiveNode{Node: models.Node{IP: "127.0.0.1", Port: 2080}})
	_ = repo2.Put(context.TODO(), constants.ActiveNodesPath+"/node1", node1)
	node2, _ := json.Marshal(models.ActiveNode{Node: models.Node{IP: "127.0.0.2", Port: 2080}})
	_ = repo2.Put(context.TODO(), constants.ActiveNodesPath+"/127.0.0.2:2080", node2)
	_ = repo2.Put(context.TODO(), constants.ActiveNodesPath+"/node3", []byte("dd"))
	time.Sleep(200 * time.Millisecond)
	c.Assert(2, check.Equals, len(cluster1.GetActiveNodes()))
	c.Assert(0, check.Equals, len(cluster2.GetActiveNodes()))

	// delete node for cluster1
	_ = repo2.Delete(context.TODO(), constants.ActiveNodesPath+"/127.0.0.2:2080")
	time.Sleep(200 * time.Millisecond)
	c.Assert(1, check.Equals, len(cluster1.GetActiveNodes()))
	// delete cluster1
	_ = repo.Delete(context.TODO(), constants.StorageClusterConfigPath+"/storage1")
	time.Sleep(200 * time.Millisecond)
	c.Check(stateMachine.GetCluster("storage1"), check.IsNil)

	// cleanup nodes when delete storage config
	c.Assert(0, check.Equals, len(cluster1.GetActiveNodes()))

	// add node into cluster2
	repo3, _ := state.NewRepo(state.Config{
		Namespace: "storage2",
		Endpoints: ts.Cluster.Endpoints,
	})
	_ = repo3.Put(context.TODO(), constants.ActiveNodesPath+"/node1", node1)
	_ = repo3.Put(context.TODO(), constants.ActiveNodesPath+"/node2", node2)
	time.Sleep(200 * time.Millisecond)
	c.Assert(2, check.Equals, len(cluster2.GetActiveNodes()))

	_ = stateMachine.Close()
}

func (ts *testClusterStateMachineSuite) TestExistNodes(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/test/exist/nodes",
		Endpoints: ts.Cluster.Endpoints,
	})
	storage1 := state.Config{
		Namespace: "/test/exist/nodes",
		Endpoints: ts.Cluster.Endpoints,
	}
	data1, _ := json.Marshal(models.StorageCluster{
		Name:   "storage1",
		Config: storage1,
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage1", data1)
	node1, _ := json.Marshal(models.ActiveNode{Node: models.Node{IP: "127.0.0.1", Port: 2080}})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node1", node1)
	node2, _ := json.Marshal(models.ActiveNode{Node: models.Node{IP: "127.0.0.2", Port: 2080}})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node2", node2)
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node3", []byte("dd"))

	stateMachine, _ := NewClusterStateMachine(context.TODO(), repo)

	c.Assert(1, check.Equals, len(stateMachine.GetAllCluster()))
	cluster1 := stateMachine.GetCluster("storage1")
	c.Assert(cluster1, check.NotNil)

	c.Assert(2, check.Equals, len(cluster1.GetActiveNodes()))
	_ = stateMachine.Close()
}
