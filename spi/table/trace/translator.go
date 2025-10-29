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
	protoTraceV1 "github.com/lindb/common/proto/gen/v1/lintrace"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TranslateResourceSpans(rs ptrace.ResourceSpans, traceID string) *protoTraceV1.CallStack {
	scopeSpans := rs.ScopeSpans()

	if scopeSpans.Len() == 0 {
		return nil
	}
	var lSpans []*protoTraceV1.Span

	for i := 0; i < scopeSpans.Len(); i++ {
		scopeSpan := scopeSpans.At(i)
		spans := scopeSpan.Spans()
		for j := 0; j < spans.Len(); j++ {
			span := spans.At(j)
			if span.TraceID().String() != traceID {
				continue
			}

			lSpans = append(lSpans, translateSpan(span))
		}
	}

	if len(lSpans) == 0 {
		return nil
	}

	return &protoTraceV1.CallStack{
		Spans:    lSpans,
		Resource: translateResource(rs.Resource()),
	}
}

func translateResource(resource pcommon.Resource) *protoTraceV1.Resource {
	attributes := resource.Attributes()
	if attributes.Len() == 0 {
		return nil
	}

	return &protoTraceV1.Resource{
		Attributes: translateAttributes(resource.Attributes()),
	}
}

func translateSpan(span ptrace.Span) *protoTraceV1.Span {
	status := span.Status()
	rs := &protoTraceV1.Span{
		TraceId:      span.TraceID().String(),
		SpanId:       span.SpanID().String(),
		ParentSpanId: span.ParentSpanID().String(),

		Name:         span.Name(),
		Kind:         span.Kind().String(),
		Status:       status.Code().String(),
		ErrorMessage: status.Message(),

		StartTime: span.StartTimestamp().AsTime().UnixNano(),
		EndTime:   span.EndTimestamp().AsTime().UnixNano(),

		Attributes: translateAttributes(span.Attributes()),

		Events: translateEvents(span.Events()),
		Links:  translateLinks(span.Links()),
	}

	return rs
}

func translateEvents(events ptrace.SpanEventSlice) []*protoTraceV1.Event {
	if events.Len() == 0 {
		return nil
	}

	rs := make([]*protoTraceV1.Event, 0, events.Len())
	for i := 0; i < events.Len(); i++ {
		event := events.At(i)

		rs = append(rs, &protoTraceV1.Event{
			Name:       event.Name(),
			Timestamp:  event.Timestamp().AsTime().UnixNano(),
			Attributes: translateAttributes(event.Attributes()),
		})
	}

	return rs
}

func translateLinks(links ptrace.SpanLinkSlice) []*protoTraceV1.Link {
	if links.Len() == 0 {
		return nil
	}

	rs := make([]*protoTraceV1.Link, 0, links.Len())
	for i := 0; i < links.Len(); i++ {
		link := links.At(i)

		rs = append(rs, &protoTraceV1.Link{
			TraceId:    link.TraceID().String(),
			SpanId:     link.SpanID().String(),
			Attributes: translateAttributes(link.Attributes()),
		})
	}

	return rs
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
