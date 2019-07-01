package elect

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/zap"
)

// Election defines the info for the election
type Election struct {
	isCurrentMaster bool
	node            discovery.Node
	key             string
	ttl             int64
	mutex           sync.Mutex
}

// NewElection returns a new election on a given key,value and heartbeat ttl
func NewElection(node discovery.Node, key string, ttl int64) *Election {
	return &Election{node: node, key: key, ttl: ttl}
}
func (e *Election) Elect(ctx context.Context) (bool, error) {
	success, err := e.tryElect(ctx)
	if err != nil {
		return false, err
	}
	// watch the change of the master
	err = e.watchMasterChange(ctx)
	if err != nil {
		log := logger.GetLogger()
		log.Error("try to watch master error ", zap.Error(err))
	}
	return success, err
}

//Resign deletes the master node if it is the master
func (e *Election) Resign(ctx context.Context) (err error) {
	log := logger.GetLogger()
	nodeBytes, err := json.Marshal(e.node)
	if err != nil {
		return err
	}
	success, err := state.GetRepo().DeleteWithValue(ctx, e.key, nodeBytes)
	if err != nil {
		return err
	}
	e.updateMasterFlag(false)
	log.Info("the master resign %s", zap.Bool("result", success))
	// close the chan of master change
	return nil
}

// IsMaster returns the flag current node is master.if false returns the current master
func (e *Election) IsMaster() (bool, *discovery.Node, error) {
	e.mutex.Lock()
	if e.isCurrentMaster {
		return true, nil, nil
	}
	e.mutex.Unlock()
	bytes, err := state.GetRepo().Get(context.TODO(), e.key)
	if err != nil {
		return false, nil, err
	}
	node := &discovery.Node{}
	err = json.Unmarshal(bytes, node)
	if err != nil {
		return false, nil, err
	}
	return false, node, nil
}

// watchMasterChange watches the changes of the master
func (e *Election) watchMasterChange(ctx context.Context) error {
	watchEventChan, err := state.GetRepo().Watch(ctx, e.key)
	if err == nil {
		go e.handlerMasterChange(ctx, watchEventChan)
	}
	return err
}

// handlerMasterChange handles the change of master.if the type is delete,it
// will try to elect master
func (e *Election) handlerMasterChange(ctx context.Context, eventChan state.WatchEventChan) {
	for event := range eventChan {
		switch event.Type {
		case state.EventTypeDelete:
			log := logger.GetLogger()
			success, _ := e.tryElect(ctx)
			log.Info("current node retries to registers as master", zap.Bool("result", success))
		case state.EventTypeModify:
			// check the value is
			node := &discovery.Node{}
			_ = json.Unmarshal(event.Value, node)
			e.updateMasterFlag(node.IP == e.node.IP && node.Port == e.node.Port)
		}
	}
}

// tryElect tries to elect as master.if the operation PutIfNotExist succeed,
// it will return success
func (e *Election) tryElect(ctx context.Context) (bool, error) {
	repo := state.GetRepo()
	bytes, err := json.Marshal(e.node)
	if err != nil {
		return false, err
	}
	success, _, err := repo.PutIfNotExist(ctx, e.key, bytes, e.ttl)
	e.updateMasterFlag(success)
	if err != nil {
		log := logger.GetLogger()
		log.Error("try to register master error ", zap.Error(err))
	}
	return success, err
}

// updateMasterFlag updates the isCurrentMaster safely
func (e *Election) updateMasterFlag(isMaster bool) {
	e.mutex.Lock()
	e.isCurrentMaster = isMaster
	e.mutex.Unlock()
}
