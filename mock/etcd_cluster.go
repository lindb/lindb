package mock

import (
	"testing"

	"gopkg.in/check.v1"

	"github.com/coreos/etcd/integration"
)

// EtcdCluster mock etcd cluster for testing
type EtcdCluster struct {
	cluster   *integration.ClusterV3
	Endpoints []string
}

// StartEtcdCluster starts integration etcd cluster
func StartEtcdCluster(t *testing.T) *EtcdCluster {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	return &EtcdCluster{
		cluster:   cluster,
		Endpoints: []string{cluster.Members[0].GRPCAddr()},
	}
}

// Terminate terminates integration etcd cluster
func (etcd *EtcdCluster) Terminate(t *testing.T) {
	etcd.cluster.Terminate(t)
}

// RepoTestSuite represents repo test suite init integration etcd cluster
type RepoTestSuite struct {
	Cluster *EtcdCluster
}

var (
	test *testing.T
)

// Test register test suite
func Test(t *testing.T) {
	test = t
	check.TestingT(t)
}

// SetUpSuite setups test suite
func (ts *RepoTestSuite) SetUpSuite(c *check.C) {
	ts.Cluster = StartEtcdCluster(test)
}

// TearDownSuite teardowns test suite, release resource
func (ts *RepoTestSuite) TearDownSuite(c *check.C) {
	ts.Cluster.Terminate(test)
}
