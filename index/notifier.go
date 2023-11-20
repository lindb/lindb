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

package index

import (
	"context"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

var (
	metaNotifierPool = sync.Pool{
		New: func() any {
			return &MetaNotifier{}
		},
	}
	fieldNotifierPool = sync.Pool{
		New: func() any {
			return &FieldNotifier{}
		},
	}
	tagNotifierPool = sync.Pool{
		New: func() any {
			return &TagNotifier{}
		},
	}
	log = logger.GetLogger("Index", "Notify")
)

// GetMetaNotifier gets meta notifier from pool
func GetMetaNotifier() *MetaNotifier {
	notifier := metaNotifierPool.Get()
	return notifier.(*MetaNotifier)
}

// PutMetaNotifier puts back meta notifier to pool.
func PutMetaNotifier(notifier *MetaNotifier) {
	metaNotifierPool.Put(notifier)
}

// GetTagNotifier gets tag notifier from pool.
func GetTagNotifier() *TagNotifier {
	notifier := tagNotifierPool.Get()
	return notifier.(*TagNotifier)
}

// PutTagNotifier puts back tag notifier to pool.
func PutTagNotifier(notifier *TagNotifier) {
	tagNotifierPool.Put(notifier)
}

// GetFieldNotifier gets field notifier from pool.
func GetFieldNotifier() *FieldNotifier {
	notifier := fieldNotifierPool.Get()
	return notifier.(*FieldNotifier)
}

// PutFieldNotifier puts back field notifier to pool.
func PutFieldNotifier(notifier *FieldNotifier) {
	fieldNotifierPool.Put(notifier)
}

// MetaNotifier represents metric meta notifier.
type MetaNotifier struct {
	Namespace  string
	MetricName string
	MetricID   metric.ID
	TagHash    uint64
	Tags       tag.Tags
	Callback   func(id uint32, err error)
}

// TagNotifier represents tag meta notifier.
type TagNotifier struct {
	metricID   metric.ID
	tags       tag.Tags
	buildIndex func(tagKeyID, tagValueID uint32)
}

// FieldNotifier represents field meta notifier.
type FieldNotifier struct {
	Namespace  string
	MetricName string
	Field      field.Meta
	Callback   func(fieldID field.ID, err error)
}

// FlushNotifier represents flush event notifier.
type FlushNotifier struct {
	Callback func(err error)
}

// NotifyWorker represents notifier event hanle worker.
type NotifyWorker struct {
	name      string
	timeout   time.Duration
	signal    chan struct{}
	buf       []Notifier
	immutable []Notifier
	lock      sync.Mutex
	handle    func(n Notifier)
}

// NewWorker creates notifier event handle worker.
func NewWorker(ctx context.Context, name string, timeout time.Duration, handle func(n Notifier)) *NotifyWorker {
	w := &NotifyWorker{
		timeout: timeout,
		name:    name,
		signal:  make(chan struct{}, 1),
		handle:  handle,
	}
	go func() {
		pprof.Do(ctx,
			pprof.Labels("type", "NotifyWorker", "name", name),
			func(ctx context.Context) {
				w.run(ctx)
			})
	}()
	return w
}

// Notify puts a notifier event into buffer, then try notify worker to process event.
func (w *NotifyWorker) Notify(n Notifier) {
	if n == nil {
		return
	}

	size := 0
	// add event to buffer
	w.lock.Lock()
	w.buf = append(w.buf, n)
	size = len(w.buf)
	w.lock.Unlock()

	// TODO: add config?
	if size >= 256 {
		select {
		// try notify buf not empty
		case w.signal <- struct{}{}:
		default:
		}
	}
}

// Shutdown shutdowns worker.
func (w *NotifyWorker) Shutdown() {
	close(w.signal)
}

func (w *NotifyWorker) run(ctx context.Context) {
	timer := time.NewTicker(w.timeout)
	defer func() {
		timer.Stop()
		w.lock.Lock()
		defer w.lock.Unlock()
		// do pending notifier
		for idx := range w.buf {
			w.handle(w.buf[idx])
		}
		w.buf = w.buf[:0]
	}()
	for {
		select {
		case <-w.signal:
			w.processEventBuf()
		case <-timer.C:
			w.processEventBuf()
		case <-ctx.Done():
			log.Warn("notify worker exist", logger.String("name", w.name))
			return
		}
	}
}

func (w *NotifyWorker) processEventBuf() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("notify worker handle event panic",
				logger.String("name", w.name), logger.Any("error", err))
		}
	}()
	w.lock.Lock()
	w.immutable = append(w.immutable, w.buf...)
	w.buf = w.buf[:0]
	w.lock.Unlock()

	if len(w.immutable) == 0 {
		return
	}

	for idx := range w.immutable {
		w.handle(w.immutable[idx])
	}
	w.immutable = w.immutable[:0]
}
