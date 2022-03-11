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

package exec

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
	sqlpkg "github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	sqlParseFn = sqlpkg.Parse
)

var (
	// ExecutePath represents lin language executor's path.
	ExecutePath             = "/exec"
	errDatabaseNameRequired = errors.New("database name cannot be empty")
)

// ExecuteAPI represent lin query language execution api.
type ExecuteAPI struct {
	deps *deps.HTTPDeps

	logger *logger.Logger
}

// NewExecuteAPI creates a lin query language execution api.
func NewExecuteAPI(deps *deps.HTTPDeps) *ExecuteAPI {
	//TODO add metric
	return &ExecuteAPI{
		deps:   deps,
		logger: logger.GetLogger("broker", "ExecuteAPI"),
	}
}

// Register adds lin language executor's path.
func (e *ExecuteAPI) Register(route gin.IRoutes) {
	// register multi http methods
	route.GET(ExecutePath, e.Execute)
	route.POST(ExecutePath, e.Execute)
	route.PUT(ExecutePath, e.Execute)
}

// Execute executes lin query language with limit.
func (e *ExecuteAPI) Execute(c *gin.Context) {
	if err := e.deps.QueryLimiter.Do(func() error {
		return e.execute(c)
	}); err != nil {
		httppkg.Error(c, err)
	}
}

// execute lin query language.
func (e *ExecuteAPI) execute(c *gin.Context) error {
	ctx, cancel := e.deps.WithTimeout()
	defer cancel()

	param := models.ExecuteParam{}
	err := c.ShouldBind(&param)
	if err != nil {
		return err
	}
	stmt, err := sqlParseFn(param.SQL)
	if err != nil {
		return err
	}

	var result interface{}
	switch s := stmt.(type) {
	case *stmtpkg.State:
		// execute state query
		result = e.execStateQuery(s)
	case *stmtpkg.Metadata:
		result, err = e.execMetadataQuery(ctx, param, s)
	case *stmtpkg.Query:
		if strings.TrimSpace(param.Database) == "" {
			return errDatabaseNameRequired
		}
		metricQuery := e.deps.QueryFactory.NewMetricQuery(ctx, param.Database, s)
		result, err = metricQuery.WaitResponse()
		if err != nil {
			return err
		}
	default:
		return errors.New("can't parse lin query language")
	}
	if err != nil {
		return err
	}
	if result == nil || reflect.ValueOf(result).IsNil() {
		httppkg.NotFound(c)
	} else {
		httppkg.OK(c, result)
	}
	return nil
}

// execStateQuery executes the state query.
func (e *ExecuteAPI) execStateQuery(stateStmt *stmtpkg.State) interface{} {
	switch stateStmt.Type {
	case stmtpkg.Master:
		return e.deps.Master.GetMaster()
	default:
		return nil
	}
}

// execMetadataQuery executes the metadata query.
func (e *ExecuteAPI) execMetadataQuery(ctx context.Context, param models.ExecuteParam, metadataStmt *stmtpkg.Metadata) (interface{}, error) {
	switch metadataStmt.Type {
	case stmtpkg.Database:
		return e.listDataBase(ctx)
	case stmtpkg.Namespace, stmtpkg.Metric, stmtpkg.Field, stmtpkg.TagKey, stmtpkg.TagValue:
		if strings.TrimSpace(param.Database) == "" {
			return nil, errDatabaseNameRequired
		}
		return e.suggest(ctx, param.Database, metadataStmt)
	default:
		return nil, nil
	}
}

// listDataBase returns database list in cluster.
func (e *ExecuteAPI) listDataBase(ctx context.Context) (*models.Metadata, error) {
	data, err := e.deps.Repo.List(ctx, constants.DatabaseConfigPath)
	if err != nil {
		return nil, err
	}
	var databaseNames []interface{}
	for _, val := range data {
		db := &models.Database{}
		err = encoding.JSONUnmarshal(val.Value, db)
		if err != nil {
			e.logger.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
			continue
		}
		databaseNames = append(databaseNames, db.Name)
	}
	return &models.Metadata{
		Type:   stmtpkg.Database.String(),
		Values: databaseNames,
	}, nil
}

// suggest executes metadata suggest query.
func (e *ExecuteAPI) suggest(ctx context.Context, database string, request *stmtpkg.Metadata) (interface{}, error) {
	metaDataQuery := e.deps.QueryFactory.NewMetadataQuery(ctx, database, request)
	values, err := metaDataQuery.WaitResponse()
	if err != nil {
		return nil, err
	}
	switch request.Type {
	case stmtpkg.Field:
		// build field result model
		result := make(map[field.Name]field.Meta)
		fields := field.Metas{}
		for _, value := range values {
			err = encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				return nil, err
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
		return &models.Metadata{
			Type:   request.Type.String(),
			Values: resultFields,
		}, nil
	default:
		return &models.Metadata{
			Type:   request.Type.String(),
			Values: values,
		}, nil
	}
}
