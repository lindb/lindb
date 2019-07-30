package parallel

import (
	"encoding/json"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

//go:generate mockgen -source=./job_manager.go -destination=./job_manager_mock.go -package=parallel

// JobManager represents the job manager for the root broker node
type JobManager interface {
	// SubmitJob submits the distribution query job based on physical plan
	SubmitJob(plan *models.PhysicalPlan) error
}

// jobManager implements the job manager for managing the query job
type jobManager struct {
	taskManager TaskManager
}

// NewJobManager creates the job manager
func NewJobManager(taskManger TaskManager) JobManager {
	return &jobManager{
		taskManager: taskManger,
	}
}

// SubmitJob submits the distribution query job based on physical plan,
// 1. if has intermediate nodes, sends the request to the intermediate nodes
// 2. else sends the request to the leaf node directly
func (j *jobManager) SubmitJob(plan *models.PhysicalPlan) error {
	planPayload, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	taskID := j.taskManager.AllocTaskID()

	// TODO need add param
	req := &pb.TaskRequest{
		ParentTaskID: taskID,
		PhysicalPlan: planPayload,
	}

	taskCtx := newTaskContext(taskID, "", "", plan.Root.NumOfTask)
	j.taskManager.Submit(taskCtx)

	if len(plan.Intermediates) > 0 {
		for _, intermediate := range plan.Intermediates {
			if err := j.sendTaskReq(intermediate.Indicator, req); err != nil {
				//TODO kill sent leaf task???
				return err
			}
		}
	} else if len(plan.Leafs) > 0 {
		for _, leaf := range plan.Leafs {
			if err := j.sendTaskReq(leaf.Indicator, req); err != nil {
				//TODO kill sent leaf task???
				return err
			}
		}
	}
	return err
}

// sendTaskReq sends the task requests to target node based on node's indicator
func (j *jobManager) sendTaskReq(indicator string, req *pb.TaskRequest) error {
	taskSender := j.taskManager.GetTaskSenderManager()
	if taskSender == nil {
		return errNoTaskSender
	}
	sendSendStream := taskSender.GetClientStream(indicator)
	if sendSendStream == nil {
		return errNoSendStream
	}
	if err := sendSendStream.Send(req); err != nil {
		//TODO kill sent leaf task???
		return errTaskSend
	}
	return nil
}
