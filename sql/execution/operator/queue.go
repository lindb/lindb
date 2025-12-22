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

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/spi/types"
)

var log = logger.GetLogger("operator", "execute")

type Queue struct {
	pageCh chan *types.Page
}

func NewQueue(pageCh chan *types.Page) *Queue {
	return &Queue{
		pageCh: pageCh,
	}
}

func (q *Queue) Produce(page *types.Page) {
	if page == nil {
		return
	}
	q.pageCh <- page
}

func (q *Queue) Consume(ctx context.Context) (*types.Page, bool) {
	select {
	case err := <-ctx.Done():
		panic(err)
	case page, ok := <-q.pageCh:
		if page != nil && page.Error != "" {
			log.Error("page has error", logger.Stack())
			panic(page.Error)
		}
		return page, ok
	}
}

func (q *Queue) GetInbound() chan *types.Page {
	return q.pageCh
}

func (q *Queue) Close() {
	close(q.pageCh)
}
