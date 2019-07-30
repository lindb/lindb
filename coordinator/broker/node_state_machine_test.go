package broker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/timeutil"
)

func TestNodeStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)
	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().WatchPrefix(gomock.Any(), constants.ActiveNodesPath).Return(eventCh)
	repo.EXPECT().Close().Return(nil)

	stateMachine, err := NewNodeStateMachine(context.TODO(), repo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))

	// wrong event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.ActiveNodesPath + "/1.1.1.1:2080", Value: nil},
		},
	})
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
	node := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 2080}, OnlineTime: timeutil.Now()}
	data, _ := json.Marshal(&node)
	// modify event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.ActiveNodesPath + "/1.1.1.1:2080", Value: data},
		},
	})
	assert.Equal(t, []models.ActiveNode{node}, stateMachine.GetActiveNodes())
	// delete event
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
		KeyValues: []state.EventKeyValue{
			{Key: constants.ActiveNodesPath + "/1.1.1.1:2080", Value: nil},
		},
	})
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))

	_ = stateMachine.Close()
	assert.Equal(t, 0, len(stateMachine.GetActiveNodes()))
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
