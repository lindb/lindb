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

package execution

import (
	"fmt"
	"sync"
	"time"

	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/sql/execution/model"
)

type RequestIDGenerator struct {
	node    string
	seq     int
	current int64
	next    int64

	lock sync.Mutex
}

func NewRequestIDGenerator(node string) *RequestIDGenerator {
	return &RequestIDGenerator{
		node: node,
		seq:  0,
	}
}

func (g *RequestIDGenerator) GenerateRequestID() model.RequestID {
	now := time.Now().UnixMilli()

	g.lock.Lock()
	defer g.lock.Unlock()

	if now >= g.next {
		// initialize next loop
		g.initialize(now)
	}

	g.seq++
	return model.RequestID(
		fmt.Sprintf("%s-%s-%08d",
			g.node,
			timeutil.FormatTimestamp(g.current, timeutil.DataTimeFormat4),
			g.seq,
		))
}

func (g *RequestIDGenerator) initialize(now int64) {
	g.current = now - now%timeutil.OneSecond
	g.next = g.current + timeutil.OneSecond
	g.seq = 0
}
