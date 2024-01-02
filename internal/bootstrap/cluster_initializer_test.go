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

package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lindb/lindb/models"
)

func TestClusterInitializer_InitDatabase(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(w http.ResponseWriter)
		wantErr bool
	}{
		{
			name: "create database successfully",
		},
		{
			name: "create database failure",
			prepare: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					if tt.prepare != nil {
						tt.prepare(w)
					}
				}))
			defer ts.Close()

			init := NewClusterInitializer(ts.URL)

			err := init.InitInternalDatabase(models.Database{})
			if (err != nil) != tt.wantErr {
				t.Errorf("InitStorageCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
