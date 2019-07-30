package parallel

import (
	pb "github.com/lindb/lindb/rpc/proto/common"
)

// rootTask represents the root task, only receives the task result from the sub tasks
type rootTask struct {
	taskManager TaskManager
}

// newRookTask creates the root task
func newRookTask(taskManager TaskManager) TaskReceiver {
	return &rootTask{taskManager: taskManager}
}

// Receive receives the task result, merges them and finally returns the final result
func (r *rootTask) Receive(resp *pb.TaskResponse) error {
	taskID := resp.TaskID
	taskCtx := r.taskManager.Get(taskID)
	if taskCtx == nil {
		return nil
	}
	//TODO impl result handler
	taskCtx.ReceiveResult()

	if taskCtx.Completed() {
		r.taskManager.Complete(taskID)

		//TODO need impl finally result build
	}
	return nil
}
