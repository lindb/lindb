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
)

func TestRequestCli_FetchRequestsByNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("content-type", "application/json")
		_, _ = w.Write([]byte(`[{"start":12314},{"start":12314}]`))
	}))
	u, err := url.Parse(svr.URL)
	assert.NoError(t, err)
	p, err := strconv.Atoi(u.Port())
	assert.NoError(t, err)
	nodes := []models.Node{&models.StatelessNode{
		HostIP:   "127.0.0.1",
		HTTPPort: uint16(p),
	}}

	cases := []struct {
		name    string
		nodes   []models.Node
		wantHas bool
	}{
		{
			name:  "show all requests, but no alive broker",
			nodes: nil,
		},
		{
			name: "show all requests, but get err from broker",
			nodes: []models.Node{&models.StatelessNode{
				HostIP:   "127.0.0.1",
				HTTPPort: 3000,
			}},
		},
		{
			name:    "show all requests successfully",
			nodes:   nodes,
			wantHas: true,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cli := NewRequestCli()
			rs := cli.FetchRequestsByNodes(tt.nodes)
			if (rs != nil) != tt.wantHas && rs == nil {
				t.Errorf("FetchRequestsByNodes() failure")
			}
		})
	}
}
