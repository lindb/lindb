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

package runtime

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/streaming/cep/sink"
)

type Listener struct {
	bridges []*sink.SinkBridge
	logger  logger.Logger
}

func NewListener(bridges []*sink.SinkBridge) *Listener {
	return &Listener{
		bridges: bridges,
		logger:  logger.GetLogger("CEP", "Listener"),
	}
}

func (l *Listener) Receive(record arrow.RecordBatch) {
	if l.logger.Enabled(logger.DebugLevel) {
		l.logger.Debug("Listener, receive record", logger.Any("record", record))
	}
	for _, s := range l.bridges {
		s.Publish(record)
	}
}
