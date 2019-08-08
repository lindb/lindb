package elect

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

type mockListener struct {
	onFailOverCount    int32
	onResignationCount int32
}

func newMockListener() *mockListener {
	return &mockListener{}
}

func (l *mockListener) OnResignation() {
	atomic.AddInt32(&l.onResignationCount, 1)
}

func (l *mockListener) OnFailOver() {
	atomic.AddInt32(&l.onFailOverCount, 1)
}

func TestElection_Initialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(nil)
	election := NewElection(context.TODO(), repo, node1, 1, listener1)
	election.Initialize()
	election.Close()

	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(eventCh)
	election = NewElection(context.TODO(), repo, node1, 1, listener1)
	election.Initialize()
	election.Close()
}

func TestElection_Elect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(nil)
	election := NewElection(context.TODO(), repo, node1, 1, listener1)
	election.Initialize()
	election.Elect()
	election.Close()
}

func TestElection_elect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(nil)
	election1 := NewElection(ctx, repo, node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)
	time.AfterFunc(700*time.Millisecond, func() {
		close(e.retryCh)
		cancel()
	})

	//fail
	repo.EXPECT().PutIfNotExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, nil, fmt.Errorf("err"))
	// success
	repo.EXPECT().PutIfNotExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	e.elect()

	election1.Close()
}

func TestElection_Handle_Event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(nil)
	election1 := NewElection(context.TODO(), repo, node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)

	eventCh := make(chan *state.Event)

	go func() {
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			KeyValues: []state.EventKeyValue{
				{Key: constants.MasterPath, Value: []byte{1, 1, 2}},
			},
		})
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeAll,
		})
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			Err:  fmt.Errorf("err"),
		})
		data := encoding.JSONMarshal(&models.Master{Node: node1})
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			KeyValues: []state.EventKeyValue{
				{Key: constants.MasterPath, Value: data},
			},
		})

		assert.Equal(t, int32(1), atomic.LoadInt32(&listener1.onFailOverCount))
		assert.True(t, e.IsMaster())

		// close chan
		close(eventCh)
	}()

	e.handleMasterChange(eventCh)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	election1.Close()
}
func TestElection_handle_event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(nil)
	election1 := NewElection(context.TODO(), repo, node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)
	data := encoding.JSONMarshal(&models.Master{Node: node1})
	e.handleEvent(&state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, e.IsMaster())

	time.AfterFunc(100*time.Millisecond, func() {
		<-e.retryCh
	})

	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	e.handleEvent(&state.Event{
		Type: state.EventTypeDelete,
	})
	assert.False(t, e.IsMaster())

	election1.Close()
}

func TestElection_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	listener1 := newMockListener()
	node1 := models.Node{IP: "127.0.0.1", Port: 2080}
	election1 := NewElection(context.TODO(), repo, node1, 1, listener1)
	election1.Close()
	e := election1.(*election)
	e.elect()
	election1.Close()
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(100 * time.Millisecond)
}
