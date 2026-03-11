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
	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/constants"
	"github.com/lindb/arrow/pkg/traces"

	"github.com/lindb/lindb/pkg/option"
)

func init() {
	RegisterDecoder(option.Trace, newTrace)
}

var spanJoinSchema = arrow.NewSchema([]arrow.Field{
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
}, nil)

var eventJoinSchema = arrow.NewSchema([]arrow.Field{
	larrow.JoinField(constants.Timestamp, arrow.FixedWidthTypes.Timestamp_ns, string(constants.DTEvents), constants.Timestamp),
	larrow.JoinField(constants.Name, arrow.BinaryTypes.String, string(constants.DTEvents), constants.Name),
	larrow.JoinDerefField(constants.Attributes, string(constants.DTEvents), constants.Attributes, string(constants.DTAttributes)),
}, nil)

type trace struct {
	reader      *traces.TraceReader
	spanJoiner  *larrow.RecordJoiner
	eventJoiner *larrow.RecordJoiner
}

func newTrace() Decoder {
	return &trace{}
}

func (t *trace) ToRecords(data []byte) ([]arrow.RecordBatch, error) {
	if t.reader == nil {
		reader, err := traces.NewTraceReader(data)
		if err != nil {
			return nil, err
		}
		t.reader = reader
	} else {
		if err := t.reader.Reset(data); err != nil {
			return nil, err
		}
	}

	if t.spanJoiner == nil {
		t.spanJoiner = larrow.NewRecordJoiner(spanJoinSchema)
		t.eventJoiner = larrow.NewRecordJoiner(eventJoinSchema)
	}

	records := t.reader.Records()
	spanRecord := t.spanJoiner.Join(records)
	eventRecord := t.eventJoiner.Join(records)
	return []arrow.RecordBatch{spanRecord, eventRecord}, nil
}
