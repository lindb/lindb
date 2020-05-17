package query

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/option"
)

func TestBrokerExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)

	nodeStateMachine := broker.NewMockNodeStateMachine(ctrl)
	dbStateMachine := database.NewMockDBStateMachine(ctrl)
	nodeStateMachine.EXPECT().GetCurrentNode().Return(currentNode.Node).AnyTimes()
	nodeStateMachine.EXPECT().GetActiveNodes().Return(nil)
	replicaStateMachine := replica.NewMockStatusStateMachine(ctrl)
	jobManager := parallel.NewMockJobManager(ctrl)

	// case 1: database not found
	exec := newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, dbStateMachine, jobManager)
	dbStateMachine.EXPECT().GetDatabaseCfg("test_db").Return(models.Database{}, false)
	exec.Execute()
	assert.NotNil(t, exec.ExecuteContext())

	// case 2: storage nodes not exist
	dbStateMachine.EXPECT().GetDatabaseCfg("test_db").
		Return(models.Database{Option: option.DatabaseOption{Interval: "10s"}}, true).AnyTimes()
	exec = newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, dbStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(nil)
	exec.Execute()
	assert.NotNil(t, exec.ExecuteContext())

	storageNodes := map[string][]int32{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 6, 9},
		"1.1.1.3:9000": {5, 7, 8},
		"1.1.1.4:9000": {10, 13, 15},
		"1.1.1.5:9000": {11, 12, 14},
	}
	brokerNodes := []models.ActiveNode{
		generateBrokerActiveNode("1.1.1.1", 8000),
		generateBrokerActiveNode("1.1.1.2", 8000),
		currentNode,
		generateBrokerActiveNode("1.1.1.4", 8000),
	}
	exec = newBrokerExecutor(context.TODO(), "test_db", "select f fro",
		replicaStateMachine, nodeStateMachine, dbStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	exec.Execute()

	exec = newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, dbStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	jobManager.EXPECT().SubmitJob(gomock.Any())
	exec.Execute()

	// submit job error
	exec = newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, dbStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	jobManager.EXPECT().SubmitJob(gomock.Any()).Return(errors.New("submit job error"))
	exec.Execute()
}
