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
	"sort"
	"strconv"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/gin-gonic/gin"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	ljson "github.com/lindb/arrow/pkg/json"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
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

	e.logger.Info("execute lin query language", logger.String("db", session.Database),
		logger.String("sql", param.SQL), logger.Any("timeRange", param.TimeRange))

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

	// Log queries include hidden cursor columns; apply broker-side pagination when present.
	if ok, err := applyBrokerPagination(c, record, &param); ok || err != nil {
		return err
	}

	if c.GetHeader("Accept") == constants.ContentTypeArrow {
		return writeArrowStream(c, record)
	}
	return writeJSONResponse(c, record)
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

// writeJSONResponse converts the RecordBatch directly to a ResultSet and writes
// it as JSON. unwrapRecord is called first to normalize any lindb-wrapped arrays.
func writeJSONResponse(c *gin.Context, record arrow.RecordBatch) error {
	unwrapped, rebuilt := unwrapRecord(record)
	if rebuilt {
		defer unwrapped.Release()
	}

	rs, err := ljson.DecodeRecordToResultSet(unwrapped)
	if err != nil {
		return fmt.Errorf("json response: decode record to result set: %w", err)
	}

	c.JSON(http.StatusOK, rs)
	return nil
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

// applyBrokerPagination handles global sorting, truncation, and cursor extraction for log
// queries.  Storage nodes append a hidden "_id" column to every raw log RecordBatch with
// format "shardID:ts:logID" where ts is the per-row log timestamp in nanoseconds.
//
// When that column is absent the record is not a paginated log result; the function
// returns (false, nil) and the caller falls through to the normal response path.
//
// When present, the function:
//  1. Parses all _id strings into (shardID, ts, logID) components.
//  2. Sorts all rows by ts DESC.
//  3. Truncates to effectiveLimit.
//  4. Records the oldest (last in sorted order) row per shard as that shard's next cursor.
//  5. Carries forward old per-shard cursors for shards absent from this page.
//  6. Serialises the cursor map as "shardID:ts:logID,…" (sorted for determinism).
//  7. Strips the hidden _id column and writes a models.LogPageResult JSON response.
//
// Returns (true, err) after writing the response, or (false, nil) when the record is not
// a paginated log result.
func applyBrokerPagination(c *gin.Context, record arrow.RecordBatch, param *models.ExecuteParam) (bool, error) {
	schema := record.Schema()
	idIdxs := schema.FieldIndices("_id")
	if len(idIdxs) == 0 {
		return false, nil
	}
	idIdx := idIdxs[0]
	idCol, ok := record.Column(idIdx).(*array.String)
	if !ok {
		return false, fmt.Errorf("_id column is not string type")
	}

	numRows := int(record.NumRows())

	// Pre-parse all _id strings to avoid repeated string splitting during sort.
	type rowID struct {
		shardID int64
		ts      int64
		logID   uint32
	}
	ids := make([]rowID, numRows)
	for i := range numRows {
		ids[i].shardID, ids[i].ts, ids[i].logID = parseID(idCol.Value(i))
	}

	// Build a sorted index array: ts DESC.
	indices := make([]int, numRows)
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return ids[indices[i]].ts > ids[indices[j]].ts
	})

	effectiveLimit := int64(1000)
	if param.Limit > 0 {
		effectiveLimit = param.Limit
	}

	hasMore := int64(numRows) > effectiveLimit
	if hasMore {
		indices = indices[:effectiveLimit]
	}

	// Build per-shard cursor: last (oldest) row per shard in the DESC-sorted result.
	// Iterating in DESC order and overwriting means the final map value is the oldest row.
	lastShardRow := make(map[int64]int)
	for _, ri := range indices {
		lastShardRow[ids[ri].shardID] = ri
	}

	newCursors := make(map[int64]spi.ShardCursor)
	for shardID, ri := range lastShardRow {
		newCursors[shardID] = spi.ShardCursor{
			Timestamp: ids[ri].ts,
			LogID:     ids[ri].logID,
		}
	}

	// Carry forward old cursors for shards not present in this page's result.
	if param.Cursor != "" {
		for shardID, cursor := range parseShardCursors(param.Cursor) {
			if _, exists := newCursors[shardID]; !exists {
				newCursors[shardID] = cursor
			}
		}
	}

	// Serialise next cursor deterministically (sorted shard IDs).
	nextCursor := ""
	if hasMore && len(newCursors) > 0 {
		parts := make([]string, 0, len(newCursors))
		for shardID, cursor := range newCursors {
			parts = append(parts, fmt.Sprintf("%d:%d:%d", shardID, cursor.Timestamp, cursor.LogID))
		}
		sort.Strings(parts)
		nextCursor = strings.Join(parts, ",")
	}

	// Determine visible (non-hidden) column indices.
	hiddenSet := map[int]bool{idIdx: true}
	numCols := schema.NumFields()
	visibleIdxs := make([]int, 0, numCols-1)
	columns := make([]string, 0, numCols-1)
	for i := range numCols {
		if !hiddenSet[i] {
			visibleIdxs = append(visibleIdxs, i)
			columns = append(columns, schema.Field(i).Name)
		}
	}

	// Extract values for the selected (sorted, truncated) rows.
	values := make([][]interface{}, len(indices))
	for rowPos, ri := range indices {
		row := make([]interface{}, len(visibleIdxs))
		for colPos, ci := range visibleIdxs {
			row[colPos] = extractArrowValue(record.Column(ci), ri)
		}
		values[rowPos] = row
	}

	c.JSON(http.StatusOK, &models.LogPageResult{
		Columns:    columns,
		Values:     values,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
	return true, nil
}

// extractArrowValue extracts a single cell from an Arrow array as a plain Go value.
// Covers the types that appear in log query output columns; falls back to ValueStr for others.
func extractArrowValue(col arrow.Array, idx int) interface{} {
	if col.IsNull(idx) {
		return nil
	}
	switch a := col.(type) {
	case *array.Int64:
		return a.Value(idx)
	case *array.String:
		return a.Value(idx)
	case *array.LargeString:
		return a.Value(idx)
	case *array.Int32:
		return int64(a.Value(idx))
	case *array.Boolean:
		return a.Value(idx)
	default:
		return col.ValueStr(idx)
	}
}

// parseShardCursors parses a per-shard cursor string "shardID:ts:logID,…" into a map.
// Malformed entries are silently skipped.
func parseShardCursors(s string) map[int64]spi.ShardCursor {
	result := make(map[int64]spi.ShardCursor)
	for _, entry := range strings.Split(s, ",") {
		if entry == "" {
			continue
		}
		shardID, ts, logID := parseID(entry)
		result[shardID] = spi.ShardCursor{Timestamp: ts, LogID: logID}
	}
	return result
}

// parseID parses a "_id" string "shardID:ts:logID" into its numeric components.
// Malformed input returns zero values.
func parseID(s string) (shardID, ts int64, logID uint32) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return
	}
	shardID, _ = strconv.ParseInt(parts[0], 10, 64)
	ts, _ = strconv.ParseInt(parts[1], 10, 64)
	v, _ := strconv.ParseInt(parts[2], 10, 64)
	logID = uint32(v)
	return
}
