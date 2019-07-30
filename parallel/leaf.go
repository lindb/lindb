package parallel

import (
	"encoding/json"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
)

// leafTask represents the leaf node's task, the leaf node is always storage node
// 1. receives the task request, and searches the data from time seres engine
// 2. sends the result to the parent node(root or intermediate)
type leafTask struct {
	currentNodeID  string
	storageService service.StorageService
	factory        ExecutorFactory
}

// NewLeafTask creates the leaf task
func NewLeafTask(currentNode models.Node,
	storageService service.StorageService,
	factory ExecutorFactory) TaskProcessor {
	return &leafTask{
		currentNodeID:  (&currentNode).Indicator(),
		storageService: storageService,
		factory:        factory,
	}
}

// Process processes the task request, searches the metric's data from time series engine
func (p *leafTask) Process(req *pb.TaskRequest) error {
	physicalPlan := models.PhysicalPlan{}
	if err := json.Unmarshal(req.PhysicalPlan, &physicalPlan); err != nil {
		return errUnmarshalPlan
	}

	foundTask := false
	var curLeaf models.Leaf
	for _, leaf := range physicalPlan.Leafs {
		if leaf.Indicator == p.currentNodeID {
			foundTask = true
			curLeaf = leaf
			break
		}
	}
	if !foundTask {
		return errWrongRequest
	}
	engine := p.storageService.GetEngine(physicalPlan.Database)
	if engine == nil {
		return errNoDatabase
	}
	//TODO impl query logic and send task result to parent node
	_ = p.factory.NewStorageExecutor(engine, curLeaf.ShardIDs, nil)
	return nil
}
