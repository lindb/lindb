package coordinator

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

func TestMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Discovery().Return(nil).AnyTimes()
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	nodeSM := broker.NewMockNodeStateMachine(ctrl)
	nodeSM.EXPECT().StartMonitoring().AnyTimes()
	nodeSM.EXPECT().StopMonitoring().AnyTimes()

	node1 := models.Node{IP: "1.1.1.1", Port: 8000}
	master1 := NewMaster(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
		BrokerSM: &BrokerStateMachines{
			NodeSM: nodeSM,
		},
	})
	master1.Start()
	data := encoding.JSONMarshal(&models.Master{Node: node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.Equal(t, node1, master1.GetMaster().Node)
	assert.True(t, master1.IsMaster())

	// re-elect
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())

	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, master1.IsMaster())

	master1.Stop()
	assert.False(t, master1.IsMaster())
}

func TestMaster_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	nodeSM := broker.NewMockNodeStateMachine(ctrl)
	nodeSM.EXPECT().StartMonitoring().AnyTimes()
	nodeSM.EXPECT().StopMonitoring().AnyTimes()

	node1 := models.Node{IP: "1.1.1.1", Port: 8000}
	master1 := NewMaster(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
		BrokerSM: &BrokerStateMachines{
			NodeSM: nodeSM,
		},
	})
	master1.Start()

	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	data := encoding.JSONMarshal(&models.Master{Node: node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())
	assert.Nil(t, master1.GetMaster())

	discovery1.EXPECT().Discovery().Return(nil)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, master1.IsMaster())
	assert.Nil(t, master1.GetMaster())

	master1.Stop()
}

func TestMaster_FlushDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh).AnyTimes()
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	discoveryFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	discovery1.EXPECT().Discovery().Return(nil).AnyTimes()
	discovery1.EXPECT().Close().AnyTimes()
	discoveryFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	nodeSM := broker.NewMockNodeStateMachine(ctrl)
	nodeSM.EXPECT().StartMonitoring().AnyTimes()
	nodeSM.EXPECT().StopMonitoring().AnyTimes()

	node1 := models.Node{IP: "1.1.1.1", Port: 8000}
	master1 := NewMaster(&MasterCfg{
		Ctx:              context.TODO(),
		Repo:             repo,
		Node:             node1,
		TTL:              1,
		DiscoveryFactory: discoveryFactory,
		BrokerSM: &BrokerStateMachines{
			NodeSM: nodeSM,
		},
	})
	err := master1.FlushDatabase("test", "test")
	assert.NoError(t, err)

	master1.Start()
	data := encoding.JSONMarshal(&models.Master{Node: node1})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, master1.IsMaster())
	err = master1.FlushDatabase("test", "test")
	assert.Error(t, err)

	m1 := master1.(*master)
	m1.mutex.Lock()
	clusterSM := storage.NewMockClusterStateMachine(ctrl)
	m1.masterCtx.StateMachine.StorageCluster = clusterSM
	m1.mutex.Unlock()

	cluster1 := storage.NewMockCluster(ctrl)
	clusterSM.EXPECT().GetCluster(gomock.Any()).Return(cluster1)
	cluster1.EXPECT().FlushDatabase(gomock.Any()).Return(nil)
	err = master1.FlushDatabase("test", "test")
	assert.NoError(t, err)
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
