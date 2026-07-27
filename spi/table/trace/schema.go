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

package trace

import (
	"encoding/hex"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/lindb/arrow/pkg/arrow/builder"
	arrowConst "github.com/lindb/arrow/pkg/constants"
	tracespkg "github.com/lindb/arrow/pkg/traces"

	"github.com/lindb/lindb/spi/types"
)

// columnReader builds an appender that reads one span row from a TraceReader into the RecordBuilder.
// Called once per query to capture the builder reference(s) in a closure.
type columnReader func(rb *builder.RecordBuilder) func(r *tracespkg.TraceReader, row int)

// traceColumn pairs an Arrow field definition with the function that reads its value from storage.
type traceColumn struct {
	Field  arrow.Field
	reader columnReader
}

// Nested Arrow type definitions reused across schema and builder construction.
var (
	attrMapType = arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String)

	resourceStructType = arrow.StructOf(
		arrow.Field{Name: arrowConst.SchemaURL, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.Attributes, Type: attrMapType, Nullable: true},
	)

	scopeStructType = arrow.StructOf(
		arrow.Field{Name: arrowConst.Name, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.Version, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.SchemaURL, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.Attributes, Type: attrMapType, Nullable: true},
	)

	eventStructType = arrow.StructOf(
		arrow.Field{Name: arrowConst.Timestamp, Type: arrow.FixedWidthTypes.Timestamp_ns},
		arrow.Field{Name: arrowConst.Name, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.Attributes, Type: attrMapType, Nullable: true},
	)

	linkStructType = arrow.StructOf(
		arrow.Field{Name: arrowConst.TraceID, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.SpanID, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.TraceState, Type: arrow.BinaryTypes.String},
		arrow.Field{Name: arrowConst.Flags, Type: arrow.PrimitiveTypes.Uint32},
		arrow.Field{Name: arrowConst.Attributes, Type: attrMapType, Nullable: true},
	)
)

// traceSchema is the canonical schema for the trace table.
// Adding a new column here automatically propagates to both broker metadata
// registration and the source connector's column-reading pipeline.
var traceSchema = []traceColumn{
	// ── Direct span fields ──────────────────────────────────────────────────
	{
		// trace_id stored as hex string so SQL WHERE trace_id = '...' comparisons work.
		Field: arrow.Field{Name: arrowConst.TraceID, Type: arrow.BinaryTypes.String},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.TraceID)
			return func(r *tracespkg.TraceReader, row int) {
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
		Field: arrow.Field{Name: arrowConst.SpanID, Type: arrow.BinaryTypes.String},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.SpanID)
			return func(r *tracespkg.TraceReader, row int) {
				v := r.SpanID(row)
				if v == nil {
					b.AppendNull()
				} else {
					b.Append(hex.EncodeToString(v))
				}
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.ParentSpanID, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.ParentSpanID)
			return func(r *tracespkg.TraceReader, row int) {
				v := r.ParentSpanID(row)
				if v == nil {
					b.AppendNull()
				} else {
					b.Append(hex.EncodeToString(v))
				}
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Name, Type: arrow.BinaryTypes.String},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.Name)
			return func(r *tracespkg.TraceReader, row int) { b.Append(r.Name(row)) }
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Kind, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.Kind)
			return func(r *tracespkg.TraceReader, row int) { b.Append(r.Kind(row)) }
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Status, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.Status)
			return func(r *tracespkg.TraceReader, row int) { b.Append(r.Status(row)) }
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.StatusMessage, Type: arrow.BinaryTypes.String, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.StringBuilder(arrowConst.StatusMessage)
			return func(r *tracespkg.TraceReader, row int) { b.Append(r.StatusMessage(row)) }
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.StartTime, Type: arrow.FixedWidthTypes.Timestamp_ns},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.TimestampBuilder(arrowConst.StartTime)
			return func(r *tracespkg.TraceReader, row int) {
				b.Append(arrow.Timestamp(r.StartTime(row)))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.Duration, Type: arrow.FixedWidthTypes.Duration_ns},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.DurationBuilder(arrowConst.Duration)
			return func(r *tracespkg.TraceReader, row int) {
				b.Append(arrow.Duration(r.Duration(row)))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.DroppedAttributesCount, Type: arrow.PrimitiveTypes.Uint16},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.Uint16Builder(arrowConst.DroppedAttributesCount)
			return func(r *tracespkg.TraceReader, row int) {
				b.Append(r.DroppedAttributesCount(row))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.DroppedEventsCount, Type: arrow.PrimitiveTypes.Uint16},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.Uint16Builder(arrowConst.DroppedEventsCount)
			return func(r *tracespkg.TraceReader, row int) {
				b.Append(r.DroppedEventsCount(row))
			}
		},
	},
	{
		Field: arrow.Field{Name: arrowConst.DroppedLinksCount, Type: arrow.PrimitiveTypes.Uint16},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			b := rb.Uint16Builder(arrowConst.DroppedLinksCount)
			return func(r *tracespkg.TraceReader, row int) {
				b.Append(r.DroppedLinksCount(row))
			}
		},
	},

	// ── FK-resolved fields ───────────────────────────────────────────────────
	{
		// attributes: resolve FK List<Uint32> → map<string,string>
		Field: arrow.Field{Name: arrowConst.Attributes, Type: attrMapType, Nullable: true},
		reader: func(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
			mb := rb.MapBuilder(arrowConst.Attributes)
			kb := mb.KeyBuilder().(*array.StringBuilder)
			vb := mb.ItemBuilder().(*array.StringBuilder)
			return func(r *tracespkg.TraceReader, row int) {
				appendAttrMap(mb, kb, vb, func(fn func(k, v string)) { r.Attributes(row, fn) })
			}
		},
	},
	{
		// resource: Struct<schema_url, attributes>
		Field:  arrow.Field{Name: arrowConst.Resource, Type: resourceStructType, Nullable: true},
		reader: buildResourceReader,
	},
	{
		// scope: Struct<name, version, schema_url, attributes>
		Field:  arrow.Field{Name: arrowConst.Scope, Type: scopeStructType, Nullable: true},
		reader: buildScopeReader,
	},
	{
		// events: List<Struct<timestamp, name, attributes>>
		Field:  arrow.Field{Name: arrowConst.Events, Type: arrow.ListOf(eventStructType), Nullable: true},
		reader: buildEventsReader,
	},
	{
		// links: List<Struct<trace_id, span_id, trace_state, flags, attributes>>
		Field:  arrow.Field{Name: arrowConst.Links, Type: arrow.ListOf(linkStructType), Nullable: true},
		reader: buildLinksReader,
	},
}

// buildResourceReader creates a StructBuilder appender for the resource column.
func buildResourceReader(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
	idx := rb.Schema().FieldIndices(arrowConst.Resource)[0]
	sb := rb.Fields()[idx].(*array.StructBuilder)
	schemaURLBuilder := sb.FieldBuilder(0).(*array.StringBuilder)
	attrMapBuilder := sb.FieldBuilder(1).(*array.MapBuilder)
	attrKeyBuilder := attrMapBuilder.KeyBuilder().(*array.StringBuilder)
	attrValBuilder := attrMapBuilder.ItemBuilder().(*array.StringBuilder)

	return func(r *tracespkg.TraceReader, row int) {
		sb.Append(true)
		r.Resource(row, func(res *tracespkg.ResourceReader, resRow int) {
			schemaURLBuilder.Append(res.SchemaURL(resRow))
			appendAttrMap(attrMapBuilder, attrKeyBuilder, attrValBuilder,
				func(fn func(k, v string)) { res.Attributes(resRow, fn) })
		})
	}
}

// buildScopeReader creates a StructBuilder appender for the scope column.
func buildScopeReader(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
	idx := rb.Schema().FieldIndices(arrowConst.Scope)[0]
	sb := rb.Fields()[idx].(*array.StructBuilder)
	nameBuilder := sb.FieldBuilder(0).(*array.StringBuilder)
	versionBuilder := sb.FieldBuilder(1).(*array.StringBuilder)
	schemaURLBuilder := sb.FieldBuilder(2).(*array.StringBuilder)
	attrMapBuilder := sb.FieldBuilder(3).(*array.MapBuilder)
	attrKeyBuilder := attrMapBuilder.KeyBuilder().(*array.StringBuilder)
	attrValBuilder := attrMapBuilder.ItemBuilder().(*array.StringBuilder)

	return func(r *tracespkg.TraceReader, row int) {
		sb.Append(true)
		r.Scope(row, func(sc *tracespkg.ScopeReader, scRow int) {
			nameBuilder.Append(sc.Name(scRow))
			versionBuilder.Append(sc.Version(scRow))
			schemaURLBuilder.Append(sc.SchemaURL(scRow))
			appendAttrMap(attrMapBuilder, attrKeyBuilder, attrValBuilder,
				func(fn func(k, v string)) { sc.Attributes(scRow, fn) })
		})
	}
}

// buildEventsReader creates a ListBuilder appender for the events column.
func buildEventsReader(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
	idx := rb.Schema().FieldIndices(arrowConst.Events)[0]
	lb := rb.Fields()[idx].(*array.ListBuilder)
	esb := lb.ValueBuilder().(*array.StructBuilder)
	tsBuilder := esb.FieldBuilder(0).(*array.TimestampBuilder)
	nameBuilder := esb.FieldBuilder(1).(*array.StringBuilder)
	attrMapBuilder := esb.FieldBuilder(2).(*array.MapBuilder)
	attrKeyBuilder := attrMapBuilder.KeyBuilder().(*array.StringBuilder)
	attrValBuilder := attrMapBuilder.ItemBuilder().(*array.StringBuilder)

	return func(r *tracespkg.TraceReader, row int) {
		lb.Append(true)
		r.Events(row, func(er *tracespkg.EventReader, erRow int) {
			esb.Append(true)
			tsBuilder.Append(arrow.Timestamp(er.Timestamp(erRow)))
			nameBuilder.Append(er.Name(erRow))
			appendAttrMap(attrMapBuilder, attrKeyBuilder, attrValBuilder,
				func(fn func(k, v string)) { er.Attributes(erRow, fn) })
		})
	}
}

// buildLinksReader creates a ListBuilder appender for the links column.
func buildLinksReader(rb *builder.RecordBuilder) func(*tracespkg.TraceReader, int) {
	idx := rb.Schema().FieldIndices(arrowConst.Links)[0]
	lb := rb.Fields()[idx].(*array.ListBuilder)
	lsb := lb.ValueBuilder().(*array.StructBuilder)
	// trace_id and span_id in links are stored as hex strings (consistent with top-level fields).
	traceIDBuilder := lsb.FieldBuilder(0).(*array.StringBuilder)
	spanIDBuilder := lsb.FieldBuilder(1).(*array.StringBuilder)
	traceStateBuilder := lsb.FieldBuilder(2).(*array.StringBuilder)
	flagsBuilder := lsb.FieldBuilder(3).(*array.Uint32Builder)
	attrMapBuilder := lsb.FieldBuilder(4).(*array.MapBuilder)
	attrKeyBuilder := attrMapBuilder.KeyBuilder().(*array.StringBuilder)
	attrValBuilder := attrMapBuilder.ItemBuilder().(*array.StringBuilder)

	return func(r *tracespkg.TraceReader, row int) {
		lb.Append(true)
		r.Links(row, func(lr *tracespkg.LinkReader, lrRow int) {
			lsb.Append(true)
			if v := lr.TraceID(lrRow); v != nil {
				traceIDBuilder.Append(hex.EncodeToString(v))
			} else {
				traceIDBuilder.AppendNull()
			}
			if v := lr.SpanID(lrRow); v != nil {
				spanIDBuilder.Append(hex.EncodeToString(v))
			} else {
				spanIDBuilder.AppendNull()
			}
			traceStateBuilder.Append(lr.TraceState(lrRow))
			flagsBuilder.Append(lr.Flags(lrRow))
			appendAttrMap(attrMapBuilder, attrKeyBuilder, attrValBuilder,
				func(fn func(k, v string)) { lr.Attributes(lrRow, fn) })
		})
	}
}

// appendAttrMap resolves attribute FK references to map<string,string> via callback.
// Since TraceReader only exposes Attributes(row, fn), we collect into a temp slice first.
func appendAttrMap(mb *array.MapBuilder, kb, vb *array.StringBuilder, iter func(fn func(k, v string))) {
	var keys, vals []string
	iter(func(k, v string) {
		keys = append(keys, k)
		vals = append(vals, v)
	})
	if len(keys) == 0 {
		mb.AppendNull()
		return
	}
	mb.Append(true)
	mb.Reserve(len(keys))
	for i := range keys {
		kb.Append(keys[i])
		vb.Append(vals[i])
	}
}

// traceSchemaIndex maps column name → traceColumn for O(1) lookup at query time.
var traceSchemaIndex map[string]*traceColumn

func init() {
	traceSchemaIndex = make(map[string]*traceColumn, len(traceSchema))
	for i := range traceSchema {
		traceSchemaIndex[traceSchema[i].Field.Name] = &traceSchema[i]
	}
	// Register the canonical Arrow fields so broker_meta can read them without
	// importing this package directly (which would create an import cycle via spi).
	types.RegisterTableArrowFields("traces", TraceTableArrowFields())
}

// TraceTableArrowFields returns the Arrow field list for the trace table.
// Used by broker_meta to register the table metadata returned to the SQL planner.
func TraceTableArrowFields() []arrow.Field {
	fields := make([]arrow.Field, len(traceSchema))
	for i, col := range traceSchema {
		fields[i] = col.Field
	}
	return fields
}

// BuildTraceColumnAppender returns an appender for the named fixed-schema column, or
// (nil, false) when the name is not in the trace schema.
// Used by buildColumnAppenders in the source connector.
func BuildTraceColumnAppender(name string, rb *builder.RecordBuilder) (func(*tracespkg.TraceReader, int), bool) {
	col, ok := traceSchemaIndex[name]
	if !ok {
		return nil, false
	}
	return col.reader(rb), true
}
