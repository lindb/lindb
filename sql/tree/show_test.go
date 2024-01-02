package tree

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestShow(t *testing.T) {
// 	cases := []struct {
// 		sql  string
// 		stmt Statement
// 	}{
// 		{
// 			"show master",
// 			&ShowMaster{},
// 		},
// 		{
// 			"show databases",
// 			&ShowDatabases{},
// 		},
// 		{
// 			"show brokers",
// 			&ShowBrokers{},
// 		},
// 		{
// 			"show requests",
// 			&ShowRequests{},
// 		},
// 		{
// 			"show limit",
// 			&ShowLimit{},
// 		},
// 		{
// 			"show metadata types",
// 			&ShowMetadataTypes{},
// 		},
// 		{
// 			"show metadatas",
// 			&ShowMetadatas{},
// 		},
// 		{
// 			"show alive",
// 			&ShowAlive{},
// 		},
// 		{
// 			"show replications",
// 			&ShowReplications{},
// 		},
// 		{
// 			"show state",
// 			&ShowState{},
// 		},
// 		{
// 			"show namespaces",
// 			&ShowNamespaces{},
// 		},
// 		{
// 			"show metrics",
// 			&ShowMetrics{},
// 		},
// 		{
// 			"show fields",
// 			&ShowFields{},
// 		},
// 		{
// 			"show tag keys",
// 			&ShowTagKeys{},
// 		},
// 		{
// 			"show tag values",
// 			&ShowTagValues{},
// 		},
// 	}
// 	for _, tt := range cases {
// 		tt := tt
// 		t.Run(tt.sql, func(t *testing.T) {
// 			stmt, err := Parse(tt.sql)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.stmt, stmt)
// 		})
// 	}
// }
