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
	"maps"
	"sync/atomic"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep"
	"github.com/lindb/lindb/streaming/decode"
)

// scheduleCmd carries a stream registration request from ScheduleStream to the
// background event loop, along with a reply channel for the result.
type scheduleCmd struct {
	stream  *models.Streaming
	replyCh chan error
}

// DataSource receives raw messages, decodes them into Arrow RecordBatches, and
// fans the records out to all registered streaming engines.
type DataSource interface {
	Name() string
	// Startup initialises the decoder and starts the internal event loop.
	Startup()
	// Shutdown stops all engines and the event loop.
	Shutdown()
	// Produce decodes msg and delivers every resulting record to all engines.
	// It is safe to call concurrently and never acquires a lock.
	Produce(msg []byte) error
	// ScheduleStream registers a new streaming engine asynchronously.
	// It returns once the engine has been started (or an error occurred).
	ScheduleStream(stream *models.Streaming) error
	// GetEngine returns the engine registered under the given stream name.
	GetEngine(stream string) (Engine, bool)
}

type dataSource struct {
	db *models.Database

	// enginesSnapshot is an atomically updated, immutable map[string]Engine.
	// Produce and GetEngine read it without holding any lock.
	enginesSnapshot atomic.Value // stores map[string]Engine

	decoder   decode.Decoder
	running   atomic.Bool
	scheduleC chan scheduleCmd // commands processed by the event loop goroutine
	stopC     chan struct{}    // closed by Shutdown to stop the event loop

	logger logger.Logger
}

// NewDataSource creates a DataSource for the given database configuration.
func NewDataSource(db *models.Database) DataSource {
	d := &dataSource{
		db:        db,
		scheduleC: make(chan scheduleCmd),
		stopC:     make(chan struct{}),
		logger:    logger.GetLogger("Streaming", "DataSource"),
	}
	// Initialise the snapshot with an empty map so Produce never sees a nil load.
	d.enginesSnapshot.Store(make(map[string]Engine))
	return d
}

// Name returns the name of the underlying database.
func (d *dataSource) Name() string {
	return d.db.Name
}

// Startup initialises the decoder and launches the background event loop that
// serialises engine registration and shutdown operations.
func (d *dataSource) Startup() {
	if !d.running.CompareAndSwap(false, true) {
		return
	}
	d.decoder = decode.GetDecoder(d.db.Option.Engine)
	go d.loop()
}

// loop is the sole goroutine that mutates the engines map.
// All other goroutines interact with it through channels, keeping the hot
// path (Produce / GetEngine) completely lock-free.
func (d *dataSource) loop() {
	// Work on a local mutable copy; publish immutable snapshots via atomic.Value.
	engines := make(map[string]Engine)

	for {
		select {
		case cmd := <-d.scheduleC:
			// Register a new engine if it is not already present.
			if _, ok := engines[cmd.stream.Name]; !ok {
				// TODO: create engine based on streaming config (add engine type)
				engine := cep.NewEngine(cmd.stream, d.db)
				if err := engine.Startup(); err != nil {
					cmd.replyCh <- err
					continue
				}
				engines[cmd.stream.Name] = engine
				// Publish an immutable copy so readers see the new engine atomically.
				d.publishSnapshot(engines)
			}
			cmd.replyCh <- nil

		case <-d.stopC:
			// Shut down all engines before exiting.
			for _, engine := range engines {
				if err := engine.Shutdown(); err != nil {
					d.logger.Error("stop engine error:", logger.Error(err))
				}
			}
			return
		}
	}
}

// publishSnapshot stores an immutable copy of the current engines map so that
// Produce and GetEngine can read it without any lock.
func (d *dataSource) publishSnapshot(engines map[string]Engine) {
	snapshot := make(map[string]Engine, len(engines))
	maps.Copy(snapshot, engines)
	d.enginesSnapshot.Store(snapshot)
}

// Produce decodes msg into one or more Arrow RecordBatches and fans each record
// out to all currently registered engines. It is entirely lock-free: it reads
// the engines snapshot atomically and performs no writes.
func (d *dataSource) Produce(msg []byte) error {
	if !d.running.Load() {
		return errors.New("data source not running")
	}
	records, err := d.decoder.ToRecords(msg)
	if err != nil {
		d.logger.Error("transfer data to event error:", logger.Error(err))
		return err
	}
	// Load the current immutable snapshot — no lock required.
	engines := d.enginesSnapshot.Load().(map[string]Engine)
	for _, record := range records {
		if record == nil || record.NumRows() == 0 {
			continue
		}
		for _, engine := range engines {
			engine.Send(record)
		}
	}
	return nil
}

// ScheduleStream sends a registration command to the event loop and blocks
// until the engine has been started (or an error is returned).
func (d *dataSource) ScheduleStream(stream *models.Streaming) error {
	if !d.running.Load() {
		d.logger.Warn("schedule stream failed, data source not running", logger.String("streaming", stream.Name))
		return nil
	}
	replyCh := make(chan error, 1)
	d.scheduleC <- scheduleCmd{stream: stream, replyCh: replyCh}
	return <-replyCh
}

// GetEngine returns the engine for the given stream name from the current
// immutable snapshot. It is lock-free.
func (d *dataSource) GetEngine(stream string) (Engine, bool) {
	engines := d.enginesSnapshot.Load().(map[string]Engine)
	e, ok := engines[stream]
	return e, ok
}

// Shutdown stops the event loop and all registered engines.
func (d *dataSource) Shutdown() {
	if d.running.CompareAndSwap(true, false) {
		close(d.stopC)
	}
}
