package parallel

import (
	"context"
	"encoding/json"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
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
func newLeafTask(
	currentNode models.Node,
	storageService service.StorageService,
	executorFactory ExecutorFactory,
	taskServerFactory rpc.TaskServerFactory,
) TaskProcessor {
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
	db, ok := p.storageService.GetDatabase(physicalPlan.Database)
	if !ok {
		return errNoDatabase
	}
	stream := p.taskServerFactory.GetStream(curLeaf.Parent)
	if stream == nil {
		return errNoSendStream
	}

	switch req.RequestType {
	case pb.RequestType_Data:
		if err := p.processDataSearch(ctx, db, physicalPlan.Namespace, curLeaf.ShardIDs, req, stream); err != nil {
			return err
		}
	case pb.RequestType_Metadata:
		if err := p.processMetadataSuggest(db, physicalPlan.Namespace, curLeaf.ShardIDs, req, stream); err != nil {
			return err
		}
	}
	return nil
}

func (p *leafTask) processMetadataSuggest(db tsdb.Database, namespace string, shardIDs []int32,
	req *pb.TaskRequest, stream pb.TaskService_HandleServer,
) error {
	payload := req.Payload
	query := &stmt.Metadata{}
	if err := encoding.JSONUnmarshal(payload, query); err != nil {
		return errUnmarshalSuggest
	}
	exec := p.executorFactory.NewMetadataStorageExecutor(db, namespace, shardIDs, query)
	result, err := exec.Execute()
	if err != nil && err != constants.ErrNotFound {
		return err
	}
	// send result to upstream
	if err := stream.Send(&pb.TaskResponse{
		JobID:     req.JobID,
		TaskID:    req.ParentTaskID,
		Completed: true,
		Payload:   encoding.JSONMarshal(&models.SuggestResult{Values: result}),
	}); err != nil {
		return err
	}
	return nil
}

func (p *leafTask) processDataSearch(ctx context.Context, db tsdb.Database, namespace string, shardIDs []int32,
	req *pb.TaskRequest, stream pb.TaskService_HandleServer,
) error {
	payload := req.Payload
	query := stmt.Query{}
	if err := encoding.JSONUnmarshal(payload, &query); err != nil {
		return errUnmarshalQuery
	}

	option := db.GetOption()
	var interval timeutil.Interval
	_ = interval.ValueOf(option.Interval)
	//TODO need get storage interval by query time if has rollup config
	timeRange, intervalRatio, queryInterval := downSamplingTimeRange(query.Interval, interval, query.TimeRange)
	// execute leaf task
	storageExecuteCtx := p.executorFactory.NewStorageExecuteContext(namespace, shardIDs, &query)
	queryFlow := NewStorageQueryFlow(ctx, storageExecuteCtx, &query, req, stream, db.ExecutorPool(), timeRange, queryInterval, intervalRatio)
	exec := p.executorFactory.NewStorageExecutor(queryFlow, db, storageExecuteCtx)
	exec.Execute()
	return nil
}
