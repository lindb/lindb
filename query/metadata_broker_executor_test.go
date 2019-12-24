package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataBrokerExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeStateMachine := broker.NewMockNodeStateMachine(ctrl)
	replicaStateMachine := replica.NewMockStatusStateMachine(ctrl)
	jobManager := parallel.NewMockJobManager(ctrl)

	exec := newMetadataBrokerExecutor(context.TODO(), "test_db", &stmt.Metadata{},
		nodeStateMachine, replicaStateMachine, jobManager)

	// no storage node
	replicaStateMachine.EXPECT().GetQueryableReplicas(gomock.Any()).Return(nil)
	rs, err := exec.Execute()
	assert.Error(t, err)
	assert.Nil(t, rs)

	// submit job err
	nodeStateMachine.EXPECT().GetCurrentNode().Return(models.Node{IP: "2.2.2.2", Port: 1234})
	replicaStateMachine.EXPECT().GetQueryableReplicas(gomock.Any()).Return(map[string][]int32{"1.1.1.1:1234": {1, 2, 3}})
	jobManager.EXPECT().SubmitMetadataJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	rs, err = exec.Execute()
	assert.Error(t, err)
	assert.Nil(t, rs)

	// execute query job
	jobManager.EXPECT().SubmitMetadataJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	e := exec.(*metadataBrokerExecutor)
	resultCh := make(chan []string)
	go func() {
		resultCh <- []string{"b", "d", "a"}
		close(resultCh)
	}()
	rs, err = e.submitJob(nil, resultCh)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "d"}, rs)
}
