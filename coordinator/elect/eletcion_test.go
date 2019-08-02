package elect

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"

	check "gopkg.in/check.v1"
)

type mockListener struct {
	onFailOverCount    int32
	onResignationCount int32
}

func newMockListener() Listener {
	return &mockListener{}
}

func (l *mockListener) OnResignation() {
	atomic.AddInt32(&l.onResignationCount, 1)
}

func (l *mockListener) OnFailOver() {
	atomic.AddInt32(&l.onFailOverCount, 1)
}

type testElectionSuite struct {
	mock.RepoTestSuite
}

func TestElection(t *testing.T) {
	check.Suite(&testElectionSuite{})
	check.TestingT(t)
}

func (ts *testElectionSuite) TestElect(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})
	listener1 := newMockListener()
	l1, _ := listener1.(*mockListener)
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	election := NewElection(repo, node1, 1, listener1)
	election.Initialize()
	election.Elect()

	time.Sleep(500 * time.Millisecond)

	c.Assert(int32(1), check.Equals, atomic.LoadInt32(&l1.onFailOverCount))

	// second node election should be false
	node2 := models.Node{IP: "127.0.0.2", Port: 2080}
	listener2 := newMockListener()
	l2, _ := listener2.(*mockListener)

	election2 := NewElection(repo, node2, 1, listener2)
	election2.Initialize()
	election2.Elect()

	// cancel first node
	election.Close()

	time.Sleep(500 * time.Millisecond)
	// second node become master
	c.Assert(int32(1), check.Equals, atomic.LoadInt32(&l2.onFailOverCount))

	election2.Close()
}
