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

package sink

import (
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/streaming/cep/annotation"
)

// Sink represents the event sink which publishes event to output transport.
type Sink interface {
	// Publish publishes the event to sink vis output transport.
	Publish(event models.Event)
	Close()
}

type createSinkFn func(props *collections.Properties) Sink

var sinks = map[string]createSinkFn{
	"lindb": newLinDBSink,
}

func CreateSink(props *collections.Properties) Sink {
	sinkType, ok := props.GetString("type")
	if !ok {
		return nil
	}
	if createFn, ok := sinks[sinkType]; ok {
		return createFn(props)
	}
	return nil
}

type SinkBridge struct {
	sink    Sink
	mappers []annotation.Mapper

	logger logger.Logger
}

func NewSinkBridge(sink Sink, sinkAnn *annotation.Annotation) *SinkBridge {
	bridge := &SinkBridge{
		sink:   sink,
		logger: logger.GetLogger("Streaming", "SinkBridge"),
	}
	for _, ann := range sinkAnn.Annotations {
		mapper := annotation.CreateMapper(ann)
		if mapper == nil {
			continue
		}
		bridge.mappers = append(bridge.mappers, mapper)
	}
	if len(bridge.mappers) == 0 {
		bridge.logger.Warn("no valid annotation mapper found for sink", logger.String("sink", sinkAnn.Name))
	}
	return bridge
}

func (b *SinkBridge) Publish(event models.Event) {
	for _, mapper := range b.mappers {
		mappedEvent := mapper.Map(event)
		if mappedEvent == nil {
			continue
		}

		// send mapped event to sink
		b.sink.Publish(mappedEvent)
	}
}
