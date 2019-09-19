package parallel

import (
	"context"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

// intermediateTask represents the intermediate node's task,
// 1. only created for group by query
// 2. exchanges leaf task
// 3. receives leaf task's result
type intermediateTask struct {
	curNode     models.Node
	curNodeID   string
	taskManager TaskManager
}

// newIntermediateTask creates the intermediate task
func newIntermediateTask(curNode models.Node, taskManger TaskManager) *intermediateTask {
	return &intermediateTask{
		curNode:     curNode,
		curNodeID:   (&curNode).Indicator(),
		taskManager: taskManger,
	}
}

// Process processes the task request, sends task request to leaf nodes based on physical plan,
// and tracks the task state
func (p *intermediateTask) Process(ctx context.Context, req *pb.TaskRequest) error {
	physicalPlan := models.PhysicalPlan{}
	if err := encoding.JSONUnmarshal(req.PhysicalPlan, &physicalPlan); err != nil {
		return errUnmarshalPlan
	}
	taskSubmitted := false
	for _, intermediate := range physicalPlan.Intermediates {
		if intermediate.Indicator == p.curNodeID {
			taskID := p.taskManager.AllocTaskID()
			//TODO set task id
			taskCtx := newTaskContext(taskID, IntermediateTask, req.ParentTaskID, intermediate.Parent,
				intermediate.NumOfTask, newResultMerger(ctx, nil))
			p.taskManager.Submit(taskCtx)
			taskSubmitted = true
			break
		}
	}
	if !taskSubmitted {
		return errWrongRequest
	}

	if err := p.sendLeafTasks(physicalPlan, req); err != nil {
		return err
	}
	return nil
}

// sendLeafTasks sends the task request to the related leaf nodes, if failure return error
func (p *intermediateTask) sendLeafTasks(physicalPlan models.PhysicalPlan, req *pb.TaskRequest) error {
	for _, leaf := range physicalPlan.Leafs {
		if leaf.Parent == p.curNodeID {
			if err := p.taskManager.SendRequest(leaf.Indicator, req); err != nil {
				//TODO kill sent leaf task???
				return err
			}
		}
	}
	return nil
}

// Receive receives the sub task's result, and merges the results
func (p *intermediateTask) Receive(resp *pb.TaskResponse) error {
	taskID := resp.TaskID
	taskCtx := p.taskManager.Get(taskID)
	if taskCtx == nil {
		return nil
	}
	//TODO impl result handler
	taskCtx.ReceiveResult(resp)

	if taskCtx.Completed() {
		p.taskManager.Complete(taskID)
		// if task complete, need send task's result to parent node, if exist parent node
		if err := p.taskManager.SendResponse(taskCtx.ParentNode(), &pb.TaskResponse{TaskID: taskCtx.ParentTaskID()}); err != nil {
			return err
		}
	}
	return nil
}
