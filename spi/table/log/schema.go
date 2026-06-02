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

package log

import (
	"encoding/hex"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/lindb/arrow/pkg/arrow/builder"
	arrowConst "github.com/lindb/arrow/pkg/constants"
	logspkg "github.com/lindb/arrow/pkg/logs"

	"github.com/lindb/lindb/spi/types"
)

// columnReader builds an appender that reads one row from a log Reader into the RecordBuilder.
// It is called once per query to capture the builder reference in a closure.
type columnReader func(rb *builder.RecordBuilder) func(r *logspkg.Reader, row int)

// logColumn pairs an Arrow field definition with the function that reads its value from storage.
// It is the single unit of schema knowledge for the log table.
type logColumn struct {
	Field  arrow.Field
	reader columnReader
}

// logSchema is the canonical schema for the log table.
// Adding a new column here automatically propagates to both broker metadata
// registration and the source connector's column-reading pipeline.
var logSchema = []logColumn{
	{
		Field: arrow.Field{Name: arrowConst.Timestamp, Type: arrow.FixedWidthTypes.Timestamp_ns, Nullable: false},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			tb := rb.TimestampBuilder(arrowConst.Timestamp)
			return func(r *logspkg.Reader, row int) {
				tb.Append(arrow.Timestamp(r.Timestamp(row)))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Message, Type: arrow.BinaryTypes.String, Nullable: false},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			sb := rb.StringBuilder(arrowConst.Message)
			return func(r *logspkg.Reader, row int) {
				sb.Append(r.Message(row))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Level, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			lb := rb.StringBuilder(arrowConst.Level)
			return func(r *logspkg.Reader, row int) {
				lb.Append(r.Level(row))
			}
		},
	},
	{
		// Attributes is stored internally as List<Uint32> FK references; exposed to users as map<string,string>.
		Field: arrow.Field{Name: arrowConst.Attributes, Type: arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String), Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			attrMap := rb.MapBuilder(arrowConst.Attributes)
			attrKey := attrMap.KeyBuilder().(*array.StringBuilder)
			attrVal := attrMap.ItemBuilder().(*array.StringBuilder)
			return func(r *logspkg.Reader, row int) {
				r.AttributesToMap(row, attrMap, attrKey, attrVal)
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.EventName, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			eb := rb.StringBuilder(arrowConst.EventName)
			return func(r *logspkg.Reader, row int) {
				eb.Append(r.EventName(row))
			}
		},
	},
	{
		// trace_id stored as hex string so SQL comparisons with string literals work.
		Field: arrow.Field{Name: arrowConst.TraceID, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			b := rb.StringBuilder(arrowConst.TraceID)
			return func(r *logspkg.Reader, row int) {
				v := r.TraceID(row)
				if v == nil {
					b.AppendNull()
				} else {
					b.Append(hex.EncodeToString(v))
				}
			}
		},
	},
	{
		// span_id stored as hex string so SQL comparisons with string literals work.
		Field: arrow.Field{Name: arrowConst.SpanID, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*logspkg.Reader, int) {
			b := rb.StringBuilder(arrowConst.SpanID)
			return func(r *logspkg.Reader, row int) {
				v := r.SpanID(row)
				if v == nil {
					b.AppendNull()
				} else {
					b.Append(hex.EncodeToString(v))
				}
			}
		},
	},
}

// logSchemaIndex maps column name → logColumn for O(1) lookup at query time.
var logSchemaIndex map[string]*logColumn

func init() {
	logSchemaIndex = make(map[string]*logColumn, len(logSchema))
	for i := range logSchema {
		logSchemaIndex[logSchema[i].Field.Name] = &logSchema[i]
	}
	// Register the canonical Arrow fields so broker_meta can read them without
	// importing this package directly (which would create an import cycle via spi).
	types.RegisterTableArrowFields("logs", LogTableArrowFields())
}

// LogTableArrowFields returns the Arrow field list for the log table.
// Used by broker_meta to register the table metadata returned to the SQL planner.
func LogTableArrowFields() []arrow.Field {
	fields := make([]arrow.Field, len(logSchema))
	for i, col := range logSchema {
		fields[i] = col.Field
	}
	return fields
}

// IsFixedLogSchemaField returns true when the column name belongs to the fixed log schema.
// Used by initializeSearchContext to distinguish fixed schema columns from
// user-defined dynamic attribute fields. This guard is necessary because the
// Dynamic extension type and String share the same underlying Arrow type.
func IsFixedLogSchemaField(name string) bool {
	_, ok := logSchemaIndex[name]
	return ok
}

// BuildLogColumnAppender returns an appender for the named fixed-schema column, or
// (nil, false) when the name is not in the fixed schema (i.e. it is a dynamic attribute).
// Used by buildColumnAppenders in the source connector.
func BuildLogColumnAppender(name string, rb *builder.RecordBuilder) (func(*logspkg.Reader, int), bool) {
	col, ok := logSchemaIndex[name]
	if !ok {
		return nil, false
	}
	return col.reader(rb), true
}
