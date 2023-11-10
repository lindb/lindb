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

package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetricCli_FetchMetricData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func() []models.Node
		wantErr   bool
	}{
		{
			name: "empty node",
			prepare: func() []models.Node {
				return nil
			},
		},
		{
			name: "fetch metric failure",
			prepare: func() []models.Node {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				return []models.Node{&models.StatelessNode{
					HostIP:   u.Hostname(),
					HTTPPort: 9000,
				}}
			},
		},
		{
			name:      "fetch metric successfully",
			statement: &stmt.State{Type: stmt.BrokerMetric, MetricNames: []string{"a", "b"}},
			prepare: func() []models.Node {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`{"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				return []models.Node{
					&models.StatelessNode{
						HostIP:   u.Hostname(),
						HTTPPort: uint16(p),
					},
					&models.StatelessNode{
						HostIP:   u.Hostname(),
						HTTPPort: uint16(p),
					},
				}
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			nodes := tt.prepare()
			cli := NewMetricCli()
			rs, err := cli.FetchMetricData(nodes, []string{"cpu"})
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("FetchMetricData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
