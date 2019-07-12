package coordinator

import (
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

type testMasterSuite struct {
	mock.RepoTestSuite
}

func TestMaster(t *testing.T) {
	check.Suite(&testMasterSuite{})
	check.TestingT(t)
}

func (ts *testMasterSuite) TestMasterElect(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/master/test",
		Endpoints: ts.Cluster.Endpoints,
	})
	node1 := models.Node{IP: "1.1.1.1", Port: 8000}
	master1 := NewMaster(repo, node1, 1)
	_ = master1.Start()
	time.Sleep(400 * time.Millisecond)
	c.Assert(true, check.Equals, master1.IsMaster())

	repo2, _ := state.NewRepo(state.Config{
		Namespace: "/master/test",
		Endpoints: ts.Cluster.Endpoints,
	})
	node2 := models.Node{IP: "1.1.1.2", Port: 8000}
	master2 := NewMaster(repo2, node2, 1)
	_ = master2.Start()
	time.Sleep(400 * time.Millisecond)
	c.Assert(false, check.Equals, master2.IsMaster())

	master1.Stop()
	time.Sleep(400 * time.Millisecond)
	c.Assert(false, check.Equals, master1.IsMaster())
	c.Assert(true, check.Equals, master2.IsMaster())
}
