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
	"strings"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/ltoml"
)

// Configuration represents node's configuration.
type Configuration interface {
	// TOML returns configuration string as toml format.
	TOML() string
}

// RepoState represents state repository config
type RepoState struct {
	Namespace   string         `env:"NAMESPACE" toml:"namespace" json:"namespace" validate:"required"`
	Endpoints   []string       `env:"ENDPOINTS" envSeparator:"," toml:"endpoints" json:"endpoints" validate:"required,gt=0"`
	LeaseTTL    ltoml.Duration `env:"LEASE_TTL" toml:"lease-ttl" json:"leaseTTL"`
	Timeout     ltoml.Duration `env:"TIMEOUT" toml:"timeout" json:"timeout"`
	DialTimeout ltoml.Duration `env:"DIAL_TIMEOUT" toml:"dial-timeout" json:"dialTimeout"`
	Username    string         `env:"USERNAME" toml:"username" json:"username"`
	Password    string         `env:"PASSWORD" toml:"password" json:"password"`
}

// String returns string value of RepoState.
func (rs *RepoState) String() string {
	return fmt.Sprintf("endpoints:[%s],leaseTTL:%s,timeout:%s,dialTimeout:%s",
		strings.Join(rs.Endpoints, ","), rs.LeaseTTL, rs.Timeout, rs.DialTimeout)
}

func (rs *RepoState) WithSubNamespace(subDir string) *RepoState {
	return &RepoState{
		Namespace:   rs.Namespace + constants.StatePathSeparator + subDir,
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
## Default: %s
## Env: COORDINATOR_NAMESPACE
namespace = "%s"
## Endpoints config list of ETCD cluster
## Default: %s
## Env: COORDINATOR_NAMESPACE  Env Separator: ,
endpoints = %s
## Lease-TTL is a number in seconds.
## It controls how long a ephemeral node like zookeeper will be removed when heartbeat fails.
## lease expiration will cause a re-elect.
## Min: 5s
## Default: %s
## Env: COORDINATOR_LEASE_TTL
lease-ttl = "%s"
## Timeout is the timeout for failing to executing a etcd command.
## Default: %s
## Env: COORDINATOR_LEASE_TTL
timeout = "%s"
## DialTimeout is the timeout for failing to establish a etcd connection.
## Default: %s
## Env: COORDINATOR_DIAL_TIMEOUT
dial-timeout = "%s"
## Username is a user name for etcd authentication.
## Default: "%s"
## Env: COORDINATOR_USERNAME
username = "%s"
## Password is a password for etcd authentication.
## Default: "%s"
## Env: COORDINATOR_PASSWORD
password = "%s"`,
		rs.Namespace,
		rs.Namespace,
		coordinatorEndpoints,
		coordinatorEndpoints,
		rs.LeaseTTL.String(),
		rs.LeaseTTL.String(),
		rs.Timeout.String(),
		rs.Timeout.String(),
		rs.DialTimeout.String(),
		rs.DialTimeout.String(),
		rs.Username,
		rs.Username,
		rs.Password,
		rs.Password,
	)
}

func NewDefaultCoordinator() *RepoState {
	return &RepoState{
		Namespace:   "/lindb-cluster",
		Endpoints:   []string{"http://localhost:2379"},
		LeaseTTL:    ltoml.Duration(time.Second * 10),
		Timeout:     ltoml.Duration(time.Second * 5),
		DialTimeout: ltoml.Duration(time.Second * 5),
	}
}

// GRPC represents grpc server config
type GRPC struct {
	Port                 uint16         `env:"PORT" toml:"port"`
	MaxConcurrentStreams int            `env:"MAX_CONCURRENT_STREAMS" toml:"max-concurrent-streams"`
	ConnectTimeout       ltoml.Duration `env:"CONNECT_TIMEOUT" toml:"connect-timeout"`
}

func (g *GRPC) TOML() string {
	return fmt.Sprintf(`
## port which the GRPC Server is listening on
## Default: %d
## Env: BROKER_GRPC_PORT
## Env: STORAGE_GRPC_PORT
port = %d
## max-concurrent-streams limits the number of concurrent streams to each ServerTransport
## Default: %d 
## Env: BROKER_GRPC_MAX_CONCURRENT_STREAMS
## Env: STORAGE_GRPC_MAX_CONCURRENT_STREAMS
max-concurrent-streams = %d
## connect-timeout sets the timeout for connection establishment.
## Default: %s
## Env: BROKER_GRPC_CONNECT_TIMEOUT
## Env: STORAGE_GRPC_CONNECT_TIMEOUT
connect-timeout = "%s"`,
		g.Port,
		g.Port,
		g.MaxConcurrentStreams,
		g.MaxConcurrentStreams,
		g.ConnectTimeout.Duration().String(),
		g.ConnectTimeout.Duration().String(),
	)
}

// BrokerCluster represents config of broker cluster.
type BrokerCluster struct {
	Config *RepoState `json:"config"`
}

// StorageCluster represents config of storage cluster.
type StorageCluster struct {
	Config *RepoState `json:"config"`
}

// Query represents query rpc config
type Query struct {
	QueryConcurrency int            `env:"CONCURRENCY" toml:"query-concurrency"`
	IdleTimeout      ltoml.Duration `env:"IDLE_TIMEOUT" toml:"idle-timeout"`
	Timeout          ltoml.Duration `env:"TIMEOUT" toml:"timeout"`
}

func (q *Query) TOML() string {
	return fmt.Sprintf(`[query]
## Number of queries allowed to execute concurrently
## Default: %d
## Env: QUERY_CONCURRENCY
query-concurrency = %d
## Idle worker will be canceled in this duration
## Default: %s
## Env: QUERY_IDLE_TIMEOUT
idle-timeout = "%s"
## Maximum timeout threshold for query.
## Default: %s
## Env: QUERY_TIMEOUT
timeout = "%s"`,
		q.QueryConcurrency,
		q.QueryConcurrency,
		q.IdleTimeout,
		q.IdleTimeout,
		q.Timeout,
		q.Timeout,
	)
}

func NewDefaultQuery() *Query {
	return &Query{
		QueryConcurrency: 1024,
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
		grpcCfg.MaxConcurrentStreams = 1024
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
