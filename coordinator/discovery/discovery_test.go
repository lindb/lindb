package discovery

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"

	"gopkg.in/check.v1"
)

type mockListener struct {
	nodes map[string][]byte
	mutex sync.Mutex
}

func newMockListener() Listener {
	return &mockListener{
		nodes: make(map[string][]byte),
	}
}

func (m *mockListener) OnCreate(key string, value []byte) {
	m.mutex.Lock()
	m.nodes[key] = value
	m.mutex.Unlock()
}

func (m *mockListener) OnDelete(key string) {
	m.mutex.Lock()
	delete(m.nodes, key)
	m.mutex.Unlock()
}

func (m *mockListener) Cleanup() {
	m.mutex.Lock()
	m.nodes = make(map[string][]byte)
	m.mutex.Unlock()
}

type testDiscoverySuite struct {
	mock.RepoTestSuite
}

var testDiscoveryPath = "/test/discovery1"
var testDiscoveryPath2 = "/test/discovery2"

func TestDiscovery(t *testing.T) {
	check.Suite(&testDiscoverySuite{})
	check.TestingT(t)
}

func (ts *testDiscoverySuite) TestNodeList(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	listener := newMockListener()
	d := NewDiscovery(repo, testDiscoveryPath, listener)
	_ = d.Discovery()

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	nodeBytes, _ := json.Marshal(node)
	_ = repo.Put(context.TODO(), "/test/discovery1/key1", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery1/key2", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery1/key3", nodeBytes)
	time.Sleep(500 * time.Millisecond)

	mockListener, _ := listener.(*mockListener)
	mockListener.mutex.Lock()
	nodes := mockListener.nodes
	c.Assert(3, check.Equals, len(nodes))
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery1/key1"])
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery1/key2"])
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery1/key3"])
	mockListener.mutex.Unlock()

	_ = repo.Delete(context.TODO(), "/test/discovery1/key2")
	time.Sleep(500 * time.Millisecond)

	mockListener.mutex.Lock()
	nodes = mockListener.nodes
	c.Assert(2, check.Equals, len(nodes))
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery1/key1"])
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery1/key3"])
	mockListener.mutex.Unlock()

	d.Close()
	mockListener.mutex.Lock()
	nodes = mockListener.nodes
	c.Assert(0, check.Equals, len(nodes))
	mockListener.mutex.Unlock()
}

func (ts *testDiscoverySuite) TestDiscoveryErr(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	listener := newMockListener()
	d := NewDiscovery(repo, "", listener)
	err := d.Discovery()
	c.Assert(err, check.NotNil)

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	nodeBytes, _ := json.Marshal(node)
	_ = repo.Put(context.TODO(), "/test/discovery1/key1", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery1/key2", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery1/key3", nodeBytes)
	time.Sleep(500 * time.Millisecond)

	mockListener, _ := listener.(*mockListener)

	mockListener.mutex.Lock()
	c.Assert(0, check.Equals, len(mockListener.nodes))
	mockListener.mutex.Unlock()

	d.Close()
}

func (ts *testDiscoverySuite) TestExistNodeList(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	listener := newMockListener()
	d := NewDiscovery(repo, testDiscoveryPath2, listener)

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	nodeBytes, _ := json.Marshal(node)
	_ = repo.Put(context.TODO(), "/test/discovery2/key1", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery2/key2", nodeBytes)
	_ = repo.Put(context.TODO(), "/test/discovery2/key3", nodeBytes)

	mockListener, _ := listener.(*mockListener)

	// no data before start discovery
	mockListener.mutex.Lock()
	nodes := mockListener.nodes
	c.Assert(0, check.Equals, len(nodes))
	mockListener.mutex.Unlock()

	_ = d.Discovery()
	time.Sleep(500 * time.Millisecond)

	mockListener.mutex.Lock()
	nodes = mockListener.nodes
	c.Assert(3, check.Equals, len(nodes))
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery2/key1"])
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery2/key2"])
	c.Assert(nodeBytes, check.DeepEquals, nodes["/test/discovery2/key3"])
	mockListener.mutex.Unlock()

	d.Close()
}
