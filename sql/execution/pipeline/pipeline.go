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

package pipeline

import (
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
)

type Pipeline struct {
	taskCtx *context.TaskContext
	root    operator.Operator
}

func NewPipeline(taskCtx *context.TaskContext, root operator.Operator) *Pipeline {
	return &Pipeline{
		taskCtx: taskCtx,
		root:    root,
	}
}

func (p *Pipeline) Run(output chan<- *types.Page) {
	fmt.Printf("run pipeline, root=>\n%s\n", renderText(p.root))
	p.root.Run(output)
}
