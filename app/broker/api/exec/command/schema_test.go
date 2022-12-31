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

	"github.com/golang/mock/gomock"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opt := &option.DatabaseOption{}
	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo: repo,
	}
	databaseCfg := `{"name":"test","storage":"cluster-test","numOfShard":12,`
	databaseCfg += `"replicaFactor":3,"option":{"intervals":[{"interval":"10s"}]}}`

	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{

			name:      "create database config unmarshal failure",
			statement: &stmt.Schema{Type: stmt.CreateDatabaseSchemaType, Value: `err`},
			wantErr:   true,
		},
		{
			name:      "create database validation failure",
			statement: &stmt.Schema{Type: stmt.CreateDatabaseSchemaType, Value: `{"name":"name"}`},
			wantErr:   true,
		},
		{
			name:      "create database, persist failure",
			statement: &stmt.Schema{Type: stmt.CreateDatabaseSchemaType, Value: databaseCfg},
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create database, option validation failure",
			statement: &stmt.Schema{
				Type: stmt.CreateDatabaseSchemaType,
				Value: string(encoding.JSONMarshal(&models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
					Option: &option.DatabaseOption{
						Intervals: option.Intervals{{Interval: 10}},
						Ahead:     "10",
					},
				})),
			},
			wantErr: true,
		},
		{
			name:      "create database successfully",
			statement: &stmt.Schema{Type: stmt.CreateDatabaseSchemaType, Value: databaseCfg},
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:      "drop database, but delete cfg failure",
			statement: &stmt.Schema{Type: stmt.DropDatabaseSchemaType, Value: "test"},
			prepare: func() {
				// delete database cfg failure
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "drop database, but delete shard assignment failure",
			statement: &stmt.Schema{Type: stmt.DropDatabaseSchemaType, Value: "test"},
			prepare: func() {
				// delete database cfg ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				// delete database shard assignment failure
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "drop database successfully",
			statement: &stmt.Schema{Type: stmt.DropDatabaseSchemaType, Value: "test"},
			prepare: func() {
				// delete database cfg ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				// delete database shard assignment ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:      "get database list err",
			statement: &stmt.Schema{Type: stmt.DatabaseNameSchemaType},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "get database successfully, with one wrong data",
			statement: &stmt.Schema{Type: stmt.DatabaseNameSchemaType},
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
					Option:        opt,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
		},
		{
			name:      "get database successfully, with one wrong data",
			statement: &stmt.Schema{Type: stmt.DatabaseNameSchemaType},
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
					Option:        opt,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
		},
		{
			name:      "schema query, unknown metadata type",
			statement: &stmt.Schema{},
		},
		{
			name:      "get all database schemas",
			statement: &stmt.Schema{Type: stmt.DatabaseSchemaType},
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
					Option:        opt,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
		},
		{
			name:      "get database list err",
			statement: &stmt.Schema{Type: stmt.DatabaseNameSchemaType},
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := SchemaCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("SchemaCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
