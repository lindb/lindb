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

package query

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	parseSQLFunc = parseSQL

	MetadataQueryPath = "/query/metadata"
)

var errWrongQueryStmt = errors.New("can't parse metadata query ql")
var errUnknownMetadataStmt = errors.New("unknown metadata statement")

// MetadataAPI represents metadata query api
type MetadataAPI struct {
	deps *deps.HTTPDeps
}

// NewMetadataAPI creates database api instance
func NewMetadataAPI(deps *deps.HTTPDeps) *MetadataAPI {
	return &MetadataAPI{
		deps: deps,
	}
}

// Register adds metadata suggest url route.
func (d *MetadataAPI) Register(route gin.IRoutes) {
	route.GET(MetadataQueryPath, d.Suggest)
}

// Suggest handles metadata suggest query by LinQL
func (d *MetadataAPI) Suggest(c *gin.Context) {
	var param struct {
		Database string `form:"db"`
		SQL      string `form:"sql" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	metaQuery, err := parseSQLFunc(param.SQL)
	if err != nil {
		http.Error(c, err)
		return
	}
	switch metaQuery.Type {
	case stmt.Database:
		d.showDatabases(c)
	case stmt.Namespace, stmt.Metric, stmt.Field, stmt.TagKey, stmt.TagValue:
		if param.Database == "" {
			http.Error(c, errors.New("database name required"))
			return
		}
		d.suggest(c, param.Database, metaQuery)
	default:
		http.Error(c, errUnknownMetadataStmt)
	}
}

// showDatabases shows all database names
func (d *MetadataAPI) showDatabases(c *gin.Context) {
	databases, err := d.deps.DatabaseSrv.List()
	if err != nil {
		http.Error(c, err)
		return
	}
	var databaseNames []string
	for _, db := range databases {
		databaseNames = append(databaseNames, db.Name)
	}
	http.OK(c, &models.Metadata{
		Type:   stmt.Database.String(),
		Values: databaseNames,
	})
}

// suggest executes the suggest query
func (d *MetadataAPI) suggest(c *gin.Context, database string, request *stmt.Metadata) {
	//TODO add timeout cfg
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	exec := d.deps.ExecutorFct.NewMetadataBrokerExecutor(ctx, database, request,
		d.deps.StateMachines.ReplicaStatusSM, d.deps.StateMachines.NodeSM, d.deps.JobManager)
	values, err := exec.Execute()
	if err != nil {
		http.Error(c, err)
		return
	}
	switch request.Type {
	case stmt.Field:
		// build field result model
		result := make(map[field.Name]field.Meta)
		fields := field.Metas{}
		for _, value := range values {
			err = encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				http.Error(c, err)
				return
			}
			for _, f := range fields {
				result[f.Name] = f
			}
		}
		var resultFields []models.Field
		for _, f := range result {
			resultFields = append(resultFields, models.Field{
				Name: string(f.Name),
				Type: f.Type.String(),
			})
		}
		http.OK(c, &models.Metadata{
			Type:   request.Type.String(),
			Values: resultFields,
		})
	default:
		http.OK(c, &models.Metadata{
			Type:   request.Type.String(),
			Values: values,
		})
	}
}

// parseSQL parses metadata query sql
func parseSQL(ql string) (*stmt.Metadata, error) {
	query, err := sql.Parse(ql)
	if err != nil {
		return nil, err
	}
	metaQuery, ok := query.(*stmt.Metadata)
	if !ok {
		return nil, errWrongQueryStmt
	}
	return metaQuery, nil
}
