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

	"github.com/lindb/lindb/app/broker/api/admin"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/app/broker/middleware"
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

	errDatabaseNameRequired = errors.New("database name required")
)

var errWrongQueryStmt = errors.New("can't parse metadata query ql")
var errUnknownMetadataStmt = errors.New("unknown metadata statement")

// MetadataAPI represents metadata query api
type MetadataAPI struct {
	deps         *deps.HTTPDeps
	ListDataBase func() ([]*models.Database, error)
}

// NewMetadataAPI creates database api instance
func NewMetadataAPI(deps *deps.HTTPDeps) *MetadataAPI {
	return &MetadataAPI{
		deps:         deps,
		ListDataBase: admin.NewDatabaseAPI(deps).ListDataBase,
	}
}

// Register adds metadata suggest url route.
func (d *MetadataAPI) Register(route gin.IRoutes) {
	route.GET(
		MetadataQueryPath,
		middleware.WithHistogram(middleware.HttHandlerTimerVec.WithTagValues(MetadataQueryPath)),
		d.Suggest,
	)
}

// Suggest handles metadata suggest query by LinQL
func (d *MetadataAPI) Suggest(c *gin.Context) {
	if err := d.deps.QueryLimiter.Do(func() error {
		return d.suggestWithLimit(c)
	}); err != nil {
		http.Error(c, err)
	}
}

// suggestWithLimit handles metadata suggest query by LinQL
func (d *MetadataAPI) suggestWithLimit(c *gin.Context) error {
	var param struct {
		Database string `form:"db"`
		SQL      string `form:"sql" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		return err
	}
	metaQuery, err := parseSQLFunc(param.SQL)
	if err != nil {
		return err
	}
	switch metaQuery.Type {
	case stmt.Database:
		if err := d.showDatabases(c); err != nil {
			return err
		}
	case stmt.Namespace, stmt.Metric, stmt.Field, stmt.TagKey, stmt.TagValue:
		if param.Database == "" {
			return errDatabaseNameRequired
		}
		if err := d.suggest(c, param.Database, metaQuery); err != nil {
			return err
		}
	default:
		return errUnknownMetadataStmt
	}
	return nil
}

// showDatabases shows all database names
func (d *MetadataAPI) showDatabases(c *gin.Context) error {
	databases, err := d.ListDataBase()
	if err != nil {
		return err
	}
	var databaseNames []string
	for _, db := range databases {
		databaseNames = append(databaseNames, db.Name)
	}
	http.OK(c, &models.Metadata{
		Type:   stmt.Database.String(),
		Values: databaseNames,
	})
	return nil
}

// suggest executes the suggest query
func (d *MetadataAPI) suggest(c *gin.Context, database string, request *stmt.Metadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.deps.BrokerCfg.Query.Timeout.Duration())
	defer cancel()

	metaDataQuery := d.deps.QueryFactory.NewMetadataQuery(ctx, database, request)
	values, err := metaDataQuery.WaitResponse()
	if err != nil {
		return err
	}
	switch request.Type {
	case stmt.Field:
		// build field result model
		result := make(map[field.Name]field.Meta)
		fields := field.Metas{}
		for _, value := range values {
			err = encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				return err
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
	return nil
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
