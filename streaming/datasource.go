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

package streaming

import (
	"sync"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep"
	"github.com/lindb/lindb/streaming/decode"
)

type DataSource interface {
	Initialize()
	Name() string
	Produce(data []byte)
	ScheduleStream(stream *models.Streaming) error
	UnscheduleStream(stream *models.Streaming) error
	GetEngine(stream string) (Engine, bool)
}

type dataSource struct {
	db *models.Database

	streamings map[string]*models.Streaming
	engines    map[string]Engine

	decoder decode.Decoder

	lock sync.Mutex

	logger logger.Logger
}

func NewDataSource(db *models.Database) DataSource {
	return &dataSource{
		db:         db,
		streamings: make(map[string]*models.Streaming),
		engines:    make(map[string]Engine),

		logger: logger.GetLogger("Streaming", "DataSource"),
	}
}

func (d *dataSource) Initialize() {
	d.decoder = decode.GetDecoder(d.db.Option.Engine)
}

func (d *dataSource) Name() string {
	return d.db.Name
}

func (d *dataSource) Produce(data []byte) {
	event, err := d.decoder.ToEvent(data)
	if err != nil {
		d.logger.Error("transfer data to event error:", logger.Error(err))
		return
	}
	if event == nil {
		return
	}
	// TODO: add lock???
	for _, engine := range d.engines {
		engine.Send(event)
	}
}

func (d *dataSource) ScheduleStream(stream *models.Streaming) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, ok := d.engines[stream.Name]
	if !ok {
		// TODO: create engine based on streaming config(add engine type)
		engine := cep.NewEngine(stream, d.db)
		if err := engine.Start(); err != nil {
			return err
		}
		d.engines[stream.Name] = engine
	}

	return nil
}

func (d *dataSource) UnscheduleStream(stream *models.Streaming) error {
	egnine, ok := d.GetEngine(stream.Name)
	if ok {
		if err := egnine.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (d *dataSource) GetEngine(stream string) (Engine, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	e, ok := d.engines[stream]
	return e, ok
}
