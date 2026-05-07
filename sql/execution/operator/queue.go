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

package operator

import (
	"context"
	"errors"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/logger"
)

var log = logger.GetLogger("operator", "execute")

type Queue struct {
	ch    chan arrow.RecordBatch
	errCh chan string // carries at most one error message from Fail()
}

func NewQueue(ch chan arrow.RecordBatch) *Queue {
	return &Queue{
		ch:    ch,
		errCh: make(chan string, 1),
	}
}

func (q *Queue) Produce(record arrow.RecordBatch) {
	if record == nil {
		return
	}
	q.ch <- record
}

func (q *Queue) Consume(ctx context.Context) (arrow.RecordBatch, bool) {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == nil || errors.Is(err, context.Canceled) {
			log.Info("queue consume canceled")
			return nil, false
		}
		panic(err)
	case record, ok := <-q.ch:
		return record, ok
	}
}

func (q *Queue) GetInbound() chan arrow.RecordBatch {
	return q.ch
}

func (q *Queue) Close() {
	close(q.ch)
}

// Fail sends an error message into errCh and closes ch so that any blocked
// Consume() call unblocks. The caller (RemoteExchangeOperator.Run) must check
// ErrMsg() after Consume returns false to decide whether to panic.
func (q *Queue) Fail(errMsg string) {
	// Non-blocking send: errCh is buffered with capacity 1, so this never blocks
	// even if called concurrently.
	select {
	case q.errCh <- errMsg:
	default:
	}
	close(q.ch)
}

// ErrMsg returns the error message set by Fail(), or "" if the queue was
// closed normally via Close(). Safe to call only after Consume() returns false.
func (q *Queue) ErrMsg() string {
	select {
	case msg := <-q.errCh:
		return msg
	default:
		return ""
	}
}
