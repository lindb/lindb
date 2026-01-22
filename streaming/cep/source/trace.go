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
	SpanStream  = "span"
	EventStream = "event"
)

func init() {
	RegisterSource(SourceType(option.Trace), func(runtime runtime.Runtime) Source {
		// register span schema
		emptyPage := NewSpanPageBuilder().Build()
		runtime.RegisterStreamBySchema(SpanStream, &types.TableSchema{Columns: emptyPage.Layout})
		// register event schema
		emptyPage = NewEventPageBuilder().Build()
		runtime.RegisterStreamBySchema(EventStream, &types.TableSchema{Columns: emptyPage.Layout})

		return &trace{
			runtime: runtime,
			logger:  logger.GetLogger("CEP", "TraceSource"),
		}
	})
}

type trace struct {
	runtime runtime.Runtime

	spans  *SpanPageBuilder
	events *EventPageBuilder

	logger logger.Logger
}

// Receive implements [Source].
func (t *trace) Receive(e models.Event) {
	traces, ok := e.(ptrace.Traces)
	if !ok {
		t.logger.Warn("trace source receive invalid event type", logger.Any("event", reflect.TypeOf(e)))
		return
	}
	// TODO: check if need create new builders per receive
	t.spans = NewSpanPageBuilder()
	t.events = NewEventPageBuilder()

	t.ToPage(traces)

	t.runtime.GetInputHandler(SpanStream).Send(t.spans.Build())
	t.runtime.GetInputHandler(EventStream).Send(t.events.Build())
}

func (t *trace) ToPage(traces ptrace.Traces) {
	resourceSpans := traces.ResourceSpans()

	if resourceSpans.Len() == 0 {
		return
	}

	for i := 0; i < resourceSpans.Len(); i++ {
		t.translateResourceSpans(resourceSpans.At(i))
	}
}

func (t *trace) translateResourceSpans(rs ptrace.ResourceSpans) {
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

			t.spans.AppendSpan(resource, span)

			t.translateSpanEvents(resource, span)
		}
	}
}

func (t *trace) translateSpanEvents(resource map[string]string, span ptrace.Span) {
	events := span.Events()
	for i := 0; i < events.Len(); i++ {
		t.events.AppendEvent(resource, span, events.At(i))
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

type SpanPageBuilder struct {
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

func NewSpanPageBuilder() *SpanPageBuilder {
	builder := &SpanPageBuilder{
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

func (b *SpanPageBuilder) AppendSpan(resource map[string]string, span ptrace.Span) {
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

func (b *SpanPageBuilder) Build() *types.Page {
	return b.page
}

type EventPageBuilder struct {
	page *types.Page

	traceID *types.Column
	spanID  *types.Column

	name       *types.Column
	timestamap *types.Column
	attributes *types.Column

	resource *types.Column
}

func NewEventPageBuilder() *EventPageBuilder {
	builder := &EventPageBuilder{
		page: types.NewPage(),

		traceID: types.NewColumn(),
		spanID:  types.NewColumn(),

		name:       types.NewColumn(),
		timestamap: types.NewColumn(),
		attributes: types.NewColumn(),

		resource: types.NewColumn(),
	}

	builder.page.AppendColumn(types.ColumnMetadata{Name: "trace_id", DataType: types.DTString}, builder.traceID)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "span_id", DataType: types.DTString}, builder.spanID)

	builder.page.AppendColumn(types.ColumnMetadata{Name: "name", DataType: types.DTString}, builder.name)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "start_time", DataType: types.DTTimestamp}, builder.timestamap)
	builder.page.AppendColumn(types.ColumnMetadata{Name: "attributes", DataType: types.DTMap}, builder.attributes)

	builder.page.AppendColumn(types.ColumnMetadata{Name: "resource", DataType: types.DTMap}, builder.resource)

	return builder
}

func (b *EventPageBuilder) AppendEvent(resource map[string]string, span ptrace.Span, event ptrace.SpanEvent) {
	// translate event
	b.traceID.Append(span.TraceID().String())
	b.spanID.Append(span.SpanID().String())

	b.name.Append(span.Name())
	b.timestamap.Append(event.Timestamp().AsTime())
	b.attributes.Append(translateAttributes(span.Attributes()))

	// set resource
	b.resource.Append(resource)
}

func (b *EventPageBuilder) Build() *types.Page {
	return b.page
}
