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

package stage

import (
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

// physicalPlanStage represents physical plan stage.
type physicalPlanStage struct {
	baseStage
	taskCtx context.TaskContext
}

func NewPhysicalPlanStage(taskCtx context.TaskContext) Stage {
	return &physicalPlanStage{
		baseStage: baseStage{
			stageType: PhysicalPlan,
		},
		taskCtx: taskCtx,
	}
}

// Plan returns sub execution tree for physical plan.
func (stage *physicalPlanStage) Plan() PlanNode {
	return NewPlanNode(operator.NewPhysicalPlan(stage.taskCtx))
}

// NextStages returns the next stages after physical plan stage completed.
func (stage *physicalPlanStage) NextStages() (rs []Stage) {
	requests := stage.taskCtx.GetRequests()
	for target, req := range requests {
		rs = append(rs, NewTaskSendStage(stage.taskCtx, target, req))
	}
	return
}

// Identifier returns identifier value of physical plan stage.
func (stage *physicalPlanStage) Identifier() string {
	return "Physical Plan"
}
