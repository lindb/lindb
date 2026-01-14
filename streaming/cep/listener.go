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

package cep

import (
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep/annotation"
	"github.com/lindb/lindb/streaming/cep/sink"
)

type Listener struct {
	mapper annotation.Mapper
	sinks  []sink.Sink
	logger logger.Logger
}

func NewListener(mapper annotation.Mapper, sinks []sink.Sink) *Listener {
	return &Listener{
		mapper: mapper,
		sinks:  sinks,
		logger: logger.GetLogger("CEP", "Listener"),
	}
}

func (l *Listener) Receive(event models.Event) {
	// TODO: add multiple mappers support?
	if l.mapper != nil {
		event = l.mapper.Map(event)
	}
	for _, s := range l.sinks {
		s.Publish(event)
	}
	l.logger.Info("Listener, receive event", logger.Any("event", event))
}
