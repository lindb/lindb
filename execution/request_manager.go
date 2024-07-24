package execution

import (
	"sync"

	"github.com/lindb/lindb/execution/model"
)

type RequestManager interface {
	CompleteRequet(id model.RequestID, err error)
}

type requestManager struct {
	requests map[model.RequestID]any
	lock     sync.RWMutex
}

func NewRequestManager() RequestManager {
	return &requestManager{
		requests: make(map[model.RequestID]any),
	}
}

func (mgr *requestManager) CompleteRequet(id model.RequestID, err error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	delete(mgr.requests, id)
}
