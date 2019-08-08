package discovery

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/state"
)

type mockListener struct {
	nodes map[string][]byte
	mutex sync.Mutex

	invokes int
}

func newMockListener() *mockListener {
	return &mockListener{
		nodes: make(map[string][]byte),
	}
}

func (m *mockListener) OnCreate(key string, value []byte) {
	m.mutex.Lock()
	m.nodes[key] = value
	m.invokes++
	m.mutex.Unlock()
}

func (m *mockListener) OnDelete(key string) {
	m.mutex.Lock()
	delete(m.nodes, key)
	m.invokes++
	m.mutex.Unlock()
}

func (m *mockListener) Cleanup() {
	m.mutex.Lock()
	m.nodes = make(map[string][]byte)
	m.invokes++
	m.mutex.Unlock()
}

var testDiscoveryPath = "/test/discovery1"

func TestDiscovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	factory := NewFactory(repo)

	d := factory.CreateDiscovery("", newMockListener())
	err := d.Discovery()
	assert.NotNil(t, err)
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	d.Close()

	listener := newMockListener()
	d = factory.CreateDiscovery(testDiscoveryPath, listener)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any()).Return(nil)
	err = d.Discovery()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	repo.EXPECT().Close().Return(nil)
	d.Close()

	eventCh := make(chan *state.Event)
	listener = newMockListener()
	d = factory.CreateDiscovery(testDiscoveryPath, listener)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any()).Return(eventCh)
	err = d.Discovery()
	if err != nil {
		t.Fatal(err)
	}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: "/test/discovery1/key1", Value: []byte{1, 1, 2}},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: "/test/discovery1/key2", Value: []byte{1, 1, 2}},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: "/test/discovery1/key3", Value: []byte{1, 1, 2}},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
		KeyValues: []state.EventKeyValue{
			{Key: "/test/discovery1/key3", Value: []byte{1, 1, 2}},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeAll,
		KeyValues: []state.EventKeyValue{
			{Key: "/test/discovery1/key4", Value: []byte{1, 1, 2}},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeAll,
		Err:  fmt.Errorf("err"),
	})

	// wait goroutine
	time.Sleep(400 * time.Millisecond)

	listener.mutex.Lock()
	nodes := listener.nodes
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, 6, listener.invokes)
	assert.Equal(t, []byte{1, 1, 2}, nodes["/test/discovery1/key4"])
	listener.mutex.Unlock()

	repo.EXPECT().Close().Return(nil)
	d.Close()
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
