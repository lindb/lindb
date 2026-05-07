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
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/gin-gonic/gin"
	larray "github.com/lindb/arrow/pkg/arrow/array"
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
		// TODO: handle panic?
		return e.execute(c)
	}); err != nil {
		// Retrieve sql/database from the parsed param if available; fall back to
		// empty strings when the error occurred before param binding succeeded.
		sql, database := "", ""
		if p, ok := c.Get(constants.CurrentSQLParams); ok {
			if param, ok := p.(*models.ExecuteParam); ok {
				sql = param.SQL
				database = param.Database
			}
		}
		e.logger.Error("execute lin query language error",
			logger.String("db", database),
			logger.String("sql", sql),
			logger.Error(err))
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
	if reqSession.Database == "" {
		reqSession.Database = param.Database
	}

	requestID := e.deps.RequestIDGen.GenerateRequestID()

	defer func() {
		e.deps.RequestMgr.CompleteRequet(requestID, nil)
	}()

	// FIXME: session?
	c.Set(constants.CurrentSQLParams, &param)
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
		Context: context.WithValue(context.WithValue(ctx, constants.ContextKeyCurrentTime, timeutil.Now()),
			constants.ContextKeyParams, &param),
		RequestID:       requestID,
		NodeIDAllocator: idAllocator,
		Database:        reqSession.Database,
		Statement:       preparedStmt,
	}

	statementType := execution.GetStatementType(stmt)
	factory := execution.GetExecutionFactory(statementType)
	record, err := factory.CreateExecution(session, preparedStmt).Start()
	if err != nil {
		// Execution error (e.g. pipeline panic, type mismatch): surface as HTTP 500.
		return err
	}
	if record == nil {
		return nil
	}
	defer record.Release()

	if c.GetHeader("Accept") == constants.ContentTypeArrow {
		return writeArrowStream(c, record)
	}
	// FIXME: impl json response
	// httppkg.OK(c, executionModel.NewResultSetFromRecord(record))
	panic("not support json response yet")
}

// writeArrowStream encodes the RecordBatch into an Arrow IPC stream and writes
// it to the HTTP response.  The encoding is done into an in-memory buffer first
// so that, if encoding fails, the HTTP status has not yet been committed and the
// caller can still return a proper error response to the client.
func writeArrowStream(c *gin.Context, record arrow.RecordBatch) error {
	unwrapped, rebuilt := unwrapRecord(record)
	if rebuilt {
		defer unwrapped.Release()
	}

	// Encode into a memory buffer first — the HTTP response must not be started
	// until we know encoding succeeded, otherwise the client would receive a
	// partial/empty Arrow stream instead of a proper error.
	var buf bytes.Buffer
	w := ipc.NewWriter(&buf)
	if err := w.Write(unwrapped); err != nil {
		return fmt.Errorf("arrow ipc: encode record batch: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("arrow ipc: close writer: %w", err)
	}

	// Encoding succeeded — commit the response.
	c.Header("Content-Type", constants.ContentTypeArrow)
	_, err := c.Writer.Write(buf.Bytes())
	return err
}

// unwrapRecord replaces columns that are lindb-wrapped arrays (e.g. *larray.Generic[T]
// produced by FilterableRecord) with their underlying standard Arrow arrays so that
// ipc.Writer can serialize them.
// Returns the (possibly new) RecordBatch and a boolean indicating whether a new
// RecordBatch was allocated. When rebuilt is true, the caller must call Release()
// on the returned record.
func unwrapRecord(record arrow.RecordBatch) (arrow.RecordBatch, bool) {
	schema := record.Schema()
	numCols := schema.NumFields()
	cols := make([]arrow.Array, numCols)
	needRebuild := false

	for i := range numCols {
		col := record.Column(i)
		if masker, ok := col.(larray.Masker); ok {
			// Retain the underlying array so the new RecordBatch holds
			// an independent reference that won't dangle if the original
			// record is released before the new one.
			underlying := masker.Storage()
			underlying.Retain()
			cols[i] = underlying
			needRebuild = true
		} else {
			cols[i] = col
		}
	}
	if !needRebuild {
		return record, false
	}
	return array.NewRecordBatch(schema, cols, record.NumRows()), true
}
