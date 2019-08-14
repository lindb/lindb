package coordinator

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestStateMachineFactory_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	discoveryFactory.EXPECT().GetRepo().Return(repo).AnyTimes()

	factory := NewStateMachineFactory(&StateMachineCfg{
		Ctx:              context.TODO(),
		CurrentNode:      models.Node{IP: "1.1.1.1", Port: 9000},
		DiscoveryFactory: discoveryFactory,
	})

	// test node state machine
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	nodeSM, err := factory.CreateNodeStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, nodeSM)

	discovery1.EXPECT().Discovery().Return(nil)
	nodeSM, err = factory.CreateNodeStateMachine()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, nodeSM)

	// test storage state machine
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	storageStateSM, err := factory.CreateStorageStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, storageStateSM)
	discovery1.EXPECT().Discovery().Return(nil)
	storageStateSM, err = factory.CreateStorageStateMachine()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, storageStateSM)

	// test replica status state machine
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	replicaStatusSM, err := factory.CreateReplicaStatusStateMachine()
	assert.NotNil(t, err)
	assert.Nil(t, replicaStatusSM)
	discovery1.EXPECT().Discovery().Return(nil)
	replicaStatusSM, err = factory.CreateReplicaStatusStateMachine()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, replicaStatusSM)
}
