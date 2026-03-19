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

package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate_Database(t *testing.T) {
	cases := []struct {
		stmt *CreateDatabase
		sql  string
	}{
		{
			sql: "create database db",
			stmt: &CreateDatabase{
				Name: "db",
			},
		},
		{
			sql: `create database db with(p1='v1',p2='v2',p3='v3')`,
			stmt: &CreateDatabase{
				Name:  "db",
				Props: []*Property{
					// "p1": "v1",
					// "p2": "v2",
					// "p3": "v3",
				},
			},
		},
		{
			sql: `create database db
				with(
				 	p1='v1',p2='v2',p3='v3'
				)
				retention(
					(r1='v1',r2='v2'),
					(r11='v1',r22='v2')
				)
			`,
			stmt: &CreateDatabase{
				Name:  "db",
				Props: []*Property{
					// "p1": "v1",
					// "p2": "v2",
					// "p3": "v3",
				},
				Retention: []*RetentionOption{
					{
						Props: []*Property{
							// "r1": "v1",
							// "r2": "v2",
						},
					},
					{
						Props: []*Property{
							// "r11": "v1",
							// "r22": "v2",
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := GetParser().CreateStatement(tt.sql, NewNodeIDAllocator())
			assert.NoError(t, err)
		})
	}
}
