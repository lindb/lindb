package parallel

import (
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

// taskReceiver represents receive the task result from the sub tasks
type taskReceiver struct {
	jobManager JobManager
}

// NewTaskReceiver creates the task receiver
func NewTaskReceiver(jobManager JobManager) rpc.TaskReceiver {
	return &taskReceiver{jobManager: jobManager}
}

// Receive receives the task result, merges them and finally returns the final result
func (r *taskReceiver) Receive(resp *pb.TaskResponse) error {
	taskID := resp.TaskID
	taskManager := r.jobManager.GetTaskManager()
	taskCtx := taskManager.Get(taskID)
	if taskCtx == nil {
		return nil
	}

	//TODO impl result handler
	taskCtx.ReceiveResult(resp)

	if taskCtx.Completed() {
		taskManager.Complete(taskID)

		if taskCtx.TaskType() == RootTask {
			jobCtx := r.jobManager.GetJob(resp.JobID)
			if jobCtx != nil && !jobCtx.Completed() {
				err := taskCtx.Error()
				if err != nil {
					jobCtx.Emit(&series.TimeSeriesEvent{Err: err})
				}
				jobCtx.Complete()
			}
		}
		//TODO need impl finally result build
	}
	return nil
}
