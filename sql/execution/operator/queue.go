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
	ch chan arrow.RecordBatch
}

func NewQueue(ch chan arrow.RecordBatch) *Queue {
	return &Queue{
		ch: ch,
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
	case page, ok := <-q.ch:
		return page, ok
	}
}

func (q *Queue) GetInbound() chan arrow.RecordBatch {
	return q.ch
}

func (q *Queue) Close() {
	close(q.ch)
}
