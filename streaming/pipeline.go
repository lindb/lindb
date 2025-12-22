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
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/query/processor"
)

type Pipeline struct {
	root processor.Processor
}

func NewPipeline(root processor.Processor) *Pipeline {
	return &Pipeline{
		root: root,
	}
}

func (p *Pipeline) Run(output chan<- any) {
	fmt.Printf("run pipeline, root=>\n%s\n", renderText(p.root))

	p.execOperator(p.root, true, output)

	p.root.Run(output)
}

func (p *Pipeline) execOperator(
	op processor.Processor,
	exclude bool,
	output chan<- any,
) {
	fmt.Printf("run operator=%T\n", op)
	children := op.Children()
	inbounds := op.GetInbounds()
	for i, child := range children {
		p.execOperator(child, false, inbounds[i])
	}

	if !exclude {
		RunAsync(op, output)
	}
}

func RunAsync(op processor.Processor, output chan<- any) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				output <- &types.Page{Error: fmt.Sprintf("%v", err)}
			}
			close(output)
		}()

		op.Run(output)
	}()
}
