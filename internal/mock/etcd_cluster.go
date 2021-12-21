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
	"fmt"
	"net/url"
	"testing"

	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap/zapcore"
)

// EtcdCluster mock etcd cluster for testing
type EtcdCluster struct {
	etcd      *embed.Etcd
	Endpoints []string
}

// StartEtcdCluster starts integration etcd cluster
func StartEtcdCluster(t *testing.T, endpoint string) *EtcdCluster {
	cfg := embed.NewConfig()
	lcurl, _ := url.Parse(endpoint)
	acurl, _ := url.Parse(fmt.Sprintf("http://localhost:1%s", lcurl.Port()))
	cfg.Dir = t.TempDir()
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.LPUrls = []url.URL{*acurl}
	cfg.LogLevel = zapcore.ErrorLevel.String()
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		panic(err)
	}
	return &EtcdCluster{
		etcd:      e,
		Endpoints: []string{endpoint},
	}
}

// Terminate terminates integration etcd cluster
func (etcd *EtcdCluster) Terminate(_ *testing.T) {
	etcd.etcd.Close()
}
