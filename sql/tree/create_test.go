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
				Name: "db",
				Props: map[string]any{
					"p1": "v1",
					"p2": "v2",
					"p3": "v3",
				},
			},
		},
		{
			sql: `create database db 
				with(
				 	p1='v1',p2='v2',p3='v3'
				)
				rollup(
					(r1='v1',r2='v2'),
					(r11='v1',r22='v2')
				)
			`,
			stmt: &CreateDatabase{
				Name: "db",
				Props: map[string]any{
					"p1": "v1",
					"p2": "v2",
					"p3": "v3",
				},
				Rollup: []RollupOption{
					{
						Options: map[string]any{
							"r1": "v1",
							"r2": "v2",
						},
					},
					{
						Options: map[string]any{
							"r11": "v1",
							"r22": "v2",
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.sql, func(t *testing.T) {
			stmt, err := GetParser().CreateStatement(tt.sql, NewNodeIDAllocator())
			assert.NoError(t, err)
			createDB := stmt.(*CreateDatabase)
			assert.Equal(t, tt.stmt.Name, createDB.Name)
			assert.Equal(t, tt.stmt.CreateOptions, createDB.CreateOptions)
			assert.Equal(t, tt.stmt.Rollup, createDB.Rollup)
		})
	}
}
