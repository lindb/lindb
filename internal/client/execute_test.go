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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/models"
)

func TestExecuteCli_Execute(t *testing.T) {
	cases := []struct {
		name    string
		param   models.ExecuteParam
		url     string
		rs      interface{}
		prepare func(rw http.ResponseWriter)
		wantErr bool
	}{
		{
			name:    "wrong url",
			url:     "http://localhost:30001",
			wantErr: true,
		},
		{
			name:    "no data return",
			wantErr: false,
		},
		{
			name: "http status no ok",
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name:  "unmarshal result failure",
			param: models.ExecuteParam{SQL: "show master"},
			rs:    &models.Master{},
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
				if tt.prepare != nil {
					tt.prepare(rw)
				}
			}))
			defer server.Close()

			cli := NewExecuteCli(server.URL)
			if len(tt.url) > 0 {
				cli = NewExecuteCli(tt.url)
			}
			err := cli.Execute(tt.param, &tt.rs)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteCli_ExecuteAsResult(t *testing.T) {
	cases := []struct {
		name    string
		param   models.ExecuteParam
		rs      interface{}
		prepare func(rw http.ResponseWriter)
		assert  func(rs string)
		wantErr bool
	}{
		{
			name: "http status no ok",
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "no data return",
			assert: func(rs string) {
				assert.True(t, strings.HasPrefix(rs, "Query OK,"))
			},
			wantErr: false,
		},
		{
			name: "not format as table",
			rs:   &[]models.ExecuteParam{},
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(encoding.JSONMarshal(&[]models.ExecuteParam{{SQL: "sql"}}))
			},
			assert: func(rs string) {
				fmt.Println(rs)
				assert.True(t, strings.Contains(rs, "0 rows"))
			},
			wantErr: false,
		},
		{
			name: "format as table",
			rs:   &models.Databases{},
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(encoding.JSONMarshal(&models.Databases{{Name: "test"}}))
			},
			assert: func(rs string) {
				_, s := (&models.Databases{{Name: "test"}}).ToTable()
				assert.True(t, strings.Contains(rs, s))
			},
			wantErr: false,
		},
		{
			name: "object format as table",
			rs:   &models.Master{},
			prepare: func(rw http.ResponseWriter) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(encoding.JSONMarshal(&models.Master{ElectTime: timeutil.Now(), Node: &models.StatelessNode{}}))
			},
			assert: func(rs string) {
				assert.True(t, strings.Contains(rs, "1 row"))
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
				if tt.prepare != nil {
					tt.prepare(rw)
				}
			}))
			defer server.Close()

			cli := NewExecuteCli(server.URL)
			rs, err := cli.ExecuteAsResult(tt.param, tt.rs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(rs)
			}
		})
	}
}
