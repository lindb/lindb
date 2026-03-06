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
	"errors"
	"sync"

	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep"
	"github.com/lindb/lindb/streaming/decode"
)

type DataSource interface {
	Initialize()
	Name() string
	Produce(data []byte) error
	ScheduleStream(stream *models.Streaming) error
	GetEngine(stream string) (Engine, bool)
	Shutdown()
}

type dataSource struct {
	db *models.Database

	streamings map[string]*models.Streaming
	engines    map[string]Engine

	decoder decode.Decoder
	running *atomic.Bool

	lock sync.Mutex

	logger logger.Logger
}

func NewDataSource(db *models.Database) DataSource {
	return &dataSource{
		db:         db,
		streamings: make(map[string]*models.Streaming),
		engines:    make(map[string]Engine),
		running:    atomic.NewBool(true),

		logger: logger.GetLogger("Streaming", "DataSource"),
	}
}

func (d *dataSource) Initialize() {
	d.decoder = decode.GetDecoder(d.db.Option.Engine)
}

func (d *dataSource) Name() string {
	return d.db.Name
}

func (d *dataSource) Produce(data []byte) error {
	if !d.running.Load() {
		return errors.New("data source not running")
	}
	record, err := d.decoder.ToRecord(data)
	if err != nil {
		d.logger.Error("transfer data to event error:", logger.Error(err))
		return err
	}
	if record == nil || record.NumRows() == 0 {
		return nil
	}
	// TODO: add lock???
	for _, engine := range d.engines {
		engine.Send(record)
	}
	return nil
}

func (d *dataSource) ScheduleStream(stream *models.Streaming) error {
	if !d.running.Load() {
		return nil
	}

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

func (d *dataSource) GetEngine(stream string) (Engine, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	e, ok := d.engines[stream]
	return e, ok
}

func (d *dataSource) Shutdown() {
	if d.running.CompareAndSwap(true, false) {
		d.lock.Lock()
		defer d.lock.Unlock()

		for _, engine := range d.engines {
			if err := engine.Stop(); err != nil {
				d.logger.Error("stop engine error:", logger.Error(err))
			}
		}
	}
}
