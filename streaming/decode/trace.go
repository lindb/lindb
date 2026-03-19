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
	"github.com/lindb/arrow/pkg/traces"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/streaming/schema"
)

func init() {
	RegisterDecoder(option.Trace, newTrace)
}

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
		t.spanJoiner = larrow.NewRecordJoiner(schema.SpanJoinSchema)
		t.eventJoiner = larrow.NewRecordJoiner(schema.EventJoinSchema)
	}

	records := t.reader.Records()
	spanRecord := t.spanJoiner.Join(records)
	eventRecord := t.eventJoiner.Join(records)
	return []arrow.RecordBatch{spanRecord, eventRecord}, nil
}
