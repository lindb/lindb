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
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestStateMachineCli_FetchStateByNode(t *testing.T) {
	cases := []struct {
		name    string
		port    int
		prepare func(rw http.ResponseWriter)
		wantErr bool
	}{
		{
			name: "fetch failure",
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: false,
		},
		{
			name: "fetch successfully",
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
			},
			wantErr: false,
		},
		{
			name:    "url wrong",
			port:    30001,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if tt.prepare != nil {
					tt.prepare(rw)
				}
			}))
			defer server.Close()
			port := strings.Split(server.URL, ":")[2]
			cli := NewStateMachineCli()
			p, _ := strconv.Atoi(port)
			if tt.port > 0 {
				p = tt.port
			}
			_, err := cli.FetchStateByNode(nil, &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: uint16(p)})
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchStateByNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStateMachineCli_FetchStateByNodes(t *testing.T) {
	cases := []struct {
		name    string
		port    int
		nodes   []models.Node
		prepare func(rw http.ResponseWriter)
	}{
		{
			name:  "url wrong",
			nodes: []models.Node{&models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: 30001}},
		},
		{
			name: "no nodes",
		},
		{
			name:  "fetch data successfully",
			nodes: []models.Node{&models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: 8080}},
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if tt.prepare != nil {
					tt.prepare(rw)
				}
			}))
			l, err := net.Listen("tcp", "127.0.0.1:8080")
			assert.NoError(t, err)
			ts.Listener = l
			ts.Start()
			defer func() {
				ts.Close()
				_ = ts.Listener.Close()
			}()

			cli := NewStateMachineCli()
			_ = cli.FetchStateByNodes(nil, tt.nodes)
		})
	}
}
