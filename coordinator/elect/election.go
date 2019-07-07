package elect

import (
	"context"
	"encoding/json"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// Election defines the info for the election
type Election struct {
	repo            state.Repository
	isCurrentMaster *atomic.Bool
	key             string
	node            models.Node
	ttl             int64
	log             *zap.Logger
}

// NewElection returns a new election on a given key,value and heartbeat ttl
func NewElection(node models.Node, key string, ttl int64) *Election {
	return &Election{
		key:             key,
		node:            node,
		ttl:             ttl,
		isCurrentMaster: atomic.NewBool(false),
		repo:            state.GetRepo(),
		log:             logger.GetLogger(),
	}
}

// Elect returns the result of register master
func (e *Election) Elect(ctx context.Context) (bool, error) {
	success, err := e.tryElect(ctx)
	if err != nil {
		return false, err
	}
	// watch the change of the master
	e.watchMasterChange(ctx)
	return success, err
}

//Resign deletes the master node if it is the master
func (e *Election) Resign(ctx context.Context) {
	e.updateMasterFlag(false)
	if !e.IsMaster() {
		e.log.Info("the node is not master ")
		return
	}
	e.log.Info("the node is master, resign")
	if err := e.repo.Delete(ctx, e.key); err != nil {
		e.log.Error("resign failed ", zap.Error(err))
	}
}

// IsMaster returns the flag current node is master.if false returns the current master
func (e *Election) IsMaster() bool {
	return e.isCurrentMaster.Load()
}

// watchMasterChange watches the changes of the master
func (e *Election) watchMasterChange(ctx context.Context) {
	watchEventChan := e.repo.Watch(ctx, e.key)
	go e.handlerMasterChange(ctx, watchEventChan)
}

// handlerMasterChange handles the change of master.if the type is delete,it
// will try to elect master
func (e *Election) handlerMasterChange(ctx context.Context, eventChan state.WatchEventChan) {
	for event := range eventChan {
		if event.Err != nil {
			continue
		}
		switch event.Type {
		case state.EventTypeDelete:
			success, err := e.tryElect(ctx)
			if err != nil {
				e.log.Error("register master failed when master node is deleted", zap.Error(err))
			}
			e.log.Info("current node retries to registers as master", zap.Bool("result", success))
		case state.EventTypeAll:
			fallthrough
		case state.EventTypeModify:
			// check the value is
			for _, kv := range event.KeyValues {
				node := &models.Node{}
				if err := json.Unmarshal(kv.Value, node); err != nil {
					e.log.Error("deserialize node values error", zap.Error(err))
				} else {
					e.updateMasterFlag(node.IP == e.node.IP && node.Port == e.node.Port)
				}
			}
		}
	}
}

// tryElect tries to elect as master.if the operation PutIfNotExist succeed,
// it will return success
func (e *Election) tryElect(ctx context.Context) (bool, error) {
	bytes, err := json.Marshal(e.node)
	if err != nil {
		return false, err
	}
	// TODO notice the changes of master
	success, _, err := e.repo.PutIfNotExist(ctx, e.key, bytes, e.ttl)
	e.updateMasterFlag(success)
	if err != nil {
		e.log.Error("try to register master error ", zap.Error(err))
	}
	return success, err
}

// updateMasterFlag updates the the flag whether current is master
func (e *Election) updateMasterFlag(isMaster bool) {
	e.isCurrentMaster.Store(isMaster)
}
