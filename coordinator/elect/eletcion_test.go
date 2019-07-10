package elect

import (
	"context"
	"testing"
	"time"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	"github.com/coreos/pkg/capnslog"
	"gopkg.in/check.v1"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

type testElectionSuite struct {
	mock.RepoTestSuite
}

var _ = check.Suite(&testElectionSuite{})

func TestElection(t *testing.T) {
	check.TestingT(t)
}

func (ts *testElectionSuite) TestElect(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})
	node := models.Node{IP: "127.0.0.1", Port: 2080}
	election := NewElection(repo, node, "test", 1)
	ctx, cancel := context.WithCancel(context.Background())
	// first node election register must be success
	success, err := election.Elect(ctx)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(success, check.Equals, true)
	node2 := models.Node{IP: "127.0.0.2", Port: 2080}

	// second node election should be false
	election2 := NewElection(repo, node2, "test", 1)
	ctx2, cancel2 := context.WithCancel(context.Background())
	success2, _ := election2.Elect(ctx2)
	c.Assert(success2, check.Equals, false)
	isMaster := election2.IsMaster()
	c.Assert(isMaster, check.Equals, false)
	cancel()
	// first node exist,the second node should be the master
	time.Sleep(2 * time.Second)
	isMaster2 := election2.IsMaster()
	c.Assert(isMaster2, check.Equals, true)

	defer cancel2()
}
