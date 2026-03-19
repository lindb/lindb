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

package streaming_test

import (
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/arrow/pkg/traces"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/streaming"
	"github.com/lindb/lindb/streaming/cep"
)

// buildTraceBytes constructs OTel trace data (2 spans, 1 event) and serialises
// it into the Arrow binary format that DataSource.Produce expects.
func buildTraceBytes(t *testing.T) []byte {
	t.Helper()

	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().PutStr("service.name", "checkout")

	ss := rs.ScopeSpans().AppendEmpty()
	ss.Scope().SetName("checkout-scope")
	ss.Scope().SetVersion("v1.0")

	now := time.Now()

	// span 1: server span with one exception event
	span1 := ss.Spans().AppendEmpty()
	traceID := pcommon.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	spanID1 := pcommon.SpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	span1.SetTraceID(traceID)
	span1.SetSpanID(spanID1)
	span1.SetName("checkout")
	span1.SetKind(ptrace.SpanKindServer)
	span1.Status().SetCode(ptrace.StatusCodeOk)
	span1.SetStartTimestamp(pcommon.NewTimestampFromTime(now))
	span1.SetEndTimestamp(pcommon.NewTimestampFromTime(now.Add(20 * time.Millisecond)))
	span1.Attributes().PutStr("http.method", "POST")
	span1.Attributes().PutStr("http.route", "/checkout")

	ev := span1.Events().AppendEmpty()
	ev.SetName("exception")
	ev.SetTimestamp(pcommon.NewTimestampFromTime(now.Add(5 * time.Millisecond)))
	ev.Attributes().PutStr("exception.type", "NullPointerException")

	// span 2: client child span, error status
	span2 := ss.Spans().AppendEmpty()
	spanID2 := pcommon.SpanID([8]byte{2, 2, 3, 4, 5, 6, 7, 8})
	span2.SetTraceID(traceID)
	span2.SetSpanID(spanID2)
	span2.SetParentSpanID(spanID1)
	span2.SetName("db.query")
	span2.SetKind(ptrace.SpanKindClient)
	span2.Status().SetCode(ptrace.StatusCodeError)
	span2.SetStartTimestamp(pcommon.NewTimestampFromTime(now.Add(2 * time.Millisecond)))
	span2.SetEndTimestamp(pcommon.NewTimestampFromTime(now.Add(15 * time.Millisecond)))
	span2.Attributes().PutStr("db.system", "postgresql")

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

// TestDataSource_TraceAnalysis exercises the full trace CEP pipeline:
//
//	OTel bytes → DataSource.Produce → decode → cep.Engine.Send
//	  → source.Receive (routes by metadata name) → runtime input handlers
//	  → LinQL jobs query "spans" and "events" tables
func TestDataSource_TraceAnalysis(t *testing.T) {
	db := &models.Database{
		Name: "trace_db",
		Option: &option.DatabaseOption{
			Engine: option.Trace,
		},
	}

	ds := streaming.NewDataSource(db)
	ds.Startup()
	defer ds.Shutdown()

	// Register a streaming engine for this database.
	stream := &models.Streaming{
		Name:     "trace_cep",
		Database: db.Name,
	}
	err := ds.ScheduleStream(stream)
	require.NoError(t, err)

	engine, ok := ds.GetEngine("trace_cep")
	require.True(t, ok)

	cepEngine := engine.(*cep.Engine)

	// Deploy a job that selects all span rows and all event rows.
	err = cepEngine.DeployJob("trace_analysis", `
		create job trace_analysis
		begin
		  select name, status, duration from spans;
		  select name, timestamp from events;
		end
	`)
	require.NoError(t, err)

	// Produce OTel trace data – triggers decode → source route → CEP query.
	traceBytes := buildTraceBytes(t)
	err = ds.Produce(traceBytes)
	require.NoError(t, err)

	// Produce a second batch to exercise the reader Reset path.
	err = ds.Produce(traceBytes)
	require.NoError(t, err)

	// Allow the async pipeline to drain.
	time.Sleep(1 * time.Second)
}
