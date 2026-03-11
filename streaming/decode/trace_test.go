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

package decode

import (
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/arrow/pkg/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// buildTestTraceBytes constructs a minimal OTel trace (1 resource span, 1 scope, 2 spans,
// one of which carries an event) and serialises it via SpanBuilder into the binary format
// that ToRecords expects.
func buildTestTraceBytes(t *testing.T) []byte {
	t.Helper()

	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().PutStr("service.name", "test-service")

	ss := rs.ScopeSpans().AppendEmpty()
	ss.Scope().SetName("test-scope")
	ss.Scope().SetVersion("v1.0")

	now := time.Now()

	// span 1 – has one event
	span1 := ss.Spans().AppendEmpty()
	traceID := pcommon.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	spanID1 := pcommon.SpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	span1.SetTraceID(traceID)
	span1.SetSpanID(spanID1)
	span1.SetName("op1")
	span1.SetKind(ptrace.SpanKindServer)
	span1.Status().SetCode(ptrace.StatusCodeOk)
	span1.SetStartTimestamp(pcommon.NewTimestampFromTime(now))
	span1.SetEndTimestamp(pcommon.NewTimestampFromTime(now.Add(10 * time.Millisecond)))
	span1.Attributes().PutStr("http.method", "GET")

	ev := span1.Events().AppendEmpty()
	ev.SetName("exception")
	ev.SetTimestamp(pcommon.NewTimestampFromTime(now.Add(5 * time.Millisecond)))
	ev.Attributes().PutStr("exception.type", "NullPointerException")

	// span 2 – no event, child of span 1
	span2 := ss.Spans().AppendEmpty()
	spanID2 := pcommon.SpanID([8]byte{2, 2, 3, 4, 5, 6, 7, 8})
	span2.SetTraceID(traceID)
	span2.SetSpanID(spanID2)
	span2.SetParentSpanID(spanID1)
	span2.SetName("op2")
	span2.SetKind(ptrace.SpanKindClient)
	span2.Status().SetCode(ptrace.StatusCodeError)
	span2.SetStartTimestamp(pcommon.NewTimestampFromTime(now.Add(1 * time.Millisecond)))
	span2.SetEndTimestamp(pcommon.NewTimestampFromTime(now.Add(8 * time.Millisecond)))

	// Build Arrow binary via SpanBuilder (same path as broker writer)
	builder := traces.NewSpanBuilder(memory.NewGoAllocator())
	defer builder.Release() //nolint:errcheck

	for i := 0; i < rs.ScopeSpans().Len(); i++ {
		scopeSpan := rs.ScopeSpans().At(i)
		for j := 0; j < scopeSpan.Spans().Len(); j++ {
			builder.Append(&model.Span{
				Span:     scopeSpan.Spans().At(j),
				Scope:    scopeSpan,
				Resource: rs,
			})
		}
	}

	data, err := builder.Bytes()
	require.NoError(t, err)
	return data
}

func TestTraceDecoder_ToRecords(t *testing.T) {
	data := buildTestTraceBytes(t)

	dec := newTrace()
	records, err := dec.ToRecords(data)
	require.NoError(t, err)

	// expect exactly 2 joined record batches: spans and events
	assert.Len(t, records, 2)

	spanRecord := records[0]
	eventRecord := records[1]

	// 2 spans were written
	assert.Equal(t, int64(2), spanRecord.NumRows())

	// 1 event was written
	assert.Equal(t, int64(1), eventRecord.NumRows())

	// span record must contain the expected column names
	spanSchema := spanRecord.Schema()
	expectedSpanCols := []string{"trace_id", "span_id", "parent_span_id", "start_time", "duration", "name", "status", "kind", "resource", "scope", "attributes"}
	for _, col := range expectedSpanCols {
		indices := spanSchema.FieldIndices(col)
		assert.NotEmpty(t, indices, "span record missing column %q", col)
	}

	// event record must contain the expected column names
	eventSchema := eventRecord.Schema()
	expectedEventCols := []string{"timestamp", "name", "attributes"}
	for _, col := range expectedEventCols {
		indices := eventSchema.FieldIndices(col)
		assert.NotEmpty(t, indices, "event record missing column %q", col)
	}
}

func TestTraceDecoder_ToRecords_Reset(t *testing.T) {
	data := buildTestTraceBytes(t)

	dec := newTrace()

	// first call – initialises reader and joiners
	records1, err := dec.ToRecords(data)
	require.NoError(t, err)
	assert.Len(t, records1, 2)

	// second call with the same data – exercises the Reset path
	records2, err := dec.ToRecords(data)
	require.NoError(t, err)
	assert.Len(t, records2, 2)
	assert.Equal(t, records1[0].NumRows(), records2[0].NumRows())
	assert.Equal(t, records1[1].NumRows(), records2[1].NumRows())
}
