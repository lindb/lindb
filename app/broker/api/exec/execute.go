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

	httppkg "github.com/lindb/common/pkg/http"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/app/broker/api/exec/command"
	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	sqlpkg "github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	sqlParseFn = sqlpkg.Parse
)

// statementExecFn represents statement execution funcation define.
type statementExecFn func(ctx context.Context,
	deps *depspkg.HTTPDeps,
	param *models.ExecuteParam,
	stmt stmtpkg.Statement) (interface{}, error)

var (
	// ExecutePath represents lin language executor's path.
	ExecutePath = "/exec"

	// register all commands for the statement of lin query language.
	commands = map[stmtpkg.StatementType]statementExecFn{
		stmtpkg.MetadataStatement:       command.MetadataCommand,
		stmtpkg.SchemaStatement:         command.SchemaCommand,
		stmtpkg.StorageStatement:        command.StorageCommand,
		stmtpkg.StateStatement:          command.StateCommand,
		stmtpkg.MetricMetadataStatement: command.MetricMetadataCommand,
		stmtpkg.QueryStatement:          command.QueryCommand,
		stmtpkg.RequestStatement:        command.RequestCommand,
		stmtpkg.LimitStatement:          command.LimitCommand,
	}
)

type ExecuteAPI struct {
	deps *depspkg.HTTPDeps

	logger logger.Logger
}

// NewExecuteAPI creates a lin query language execution api.
func NewExecuteAPI(deps *depspkg.HTTPDeps) *ExecuteAPI {
	// TODO add metric
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

// Execute executes lin query language with rate limit.
// 1. metric data/metadata query statement;
// 2. cluster metadata/state query statement;
// 3. database/storage management statement;
//
// @Summary execute lin query language
// @Description Execute lin query language with rate limit, then return different response based on execution statement.
// @Description 1. metric data/metadata query statement;
// @Description 2. cluster metadata/state query statement;
// @Description 3. database/storage management statement;
// @Tags LinQL
// @Accept json
// @Param param body models.ExecuteParam ture "param data"
// @Produce json
// @Success 200 {object} models.ResultSet
// @Success 200 {object} models.Metadata
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "can't parse lin query language"
// @Failure 500 {string} string "internal error"
// @Router /exec [get]
// @Router /exec [put]
// @Router /exec [post]
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
	c.Set(constants.CurrentSQL, &param)
	stmt, err := sqlParseFn(param.SQL)
	if err != nil {
		return err
	}

	if stmt == nil {
		return errors.New("can't parse lin query language")
	}

	if commandFn, ok := commands[stmt.StatementType()]; ok {
		result, err := commandFn(ctx, e.deps, &param, stmt)
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
	return errors.New("can't parse lin query language")
}
