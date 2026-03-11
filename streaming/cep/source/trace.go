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
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/arrow/pkg/constants"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/streaming/cep/runtime"
	"github.com/lindb/lindb/streaming/schema"
)

func init() {
	RegisterSource(SourceType(option.Trace), func(rt runtime.Runtime) Source {
		for _, s := range []*arrow.Schema{schema.SpanJoinSchema, schema.EventJoinSchema} {
			if name, ok := s.Metadata().GetValue(constants.MetadataNameKey); ok {
				_ = rt.RegisterStream(name, s)
			}
		}
		return &trace{
			runtime: rt,
			logger:  logger.GetLogger("CEP", "TraceSource"),
		}
	})
}

type trace struct {
	runtime runtime.Runtime
	logger  logger.Logger
}

// Receive implements [Source].
func (t *trace) Receive(record arrow.RecordBatch) {
	name, ok := record.Schema().Metadata().GetValue(constants.MetadataNameKey)
	if !ok {
		return
	}
	t.runtime.GetInputHandler(name).Send(record)
}
