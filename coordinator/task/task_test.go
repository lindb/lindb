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

package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const kindDummy Kind = "you-guess"

type dummyParams struct{}

func (p dummyParams) Bytes() []byte { return []byte("{}") }

type dummyProcessor struct{ callcnt int32 }

func (p *dummyProcessor) Kind() Kind                  { return kindDummy }
func (p *dummyProcessor) RetryCount() int             { return 0 }
func (p *dummyProcessor) RetryBackOff() time.Duration { return 0 }
func (p *dummyProcessor) Concurrency() int            { return 1 }
func (p *dummyProcessor) Process(ctx context.Context, task Task) error {
	atomic.AddInt32(&p.callcnt, 1)
	return nil
}
func (p *dummyProcessor) CallCount() int { return int(atomic.LoadInt32(&p.callcnt)) }

func TestTask(t *testing.T) {
	assert.Equal(t, "StateCreated", StateCreated.String())
	assert.Equal(t, "StateRunning", StateRunning.String())
	assert.Equal(t, "StateDoneErr", StateDoneErr.String())
	assert.Equal(t, "StateDoneOK", StateDoneOK.String())
}
