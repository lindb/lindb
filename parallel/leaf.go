package parallel

import (
	"context"
	"encoding/json"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
)

// leafTask represents the leaf node's task, the leaf node is always storage node
// 1. receives the task request, and searches the data from time seres engine
// 2. sends the result to the parent node(root or intermediate)
type leafTask struct {
	currentNodeID     string
	storageService    service.StorageService
	executorFactory   ExecutorFactory
	taskServerFactory rpc.TaskServerFactory
}

// newLeafTask creates the leaf task
func newLeafTask(currentNode models.Node,
	storageService service.StorageService,
	executorFactory ExecutorFactory, taskServerFactory rpc.TaskServerFactory) TaskProcessor {
	return &leafTask{
		currentNodeID:     (&currentNode).Indicator(),
		storageService:    storageService,
		executorFactory:   executorFactory,
		taskServerFactory: taskServerFactory,
	}
}

// Process processes the task request, searches the metric's data from time series engine
func (p *leafTask) Process(ctx context.Context, req *pb.TaskRequest) error {
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

	payload := req.Payload
	query := stmt.Query{}
	if err := encoding.JSONUnmarshal(payload, &query); err != nil {
		return errUnmarshalQuery
	}

	stream := p.taskServerFactory.GetStream(curLeaf.Parent)
	if stream == nil {
		return errNoSendStream
	}

	exec := p.executorFactory.NewStorageExecutor(ctx, engine, curLeaf.ShardIDs, &query)
	groupedTimeSeries := exec.Execute()

	var err error
	var data []byte
	if groupedTimeSeries != nil {
		for result := range groupedTimeSeries {
			if err != nil {
				break
			}
			series := result.Series
			fields := make(map[string][]byte)
			for series.HasNext() {
				field := series.Next()
				data, err = field.Bytes()
				if err != nil {
					break
				}
				fields[field.FieldName()] = data
			}
			timeSeries := &pb.TimeSeries{
				Tags:   series.Tags(),
				Fields: fields,
			}
			data, _ := timeSeries.Marshal()
			_ = stream.Send(&pb.TaskResponse{
				JobID:     req.JobID,
				TaskID:    req.ParentTaskID,
				Completed: false,
				Payload:   data,
			})
		}
	}
	err = exec.Error()
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	_ = stream.Send(&pb.TaskResponse{
		JobID:     req.JobID,
		TaskID:    req.ParentTaskID,
		Completed: true,
		ErrMsg:    errMsg,
	})

	return nil
}
