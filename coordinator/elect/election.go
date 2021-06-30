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

package elect

import (
	"context"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./election.go -destination=./election_mock.go -package=elect

var log = logger.GetLogger("coordinator", "Election")

// Listener represent master change callback interface
type Listener interface {
	// OnFailOver triggers master fail-over, current node become master,
	// if fail over is failure return err
	OnFailOver() error
	// OnResignation triggers master resignation
	OnResignation()
}

// Election represents master elect
type Election interface {
	// Initialize initializes election, such as master change watch
	Initialize()
	// Elect elects master, include retry elect when elect fail
	Elect()
	// Close closes master elect
	Close()
	// IsMaster returns current node if is master
	IsMaster() bool
	// GetMaster returns the current master info
	GetMaster() *models.Master
}

// election implements election interface for master elect
type election struct {
	repo     state.Repository
	isMaster *atomic.Bool
	master   atomic.Value
	node     models.Node
	ttl      int64

	listener Listener

	ctx    context.Context
	cancel context.CancelFunc

	retryCh chan int
}

// NewElection returns a new election
func NewElection(ctx context.Context, repo state.Repository, node models.Node, ttl int64, listener Listener) Election {
	c, cancel := context.WithCancel(ctx)
	return &election{
		node:     node,
		ttl:      ttl,
		isMaster: atomic.NewBool(false),
		repo:     repo,
		listener: listener,
		ctx:      c,
		cancel:   cancel,
		retryCh:  make(chan int),
	}
}

// Initialize initializes election, such as master change watch
func (e *election) Initialize() {
	// watch master change event
	watchEventChan := e.repo.Watch(e.ctx, constants.MasterPath, true)

	go func() {
		e.handleMasterChange(watchEventChan)
		log.Info("exit master change event watch loop", logger.Any("node", e.node))
	}()
}

// Elect elects master,start goroutine do elect logic
func (e *election) Elect() {
	go func() {
		// wait init
		time.Sleep(10 * time.Millisecond)
		e.elect()
		log.Warn("exit master elect loop", logger.Any("node", e.node))
	}()
}

// IsMaster returns current node if is master
func (e *election) IsMaster() bool {
	return e.isMaster.Load()
}

// GetMaster returns the current master
func (e *election) GetMaster() *models.Master {
	m := e.master.Load()
	master, ok := m.(*models.Master)
	if ok {
		return master
	}
	return nil
}

// elect elects master, start elect loop for retry when failure
func (e *election) elect() {
	for {
		if e.ctx.Err() != nil {
			log.Error("context canceled, exit elect loop")
			return
		}
		log.Info("try elect master", logger.Any("node", e.node))

		master := models.Master{Node: e.node, ElectTime: timeutil.Now()}
		masterBytes := encoding.JSONMarshal(master)
		result, _, err := e.repo.Elect(e.ctx, constants.MasterPath, masterBytes, e.ttl)

		if err != nil {
			log.Warn("got an error when master elect, sleep 500ms then retry",
				logger.Error(err), logger.Any("node", e.node))
			// sleep, then try again
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if result {
			log.Info("finished election, i'm master now", logger.Any("self", e.node))
		} else {
			log.Info("finished election, i'm follower now", logger.Any("self", e.node))
		}
		select {
		case <-e.ctx.Done():
			return
		case <-e.retryCh:
			// wait retry signal
		}
	}
}

// Close closes master elect
func (e *election) Close() {
	e.resign()
	e.cancel()
}

// resign resigns master role, delete master elect node
func (e *election) resign() {
	if e.isMaster.Load() {
		log.Info("do master resign because current node is master")
		if err := e.repo.Delete(e.ctx, constants.MasterPath); err != nil {
			log.Error("delete master path failed", logger.Error(err))
		}
		e.isMaster.Store(false)
		e.master.Store(&models.Master{}) // store empty master
	}
}

// handlerMasterChange handles the event of master change,
// if master node is deleted, retry elect master
func (e *election) handleMasterChange(eventChan state.WatchEventChan) {
	for event := range eventChan {
		e.handleEvent(event)
	}
}

func (e *election) handleEvent(event *state.Event) {
	if event.Err != nil {
		log.Error("get error master change event", logger.Error(event.Err))
		return
	}
	log.Info("receive master change event", logger.String("type", event.Type.String()))
	switch event.Type {
	case state.EventTypeDelete:
		log.Info("master node lost, retry elect new master")
		if e.isMaster.Load() {
			// current node is master, do resignation when master delete is deleted
			log.Info("current node is master, do resign when master node is deleted")
			e.listener.OnResignation()
		}
		e.reElect()
	case state.EventTypeAll:
		fallthrough
	case state.EventTypeModify:
		// check the value is
		for _, kv := range event.KeyValues {
			master := models.Master{}
			if err := encoding.JSONUnmarshal(kv.Value, &master); err != nil {
				log.Error("unmarshal master value error", logger.Error(err))
			}
			log.Info("current master is", logger.Any("master", master))
			// check current node if is master
			if master.Node.IP == e.node.IP && master.Node.Port == e.node.Port {
				// current node become master
				if err := e.listener.OnFailOver(); err != nil {
					e.reElect()
					log.Error("master fail over", logger.Error(err))
					return
				}
				e.isMaster.Store(true)
			}
			// cache master info
			e.master.Store(&master)
		}
	}
}

// reElect re-elects master
func (e *election) reElect() {
	e.resign()
	// notify try elect master
	e.retryCh <- 1
}
