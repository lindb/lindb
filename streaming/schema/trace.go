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

package schema

import (
	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/constants"
)

var spanJoinMeta = arrow.MetadataFrom(map[string]string{
	constants.MetadataNameKey: string(constants.DTSpans),
})

// SpanJoinSchema is the join schema for span records.
var SpanJoinSchema = arrow.NewSchema([]arrow.Field{
	larrow.JoinField(constants.TraceID, &arrow.FixedSizeBinaryType{ByteWidth: 16}, string(constants.DTSpans), constants.TraceID),
	larrow.JoinField(constants.SpanID, &arrow.FixedSizeBinaryType{ByteWidth: 8}, string(constants.DTSpans), constants.SpanID),
	larrow.JoinField(constants.ParentSpanID, &arrow.FixedSizeBinaryType{ByteWidth: 8}, string(constants.DTSpans), constants.ParentSpanID),
	larrow.JoinField(constants.StartTime, arrow.FixedWidthTypes.Timestamp_ns, string(constants.DTSpans), constants.StartTime),
	larrow.JoinField(constants.Duration, arrow.FixedWidthTypes.Duration_ns, string(constants.DTSpans), constants.Duration),
	larrow.JoinField(constants.Name, arrow.BinaryTypes.String, string(constants.DTSpans), constants.Name),
	larrow.JoinField(constants.Status, arrow.BinaryTypes.String, string(constants.DTSpans), constants.Status),
	larrow.JoinField(constants.Kind, arrow.BinaryTypes.String, string(constants.DTSpans), constants.Kind),
	larrow.JoinDerefField(constants.Resource, string(constants.DTSpans), constants.Resource, string(constants.DTResources)),
	larrow.JoinDerefField(constants.Scope, string(constants.DTSpans), constants.Scope, string(constants.DTScopes)),
	larrow.JoinDerefField(constants.Attributes, string(constants.DTSpans), constants.Attributes, string(constants.DTAttributes)),
}, &spanJoinMeta)

var eventJoinMeta = arrow.MetadataFrom(map[string]string{
	constants.MetadataNameKey: string(constants.DTEvents),
})

// EventJoinSchema is the join schema for event records.
var EventJoinSchema = arrow.NewSchema([]arrow.Field{
	larrow.JoinField(constants.Timestamp, arrow.FixedWidthTypes.Timestamp_ns, string(constants.DTEvents), constants.Timestamp),
	larrow.JoinField(constants.Name, arrow.BinaryTypes.String, string(constants.DTEvents), constants.Name),
	larrow.JoinDerefField(constants.Attributes, string(constants.DTEvents), constants.Attributes, string(constants.DTAttributes)),
}, &eventJoinMeta)
