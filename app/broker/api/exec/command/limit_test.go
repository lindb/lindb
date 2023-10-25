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
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/golang/mock/gomock"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo: repo,
	}
	cases := []struct {
		name      string
		db        string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{
			name:    "database name not input",
			wantErr: true,
		},
		{
			name:      "unknow limit op",
			db:        "test",
			statement: &stmt.Limit{},
		},
		{
			name:      "invalid toml",
			db:        "test",
			statement: &stmt.Limit{Limit: "test", Type: stmt.SetLimit},
			prepare: func() {
				tomlDecodeFn = func(data string, v interface{}) (toml.MetaData, error) {
					return toml.MetaData{}, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name:      "save limit failure",
			db:        "test",
			statement: &stmt.Limit{Limit: "test", Type: stmt.SetLimit},
			prepare: func() {
				tomlDecodeFn = func(data string, v interface{}) (toml.MetaData, error) {
					return toml.MetaData{}, nil
				}
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "save limit successfully",
			db:        "test",
			statement: &stmt.Limit{Limit: "test", Type: stmt.SetLimit},
			prepare: func() {
				tomlDecodeFn = func(data string, v interface{}) (toml.MetaData, error) {
					return toml.MetaData{}, nil
				}
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:      "get default limits",
			db:        "test",
			statement: &stmt.Limit{Type: stmt.ShowLimit},
			prepare: func() {
				repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
			},
		},
		{
			name:      "get limits failure",
			db:        "test",
			statement: &stmt.Limit{Type: stmt.ShowLimit},
			prepare: func() {
				repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "get limits successfully",
			db:        "test",
			statement: &stmt.Limit{Type: stmt.ShowLimit},
			prepare: func() {
				repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("test"), nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				tomlDecodeFn = toml.Decode
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := LimitCommand(context.TODO(), deps, &models.ExecuteParam{Database: tt.db}, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("LimitCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
