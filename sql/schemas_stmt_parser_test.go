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

package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestShowSchemasStatement(t *testing.T) {
	q, err := Parse("show schemas")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Schema{Type: stmt.DatabaseSchemaType}, q)
}

func TestDropDatabaseStatement(t *testing.T) {
	q, err := Parse("drop database 'test'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Schema{Type: stmt.DropDatabaseSchemaType, Value: "test"}, q)
}

func TestCreateDatabase(t *testing.T) {
	cfg := `{\"name\":\"test\"}`
	sql := `create database ` + cfg
	q, err := Parse(sql)
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Schema{
		Type:  stmt.CreateDatabaseSchemaType,
		Value: `{"name":"test"}`,
	}, q)
}

func TestCreateDatabaseWith(t *testing.T) {
	sql := `create database cpu with (storage: "/lind-cluster", numofshard:2, replicafactor:11, 
       behead: 3h, ahead: 4h, autocreatens: true) 
       rollup ((interval: 5s, retention: 2M), (interval: 10m, retention: 2y))`
	_, err := Parse(sql)
	assert.NoError(t, err)
}
