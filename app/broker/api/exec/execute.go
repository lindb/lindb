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
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	httppkg "github.com/lindb/common/pkg/http"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/execution"
	"github.com/lindb/lindb/sql/tree"
)

// ExecutePath represents lin language executor's path.
var ExecutePath = "/exec"

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
		// FIXME: move to common pkg
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%v", err)
				_ = c.Error(errors.New(msg))
				c.Header("Content-Type", "text/plain")
				c.String(http.StatusInternalServerError, msg)
			}
		}()
		return e.execute(c)
	}); err != nil {
		_ = c.Error(err)
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// execute lin query language.
func (e *ExecuteAPI) execute(c *gin.Context) error {
	ctx, cancel := e.deps.WithTimeout()
	defer cancel()

	reqSession := &models.Session{}
	if err := c.BindHeader(reqSession); err != nil {
		return err
	}

	param := models.ExecuteParam{}
	err := c.ShouldBind(&param)
	if err != nil {
		return err
	}
	if reqSession.Databases == "" {
		reqSession.Databases = param.Database
	}

	requestID := e.deps.RequestIDGen.GenerateRequestID()

	defer func() {
		e.deps.RequestMgr.CompleteRequet(requestID, nil)
	}()

	fmt.Println(param)

	// FIXME: session?
	c.Set(constants.CurrentSQL, &param)
	idAllocator := tree.NewNodeIDAllocator()
	stmt, err := tree.GetParser().CreateStatement(param.SQL, idAllocator)
	if err != nil {
		return err
	}

	if stmt == nil {
		return errors.New("can't parse lin query language")
	}
	// set query session context
	preparedStmt := &tree.PreparedStatement{
		Statement:  stmt,
		PrepareSQL: param.SQL,
	}
	session := &execution.Session{
		// the current system time when creating this session,
		// make sure the same current time is used in the session.
		Context:         context.WithValue(ctx, constants.ContextKeyCurrentTime, timeutil.Now()),
		RequestID:       requestID,
		NodeIDAllocator: idAllocator,
		Database:        reqSession.Databases,
		Statement:       preparedStmt,
	}

	statementType := execution.GetStatementType(stmt)
	factory := execution.GetExecutionFactory(statementType)
	exec := factory.CreateExecution(session, preparedStmt)
	result := exec.Start()
	if result == nil || reflect.ValueOf(result).IsNil() {
		httppkg.NotFound(c)
	} else {
		httppkg.OK(c, result)
	}

	// TODO: resource group
	return nil
}
