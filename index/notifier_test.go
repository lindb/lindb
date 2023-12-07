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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotifier_Pool(t *testing.T) {
	mn := GetMetaNotifier()
	assert.NotNil(t, mn)
	PutMetaNotifier(mn)
	tn := GetTagNotifier()
	assert.NotNil(t, tn)
	PutTagNotifier(tn)
	fn := GetFieldNotifier()
	assert.NotNil(t, fn)
	PutFieldNotifier(fn)
}

func TestNotifyWorker_Notify(t *testing.T) {
	c := 0
	ch := make(chan struct{})
	w := NewWorker(context.TODO(), "test", time.Second, func(_ Notifier) {
		c++
		if c == 256 {
			ch <- struct{}{}
		}
	})
	w.Notify(nil)
	w.Notify(GetMetaNotifier())
	w.lock.Lock()
	assert.Len(t, w.buf, 1)
	w.lock.Unlock()
	for i := 0; i < 255; i++ {
		w.Notify(GetMetaNotifier())
	}
	<-ch
	w.lock.Lock()
	assert.Empty(t, w.buf)
	w.lock.Unlock()
	w.Shutdown()
}

func TestNotifyWorker_Notify_Timeout(t *testing.T) {
	ch := make(chan struct{})
	w := NewWorker(context.TODO(), "test", 10*time.Millisecond, func(_ Notifier) {
		ch <- struct{}{}
	})
	w.Notify(GetMetaNotifier())
	<-ch
	w.lock.Lock()
	assert.Empty(t, w.buf)
	w.lock.Unlock()
	w.Shutdown()
}

func TestNotifyWorker_Notify_Pending(t *testing.T) {
	ch := make(chan struct{})
	ctx, cancel := context.WithCancel(context.TODO())
	w := NewWorker(ctx, "test", time.Second, func(_ Notifier) {
		ch <- struct{}{}
	})
	w.Notify(GetMetaNotifier())
	cancel()
	<-ch
	w.lock.Lock()
	assert.Empty(t, w.buf)
	w.lock.Unlock()
}

func TestNotifyWorker_Notify_NonBlocking(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	w := NewWorker(ctx, "test", time.Second, func(_ Notifier) {})
	go func() { cancel() }()
	w.signal <- struct{}{}

	// non blocking
	for i := 0; i < 600; i++ {
		w.Notify(GetMetaNotifier())
	}
}

func TestNotifyWorker_Handle_Panic(t *testing.T) {
	ch := make(chan struct{})
	w := NewWorker(context.TODO(), "test", 10*time.Millisecond, func(_ Notifier) {
		ch <- struct{}{}
		panic("err")
	})
	w.Notify(GetMetaNotifier())
	<-ch
}
