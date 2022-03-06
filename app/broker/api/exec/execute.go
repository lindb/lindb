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

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
	sqlpkg "github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var (
	sqlParseFn = sqlpkg.Parse
)

var (
	// ExecutePath represents lin language executor's path.
	ExecutePath = "/exec"
)

type ExecuteAPI struct {
	deps *deps.HTTPDeps

	logger *logger.Logger
}

func NewExecuteAPI(deps *deps.HTTPDeps) *ExecuteAPI {
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

// Execute executes lin language.
func (e *ExecuteAPI) Execute(c *gin.Context) {
	ctx, cancel := e.deps.WithTimeout()
	defer cancel()

	var param struct {
		Database string `form:"db"`
		SQL      string `form:"sql" binding:"required"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	stmt, err := sqlParseFn(param.SQL)
	if err != nil {
		httppkg.Error(c, err)
		return
	}

	var result interface{}
	switch s := stmt.(type) {
	case *stmtpkg.State:
		// execute state query
		result = e.execStateQuery(s)
	case *stmtpkg.Metadata:
		result, err = e.execMetadataQuery(ctx, s)
	default:
		httppkg.Error(c, errors.New("can't parse lin language"))
		return
	}
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	if result == nil || reflect.ValueOf(result).IsNil() {
		httppkg.NotFound(c)
	} else {
		httppkg.OK(c, result)
	}
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
func (e *ExecuteAPI) execMetadataQuery(ctx context.Context, metadataStmt *stmtpkg.Metadata) (interface{}, error) {
	switch metadataStmt.Type {
	case stmtpkg.Database:
		return e.listDataBase(ctx)
	default:
		return nil, nil
	}
}

// listDataBase returns database list in cluster.
func (e *ExecuteAPI) listDataBase(ctx context.Context) ([]*models.Database, error) {
	var result []*models.Database
	data, err := e.deps.Repo.List(ctx, constants.DatabaseConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		db := &models.Database{}
		err = encoding.JSONUnmarshal(val.Value, db)
		if err != nil {
			e.logger.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			db.Desc = db.String()
			result = append(result, db)
		}
	}
	return result, nil
}
