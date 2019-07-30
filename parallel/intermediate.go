package parallel

import (
	"encoding/json"

	"github.com/lindb/lindb/models"
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
func (p *intermediateTask) Process(req *pb.TaskRequest) error {
	physicalPlan := models.PhysicalPlan{}
	if err := json.Unmarshal(req.PhysicalPlan, &physicalPlan); err != nil {
		return errUnmarshalPlan
	}
	taskSubmitted := false
	for _, intermediate := range physicalPlan.Intermediates {
		if intermediate.Indicator == p.curNodeID {
			taskID := p.taskManager.AllocTaskID()
			//TODO set task id
			taskCtx := newTaskContext(taskID, req.ParentTaskID, intermediate.Parent, intermediate.NumOfTask)
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
	taskSender := p.taskManager.GetTaskSenderManager()
	if taskSender == nil {
		return errNoTaskSender
	}
	for _, leaf := range physicalPlan.Leafs {
		if leaf.Parent == p.curNodeID {

			sendSendStream := taskSender.GetClientStream(leaf.Indicator)
			if sendSendStream == nil {
				return errNoSendStream
			}
			if err := sendSendStream.Send(req); err != nil {
				//TODO kill sent leaf task???
				return errTaskSend
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
	taskCtx.ReceiveResult()

	if taskCtx.Completed() {
		p.taskManager.Complete(taskID)
		// if task complete, need send task's result to parent node, if exist parent node
		if err := p.sendTaskResult(taskCtx); err != nil {
			return err
		}
	}
	return nil
}

// sendTaskResult sends the task result to the parent node if exist
func (p *intermediateTask) sendTaskResult(taskCtx TaskContext) error {
	serverStream := p.taskManager.GetTaskSenderManager().GetServerStream(taskCtx.ParentNode())
	if serverStream != nil {
		// todo need add task result
		if err := serverStream.Send(&pb.TaskResponse{TaskID: taskCtx.ParentTaskID()}); err != nil {
			return err
		}
	}
	//TODO add warn log???
	return nil
}
