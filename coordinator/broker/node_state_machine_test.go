package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
)

var currentNode = models.Node{IP: "1.1.1.2", Port: 2080}

func TestNodeStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err := NewNodeStateMachine(context.TODO(), currentNode, factory)
	assert.NotNil(t, err)

	// normal case
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewNodeStateMachine(context.TODO(), currentNode, factory)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
	assert.Equal(t, currentNode, stateMachine.GetCurrentNode())
}

func TestNodeStateMachine_Listener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	// normal case
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewNodeStateMachine(context.TODO(), currentNode, factory)
	if err != nil {
		t.Fatal(err)
	}
	activeNode := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}}
	data, _ := json.Marshal(&activeNode)
	stateMachine.OnCreate("/data/test", data)

	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))
	assert.Equal(t, activeNode, stateMachine.GetActiveNodes()[0])
	assert.Equal(t, currentNode, stateMachine.GetCurrentNode())

	stateMachine.OnCreate("/data/test2", []byte{1, 1})
	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))

	stateMachine.OnDelete("/data/test")
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))

	// add
	stateMachine.OnCreate("/data/test", data)
	assert.Equal(t, 1, len(stateMachine.GetActiveNodes()))

	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
}
