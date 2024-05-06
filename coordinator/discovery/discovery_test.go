// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package discovery

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lindb/common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

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
	assert.Equal(t, repo, factory.GetRepo())

	d := factory.CreateDiscovery("", newMockListener())
	err := d.Discovery(false)
	assert.NotNil(t, err)
	d.Close()

	listener := newMockListener()
	d = factory.CreateDiscovery(testDiscoveryPath, listener)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), false).Return(nil)
	err = d.Discovery(false)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	d.Close()

	eventCh := make(chan *state.Event)
	listener = newMockListener()
	d = factory.CreateDiscovery(testDiscoveryPath, listener)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), false).Return(eventCh)
	err = d.Discovery(false)
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
			{Key: "/test/discovery1/key3"},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		Err:  fmt.Errorf("err"),
	})

	// wait goroutine
	time.Sleep(400 * time.Millisecond)

	listener.mutex.Lock()
	nodes := listener.nodes
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 4, listener.invokes)
	assert.Equal(t, []byte{1, 1, 2}, nodes["/test/discovery1/key2"])
	listener.mutex.Unlock()

	d.Close()
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}

func TestDiscovery_Discovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	listener := NewMockListener(ctrl)
	d := &discovery{prefix: "/test", repo: repo, listener: listener, logger: logger.GetLogger("Coordinator", "Test")}

	// case 1: list err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := d.Discovery(true)
	assert.Error(t, err)

	// case 2: no data
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	eventCh := make(chan *state.Event)
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), false).Return(eventCh)
	err = d.Discovery(true)
	assert.NoError(t, err)
	close(eventCh)

	// case 3: find data
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{{}, {}}, nil)
	listener.EXPECT().OnCreate(gomock.Any(), gomock.Any()).MaxTimes(2)
	eventCh = make(chan *state.Event)
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), false).Return(eventCh)
	err = d.Discovery(true)
	assert.NoError(t, err)
	close(eventCh)
}
