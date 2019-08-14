package parallel

import pb "github.com/lindb/lindb/rpc/proto/common"

// TaskDispatcher represents the task dispatcher
type TaskDispatcher interface {
	// Dispatch dispatches the task request based on task type
	Dispatch(req *pb.TaskRequest)
}

// TaskProcessor represents the task processor, all task processors are async
type TaskProcessor interface {
	// Process processes the task request
	Process(req *pb.TaskRequest) error
}
