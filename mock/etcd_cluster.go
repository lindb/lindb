// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mock

import (
	"testing"

	"go.etcd.io/etcd/integration"
	"gopkg.in/check.v1"
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
