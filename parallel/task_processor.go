package parallel

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./task_processor.go -destination=./task_processor_mock.go -package=parallel

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

// leafTaskDispatcher represents leaf task dispatcher for storage
type leafTaskDispatcher struct {
	processor TaskProcessor
}

// NewLeafTaskDispatcher creates a leaf task dispatcher
func NewLeafTaskDispatcher(currentNode models.Node,
	storageService service.StorageService,
	executorFactory ExecutorFactory, taskServerFactory rpc.TaskServerFactory) TaskDispatcher {
	return &leafTaskDispatcher{
		processor: newLeafTask(currentNode, storageService, executorFactory, taskServerFactory),
	}
}

// Dispatch dispatches the request to storage engine query processor
func (d *leafTaskDispatcher) Dispatch(req *pb.TaskRequest) {
	//FIXME(stone1100) need remove err
	go func() {
		_ = d.processor.Process(req)
	}()
}

// intermediateTaskDispatcher represents intermediate task dispatcher for broker
type intermediateTaskDispatcher struct {
}

// NewIntermediateTaskDispatcher create an intermediate task dispatcher
func NewIntermediateTaskDispatcher() TaskDispatcher {
	return &intermediateTaskDispatcher{}
}

// Dispatch dispatches the request to distribution query processor, merges the results
func (d *intermediateTaskDispatcher) Dispatch(req *pb.TaskRequest) {

}
