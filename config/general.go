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
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// RepoState represents state repository config
type RepoState struct {
	Namespace   string         `toml:"namespace" json:"namespace"`
	Endpoints   []string       `toml:"endpoints" json:"endpoints"`
	Timeout     ltoml.Duration `toml:"timeout" json:"timeout"`
	DialTimeout ltoml.Duration `toml:"dial-timeout" json:"dialTimeout"`
	Username    string         `toml:"username" json:"username"`
	Password    string         `toml:"password" json:"password"`
}

// TOML returns RepoState's toml config string
func (rs *RepoState) TOML() string {
	coordinatorEndpoints, _ := json.Marshal(rs.Endpoints)
	return fmt.Sprintf(`
    ## Coordinator coordinates reads/writes operations between different nodes
    ## namespace organizes etcd keys into a isolated complete keyspaces for coordinator
    namespace = "%s"
    ## endpoints config list of ETCD cluster
    endpoints = %s
	## Timeout is the timeout for failing to executing a etcd command.
	timeout = "%s"
	## DialTimeout is the timeout for failing to establish a etcd connection.
    dial-timeout = "%s"
	## Username is a user name for etcd authentication.
	username = "%s"
	## Password is a password for etcd authentication.
	password = "%s"`,
		rs.Namespace,
		coordinatorEndpoints,
		rs.Timeout.String(),
		rs.DialTimeout.String(),
		rs.Username,
		rs.Password,
	)
}

// GRPC represents grpc server config
type GRPC struct {
	Port uint16         `toml:"port"`
	TTL  ltoml.Duration `toml:"ttl"`
}

func (g *GRPC) TOML() string {
	return fmt.Sprintf(`
    port = %d
    ttl = "%s"`,
		g.Port,
		g.TTL.String(),
	)
}

// StorageCluster represents config of storage cluster
type StorageCluster struct {
	Name   string    `json:"name" binding:"required"`
	Config RepoState `json:"config"`
}

// Query represents query rpc config
type Query struct {
	MaxWorkers  int            `toml:"max-workers"`
	IdleTimeout ltoml.Duration `toml:"idle-timeout"`
	Timeout     ltoml.Duration `toml:"timeout"`
}

func (q *Query) TOML() string {
	return fmt.Sprintf(`
    ## max concurrentcy number of workers in the executor pool,
    ## each worker is only responsible for execting querying task,
    ## and idle workers will be recycled.
    max-workers = %d

    ## idle worker will be canceled in this duration
    idle-timeout = "%s"

    ## maximum timeout threshold for the task performed
    timeout = "%s"`,
		q.MaxWorkers,
		q.IdleTimeout,
		q.Timeout,
	)
}

func NewDefaultQuery() *Query {
	return &Query{
		MaxWorkers:  30,
		IdleTimeout: ltoml.Duration(5 * time.Second),
		Timeout:     ltoml.Duration(30 * time.Second),
	}
}
