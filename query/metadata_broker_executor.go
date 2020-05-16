package query

import (
	"context"
	"sort"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
)

// metadataBrokerExecutor represents the executor which executes metric metadata suggest in broker side.
// 1. builds distribution physical execute plan
// 2. submits distribution query job
// 3. merges the result from distribution task execute result set
type metadataBrokerExecutor struct {
	database            string
	request             *stmt.Metadata
	replicaStateMachine replica.StatusStateMachine
	nodeStateMachine    broker.NodeStateMachine
	jobManager          parallel.JobManager

	ctx context.Context
}

// newMetadataBrokerExecutor creates a metadata suggest executor in broker side
func newMetadataBrokerExecutor(ctx context.Context, database string, request *stmt.Metadata,
	nodeStateMachine broker.NodeStateMachine, replicaStateMachine replica.StatusStateMachine,
	jobManager parallel.JobManager) parallel.MetadataExecutor {
	return &metadataBrokerExecutor{
		ctx:                 ctx,
		database:            database,
		request:             request,
		replicaStateMachine: replicaStateMachine,
		nodeStateMachine:    nodeStateMachine,
		jobManager:          jobManager,
	}
}

// Execute builds the execute plan, then submits the distribution query job
func (e *metadataBrokerExecutor) Execute() (result []string, err error) {
	physicalPlan, err := e.buildPhysicalPlan()
	if err != nil {
		return nil, err
	}

	resultCh := make(chan []string)
	// submit execute job
	result, err = e.submitJob(physicalPlan, resultCh)
	return
}

// submitJob submits the metadata suggest query job
func (e *metadataBrokerExecutor) submitJob(physicalPlan *models.PhysicalPlan, resultCh chan []string) (result []string, err error) {
	if err := e.jobManager.SubmitMetadataJob(e.ctx, physicalPlan, e.request, resultCh); err != nil {
		close(resultCh)
		return nil, err
	}
	resultMap := make(map[string]struct{})
	for rs := range resultCh {
		for _, value := range rs {
			resultMap[value] = struct{}{}
		}
	}
	result = []string{}
	for value := range resultMap {
		result = append(result, value)
	}
	sort.Strings(result)
	return
}

// buildPhysicalPlan builds distribution physical execute plan
func (e *metadataBrokerExecutor) buildPhysicalPlan() (*models.PhysicalPlan, error) {
	//FIXME need using storage's replica state ???
	storageNodes := e.replicaStateMachine.GetQueryableReplicas(e.database)
	storageNodesLen := len(storageNodes)
	if storageNodesLen == 0 {
		return nil, errNoAvailableStorageNode
	}
	curBroker := e.nodeStateMachine.GetCurrentNode()
	curBrokerIndicator := (&curBroker).Indicator()
	physicalPlan := &models.PhysicalPlan{
		Database: e.database,
		Root: models.Root{
			Indicator: curBrokerIndicator,
			NumOfTask: int32(storageNodesLen),
		},
	}
	receivers := []models.Node{curBroker}
	for storageNode, shardIDs := range storageNodes {
		physicalPlan.AddLeaf(models.Leaf{
			BaseNode: models.BaseNode{
				Parent:    curBrokerIndicator,
				Indicator: storageNode,
			},
			ShardIDs:  shardIDs,
			Receivers: receivers,
		})
	}
	return physicalPlan, nil
}
