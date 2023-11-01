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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestSchemaCommand_NotFound(t *testing.T) {
	rs, err := SchemaCommand(context.TODO(), nil, nil, &stmt.Schema{Type: stmt.SchemaType(-1)})
	assert.Nil(t, rs)
	assert.Nil(t, err)
}

func TestSchemaCommand_listDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo: repo,
	}

	t.Run("list database failure", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
		rs, err := SchemaCommand(context.TODO(), deps, nil, &stmt.Schema{Type: stmt.DatabaseSchemaType})
		assert.Error(t, err)
		assert.Nil(t, rs)
	})

	t.Run("list database successfully", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).
			Return([]state.KeyValue{{Value: []byte("1")}, {Value: []byte("{}")}}, nil)
		rs, err := SchemaCommand(context.TODO(), deps, nil, &stmt.Schema{Type: stmt.DatabaseSchemaType})
		assert.NoError(t, err)
		assert.NotNil(t, rs)
	})
}

func TestSchemaCommand_listDatabaseNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo: repo,
	}

	t.Run("list database failure", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
		rs, err := SchemaCommand(context.TODO(), deps, nil, &stmt.Schema{Type: stmt.DatabaseNameSchemaType})
		assert.Error(t, err)
		assert.Nil(t, rs)
	})

	t.Run("list database successfully", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).
			Return([]state.KeyValue{{Value: []byte("1")}, {Value: []byte(`{"name":"test"}`)}}, nil)
		rs, err := SchemaCommand(context.TODO(), deps, nil, &stmt.Schema{Type: stmt.DatabaseNameSchemaType})
		assert.NoError(t, err)
		assert.NotNil(t, rs)
	})
}

func TestSchemaCommand_saveDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo: repo,
	}
	cfg := `{"name":"test","routers":[]}`
	cases := []struct {
		name    string
		body    string
		prepare func()
		wantErr bool
	}{
		{
			name:    "unmarshal err",
			body:    "a",
			wantErr: true,
		},
		{
			name:    "validator failure",
			body:    "{}",
			wantErr: true,
		},
		{
			name: "save config failure",
			body: cfg,
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "save config successfually",
			body: cfg,
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := SchemaCommand(context.TODO(), deps, nil, &stmt.Schema{Type: stmt.CreateDatabaseSchemaType, Value: tt.body})
			if ((err != nil) != tt.wantErr && rs == nil) || (!tt.wantErr && rs == nil) {
				t.Errorf("SchemaCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
