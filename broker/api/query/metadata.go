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
	ctx, cancel := context.WithTimeout(context.Background(), d.deps.BrokerCfg.Query.Timeout.Duration())
	defer cancel()

	metaDataQuery := d.deps.QueryFactory.NewMetadataQuery(ctx, database, request)
	values, err := metaDataQuery.WaitResponse()
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
		// HistogramSum(sum), HistogramCount(sum), HistogramMin(min), HistogramMax(max) is visible
		// __bucket_{id}(HistogramField) is not visible for api,
		// underlying histogram data is only restricted access by user via quantile function
		// furthermore, we suggest some quantile functions for user in field names, such as quantile(0.99)
		var (
			resultFields []models.Field
			hasHistogram bool
		)
		for _, f := range result {
			if f.Type != field.HistogramField {
				resultFields = append(resultFields, models.Field{
					Name: string(f.Name),
					Type: f.Type.String(),
				})
			} else {
				hasHistogram = true
			}
		}
		//
		if hasHistogram {
			resultFields = append(resultFields,
				models.Field{Name: "quantile(0.99)", Type: field.HistogramField.String()},
				models.Field{Name: "quantile(0.95)", Type: field.HistogramField.String()},
				models.Field{Name: "quantile(0.90)", Type: field.HistogramField.String()},
			)
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
