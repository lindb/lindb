package elect

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/timeutil"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// masterPath represents master elect path
const masterPath = "/master/node"

// Listener represent master change callback interface
type Listener interface {
	// OnFailOver triggers master fail-over, current node become master
	OnFailOver()
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
}

// election implements election interface for master elect
type election struct {
	repo     state.Repository
	isMaster *atomic.Bool
	node     models.Node
	ttl      int64

	listener Listener

	ctx    context.Context
	cancel context.CancelFunc

	retryCh chan int

	log *logger.Logger
}

// NewElection returns a new election
func NewElection(repo state.Repository, node models.Node, ttl int64, listener Listener) Election {
	ctx, cancel := context.WithCancel(context.Background())
	return &election{
		node:     node,
		ttl:      ttl,
		isMaster: atomic.NewBool(false),
		repo:     repo,
		listener: listener,
		ctx:      ctx,
		cancel:   cancel,
		retryCh:  make(chan int),
		log:      logger.GetLogger("coordinator/elect"),
	}
}

// Initialize initializes election, such as master change watch
func (e *election) Initialize() {
	// watch master change event
	watchEventChan := e.repo.Watch(e.ctx, masterPath)

	go func() {
		e.handlerMasterChange(watchEventChan)
		e.log.Info("exit master change event watch loop", logger.Any("node", e.node))
	}()
}

// Elect elects master,start goroutine do elect logic
func (e *election) Elect() {
	go func() {
		e.elect()
		e.log.Warn("exit master elect loop", zap.Any("node", e.node))
	}()
}

// IsMaster returns current node if is master
func (e *election) IsMaster() bool {
	return e.isMaster.Load()
}

// elect elects master, start elect loop for retry when failure
func (e *election) elect() {
	log := e.log
	for {
		if e.ctx.Err() != nil {
			log.Error("context canceled, exit elect loop")
			return
		}
		log.Info("starting try elect master", zap.Any("node", e.node))

		master := models.Master{Node: e.node, ElectTime: timeutil.Now()}
		masterBytes, err := json.Marshal(master)
		var result bool
		if err == nil {
			result, _, err = e.repo.PutIfNotExist(e.ctx, masterPath, masterBytes, e.ttl)
		}
		if err != nil {
			log.Warn("got an error when master elect, sleep 500ms then retry",
				zap.Error(err), zap.Any("node", e.node))
			// sleep, then try again
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if result {
			log.Info("become master", zap.Any("node", e.node))
		}
		log.Info("finish master elect....")

		// wait retry signal
		<-e.retryCh
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
		if err := e.repo.Delete(e.ctx, masterPath); err != nil {
			e.log.Error("delete master path failed", zap.Error(err))
		}
		e.isMaster.Store(false)
	}
}

// handlerMasterChange handles the event of master change,
// if master node is deleted, retry elect master
func (e *election) handlerMasterChange(eventChan state.WatchEventChan) {
	log := e.log
	for event := range eventChan {
		if event.Err != nil {
			log.Error("get error master change event", zap.Error(event.Err))
			continue
		}
		switch event.Type {
		case state.EventTypeDelete:
			log.Info("master node lost, retry elect new master")
			if e.isMaster.Load() {
				// current node is master, do resignation when master delete is deleted
				log.Info("current node is master, do resig when master node is deleted")
				e.listener.OnResignation()
			}
			e.resign()
			// notify try elect master
			e.retryCh <- 1
		case state.EventTypeAll:
			fallthrough
		case state.EventTypeModify:
			// check the value is
			for _, kv := range event.KeyValues {
				master := models.Master{}
				if err := json.Unmarshal(kv.Value, &master); err != nil {
					e.log.Error("unmarshal master value error", zap.Error(err))
				} else if master.Node.IP == e.node.IP && master.Node.Port == e.node.Port {
					e.isMaster.Store(true)
					// current node become master
					e.listener.OnFailOver()
				}
			}
		}
	}
}
