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

package source

import (
	"reflect"

	"github.com/lindb/common/pkg/logger"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/cep/runtime"
)

const (
	SpanStream = "span"
)

func init() {
	RegisterSource(SourceType(option.Trace), func(runtime runtime.Runtime) Source {
		// register span schema
		emptyPage := NewTracePageBuilder().Build()
		runtime.RegisterStreamBySchema(SpanStream, &types.TableSchema{Columns: emptyPage.Layout})

		return &trace{
			runtime: runtime,
			logger:  logger.GetLogger("CEP", "TraceSource"),
		}
	})
}

type trace struct {
	runtime runtime.Runtime

	logger logger.Logger
}

// Receive implements [Source].
func (t *trace) Receive(e models.Event) {
	traces, ok := e.(ptrace.Traces)
	if !ok {
		t.logger.Warn("trace source receive invalid event type", logger.Any("event", reflect.TypeOf(e)))
		return
	}
	page, err := ToPage(traces)
	if err != nil {
		return
	}
	t.runtime.GetInputHandler(SpanStream).Send(page)
}

func ToPage(traces ptrace.Traces) (*types.Page, error) {
	resourceSpans := traces.ResourceSpans()

	if resourceSpans.Len() == 0 {
		return nil, nil
	}

	builder := NewTracePageBuilder()

	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		translateResourceSpans(rs, builder)
	}

	return builder.Build(), nil
}

func translateResourceSpans(rs ptrace.ResourceSpans, builder *TracePageBuilder) {
	scopeSpans := rs.ScopeSpans()

	if scopeSpans.Len() == 0 {
		return
	}
	resource := translateResource(rs.Resource())

	for i := 0; i < scopeSpans.Len(); i++ {
		scopeSpan := scopeSpans.At(i)
		spans := scopeSpan.Spans()
		for j := 0; j < spans.Len(); j++ {
			span := spans.At(j)

			builder.AppendSpan(resource, span)

			// TODO: add events
		}
	}
}

func translateResource(resource pcommon.Resource) map[string]string {
	attributes := resource.Attributes()
	if attributes.Len() == 0 {
		return nil
	}
	return translateAttributes(attributes)
}

func translateAttributes(attrs pcommon.Map) map[string]string {
	if attrs.Len() == 0 {
		return nil
	}
	rs := make(map[string]string)
	attrs.Range(func(k string, v pcommon.Value) bool {
		rs[k] = v.AsString()
		return true
	})
	return rs
}

type TracePageBuilder struct {
	page *types.Page

	traceID      *types.Column
	spanID       *types.Column
	parentSpanID *types.Column

	name         *types.Column
	kind         *types.Column
	status       *types.Column
	errorMessage *types.Column
	startTime    *types.Column
	endTime      *types.Column
	duration     *types.Column
	attributes   *types.Column

	resource *types.Column
}

func NewTracePageBuilder() *TracePageBuilder {
	builder := &TracePageBuilder{
		page: types.NewPage(),

		traceID:      types.NewColumn(),
		spanID:       types.NewColumn(),
		parentSpanID: types.NewColumn(),

		name:         types.NewColumn(),
		kind:         types.NewColumn(),
		status:       types.NewColumn(),
		errorMessage: types.NewColumn(),
		startTime:    types.NewColumn(),
		endTime:      types.NewColumn(),
		duration:     types.NewColumn(),
		attributes:   types.NewColumn(),
		resource:     types.NewColumn(),
	}

	builder.page.AppendColumn(types.ColumnMetadata{Name: "trace_id", DataType: types.DTString}, builder.traceID)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "span_id", DataType: types.DTString}, builder.spanID)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "parent_span_id", DataType: types.DTString}, builder.parentSpanID)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "name", DataType: types.DTString}, builder.name)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "kind", DataType: types.DTString}, builder.kind)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "status", DataType: types.DTString}, builder.status)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "error_message", DataType: types.DTString}, builder.errorMessage)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "start_time", DataType: types.DTTimestamp}, builder.startTime)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "end_time", DataType: types.DTTimestamp}, builder.endTime)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "duration", DataType: types.DTDuration}, builder.duration)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "attributes", DataType: types.DTMap}, builder.attributes)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "resource", DataType: types.DTMap}, builder.resource)

	return builder
}

func (b *TracePageBuilder) AppendSpan(resource map[string]string, span ptrace.Span) {
	// translate span
	b.traceID.Append(span.TraceID().String())
	b.spanID.Append(span.SpanID().String())
	b.parentSpanID.Append(span.ParentSpanID().String())
	b.name.Append(span.Name())
	b.kind.Append(span.Kind().String())
	b.attributes.Append(translateAttributes(span.Attributes()))

	status := span.Status()
	b.status.Append(status.Code().String())
	b.errorMessage.Append(status.Message())

	start := span.StartTimestamp().AsTime()
	end := span.EndTimestamp().AsTime()
	b.startTime.Append(start)
	b.endTime.Append(end)
	b.duration.Append(end.Sub(start))

	// set resource
	b.resource.Append(resource)
}

func (b *TracePageBuilder) Build() *types.Page {
	return b.page
}
