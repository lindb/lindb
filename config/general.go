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

package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// Configuration represents node's configuration.
type Configuration interface {
	// TOML returns configuration string as toml format.
	TOML() string
}

// RepoState represents state repository config
type RepoState struct {
	Namespace   string         `toml:"namespace" json:"namespace"`
	Endpoints   []string       `toml:"endpoints" json:"endpoints"`
	LeaseTTL    int64          `toml:"lease-ttl" json:"leaseTTL"`
	Timeout     ltoml.Duration `toml:"timeout" json:"timeout"`
	DialTimeout ltoml.Duration `toml:"dial-timeout" json:"dialTimeout"`
	Username    string         `toml:"username" json:"username"`
	Password    string         `toml:"password" json:"password"`
}

func (rs *RepoState) WithSubNamespace(subDir string) RepoState {
	return RepoState{
		Namespace:   filepath.Join(rs.Namespace, subDir),
		Endpoints:   rs.Endpoints,
		Timeout:     rs.Timeout,
		DialTimeout: rs.DialTimeout,
		Username:    rs.Username,
		Password:    rs.Password,
	}
}

// TOML returns RepoState's toml config string
func (rs *RepoState) TOML() string {
	coordinatorEndpoints, _ := json.Marshal(rs.Endpoints)
	return fmt.Sprintf(`[coordinator]
## Coordinator coordinates reads/writes operations between different nodes
## namespace organizes etcd keys into a isolated complete keyspaces for coordinator
namespace = "%s"
## Endpoints config list of ETCD cluster
endpoints = %s
## Lease-TTL is a number in seconds.
## It controls how long a ephemeral node like zookeeper will be removed when heartbeat fails.
## lease expiration will cause a re-elect.
## Min: 5； Default: 10
lease-ttl = %d
## Timeout is the timeout for failing to executing a etcd command.
## Default: 5s
timeout = "%s"
## DialTimeout is the timeout for failing to establish a etcd connection.
## Default: 5s
dial-timeout = "%s"
## Username is a user name for etcd authentication.
username = "%s"
## Password is a password for etcd authentication.
password = "%s"`,
		rs.Namespace,
		coordinatorEndpoints,
		rs.LeaseTTL,
		rs.Timeout.String(),
		rs.DialTimeout.String(),
		rs.Username,
		rs.Password,
	)
}

func NewDefaultCoordinator() *RepoState {
	return &RepoState{
		Namespace:   "/lindb-cluster",
		Endpoints:   []string{"http://localhost:2379"},
		LeaseTTL:    10,
		Timeout:     ltoml.Duration(time.Second * 5),
		DialTimeout: ltoml.Duration(time.Second * 5),
	}
}

// GRPC represents grpc server config
type GRPC struct {
	Port                 uint16         `toml:"port"`
	MaxConcurrentStreams int            `toml:"max-concurrent-streams"`
	ConnectTimeout       ltoml.Duration `toml:"connect-timeout"`
}

func (g *GRPC) TOML() string {
	return fmt.Sprintf(`
port = %d
## max-concurrent-streams limits the number of concurrent streams to each ServerTransport
## Default: runtime.GOMAXPROCS(-1) * 2
max-concurrent-streams = %d
## connect-timeout sets the timeout for connection establishment.
## Default: 3s
connect-timeout = "%s"`,
		g.Port,
		g.MaxConcurrentStreams,
		g.ConnectTimeout.Duration().String(),
	)
}

// StorageCluster represents config of storage cluster
type StorageCluster struct {
	Name   string    `json:"name" binding:"required"`
	Config RepoState `json:"config"`
}

// Query represents query rpc config
type Query struct {
	QueryConcurrency int            `toml:"query-concurrency"`
	IdleTimeout      ltoml.Duration `toml:"idle-timeout"`
	Timeout          ltoml.Duration `toml:"timeout"`
}

func (q *Query) TOML() string {
	return fmt.Sprintf(`
[query]
## Number of queries allowed to execute concurrently
## Default: runtime.GOMAXPROCS(-1) * 2
query-concurrency = %d
## Idle worker will be canceled in this duration
## Default: 5s
idle-timeout = "%s"
## Maximum timeout threshold for query.
## Default: 5s
timeout = "%s"`,
		q.QueryConcurrency,
		q.IdleTimeout,
		q.Timeout,
	)
}

func NewDefaultQuery() *Query {
	return &Query{
		QueryConcurrency: runtime.GOMAXPROCS(-1) * 2,
		IdleTimeout:      ltoml.Duration(5 * time.Second),
		Timeout:          ltoml.Duration(5 * time.Second),
	}
}

func checkCoordinatorCfg(state *RepoState) error {
	if state.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if state.LeaseTTL < 5 {
		state.LeaseTTL = 5
	}
	if len(state.Endpoints) == 0 {
		return fmt.Errorf("endpoints cannot be empty")
	}
	if state.Timeout <= 0 {
		state.Timeout = ltoml.Duration(time.Second * 5)
	}
	if state.DialTimeout <= 0 {
		state.Timeout = ltoml.Duration(time.Second * 5)
	}
	return nil
}

func checkGRPCCfg(grpcCfg *GRPC) error {
	if grpcCfg.Port == 0 {
		return fmt.Errorf("grpc endpoint cannot be empty")
	}
	if grpcCfg.MaxConcurrentStreams <= 0 {
		grpcCfg.MaxConcurrentStreams = runtime.GOMAXPROCS(-1) * 2
	}
	if grpcCfg.ConnectTimeout <= 0 {
		grpcCfg.ConnectTimeout = ltoml.Duration(time.Second * 3)
	}
	return nil
}

func checkQueryCfg(queryCfg *Query) {
	defaultQuery := NewDefaultQuery()
	if queryCfg.QueryConcurrency <= 0 {
		queryCfg.QueryConcurrency = defaultQuery.QueryConcurrency
	}
	if queryCfg.Timeout <= 0 {
		queryCfg.Timeout = defaultQuery.Timeout
	}
	if queryCfg.IdleTimeout <= 0 {
		queryCfg.IdleTimeout = defaultQuery.IdleTimeout
	}
}
