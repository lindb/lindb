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

package ingest

import (
	"bytes"
	"context"
	"github.com/lindb/common/series"
	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/ingestion/flat"
	"github.com/lindb/lindb/pkg/strutil"
	"io"
	"sync"
	"time"
)

// Writer represents writer for writing time series data asynchronously.
type Writer interface {
	// AddPoint adds a time series point into buffer.
	AddPoint(ctx context.Context, point *Point)
	// Errors watches error in background goroutine.
	Errors() <-chan error
	// Close closes writer client, before close try to writer pending points.
	Close()
}

// writer implements Writer interface.
type writer struct {
	deps *depspkg.HTTPDeps

	namespace    string
	database     string
	writeOptions *WriteOptions

	bufferCh    chan *Point
	writeCh     chan []byte
	errCh       chan error
	stopBatchCh chan struct{}
	stopWriteCh chan struct{}
	doneCh      chan struct{}

	builder     *series.RowBuilder
	buf         *bytes.Buffer
	batchedSize int

	closed bool
	mutex  sync.Mutex
}

// NewWriter creates a writer
func NewWriter(deps *depspkg.HTTPDeps, writeOptions *WriteOptions) Writer {
	w := &writer{
		deps:         deps,
		writeOptions: writeOptions,
		bufferCh:     make(chan *Point, writeOptions.BatchSize()+1),
		writeCh:      make(chan []byte),
		errCh:        make(chan error),
		stopBatchCh:  make(chan struct{}),
		stopWriteCh:  make(chan struct{}),
		doneCh:       make(chan struct{}),
		builder:      series.CreateRowBuilder(),
		buf:          &bytes.Buffer{},
	}
	go w.bufferProc() // process point->data([]byte)
	go w.writeProc()  // writer data to broker
	return w
}

// AddPoint adds a time series point into buffer.
func (w *writer) AddPoint(ctx context.Context, point *Point) {
	if point == nil || !point.Valid() {
		return
	}
	select {
	case <-ctx.Done():
	case w.bufferCh <- point:
	}
}

// Errors watches error in background goroutine.
func (w *writer) Errors() <-chan error {
	return w.errCh
}

// Close closes writer, before close try to write pending points.
func (w *writer) Close() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// double check
	if w.closed {
		return
	}

	close(w.stopBatchCh)
	close(w.bufferCh)
	<-w.doneCh // wait buffer process completed

	close(w.stopWriteCh)
	close(w.writeCh)
	<-w.doneCh // wait writer process completed

	close(w.errCh)
	w.closed = true
}

// bufferProc consumes time series point from buffer chan, marshals point then put data into write buffer.
func (w *writer) bufferProc() {
	batchSize := w.writeOptions.BatchSize()
	ticker := time.NewTicker(time.Duration(w.writeOptions.flushInterval) * time.Millisecond)

	defer func() {
		ticker.Stop()
		w.doneCh <- struct{}{}
	}()

	for {
		select {
		case point := <-w.bufferCh:
			if err := w.batchPoint(point); err != nil {
				w.emitErr(err)
				continue
			}
			// check batch buffer is full
			if w.batchedSize >= batchSize {
				w.flushBuffer()
			}
		case <-ticker.C:
			w.flushBuffer()
		case <-w.stopBatchCh:
			// try to batch pending points
			for point := range w.bufferCh {
				if err := w.batchPoint(point); err != nil {
					w.emitErr(err)
				}
			}
			w.flushBuffer()
			return
		}
	}
}

// flushBuffer flushes buffer data, put data into write chan, then clear buffer.
func (w *writer) flushBuffer() {
	if w.batchedSize == 0 {
		return
	}
	data := w.buf.Bytes()
	w.buf.Reset() // reset batch buf
	w.batchedSize = 0

	// copy data
	dst := make([]byte, len(data))
	copy(dst, data)

	// put data into writer chan
	w.writeCh <- dst
}

// batchPoint marshals point, if success put data into buffer.
func (w *writer) batchPoint(point *Point) error {
	if point == nil {
		return nil
	}
	defer w.builder.Reset()

	builder := w.builder

	builder.AddNameSpace(strutil.String2ByteSlice(point.namespace))
	builder.AddMetricName(strutil.String2ByteSlice(point.MetricName()))
	builder.AddTimestamp(point.Timestamp().UnixMilli())

	addTag := func(tags map[string]string) error {
		for k, v := range tags {
			if err := builder.AddTag(strutil.String2ByteSlice(k), strutil.String2ByteSlice(v)); err != nil {
				return err
			}
		}
		return nil
	}
	// add default tags
	if err := addTag(w.writeOptions.DefaultTags()); err != nil {
		return err
	}
	// add tags of current point
	if err := addTag(point.Tags()); err != nil {
		return err
	}

	// writer field
	fields := point.Fields()
	for _, f := range fields {
		if err := f.write(builder); err != nil {
			return err
		}
	}

	// put point into buffer
	data, err := builder.Build()
	if err != nil {
		return err
	}
	_, err = w.buf.Write(data)
	if err != nil {
		return err
	}
	w.batchedSize++
	return nil
}

// writeProc consumes batched writer data, then writer it to broker.
func (w *writer) writeProc() {
	defer func() {
		// invoke when writer goroutine exit.
		w.doneCh <- struct{}{}
	}()
	// write writes data
	write := func(data []byte) {
		if len(data) == 0 {
			return
		}
		reader := bytes.NewReader(data)
		if err := w.write(reader); err != nil {
			w.emitErr(err)
		}
	}
	for {
		select {
		case data := <-w.writeCh:
			write(data)
		case <-w.stopWriteCh:
			// try to write pending messages
			for data := range w.writeCh {
				write(data)
			}
			return
		}
	}
}

// write writes data to broker.
func (w *writer) write(data io.Reader) error {
	namespace, database := w.deps.BrokerCfg.Prometheus.Namespace, w.deps.BrokerCfg.Prometheus.Database
	ctx, cancel := context.WithTimeout(
		context.Background(),
		w.deps.BrokerCfg.BrokerBase.Ingestion.IngestTimeout.Duration())
	defer cancel()

	limits := w.deps.StateMgr.GetDatabaseLimits(database)
	if limits.EnableNamespaceLengthCheck() && len(namespace) > limits.MaxNamespaceLength {
		return constants.ErrNamespaceTooLong
	}
	rows, err := flat.ParseReader(data, nil, namespace, limits)
	if err != nil {
		return err
	}
	if err := w.deps.CM.Write(ctx, database, rows); err != nil {
		return err
	}
	return nil
}

// emitErr emits error into chan.
func (w *writer) emitErr(err error) {
	select {
	case w.errCh <- err:
	default:
		// no err read, cannot put err into chan
	}
}
