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

package input

import (
	"fmt"

	"github.com/lindb/lindb/streaming/query/context"
	"github.com/lindb/lindb/streaming/query/processor"
)

type SingleStreamReceiver struct {
	ctx    *context.QueryContext
	stream string

	inbound *processor.Queue
}

func NewSingleStreamReceiver(ctx *context.QueryContext, stream string) *SingleStreamReceiver {
	return &SingleStreamReceiver{
		ctx:     ctx,
		stream:  stream,
		inbound: processor.NewQueue(make(chan any)),
	}
}

func (r *SingleStreamReceiver) Receive(event any) {
	fmt.Printf("single stream receiver, current:%v, receive event:%v\n", r.ctx.Query, event)
	r.inbound.Produce(event)
}

func (r *SingleStreamReceiver) Run(output chan<- any) {
	for {
		source, ok := r.inbound.Consume(r.ctx.Context)
		if !ok {
			break
		}
		output <- source
	}
}

func (r *SingleStreamReceiver) GetInbounds() []chan any {
	return []chan any{r.inbound.GetInbound()}
}

func (r *SingleStreamReceiver) Children() []processor.Processor {
	return nil
}

func (r *SingleStreamReceiver) String() string {
	return fmt.Sprintf("SingleStreamReceiver[stream=%s]", r.stream)
}
