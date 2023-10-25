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

package command

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

func TestState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := root.NewMockStateManager(ctrl)
	deps := &depspkg.HTTPDeps{
		StateMgr: stateMgr,
	}

	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "unknown metadata statement type",
			statement: &stmt.State{},
		},
		{
			name:      "show root alive node",
			statement: &stmt.State{Type: stmt.RootAlive},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "1.1.1.1",
					HTTPPort: 8080,
				}})
			},
		},
		{
			name:      "show broker alive node",
			statement: &stmt.State{Type: stmt.BrokerAlive},
			prepare: func() {
				stateMgr.EXPECT().GetBrokerStates().Return([]models.BrokerState{{
					Name: "test",
				}})
			},
		},
		{
			name:      "show root metric, no alive node",
			statement: &stmt.State{Type: stmt.RootMetric, MetricNames: []string{"a", "b"}},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return(nil)
			},
		},
		{
			name:      "show root metric successfully",
			statement: &stmt.State{Type: stmt.RootMetric, MetricNames: []string{"a", "b"}},
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`{"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}, {
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}})
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := StateCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("StateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
