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

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

func TestRequst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMgr := broker.NewMockStateManager(ctrl)
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
			name:      "show all requests, but no alive broker",
			statement: &stmt.Request{},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return(nil)
			},
		},
		{
			name:      "show all requests, but get err from broker",
			statement: &stmt.Request{},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "127.0.0.1",
					HTTPPort: 3000,
				}})
			},
		},
		{
			name:      "show all requests successfully",
			statement: &stmt.Request{},
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`[{"start":12314},{"start":12314}]`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "127.0.0.1",
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
			rs, err := RequestCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("RequestCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
